package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/support"
)

type ContactService struct {
	store repository.Store
}

func NewContactService(store repository.Store) *ContactService {
	return &ContactService{store: store}
}

type ContactFilter struct {
	Query string
	Tag   string
	Phone string
}

type ProfileContactCategoryResult struct {
	Contact        domain.Contact `json:"contact"`
	AutoConsultant bool           `json:"autoConsultant"`
	Existing       bool           `json:"existing"`
}

func (s *ContactService) SystemTags() []string {
	return domain.SystemContactTags()
}

func (s *ContactService) List(ctx context.Context, userID, businessID string, filter ContactFilter) ([]domain.Contact, error) {
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return nil, err
	}
	items, err := s.store.ListContacts(ctx, businessID)
	if err != nil {
		return nil, err
	}
	filter.Query = normalizeContactText(filter.Query)
	filter.Tag = strings.TrimSpace(filter.Tag)
	filter.Phone = normalizeContactText(filter.Phone)
	result := []domain.Contact{}
	for _, item := range items {
		if filter.Tag != "" && !contactHasTag(item, filter.Tag) {
			continue
		}
		if filter.Phone != "" && !contactHasPhone(item, filter.Phone) {
			continue
		}
		if filter.Query != "" && !contactMatchesQuery(item, filter.Query) {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *ContactService) Create(ctx context.Context, userID, businessID string, input domain.Contact) (domain.Contact, error) {
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return domain.Contact{}, err
	}
	contact, err := normalizeContact(input)
	if err != nil {
		return domain.Contact{}, err
	}
	if err := rejectManualDoneRequests(contact.Requests, nil); err != nil {
		return domain.Contact{}, err
	}
	initializeRequestHistories(contact.Requests, userID)
	contact.BusinessID = businessID
	contact.CreatedByID = userID
	return s.store.CreateContact(ctx, contact)
}

func (s *ContactService) Update(ctx context.Context, userID, businessID, contactID string, input domain.Contact) (domain.Contact, error) {
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return domain.Contact{}, err
	}
	existing, err := s.store.GetContact(ctx, businessID, contactID)
	if err != nil {
		return domain.Contact{}, err
	}
	contact, err := normalizeContact(input)
	if err != nil {
		return domain.Contact{}, err
	}
	if err := rejectManualDoneRequests(contact.Requests, existing.Requests); err != nil {
		return domain.Contact{}, err
	}
	if err := mergeRequestHistories(contact.Requests, existing.Requests, userID); err != nil {
		return domain.Contact{}, err
	}
	contact.ID = contactID
	contact.BusinessID = businessID
	contact.CreatedByID = existing.CreatedByID
	return s.store.UpdateContact(ctx, contact)
}

func (s *ContactService) AddProfileToCategories(ctx context.Context, userID, businessID, phone, displayName string, tags []string) (ProfileContactCategoryResult, error) {
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return ProfileContactCategoryResult{}, err
	}
	normalized, err := NormalizePhone(phone)
	if err != nil {
		return ProfileContactCategoryResult{}, err
	}
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		displayName = normalized
	}
	tags = normalizeTags(tags)
	autoConsultant := s.isRealEstateUser(ctx, normalized)
	if autoConsultant {
		tags = mergeTags(tags, []string{domain.ContactTagRealEstateConsultant})
	}
	if len(tags) == 0 {
		return ProfileContactCategoryResult{}, errors.New("حداقل یک دسته‌بندی مخاطب باید انتخاب شود")
	}
	items, err := s.store.ListContacts(ctx, businessID)
	if err != nil {
		return ProfileContactCategoryResult{}, err
	}
	for _, existing := range items {
		if !contactHasExactPhone(existing, normalized) {
			continue
		}
		existing.DisplayName = firstNonEmpty(existing.DisplayName, displayName)
		existing.Tags = mergeTags(existing.Tags, tags)
		if !contactHasExactPhone(existing, normalized) {
			existing.Phones = append(existing.Phones, domain.ContactPhone{Label: "موبایل", Value: normalized})
		}
		updated, err := s.Update(ctx, userID, businessID, existing.ID, existing)
		if err != nil {
			return ProfileContactCategoryResult{}, err
		}
		return ProfileContactCategoryResult{Contact: updated, AutoConsultant: autoConsultant, Existing: true}, nil
	}
	contact, err := s.Create(ctx, userID, businessID, domain.Contact{
		DisplayName: displayName,
		Phones: []domain.ContactPhone{{
			Label: "موبایل",
			Value: normalized,
		}},
		Tags: tags,
	})
	if err != nil {
		return ProfileContactCategoryResult{}, err
	}
	return ProfileContactCategoryResult{Contact: contact, AutoConsultant: autoConsultant}, nil
}

