package i18n

import (
	"os"

	"golang.org/x/text/language"
)

// catalogs maps full canonical locale codes to their Catalog. Keys MUST stay in
// lockstep with installer.SupportedLocales codes. To add a new language: add the
// catalog var, register it here, and extend normalize() to funnel into the key.
var catalogs = map[string]Catalog{
	"en":     English,
	"es-419": SpanishNeutral,
	"es-ES":  SpanishSpain,
}

// normalize funnels any locale code (canonical, legacy, or partial) into one of
// the three registered catalog keys: "en", "es-419", or "es-ES". It is the
// SINGLE backward-compat boundary: every legacy shape (bare "es", generic
// es-<region>, "und", unknown, unparsable) collapses here.
//
//   - base != "es"                         → "en"
//   - base "es" with EXPLICIT region "ES"  → "es-ES"
//   - base "es", any other / inferred region (incl. bare "es", es-419, es-AR…)
//     → "es-419"
//
// The region-confidence check distinguishes an explicit "es-ES" (Exact) from a
// bare "es" whose ES region is merely inferred (Low) — the latter must resolve
// to neutral es-419 for backward compatibility, never peninsular.
func normalize(code string) string {
	tag, err := language.Parse(code)
	if err != nil {
		return "en"
	}
	return normalizeTag(tag)
}

// normalizeTag is the language.Tag variant of normalize, used by callers that
// already hold a parsed tag (catalogFor, ActiveCode) so they avoid a redundant
// String()/Parse round-trip.
func normalizeTag(tag language.Tag) string {
	base, _ := tag.Base()
	if base.String() != "es" {
		return "en"
	}
	if region, conf := tag.Region(); region.String() == "ES" && conf != language.Low {
		return "es-ES"
	}
	return "es-419"
}

// resolve returns the language tag for the current environment. ASDT_LANG takes
// priority over system locale detection so it can be overridden without changing
// the system locale (e.g. ASDT_LANG=es).
//
// The override returns the PARSED tag verbatim (not supported.Match'd): the
// matcher's region-distance heuristic would canonicalize a bare "es" into an
// es-ES tag with an EXACT region, erasing the inferred-region signal that
// normalize() relies on to route legacy "es" to neutral es-419. Normalization
// happens downstream at the catalogFor/ActiveCode boundary, consistent with
// detect()'s region pre-switch.
func resolve() language.Tag {
	if override := os.Getenv("ASDT_LANG"); override != "" {
		if tag, err := language.Parse(override); err == nil {
			return tag
		}
	}
	return detect()
}

// Active returns the Catalog for the user's locale.
func Active() Catalog {
	return catalogFor(resolve())
}

// ActiveCode returns the resolved full canonical locale code for the user's
// locale (e.g. "en", "es-419", "es-ES"). All inputs are funneled through
// normalize, so the result is always a registered catalog key.
func ActiveCode() string {
	return normalizeTag(resolve())
}

// ForCode returns the Catalog registered for the given locale code, falling back
// to English for unknown codes. The code is funneled through normalize, so
// legacy shapes (bare "es", generic es-<region>) resolve correctly. Useful when
// the language choice comes from persisted state or UI selection instead of env
// detection.
func ForCode(code string) Catalog {
	if c, ok := catalogs[normalize(code)]; ok {
		return c
	}
	return English
}

func catalogFor(tag language.Tag) Catalog {
	if c, ok := catalogs[normalizeTag(tag)]; ok {
		return c
	}
	return English
}
