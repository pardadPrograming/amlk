package searchengine

import (
	"math"
	"strings"

	"amlakcrm/backend/internal/domain"
)

const (
	rentLeaseDepositStep = int64(100_000_000)
	rentLeaseRentStep    = int64(3_000_000)
)

type matchCounter struct {
	total   int
	matched int
	missed  int
	soft    bool
	hard    bool
	reasons []string
	misses  []string
}

func (c *matchCounter) add(ok bool, reason, miss string) {
	c.total++
	if ok {
		c.matched++
		if reason != "" {
			c.reasons = append(c.reasons, reason)
		}
		return
	}
	c.missed++
	if miss != "" {
		c.misses = append(c.misses, miss)
	}
}

func matchProperty(request domain.ContactRequest, file domain.PropertyFile) domain.PropertyMatchResult {
	counter := &matchCounter{}
	if !propertyTypeMatches(request.Type, file.Types, file.Type) {
		return domain.PropertyMatchResult{}
	}
	counter.add(true, "نوع فایل با درخواست سازگار است", "")
	addFinancialMatch(counter, request, file)
	if request.MinAreaM2 > 0 {
		counter.add(file.HouseInfo.AreaM2 >= request.MinAreaM2, "متراژ حداقل درخواست را دارد", "متراژ کمتر از حداقل درخواست است")
	}
	if request.MaxAgeYears > 0 {
		counter.add(file.HouseInfo.AgeYears <= request.MaxAgeYears, "سن بنا در بازه قابل قبول است", "سن بنا بیشتر از حداکثر درخواست است")
	}
	if request.LandMinAreaM2 > 0 {
		counter.add(file.HouseInfo.AreaM2 >= request.LandMinAreaM2, "متراژ زمین برای مشارکت مناسب است", "متراژ زمین کمتر از حداقل مشارکت است")
	}
	if request.BuildingMinAreaM2 > 0 {
		counter.add(file.HouseInfo.GardenBuildingAreaM2 >= request.BuildingMinAreaM2 || file.HouseInfo.AreaM2 >= request.BuildingMinAreaM2, "بنای موجود با درخواست سازگار است", "بنای موجود کمتر از حداقل درخواست است")
	}
	addLocationMatch(counter, request, file)
	addFloorRuleMatch(counter, request, file)
	addOptionMatches(counter, request, file)
	addBooleanMatches(counter, request, file)
	addNumberMatches(counter, request, file)
	if counter.total == 0 {
		return domain.PropertyMatchResult{}
	}
	score := counter.matched * 100 / counter.total
	if counter.hard && score < 80 {
		return domain.PropertyMatchResult{}
	}
	tier := ""
	switch {
	case score == 100 && counter.missed == 0 && !counter.soft:
		tier = "green"
	case score >= 80 && counter.missed > 0:
		tier = "orange"
	case score > 0:
		tier = "yellow"
	}
	if tier == "" {
		return domain.PropertyMatchResult{}
	}
	return domain.PropertyMatchResult{
		PropertyFile:   file,
		Score:          score,
		Tier:           tier,
		MatchedReasons: counter.reasons,
		MissedReasons:  counter.misses,
	}
}

func MatchProperty(request domain.ContactRequest, file domain.PropertyFile) domain.PropertyMatchResult {
	return matchProperty(request, file)
}

func propertyTypeMatches(requestType string, fileTypes []domain.PropertyFileType, primary domain.PropertyFileType) bool {
	target := domain.PropertyFileType(requestType)
	if target == "" {
		return true
	}
	if primary == target {
		return true
	}
	for _, item := range fileTypes {
		if item == target {
			return true
		}
	}
	return false
}

func addFinancialMatch(counter *matchCounter, request domain.ContactRequest, file domain.PropertyFile) {
	switch request.Type {
	case "rent_lease":
		addRentLeaseFinancialMatch(counter, request, file)
		if request.RentWithOwner {
			counter.add(file.RentWithOwner, "همراه مالک بودن با درخواست سازگار است", "همراه مالک نیست")
		}
	case "partnership":
		price := file.SalePrice
		if request.PartnershipMin > 0 {
			counter.add(price >= request.PartnershipMin, "ارزش مشارکت از حداقل درخواست بالاتر است", "ارزش مشارکت کمتر از حداقل درخواست است")
		}
		if request.PartnershipMax > 0 {
			counter.add(price <= request.PartnershipMax, "ارزش مشارکت در سقف درخواست است", "ارزش مشارکت از سقف درخواست بیشتر است")
		}
	default:
		price := file.FinalPrice
		if price == 0 {
			price = file.SalePrice
		}
		if request.PurchaseMin > 0 {
			counter.add(price >= request.PurchaseMin, "قیمت از حداقل درخواست بالاتر است", "قیمت کمتر از حداقل درخواست است")
		}
		if request.PurchaseMax > 0 {
			counter.add(price <= request.PurchaseMax, "قیمت در سقف درخواست است", "قیمت از سقف درخواست بیشتر است")
		}
	}
}