func (s *ContactService) authorize(ctx context.Context, userID, businessID string) (domain.BusinessMember, error) {
	member, err := s.store.GetMemberByUser(ctx, businessID, userID)
	if err != nil || member.Status != domain.MemberActive {
		return domain.BusinessMember{}, errors.New("دسترسی به کسب‌وکار وجود ندارد")
	}
	return member, nil
}

func (s *ContactService) isRealEstateUser(ctx context.Context, phone string) bool {
	user, err := s.store.GetUserByPhone(ctx, phone)
	if err != nil {
		return false
	}
	businesses, err := s.store.ListBusinessesForUser(ctx, user.ID)
	return err == nil && len(businesses) > 0
}

func normalizeContact(input domain.Contact) (domain.Contact, error) {
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Company = strings.TrimSpace(input.Company)
	input.Note = strings.TrimSpace(input.Note)
	if input.DisplayName == "" {
		input.DisplayName = strings.TrimSpace(input.FirstName + " " + input.LastName)
	}
	if input.DisplayName == "" && input.Company != "" {
		input.DisplayName = input.Company
	}
	if input.DisplayName == "" {
		return domain.Contact{}, errors.New("نام مخاطب الزامی است")
	}
	phones := []domain.ContactPhone{}
	for _, phone := range input.Phones {
		phone.Label = strings.TrimSpace(phone.Label)
		phone.Value = strings.TrimSpace(phone.Value)
		if phone.Value == "" {
			continue
		}
		if phone.Label == "" {
			phone.Label = "موبایل"
		}
		phones = append(phones, phone)
	}
	if len(phones) == 0 {
		return domain.Contact{}, errors.New("حداقل یک شماره تلفن الزامی است")
	}
	input.Phones = phones
	input.Tags = normalizeTags(input.Tags)
	input.Requests = normalizeContactRequests(input.Requests)
	return input, nil
}

func normalizeContactRequests(requests []domain.ContactRequest) []domain.ContactRequest {
	result := make([]domain.ContactRequest, 0, len(requests))
	for _, req := range requests {
		req.Title = strings.TrimSpace(req.Title)
		req.Type = normalizeRequestType(req.Type)
		req.Status = normalizeRequestStatus(req.Status)
		req.Note = strings.TrimSpace(req.Note)
		if req.Title == "" {
			req.Title = requestTypeLabel(req.Type)
		}
		if req.Type == "" && req.Title == "" {
			continue
		}
		if req.MinAreaM2 < 0 {
			req.MinAreaM2 = 0
		}
		if req.MaxAgeYears < 0 {
			req.MaxAgeYears = 0
		}
		if req.LandMinAreaM2 < 0 {
			req.LandMinAreaM2 = 0
		}
		if req.BuildingMinAreaM2 < 0 {
			req.BuildingMinAreaM2 = 0
		}
		if req.PermitFloorsMin < 0 {
			req.PermitFloorsMin = 0
		}
		req = normalizeRequestFinancials(req)
		req.Locations = normalizeRequestLocations(req.Locations)
		req.FloorRules = normalizeRequestFloorRules(req.FloorRules)
		req.OptionFilters = normalizeRequestOptionFilters(req.OptionFilters)
		req.BooleanFilters = normalizeRequestBooleanFilters(req.BooleanFilters)
		req.NumberFilters = normalizeRequestNumberFilters(req.NumberFilters)
		result = append(result, req)
	}
	return result
}

