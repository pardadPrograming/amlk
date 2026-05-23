package searchengine

import (
	"context"
	"errors"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
)

type OptimizedMatchingService struct {
	store        repository.Store
	ttl          time.Duration
	mu           sync.RWMutex
	indexes      map[string]propertySearchIndex
	resultCaches map[string]matchResultCache
}

type MatchPage struct {
	Items  []domain.PropertyMatchResult
	Total  int
	Limit  int
	Offset int
}

type matchResultCache struct {
	expiresAt time.Time
	results   []domain.PropertyMatchResult
}

type propertySearchIndex struct {
	expiresAt       time.Time
	files           []domain.PropertyFile
	accessByFileID  map[string][]domain.PropertyMatchAccess
	areaIDs         map[string][]int
	streetIDs       map[string][]int
	neighborhoodIDs map[string][]int
	areaNames       map[string][]int
	streetNames     map[string][]int
	neighborhoods   map[string][]int
}

func NewOptimizedMatchingService(store repository.Store, ttl time.Duration) *OptimizedMatchingService {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &OptimizedMatchingService{
		store:        store,
		ttl:          ttl,
		indexes:      map[string]propertySearchIndex{},
		resultCaches: map[string]matchResultCache{},
	}
}

func (s *OptimizedMatchingService) RequestMatches(ctx context.Context, userID, businessID, contactID, requestID string, limit, offset int) (MatchPage, error) {
	limit, offset = normalizePage(limit, offset)
	cacheKey := businessID + ":" + userID + ":" + contactID + ":" + requestID
	if cached, ok := s.cachedResults(cacheKey); ok {
		items := paginateMatches(cached, limit, offset)
		return MatchPage{Items: items, Total: len(cached), Limit: limit, Offset: offset}, nil
	}
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return MatchPage{}, err
	}
	contact, err := s.store.GetContact(ctx, businessID, contactID)
	if err != nil {
		return MatchPage{}, err
	}
	var request domain.ContactRequest
	for _, item := range contact.Requests {
		if item.ID == requestID {
			request = item
			break
		}
	}
	if request.ID == "" {
		return MatchPage{}, repository.ErrNotFound
	}
	index, err := s.propertyIndex(ctx, businessID, userID)
	if err != nil {
		return MatchPage{}, err
	}
	results := sortMatchResults(scorePropertyCandidates(request, candidatePropertyFiles(index, request), index.accessByFileID))
	s.storeCachedResults(cacheKey, results)
	items := paginateMatches(results, limit, offset)
	return MatchPage{Items: items, Total: len(results), Limit: limit, Offset: offset}, nil
}

func (s *OptimizedMatchingService) Invalidate(businessID, userID string) {
	prefix := businessID + ":" + userID + ":"
	indexKey := businessID + ":" + userID
	s.mu.Lock()
	delete(s.indexes, indexKey)
	for key := range s.resultCaches {
		if strings.HasPrefix(key, prefix) {
			delete(s.resultCaches, key)
		}
	}
	s.mu.Unlock()
}

func (s *OptimizedMatchingService) cachedResults(key string) ([]domain.PropertyMatchResult, bool) {
	now := time.Now()
	s.mu.RLock()
	cache, ok := s.resultCaches[key]
	s.mu.RUnlock()
	if !ok || now.After(cache.expiresAt) {
		return nil, false
	}
	return cache.results, true
}

func (s *OptimizedMatchingService) storeCachedResults(key string, results []domain.PropertyMatchResult) {
	s.mu.Lock()
	s.resultCaches[key] = matchResultCache{expiresAt: time.Now().Add(s.ttl), results: results}
	s.mu.Unlock()
}

func normalizePage(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func paginateMatches(results []domain.PropertyMatchResult, limit, offset int) []domain.PropertyMatchResult {
	if offset >= len(results) {
		return []domain.PropertyMatchResult{}
	}
	end := offset + limit
	if end > len(results) {
		end = len(results)
	}
	return results[offset:end]
}

func sortMatchResults(results []domain.PropertyMatchResult) []domain.PropertyMatchResult {
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].PropertyFile.UpdatedAt.After(results[j].PropertyFile.UpdatedAt)
		}
		return results[i].Score > results[j].Score
	})
	return results
}

func (s *OptimizedMatchingService) authorize(ctx context.Context, userID, businessID string) (domain.BusinessMember, error) {
	member, err := s.store.GetMemberByUser(ctx, businessID, userID)
	if err != nil || member.Status != domain.MemberActive {
		return domain.BusinessMember{}, errors.New("دسترسی به کسب و کار وجود ندارد")
	}
	return member, nil
}