func addRentLeaseFinancialMatch(counter *matchCounter, request domain.ContactRequest, file domain.PropertyFile) {
	preferred := rentLeaseRange{
		depositMin: request.DepositMin,
		depositMax: request.DepositMax,
		rentMin:    request.RentMin,
		rentMax:    request.RentMax,
	}
	suggested := rentLeaseRange{
		depositMin: request.SuggestedDepositMin,
		depositMax: request.SuggestedDepositMax,
		rentMin:    request.SuggestedRentMin,
		rentMax:    request.SuggestedRentMax,
	}
	if !preferred.active() && !suggested.active() {
		return
	}
	pricing := rentLeasePricingForFile(file)
	if pricing.matches(preferred) {
		counter.add(true, "رهن و اجاره با بازه اصلی درخواست مچ است", "")
		return
	}
	if pricing.matches(suggested) {
		counter.soft = true
		counter.add(true, "رهن و اجاره با بازه پیشنهادی مشاور مچ است", "")
		return
	}
	counter.hard = true
	counter.add(false, "", "رهن و اجاره با بازه اصلی یا پیشنهادی درخواست مچ نیست")
}

type rentLeaseRange struct {
	depositMin int64
	depositMax int64
	rentMin    int64
	rentMax    int64
}

func (r rentLeaseRange) active() bool {
	return r.depositMin > 0 || r.depositMax > 0 || r.rentMin > 0 || r.rentMax > 0
}

func (r rentLeaseRange) contains(deposit, rent int64) bool {
	if !r.active() {
		return false
	}
	if r.depositMin > 0 && deposit < r.depositMin {
		return false
	}
	if r.depositMax > 0 && deposit > r.depositMax {
		return false
	}
	if r.rentMin > 0 && rent < r.rentMin {
		return false
	}
	if r.rentMax > 0 && rent > r.rentMax {
		return false
	}
	return true
}

type rentLeasePricing struct {
	depositStart int64
	rentStart    int64
	depositEnd   int64
	rentEnd      int64
}

func rentLeasePricingForFile(file domain.PropertyFile) rentLeasePricing {
	pricing := rentLeasePricing{
		depositStart: file.DepositPrice,
		rentStart:    file.RentPrice,
		depositEnd:   file.DepositPrice,
		rentEnd:      file.RentPrice,
	}
	if !file.Convertible || file.MaxConvertibleDeposit <= file.DepositPrice {
		return pricing
	}
	pricing.depositEnd = file.MaxConvertibleDeposit
	pricing.rentEnd = convertedRent(file.RentPrice, file.MaxConvertibleDeposit-file.DepositPrice)
	return pricing
}

func convertedRent(baseRent, depositIncrease int64) int64 {
	if depositIncrease <= 0 {
		return baseRent
	}
	reduction := depositIncrease * rentLeaseRentStep / rentLeaseDepositStep
	if reduction >= baseRent {
		return 0
	}
	return baseRent - reduction
}

func (p rentLeasePricing) matches(r rentLeaseRange) bool {
	if !r.active() {
		return false
	}
	if r.contains(p.depositStart, p.rentStart) || r.contains(p.depositEnd, p.rentEnd) {
		return true
	}
	if p.depositStart == p.depositEnd {
		return false
	}
	minDeposit, maxDeposit := int64Bounds(p.depositStart, p.depositEnd)
	candidateMin := maxInt64(minDeposit, positiveOrMin(r.depositMin, minDeposit))
	candidateMax := minInt64(maxDeposit, positiveOrMax(r.depositMax, maxDeposit))
	if candidateMin > candidateMax {
		return false
	}
	for _, deposit := range candidateDepositsForRentRange(p, r, candidateMin, candidateMax) {
		if deposit < candidateMin || deposit > candidateMax {
			continue
		}
		if r.contains(deposit, p.rentAt(deposit)) {
			return true
		}
	}
	return false
}

func candidateDepositsForRentRange(p rentLeasePricing, r rentLeaseRange, minDeposit, maxDeposit int64) []int64 {
	candidates := []int64{minDeposit, maxDeposit, p.depositStart, p.depositEnd}
	if r.rentMin > 0 {
		candidates = append(candidates, p.depositAtRentBoundary(r.rentMin)...)
	}
	if r.rentMax > 0 {
		candidates = append(candidates, p.depositAtRentBoundary(r.rentMax)...)
	}
	return candidates
}

func (p rentLeasePricing) rentAt(deposit int64) int64 {
	if p.depositEnd == p.depositStart {
		return p.rentStart
	}
	ratio := float64(deposit-p.depositStart) / float64(p.depositEnd-p.depositStart)
	rent := float64(p.rentStart) + ratio*float64(p.rentEnd-p.rentStart)
	if rent < 0 {
		return 0
	}
	return int64(math.Round(rent))
}

func (p rentLeasePricing) depositAtRentBoundary(rent int64) []int64 {
	if p.rentEnd == p.rentStart {
		return nil
	}
	ratio := float64(rent-p.rentStart) / float64(p.rentEnd-p.rentStart)
	deposit := float64(p.depositStart) + ratio*float64(p.depositEnd-p.depositStart)
	rounded := int64(math.Round(deposit))
	return []int64{rounded - 1, rounded, rounded + 1}
}