func rejectManualDoneRequests(requests, existing []domain.ContactRequest) error {
	existingStatus := map[string]string{}
	for _, req := range existing {
		if req.ID != "" {
			existingStatus[req.ID] = req.Status
		}
	}
	for _, req := range requests {
		if req.Status != "done" {
			continue
		}
		if req.ID != "" && existingStatus[req.ID] == "done" {
			continue
		}
		return errors.New("وضعیت انجام شده فقط بعد از ثبت قرارداد قابل اعمال است")
	}
	return nil
}

func initializeRequestHistories(requests []domain.ContactRequest, userID string) {
	now := time.Now().UTC()
	for i := range requests {
		requests[i].ChangeDescription = ""
		if requests[i].CreatedAt.IsZero() {
			requests[i].CreatedAt = now
		}
		if len(requests[i].History) == 0 {
			requests[i].History = []domain.ContactRequestHistoryEntry{{
				ID:          support.NewID(),
				ChangedByID: userID,
				ChangedAt:   now,
				Description: "ثبت اولیه درخواست",
				Changes: []domain.ContactRequestHistoryChange{{
					Field: "request",
					To:    "created",
				}},
			}}
		}
	}
}

func mergeRequestHistories(requests []domain.ContactRequest, existing []domain.ContactRequest, userID string) error {
	existingByID := map[string]domain.ContactRequest{}
	for _, req := range existing {
		if req.ID != "" {
			existingByID[req.ID] = req
		}
	}
	now := time.Now().UTC()
	for i := range requests {
		current := &requests[i]
		if current.ID == "" {
			if strings.TrimSpace(current.ChangeDescription) == "" {
				current.ChangeDescription = "ثبت درخواست جدید"
			}
			current.History = []domain.ContactRequestHistoryEntry{{
				ID:          support.NewID(),
				ChangedByID: userID,
				ChangedAt:   now,
				Description: strings.TrimSpace(current.ChangeDescription),
				Changes: []domain.ContactRequestHistoryChange{{
					Field: "request",
					To:    "created",
				}},
			}}
			current.ChangeDescription = ""
			continue
		}
		prev, ok := existingByID[current.ID]
		if !ok {
			if strings.TrimSpace(current.ChangeDescription) == "" {
				current.ChangeDescription = "ثبت درخواست جدید"
			}
			current.History = []domain.ContactRequestHistoryEntry{{
				ID:          support.NewID(),
				ChangedByID: userID,
				ChangedAt:   now,
				Description: strings.TrimSpace(current.ChangeDescription),
				Changes: []domain.ContactRequestHistoryChange{{
					Field: "request",
					To:    "created",
				}},
			}}
			current.ChangeDescription = ""
			continue
		}
		current.CreatedAt = prev.CreatedAt
		changes := requestChanges(prev, *current)
		current.History = prev.History
		if len(changes) > 0 {
			description := strings.TrimSpace(current.ChangeDescription)
			if description == "" {
				return fmt.Errorf("توضیح تغییرات برای درخواست %s الزامی است", current.Title)
			}
			current.History = append(current.History, domain.ContactRequestHistoryEntry{
				ID:          support.NewID(),
				ChangedByID: userID,
				ChangedAt:   now,
				Description: description,
				Changes:     changes,
			})
		}
		current.ChangeDescription = ""
	}
	return nil
}

