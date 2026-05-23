package searchengine

import (
	"strings"

	"amlakcrm/backend/internal/domain"
)

type matchCounter struct {
	total   int
	matched int
	missed  int
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
	tier := ""
	switch {
	case score == 100 && counter.missed == 0:
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
		if request.DepositMin > 0 {
			counter.add(file.DepositPrice >= request.DepositMin, "رهن از حداقل درخواست بالاتر است", "رهن کمتر از حداقل درخواست است")
		}
		if request.DepositMax > 0 {
			counter.add(file.DepositPrice <= request.DepositMax || (request.Convertible && file.Convertible), "رهن در بازه یا قابل تبدیل است", "رهن از سقف درخواست بیشتر است")
		}
		if request.RentMin > 0 {
			counter.add(file.RentPrice >= request.RentMin, "اجاره از حداقل درخواست بالاتر است", "اجاره کمتر از حداقل درخواست است")
		}
		if request.RentMax > 0 {
			counter.add(file.RentPrice <= request.RentMax, "اجاره در سقف درخواست است", "اجاره از سقف درخواست بیشتر است")
		}
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