func (s *OptimizedMatchingService) propertyIndex(ctx context.Context, businessID, userID string) (propertySearchIndex, error) {
	key := businessID + ":" + userID
	now := time.Now()

	s.mu.RLock()
	index, ok := s.indexes[key]
	s.mu.RUnlock()
	if ok && now.Before(index.expiresAt) {
		return index, nil
	}

	files, accessByFileID, err := s.accessiblePropertyFiles(ctx, businessID, userID)
	if err != nil {
		return propertySearchIndex{}, err
	}
	index = buildPropertySearchIndex(files, accessByFileID, now.Add(s.ttl))

	s.mu.Lock()
	s.indexes[key] = index
	s.mu.Unlock()
	return index, nil
}

func (s *OptimizedMatchingService) accessiblePropertyFiles(ctx context.Context, businessID, userID string) ([]domain.PropertyFile, map[string][]domain.PropertyMatchAccess, error) {
	user, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	businesses, err := s.store.ListBusinessesForUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	businessIDs := make([]string, 0, 1)
	for _, business := range businesses {
		if business.ID == businessID {
			businessIDs = append(businessIDs, business.ID)
			break
		}
	}
	channels, err := s.store.ListChannelsForUser(ctx, userID, []string{user.Phone}, businessIDs)
	if err != nil {
		return nil, nil, err
	}
	accessibleVaults := map[string]domain.ChannelVault{}
	accessibleVaultIDs := []string{}
	for _, channel := range channels {
		if channel.VaultID == "" {
			continue
		}
		if channel.BusinessID != "" && channel.BusinessID != businessID {
			continue
		}
		vault, err := s.store.GetChannelVault(ctx, channel.VaultID)
		if err != nil {
			continue
		}
		if _, exists := accessibleVaults[vault.ID]; !exists {
			accessibleVaultIDs = append(accessibleVaultIDs, vault.ID)
		}
		accessibleVaults[vault.ID] = vault
	}
	if member, err := s.store.GetMemberByUser(ctx, businessID, userID); err == nil &&
		(member.Role == domain.RoleOwner || member.Role == domain.RoleManager || domain.HasPermission(member, domain.PermBusinessUpdate)) {
		if vaults, err := s.store.ListBusinessVaults(ctx, businessID); err == nil {
			for _, vault := range vaults {
				if vault.ID == "" {
					continue
				}
				if _, exists := accessibleVaults[vault.ID]; !exists {
					accessibleVaultIDs = append(accessibleVaultIDs, vault.ID)
				}
				accessibleVaults[vault.ID] = vault
			}
		}
	}

	allFiles, err := s.store.ListPropertyFilesForAccess(ctx, businessID, userID, accessibleVaultIDs)
	if err != nil {
		return nil, nil, err
	}
	files := make([]domain.PropertyFile, 0, len(allFiles))
	accessByFileID := map[string][]domain.PropertyMatchAccess{}
	for _, file := range allFiles {
		access := propertyMatchAccessForFile(file, userID, accessibleVaults)
		if len(access) == 0 {
			continue
		}
		files = append(files, file)
		accessByFileID[file.ID] = access
	}
	return files, accessByFileID, nil
}

func propertyMatchAccessForFile(file domain.PropertyFile, userID string, accessibleVaults map[string]domain.ChannelVault) []domain.PropertyMatchAccess {
	result := []domain.PropertyMatchAccess{}
	if file.OwnerUserID == userID {
		source := "own"
		if file.IsPartnershipCopy || file.SharedFromFileID != "" {
			source = "collaboration"
		}
		result = append(result, domain.PropertyMatchAccess{
			Source:        source,
			Collaboration: source == "collaboration",
		})
	}
	commissionByVault := map[string]float64{}
	for _, placement := range file.VaultPlacements {
		if placement.VaultID != "" {
			commissionByVault[placement.VaultID] = placement.CommissionPercent
		}
	}
	for _, vaultID := range file.VaultIDs {
		vault, ok := accessibleVaults[vaultID]
		if !ok {
			continue
		}
		if file.OwnerUserID == userID || file.IsPartnershipCopy {
			continue
		}
		title := strings.TrimSpace(vault.Title)
		if title == "" {
			title = "صندوقچه"
		}
		result = append(result, domain.PropertyMatchAccess{
			Source:            "vault",
			VaultID:           vault.ID,
			VaultTitle:        title,
			CommissionPercent: commissionByVault[vaultID],
			Collaboration:     false,
		})
	}
	return result
}

