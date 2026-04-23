package i18n

import "strings"

// T looks up a translation key of the form "file:dotted.path" in the given
// language. It falls back to the default language (English) if the key is
// missing, and returns the key itself if still not found — so omissions are
// visible during review but never crash rendering.
func T(lang Lang, key string) string {
	if v := lookup(lang, key); v != "" {
		return v
	}
	if lang != DefaultLang {
		if v := lookup(DefaultLang, key); v != "" {
			return v
		}
	}
	return key
}

func lookup(lang Lang, key string) string {
	colon := strings.IndexByte(key, ':')
	if colon <= 0 || colon == len(key)-1 {
		return ""
	}
	file := key[:colon]
	path := key[colon+1:]

	langMap, ok := translations[lang]
	if !ok {
		return ""
	}
	fileTree, ok := langMap[file].(map[string]interface{})
	if !ok {
		return ""
	}

	var current interface{} = fileTree
	for _, seg := range strings.Split(path, ".") {
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		current, ok = m[seg]
		if !ok {
			return ""
		}
	}
	if s, ok := current.(string); ok {
		return s
	}
	return ""
}