func int64Bounds(a, b int64) (int64, int64) {
	if a < b {
		return a, b
	}
	return b, a
}

func positiveOrMin(value, fallback int64) int64 {
	if value > 0 {
		return value
	}
	return fallback
}

func positiveOrMax(value, fallback int64) int64 {
	if value > 0 {
		return value
	}
	return fallback
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func addLocationMatch(counter *matchCounter, request domain.ContactRequest, file domain.PropertyFile) {
	if len(request.Locations) == 0 {
		return
	}
	ok := false
	for _, wanted := range request.Locations {
		for _, address := range file.Addresses {
			if locationMatches(wanted, address) {
				ok = true
				break
			}
		}
	}
	counter.add(ok, "محدوده مکانی با درخواست همخوان است", "محدوده مکانی با درخواست تناقض دارد")
}

func locationMatches(wanted domain.ContactRequestLocation, address domain.PropertyAddress) bool {
	name := NormalizePersianName(wanted.Name)
	switch wanted.Level {
	case "neighborhood":
		return (wanted.ID != "" && wanted.ID == address.NeighborhoodID) || (name != "" && name == NormalizePersianName(address.NeighborhoodName))
	case "street":
		return (wanted.ID != "" && wanted.ID == address.StreetID) || (wanted.StreetID != "" && wanted.StreetID == address.StreetID) || (name != "" && name == NormalizePersianName(address.StreetName))
	default:
		return (wanted.ID != "" && wanted.ID == address.AreaID) || (wanted.AreaID != "" && wanted.AreaID == address.AreaID) || (name != "" && name == NormalizePersianName(address.AreaName))
	}
}

func addFloorRuleMatch(counter *matchCounter, request domain.ContactRequest, file domain.PropertyFile) {
	if len(request.FloorRules) == 0 {
		return
	}
	ok := false
	for _, rule := range request.FloorRules {
		minOK := rule.FloorMin == 0 || file.HouseInfo.Floor >= rule.FloorMin
		maxOK := rule.FloorMax == 0 || file.HouseInfo.Floor <= rule.FloorMax
		elevatorOK := !rule.Elevator || file.HouseInfo.Elevator
		if minOK && maxOK && elevatorOK {
			ok = true
			break
		}
	}
	counter.add(ok, "قانون طبقه و آسانسور مچ است", "طبقه یا آسانسور با قوانین درخواست نمی‌خواند")
}

func addOptionMatches(counter *matchCounter, request domain.ContactRequest, file domain.PropertyFile) {
	for _, filter := range request.OptionFilters {
		value := optionValue(filter.Key, file)
		if value == "" {
			continue
		}
		ok := false
		for _, wanted := range filter.Values {
			if NormalizePersianName(wanted) == NormalizePersianName(value) {
				ok = true
				break
			}
		}
		counter.add(ok, filter.Key+" با گزینه‌های درخواست سازگار است", filter.Key+" با گزینه‌های درخواست تناقض دارد")
	}
}

func optionValue(key string, file domain.PropertyFile) string {
	switch key {
	case "cabinetType":
		return file.HouseInfo.CabinetType
	case "flooring":
		return file.HouseInfo.Flooring
	case "documentType":
		return file.HouseInfo.DocumentType
	default:
		return ""
	}
}

func addBooleanMatches(counter *matchCounter, request domain.ContactRequest, file domain.PropertyFile) {
	for _, filter := range request.BooleanFilters {
		actual, ok := boolValue(filter.Key, file)
		if !ok {
			continue
		}
		counter.add(actual == filter.Value, filter.Key+" مطابق درخواست است", filter.Key+" با درخواست تناقض دارد")
	}
}

func boolValue(key string, file domain.PropertyFile) (bool, bool) {
	switch key {
	case "parking":
		return file.HouseInfo.Parking, true
	case "warehouse", "storage":
		return file.HouseInfo.Storage, true
	case "elevator":
		return file.HouseInfo.Elevator, true
	case "renovated":
		return file.HouseInfo.Renovated, true
	default:
		return false, false
	}
}

func addNumberMatches(counter *matchCounter, request domain.ContactRequest, file domain.PropertyFile) {
	for _, filter := range request.NumberFilters {
		actual, ok := numberValue(filter.Key, file)
		if !ok || filter.Min == 0 {
			continue
		}
		counter.add(actual >= filter.Min, filter.Key+" حداقل درخواست را دارد", filter.Key+" کمتر از حداقل درخواست است")
	}
}

func numberValue(key string, file domain.PropertyFile) (int, bool) {
	switch key {
	case "bedrooms":
		return file.HouseInfo.Bedrooms, true
	case "masterServiceCount":
		return file.HouseInfo.MasterServiceCount, true
	case "terraceCount":
		return file.HouseInfo.TerraceCount, true
	case "bathrooms":
		return file.HouseInfo.MasterServiceCount, true
	default:
		if strings.TrimSpace(key) == "" {
			return 0, false
		}
		return 0, false
	}
}

func NormalizePersianName(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "?", "?")
	value = strings.ReplaceAll(value, "?", "?")
	return strings.Join(strings.Fields(value), " ")
}
