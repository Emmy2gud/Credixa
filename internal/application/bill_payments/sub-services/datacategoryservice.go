package sub_services
import "strings"

func GetGloDataCategory(variationCode string) string {
	switch {
	case strings.HasPrefix(variationCode, "glo-social-") ||
		strings.HasPrefix(variationCode, "glo-telegram-") ||
		strings.HasPrefix(variationCode, "glo-insta-") ||
		strings.HasPrefix(variationCode, "glo-tiktok-") ||
		strings.HasPrefix(variationCode, "glo-opera-") ||
		strings.HasPrefix(variationCode, "glo-youtube-"):
		return "Social Media"
	case strings.HasPrefix(variationCode, "glo-daily-"):
		return "Daily"
	case strings.HasPrefix(variationCode, "glo-2days-"),
		strings.HasPrefix(variationCode, "glo-2weeks-"):
		return "Short-Term"
	case strings.HasPrefix(variationCode, "glo-monthly-"):
		return "Monthly (Day + Night)"
	case strings.HasPrefix(variationCode, "glo-weekend-") || strings.HasPrefix(variationCode, "glo-sunday"):
		return "Weekend"
	case strings.HasPrefix(variationCode, "glo-mega-"):
		return "Mega"
	case strings.HasPrefix(variationCode, "glo-tv-"):
		return "Glo TV"
	case strings.HasPrefix(variationCode, "glo-wtf-"):
		return "WTF"
	case strings.HasPrefix(variationCode, "glo-dg-"):
		return "SME"
	case strings.HasPrefix(variationCode, "glo-special-"):
		return "Special"
	default:
		return "Standard"
	}
}

func GetMTNDataCategory(variationCode string) string {
	switch {
	case variationCode == "mtn-10mb-100" || variationCode == "mtn-50mb-200" ||
		variationCode == "mtn-2-5gb-600" || variationCode == "mtn-3gb-800" || variationCode == "mtn-230mb-200":
		return "Daily / 2-Day"
	case variationCode == "mtn-20hrs-1500" || variationCode == "mtn-7gb-2000" || variationCode == "mtn-1500mb-1000":
		return "Weekly"
	case variationCode == "mtn-100gb-20000" || variationCode == "mtn-160gb-30000" ||
		variationCode == "mtn-400gb-50000" || variationCode == "mtn-600gb-75000" ||
		variationCode == "mtn-120gb-22000" || variationCode == "monthly":
		return "Multi-Month"
	case variationCode == "mtn-4-5tb-450000" || variationCode == "mtn-1tb-110000":
		return "Yearly"
	default:
		return "Monthly"
	}
}
func GetAirtelDataCategory(variationCode string) string {
	switch {

	case strings.HasPrefix(variationCode, "social"):
		return "social"
	case strings.HasPrefix(variationCode, "mifi"):
		return "mifi"
	case strings.HasSuffix(variationCode, "-1"),
		strings.HasSuffix(variationCode, "-2"),
		strings.HasSuffix(variationCode, "-7"):
		return "Short-Term"

	case strings.Contains(variationCode, "6000-30"),
		strings.Contains(variationCode, "10000"),
		strings.Contains(variationCode, "15000"),
		strings.Contains(variationCode, "20000"):
		return "Mega"

	case strings.Contains(variationCode, "1500-2"):
		return "Binge"

	case strings.Contains(variationCode, "50"),
		strings.Contains(variationCode, "100"),
		strings.Contains(variationCode, "200"),
		strings.Contains(variationCode, "300"):

		return "Daily"

	default:
		return "Monthly"
	}
}
func Get9mobileDataCategory(variationCode string) string {
	switch {
	case strings.Contains(variationCode, "150-1") ||
		strings.Contains(variationCode, "100") ||
		strings.Contains(variationCode, "200") ||
		strings.Contains(variationCode, "300"):
		return "Daily"

	case strings.Contains(variationCode, "1500-7"):
		return "Weekly"

	case strings.Contains(variationCode, "2500"):
		return "Night + Weekend"

	case strings.Contains(variationCode, "27500") ||
		strings.Contains(variationCode, "55000") ||
		strings.Contains(variationCode, "110000"):
		return "Long-Term"

	default:
		return "Monthly"
	}
}