func buildPropertySearchIndex(files []domain.PropertyFile, accessByFileID map[string][]domain.PropertyMatchAccess, expiresAt time.Time) propertySearchIndex {
	index := propertySearchIndex{
		expiresAt:       expiresAt,
		files:           make([]domain.PropertyFile, 0, len(files)),
		accessByFileID:  accessByFileID,
		areaIDs:         map[string][]int{},
		streetIDs:       map[string][]int{},
		neighborhoodIDs: map[string][]int{},
		areaNames:       map[string][]int{},
		streetNames:     map[string][]int{},
		neighborhoods:   map[string][]int{},
	}
	for _, file := range files {
		idx := len(index.files)
		index.files = append(index.files, file)
		for _, address := range file.Addresses {
			addIndexValue(index.areaIDs, address.AreaID, idx)
			addIndexValue(index.streetIDs, address.StreetID, idx)
			addIndexValue(index.neighborhoodIDs, address.NeighborhoodID, idx)
			addIndexValue(index.areaNames, NormalizePersianName(address.AreaName), idx)
			addIndexValue(index.streetNames, NormalizePersianName(address.StreetName), idx)
			addIndexValue(index.neighborhoods, NormalizePersianName(address.NeighborhoodName), idx)
		}
	}
	return index
}

func addIndexValue(values map[string][]int, key string, idx int) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}
	values[key] = append(values[key], idx)
}

func candidatePropertyFiles(index propertySearchIndex, request domain.ContactRequest) []domain.PropertyFile {
	if len(request.Locations) == 0 {
		return index.files
	}
	seen := map[int]struct{}{}
	for _, location := range request.Locations {
		for _, idx := range candidateIndexesForLocation(index, location) {
			seen[idx] = struct{}{}
		}
	}
	if len(seen) == 0 {
		return index.files
	}
	candidates := make([]domain.PropertyFile, 0, len(seen))
	for idx := range seen {
		if idx >= 0 && idx < len(index.files) {
			candidates = append(candidates, index.files[idx])
		}
	}
	return candidates
}

func candidateIndexesForLocation(index propertySearchIndex, location domain.ContactRequestLocation) []int {
	name := NormalizePersianName(location.Name)
	switch location.Level {
	case "neighborhood":
		if location.ID != "" {
			return index.neighborhoodIDs[location.ID]
		}
		return index.neighborhoods[name]
	case "street":
		if location.ID != "" {
			return index.streetIDs[location.ID]
		}
		if location.StreetID != "" {
			return index.streetIDs[location.StreetID]
		}
		return index.streetNames[name]
	default:
		if location.ID != "" {
			return index.areaIDs[location.ID]
		}
		if location.AreaID != "" {
			return index.areaIDs[location.AreaID]
		}
		return index.areaNames[name]
	}
}

func scorePropertyCandidates(request domain.ContactRequest, files []domain.PropertyFile, accessByFileID map[string][]domain.PropertyMatchAccess) []domain.PropertyMatchResult {
	if len(files) < 128 {
		return scorePropertyChunk(request, files, accessByFileID)
	}
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	if workers > len(files) {
		workers = len(files)
	}
	chunkSize := (len(files) + workers - 1) / workers
	out := make(chan []domain.PropertyMatchResult, workers)
	var wg sync.WaitGroup
	for start := 0; start < len(files); start += chunkSize {
		end := start + chunkSize
		if end > len(files) {
			end = len(files)
		}
		chunk := files[start:end]
		wg.Add(1)
		go func() {
			defer wg.Done()
			out <- scorePropertyChunk(request, chunk, accessByFileID)
		}()
	}
	wg.Wait()
	close(out)

	results := make([]domain.PropertyMatchResult, 0, len(files))
	for items := range out {
		results = append(results, items...)
	}
	return results
}

func scorePropertyChunk(request domain.ContactRequest, files []domain.PropertyFile, accessByFileID map[string][]domain.PropertyMatchAccess) []domain.PropertyMatchResult {
	results := make([]domain.PropertyMatchResult, 0, len(files))
	for _, file := range files {
		match := matchProperty(request, file)
		if match.Score == 0 || match.Tier == "" {
			continue
		}
		match.Access = accessByFileID[file.ID]
		results = append(results, match)
	}
	return results
}
