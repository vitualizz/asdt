package i18n

import (
	"os"
	"strings"

	"golang.org/x/text/language"
)

// supported lists the locales the application ships catalogs for.
// English must be first — language.NewMatcher uses position 0 as the fallback
// (NFR-i18n-constraint). x/text has no es-419 constant, so the regional Spanish
// tags are built with MustParse.
var supported = language.NewMatcher([]language.Tag{
	language.English,
	language.MustParse("es-419"),
	language.MustParse("es-ES"),
})

// detect returns the best-matching supported language tag from the system
// locale environment variables. Reads LC_ALL, LANG, and LANGUAGE in that
// order; falls back to English if none is parseable or supported.
func detect() language.Tag {
	for _, env := range []string{"LC_ALL", "LANG", "LANGUAGE"} {
		val := os.Getenv(env)
		if val == "" || val == "C" || val == "POSIX" {
			continue
		}
		// LANGUAGE can be a colon-separated priority list; take the first tag.
		if idx := strings.IndexByte(val, ':'); idx != -1 {
			val = val[:idx]
		}
		// Strip encoding suffix: "es_AR.UTF-8" → "es_AR".
		if idx := strings.IndexByte(val, '.'); idx != -1 {
			val = val[:idx]
		}
		// Normalize POSIX underscores to BCP 47 hyphens: "es_AR" → "es-AR".
		val = strings.ReplaceAll(val, "_", "-")
		tag, err := language.Parse(val)
		if err != nil {
			continue
		}
		// Deterministic region pre-switch: route Spanish locales by region
		// regardless of the matcher's region-distance heuristic. An explicit
		// "ES" region resolves to peninsular es-ES; every other Spanish region
		// (419/AR/MX/CO/…) and bare "es" resolves to neutral es-419. Non-Spanish
		// bases fall through to Match, keeping English (index 0) as the fallback.
		base, _ := tag.Base()
		if base.String() == "es" {
			if region, conf := tag.Region(); region.String() == "ES" && conf != language.Low {
				return language.MustParse("es-ES")
			}
			return language.MustParse("es-419")
		}
		matched, _, _ := supported.Match(tag)
		return matched
	}
	return language.English
}