func requestChanges(prev, next domain.ContactRequest) []domain.ContactRequestHistoryChange {
	prev.ChangeDescription = ""
	next.ChangeDescription = ""
	prev.History = nil
	next.History = nil
	prev.CreatedAt = time.Time{}
	next.CreatedAt = time.Time{}
	fields := []struct {
		name string
		from interface{}
		to   interface{}
	}{
		{"title", prev.Title, next.Title},
		{"type", prev.Type, next.Type},
		{"status", prev.Status, next.Status},
		{"budgetMin", prev.BudgetMin, next.BudgetMin},
		{"budgetMax", prev.BudgetMax, next.BudgetMax},
		{"purchaseMin", prev.PurchaseMin, next.PurchaseMin},
		{"purchaseMax", prev.PurchaseMax, next.PurchaseMax},
		{"suggestedPurchaseMin", prev.SuggestedPurchaseMin, next.SuggestedPurchaseMin},
		{"suggestedPurchaseMax", prev.SuggestedPurchaseMax, next.SuggestedPurchaseMax},
		{"partnershipMin", prev.PartnershipMin, next.PartnershipMin},
		{"partnershipMax", prev.PartnershipMax, next.PartnershipMax},
		{"shareMin", prev.ShareMin, next.ShareMin},
		{"shareMax", prev.ShareMax, next.ShareMax},
		{"depositMin", prev.DepositMin, next.DepositMin},
		{"depositMax", prev.DepositMax, next.DepositMax},
		{"suggestedDepositMin", prev.SuggestedDepositMin, next.SuggestedDepositMin},
		{"suggestedDepositMax", prev.SuggestedDepositMax, next.SuggestedDepositMax},
		{"rentMin", prev.RentMin, next.RentMin},
		{"rentMax", prev.RentMax, next.RentMax},
		{"suggestedRentMin", prev.SuggestedRentMin, next.SuggestedRentMin},
		{"suggestedRentMax", prev.SuggestedRentMax, next.SuggestedRentMax},
		{"minAreaM2", prev.MinAreaM2, next.MinAreaM2},
		{"maxAgeYears", prev.MaxAgeYears, next.MaxAgeYears},
		{"convertible", prev.Convertible, next.Convertible},
		{"maxConvertibleDeposit", prev.MaxConvertibleDeposit, next.MaxConvertibleDeposit},
		{"rentWithOwner", prev.RentWithOwner, next.RentWithOwner},
		{"landMinAreaM2", prev.LandMinAreaM2, next.LandMinAreaM2},
		{"buildingMinAreaM2", prev.BuildingMinAreaM2, next.BuildingMinAreaM2},
		{"permitFloorsMin", prev.PermitFloorsMin, next.PermitFloorsMin},
		{"locations", prev.Locations, next.Locations},
		{"floorRules", prev.FloorRules, next.FloorRules},
		{"optionFilters", prev.OptionFilters, next.OptionFilters},
		{"booleanFilters", prev.BooleanFilters, next.BooleanFilters},
		{"numberFilters", prev.NumberFilters, next.NumberFilters},
		{"note", prev.Note, next.Note},
	}
	changes := []domain.ContactRequestHistoryChange{}
	for _, field := range fields {
		from := historyValue(field.from)
		to := historyValue(field.to)
		if from == to {
			continue
		}
		changes = append(changes, domain.ContactRequestHistoryChange{
			Field: field.name,
			From:  from,
			To:    to,
		})
	}
	return changes
}

func historyValue(value interface{}) string {
	body, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(body)
}

