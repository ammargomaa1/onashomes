package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const languageHeader = "Accept-Language"

func GetRequestLang(c *gin.Context) string {
	lang := strings.TrimSpace(c.GetHeader(languageHeader))
	return strings.ToLower(lang)
}

// GetCurrentLang retrieves the language from the goroutine-local context
func GetCurrentLang() string {
	c := GetContextForGoroutine()
	if c == nil {
		return ""
	}
	return GetRequestLang(c)
}

// SelectLocalizedString returns the Arabic or English value based on the provided
// language string. It prefers Arabic when lang starts with "ar", otherwise it
// falls back to English. If the chosen value is empty, it gracefully falls back
// to the other one.
func SelectLocalizedString(lang, valueAr, valueEn string) string {
	if lang == "" {
		if valueEn != "" {
			return valueEn
		}
		return valueAr
	}

	if strings.HasPrefix(lang, "ar") {
		if valueAr != "" {
			return valueAr
		}
		return valueEn
	}

	if valueEn != "" {
		return valueEn
	}
	return valueAr
}

// GetLocalizedStringFromContext retrieves the Arabic or English value based on
// the goroutine-local context. It does not require passing gin.Context explicitly.
func GetLocalizedStringFromContext(valueAr, valueEn string) string {
	lang := GetCurrentLang()
	return SelectLocalizedString(lang, valueAr, valueEn)
}
