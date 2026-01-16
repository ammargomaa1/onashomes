package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// languageCache stores loaded language files to avoid repeated disk reads
	languageCache = make(map[string]map[string]interface{})
	// cacheMutex protects concurrent access to languageCache
	cacheMutex sync.RWMutex
)

// Translate retrieves a translation value from the language file based on a nested key.
// The key format is dot-separated, e.g., "PROCUREMENT.SUPPLIERS.NAME"
// It uses the language from the goroutine-local context.
// Falls back to English if the requested language is not found.
func Translate(key string) string {
	lang := GetCurrentLang()
	return TranslateWithLang(lang, key)
}

// TranslateWithLang retrieves a translation value from the language file for a specific language.
// The key format is dot-separated, e.g., "PROCUREMENT.SUPPLIERS.NAME"
// Falls back to English if the requested language is not found or key is missing.
func TranslateWithLang(lang, key string) string {
	if lang == "" {
		lang = "en"
	}

	// Normalize language code (e.g., "ar-SA" -> "ar")
	baseLang := strings.Split(lang, "-")[0]

	value := getValueFromLanguageFile(baseLang, key)
	if value != "" {
		return value
	}

	// Fallback to English if the language or key is not found
	if baseLang != "en" {
		return getValueFromLanguageFile("en", key)
	}

	// If key is not found in any language, return the key itself
	return key
}

// getValueFromLanguageFile loads a language file and retrieves a nested value
func getValueFromLanguageFile(lang, key string) string {
	langData := loadLanguageFile(lang)
	if langData == nil {
		return ""
	}

	return getNestedValue(langData, key)
}

// loadLanguageFile loads a language file from resources/langs/ directory
func loadLanguageFile(lang string) map[string]interface{} {
	cacheMutex.RLock()
	if cached, ok := languageCache[lang]; ok {
		cacheMutex.RUnlock()
		return cached
	}
	cacheMutex.RUnlock()

	// Construct the file path: resources/langs/[lang].json
	filePath := filepath.Join("resources", "langs", fmt.Sprintf("%s.json", lang))

	data, err := os.ReadFile(filePath)
	if err != nil {
		// Language file not found
		return nil
	}

	var langData map[string]interface{}
	if err := json.Unmarshal(data, &langData); err != nil {
		// Failed to parse JSON
		return nil
	}

	// Cache the language data
	cacheMutex.Lock()
	languageCache[lang] = langData
	cacheMutex.Unlock()

	return langData
}

// getNestedValue retrieves a value from a nested map using a dot-separated key
// Example: "PROCUREMENT.SUPPLIERS.NAME" -> nested map access
func getNestedValue(data map[string]interface{}, key string) string {
	keys := strings.Split(key, ".")

	var current interface{} = data
	for _, k := range keys {
		k = strings.TrimSpace(k)

		switch v := current.(type) {
		case map[string]interface{}:
			if next, ok := v[k]; ok {
				current = next
			} else {
				return ""
			}
		default:
			return ""
		}
	}

	// Convert the final value to string
	if result, ok := current.(string); ok {
		return result
	}

	return ""
}

// ClearLanguageCache clears the in-memory language cache
func ClearLanguageCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	languageCache = make(map[string]map[string]interface{})
}