func normalizeRequestFinancials(req domain.ContactRequest) domain.ContactRequest {
	req.BudgetMin = nonNegativeInt64(req.BudgetMin)
	req.BudgetMax = nonNegativeInt64(req.BudgetMax)
	req.PurchaseMin = nonNegativeInt64(req.PurchaseMin)
	req.PurchaseMax = nonNegativeInt64(req.PurchaseMax)
	req.SuggestedPurchaseMin = nonNegativeInt64(req.SuggestedPurchaseMin)
	req.SuggestedPurchaseMax = nonNegativeInt64(req.SuggestedPurchaseMax)
	req.PartnershipMin = nonNegativeInt64(req.PartnershipMin)
	req.PartnershipMax = nonNegativeInt64(req.PartnershipMax)
	req.DepositMin = nonNegativeInt64(req.DepositMin)
	req.DepositMax = nonNegativeInt64(req.DepositMax)
	req.SuggestedDepositMin = nonNegativeInt64(req.SuggestedDepositMin)
	req.SuggestedDepositMax = nonNegativeInt64(req.SuggestedDepositMax)
	req.RentMin = nonNegativeInt64(req.RentMin)
	req.RentMax = nonNegativeInt64(req.RentMax)
	req.SuggestedRentMin = nonNegativeInt64(req.SuggestedRentMin)
	req.SuggestedRentMax = nonNegativeInt64(req.SuggestedRentMax)
	req.MaxConvertibleDeposit = nonNegativeInt64(req.MaxConvertibleDeposit)
	req.ShareMin = clampPercent(req.ShareMin)
	req.ShareMax = clampPercent(req.ShareMax)
	req.BudgetMin, req.BudgetMax = orderedInt64(req.BudgetMin, req.BudgetMax)
	req.PurchaseMin, req.PurchaseMax = orderedInt64(req.PurchaseMin, req.PurchaseMax)
	req.SuggestedPurchaseMin, req.SuggestedPurchaseMax = orderedInt64(req.SuggestedPurchaseMin, req.SuggestedPurchaseMax)
	req.PartnershipMin, req.PartnershipMax = orderedInt64(req.PartnershipMin, req.PartnershipMax)
	req.DepositMin, req.DepositMax = orderedInt64(req.DepositMin, req.DepositMax)
	req.SuggestedDepositMin, req.SuggestedDepositMax = orderedInt64(req.SuggestedDepositMin, req.SuggestedDepositMax)
	req.RentMin, req.RentMax = orderedInt64(req.RentMin, req.RentMax)
	req.SuggestedRentMin, req.SuggestedRentMax = orderedInt64(req.SuggestedRentMin, req.SuggestedRentMax)
	req.ShareMin, req.ShareMax = orderedInt(req.ShareMin, req.ShareMax)
	switch req.Type {
	case "sale":
		if req.PurchaseMin == 0 {
			req.PurchaseMin = req.BudgetMin
		}
		if req.PurchaseMax == 0 {
			req.PurchaseMax = req.BudgetMax
		}
		if req.BudgetMin == 0 {
			req.BudgetMin = req.PurchaseMin
		}
		if req.BudgetMax == 0 {
			req.BudgetMax = req.PurchaseMax
		}
	case "partnership":
		if req.PartnershipMin == 0 {
			req.PartnershipMin = req.BudgetMin
		}
		if req.PartnershipMax == 0 {
			req.PartnershipMax = req.BudgetMax
		}
		if req.BudgetMin == 0 {
			req.BudgetMin = req.PartnershipMin
		}
		if req.BudgetMax == 0 {
			req.BudgetMax = req.PartnershipMax
		}
	}
	return req
}

func nonNegativeInt64(value int64) int64 {
	if value < 0 {
		return 0
	}
	return value
}

