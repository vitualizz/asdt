package installer

// Locale describes one selectable installer language. A single TUI selection
// resolves to one Locale: its Code is persisted verbatim into
// InstallMeta.Language, the same full code feeds the i18n catalog (the i18n
// package normalizes it), and its Directive is injected into AGENTS.md as
// {{language_directive}}.
type Locale struct {
	Code      string // full BCP-47 locale stored in InstallMeta.Language (en/es-419/es-ES)
	Label     string // friendly native name shown in the single TUI selector
	Directive string // English instruction injected as {{language_directive}}
}

// SupportedLocales is the single source of truth for the installer's selectable
// languages. The TUI selector and renderAgentConfig both consume it.
//
// The neutral-register rule (clear register, no strong regional dialect or slang)
// lives in the agents-template.md framing, not here — Directives name only the
// spoken language so the neutrality guidance stays single-sourced across locales.
var SupportedLocales = []Locale{
	{Code: "en", Label: "🇬🇧 English", Directive: "Respond in English, keeping the assistant's own defined voice and style."},
	{Code: "es-419", Label: "🌎 Español (Latinoamérica)", Directive: "Respond in Latin American Spanish, keeping the assistant's own defined voice and style."},
	{Code: "es-ES", Label: "🇪🇸 Español (España)", Directive: "Respond in peninsular Spanish (Spain), keeping the assistant's own defined voice and style."},
}

// DefaultLocaleCode is the fallback locale used for empty or unknown codes.
const DefaultLocaleCode = "en"

// localeFor returns the SupportedLocales row whose Code matches code exactly,
// or the default locale row when there is no exact match.
func localeFor(code string) Locale {
	for _, l := range SupportedLocales {
		if l.Code == code {
			return l
		}
	}
	for _, l := range SupportedLocales {
		if l.Code == DefaultLocaleCode {
			return l
		}
	}
	return Locale{Code: DefaultLocaleCode, Label: "🇬🇧 English", Directive: "Respond in English, keeping the assistant's own defined voice and style."}
}

// LocaleByCode resolves code to a Locale: exact match first, then the legacy
// bare "es" maps to es-419 for backward compatibility, then the default locale
// for any empty or unknown code. The returned Directive is always non-empty.
func LocaleByCode(code string) Locale {
	for _, l := range SupportedLocales {
		if l.Code == code {
			return l
		}
	}
	if code == "es" {
		return localeFor("es-419") // legacy bare-es backward-compat
	}
	return localeFor(DefaultLocaleCode)
}
