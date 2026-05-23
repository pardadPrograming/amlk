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
	results := sortMatchResults(scorePropertyCandidates(request, candidatePropertyFiles(index, request)))
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

	files, err := s.store.ListPropertyFilesForOwner(ctx, businessID, userID)
	if err != nil {
		return propertySearchIndex{}, err
	}
	index = buildPropertySearchIndex(files, now.Add(s.ttl))

	s.mu.Lock()
	s.indexes[key] = index
	s.mu.Unlock()
	return index, nil
}

func buildPropertySearchIndex(files []domain.PropertyFile, expiresAt time.Time) propertySearchIndex {
	index := propertySearchIndex{
		expiresAt:       expiresAt,
		files:           make([]domain.PropertyFile, 0, len(files)),
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

func scorePropertyCandidates(request domain.ContactRequest, files []domain.PropertyFile) []domain.PropertyMatchResult {
	if len(files) < 128 {
		return scorePropertyChunk(request, files)
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
			out <- scorePropertyChunk(request, chunk)
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

func scorePropertyChunk(request domain.ContactRequest, files []domain.PropertyFile) []domain.PropertyMatchResult {
	results := make([]domain.PropertyMatchResult, 0, len(files))
	for _, file := range files {
		match := matchProperty(request, file)
		if match.Score == 0 || match.Tier == "" {
			continue
		}
		results = append(results, match)
	}
	return results
}