func clampPercent(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func orderedInt64(minValue, maxValue int64) (int64, int64) {
	if maxValue > 0 && minValue > maxValue {
		return maxValue, minValue
	}
	return minValue, maxValue
}

func orderedInt(minValue, maxValue int) (int, int) {
	if maxValue > 0 && minValue > maxValue {
		return maxValue, minValue
	}
	return minValue, maxValue
}

func normalizeRequestType(value string) string {
	switch strings.TrimSpace(value) {
	case "purchase", "buy", "sale":
		return "sale"
	case "partnership":
		return "partnership"
	case "rent", "rent_lease", "lease":
		return "rent_lease"
	default:
		return strings.TrimSpace(value)
	}
}

func normalizeRequestStatus(value string) string {
	switch strings.TrimSpace(value) {
	case "done", "finished", "suspended", "active":
		return strings.TrimSpace(value)
	default:
		return "active"
	}
}

func requestTypeLabel(value string) string {
	switch value {
	case "sale":
		return "درخواست خرید"
	case "partnership":
		return "درخواست مشارکت"
	case "rent_lease":
		return "درخواست رهن و اجاره"
	default:
		return "درخواست مشتری"
	}
}

func normalizePreference(value string) string {
	if strings.TrimSpace(value) == "suggested" {
		return "suggested"
	}
	return "preferred"
}

func normalizeRequestLocations(items []domain.ContactRequestLocation) []domain.ContactRequestLocation {
	result := make([]domain.ContactRequestLocation, 0, len(items))
	for _, item := range items {
		item.Level = strings.TrimSpace(item.Level)
		item.Name = strings.TrimSpace(item.Name)
		item.AreaID = strings.TrimSpace(item.AreaID)
		item.StreetID = strings.TrimSpace(item.StreetID)
		item.Description = strings.TrimSpace(item.Description)
		item.Preference = normalizePreference(item.Preference)
		if item.Level == "" {
			item.Level = "area"
		}
		if item.Name == "" && item.AreaID == "" && item.StreetID == "" {
			continue
		}
		result = append(result, item)
	}
	return result
}

func normalizeRequestFloorRules(items []domain.ContactRequestFloorRule) []domain.ContactRequestFloorRule {
	result := make([]domain.ContactRequestFloorRule, 0, len(items))
	for _, item := range items {
		item.Preference = normalizePreference(item.Preference)
		if item.FloorMin < 0 {
			item.FloorMin = 0
		}
		if item.FloorMax < 0 {
			item.FloorMax = 0
		}
		if item.FloorMax > 0 && item.FloorMin > item.FloorMax {
			item.FloorMin, item.FloorMax = item.FloorMax, item.FloorMin
		}
		if item.FloorMin == 0 && item.FloorMax == 0 {
			continue
		}
		result = append(result, item)
	}
	return result
}

func normalizeRequestOptionFilters(items []domain.ContactRequestOptionFilter) []domain.ContactRequestOptionFilter {
	result := make([]domain.ContactRequestOptionFilter, 0, len(items))
	for _, item := range items {
		item.Key = strings.TrimSpace(item.Key)
		item.Preference = normalizePreference(item.Preference)
		values := make([]string, 0, len(item.Values))
		seen := map[string]bool{}
		for _, value := range item.Values {
			value = strings.TrimSpace(value)
			if value == "" || seen[value] {
				continue
			}
			seen[value] = true
			values = append(values, value)
		}
		if item.Key == "" || len(values) == 0 {
			continue
		}
		item.Values = values
		result = append(result, item)
	}
	return result
}

func normalizeRequestBooleanFilters(items []domain.ContactRequestBooleanFilter) []domain.ContactRequestBooleanFilter {
	result := make([]domain.ContactRequestBooleanFilter, 0, len(items))
	for _, item := range items {
		item.Key = strings.TrimSpace(item.Key)
		item.Preference = normalizePreference(item.Preference)
		if item.Key == "" {
			continue
		}
		result = append(result, item)
	}
	return result
}

func normalizeRequestNumberFilters(items []domain.ContactRequestNumberFilter) []domain.ContactRequestNumberFilter {
	result := make([]domain.ContactRequestNumberFilter, 0, len(items))
	for _, item := range items {
		item.Key = strings.TrimSpace(item.Key)
		item.Preference = normalizePreference(item.Preference)
		if item.Min < 0 {
			item.Min = 0
		}
		if item.Key == "" || item.Min == 0 {
			continue
		}
		result = append(result, item)
	}
	return result
}

func normalizeTags(tags []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, tag := range tags {
		tag = strings.TrimSpace(strings.TrimPrefix(tag, "#"))
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		result = append(result, tag)
	}
	return result
}

func contactHasTag(contact domain.Contact, tag string) bool {
	for _, item := range contact.Tags {
		if item == tag {
			return true
		}
	}
	return false
}

func contactHasPhone(contact domain.Contact, phone string) bool {
	for _, item := range contact.Phones {
		if strings.Contains(normalizeContactText(item.Value), phone) {
			return true
		}
	}
	return false
}

func contactHasExactPhone(contact domain.Contact, phone string) bool {
	for _, item := range contact.Phones {
		if normalized, err := NormalizePhone(item.Value); err == nil && normalized == phone {
			return true
		}
		if normalizeContactText(item.Value) == normalizeContactText(phone) {
			return true
		}
	}
	return false
}

func mergeTags(current, next []string) []string {
	merged := make([]string, 0, len(current)+len(next))
	merged = append(merged, current...)
	merged = append(merged, next...)
	return normalizeTags(merged)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func contactMatchesQuery(contact domain.Contact, query string) bool {
	values := []string{contact.FirstName, contact.LastName, contact.DisplayName, contact.Company, contact.Note}
	values = append(values, contact.Tags...)
	for _, phone := range contact.Phones {
		values = append(values, phone.Value)
	}
	for _, item := range values {
		if strings.Contains(normalizeContactText(item), query) {
			return true
		}
	}
	return false
}

func normalizeContactText(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "ي", "ی")
	value = strings.ReplaceAll(value, "ك", "ک")
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "-", "")
	return strings.ToLower(value)
}
