package i18n_test

import (
	"reflect"
	"testing"

	"github.com/vitualizz/asdt/internal/i18n"
)

func TestActive_DefaultsToEnglish(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANG", "")
	t.Setenv("LANGUAGE", "")
	t.Setenv("ASDT_LANG", "")

	got := i18n.Active()
	if got.Installer.BtnContinue != i18n.English.Installer.BtnContinue {
		t.Errorf("empty locale: expected English, got BtnContinue=%q", got.Installer.BtnContinue)
	}
}

func TestActive_SpanishFromLANG(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANG", "es_AR.UTF-8")
	t.Setenv("ASDT_LANG", "")

	got := i18n.Active()
	if got.Installer.BodyEmojiPrefSubtitle != i18n.SpanishNeutral.Installer.BodyEmojiPrefSubtitle {
		t.Errorf("LANG=es_AR.UTF-8: expected SpanishNeutral (es-419), got BodyEmojiPrefSubtitle=%q", got.Installer.BodyEmojiPrefSubtitle)
	}
}

func TestActive_SpanishSpainFromLCAll(t *testing.T) {
	t.Setenv("LC_ALL", "es_ES.UTF-8")
	t.Setenv("ASDT_LANG", "")

	got := i18n.Active()
	if got.Installer.TitlePreflightCheck != i18n.SpanishSpain.Installer.TitlePreflightCheck {
		t.Errorf("LC_ALL=es_ES: expected SpanishSpain (es-ES), got TitlePreflightCheck=%q", got.Installer.TitlePreflightCheck)
	}
}

func TestActive_ASDTLANGOverridesSystem(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANG", "en_US.UTF-8")
	t.Setenv("ASDT_LANG", "es")

	got := i18n.Active()
	if got.Installer.BodyEmojiPrefSubtitle != i18n.SpanishNeutral.Installer.BodyEmojiPrefSubtitle {
		t.Errorf("ASDT_LANG=es: expected SpanishNeutral (es-419), got BodyEmojiPrefSubtitle=%q", got.Installer.BodyEmojiPrefSubtitle)
	}
}

func TestActive_FallsBackToEnglishOnUnknownLocale(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANG", "xx_XX.UTF-8")
	t.Setenv("ASDT_LANG", "")

	got := i18n.Active()
	if got.Installer.BtnContinue != i18n.English.Installer.BtnContinue {
		t.Errorf("unknown locale: expected English fallback, got BtnContinue=%q", got.Installer.BtnContinue)
	}
}

func TestActiveCode_ResolvesFullCode(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LANG", "es_AR.UTF-8")
	t.Setenv("LANGUAGE", "")
	t.Setenv("ASDT_LANG", "")
	if got := i18n.ActiveCode(); got != "es-419" {
		t.Errorf("ActiveCode() with LANG=es_AR = %q, want %q", got, "es-419")
	}

	t.Setenv("LANG", "es_ES.UTF-8")
	if got := i18n.ActiveCode(); got != "es-ES" {
		t.Errorf("ActiveCode() with LANG=es_ES = %q, want %q", got, "es-ES")
	}

	t.Setenv("LANG", "")
	if got := i18n.ActiveCode(); got != "en" {
		t.Errorf("ActiveCode() with empty locale = %q, want %q", got, "en")
	}
}

func TestForCode_KnownAndUnknown(t *testing.T) {
	// Legacy bare "es" and generic Spanish funnel to es-419 (SpanishNeutral).
	if got := i18n.ForCode("es"); got.Installer.BodyEmojiPrefSubtitle != i18n.SpanishNeutral.Installer.BodyEmojiPrefSubtitle {
		t.Errorf("ForCode(\"es\"): expected SpanishNeutral (es-419), got BodyEmojiPrefSubtitle=%q", got.Installer.BodyEmojiPrefSubtitle)
	}
	if got := i18n.ForCode("es-419"); got.Installer.BodyEmojiPrefSubtitle != i18n.SpanishNeutral.Installer.BodyEmojiPrefSubtitle {
		t.Errorf("ForCode(\"es-419\"): expected SpanishNeutral, got BodyEmojiPrefSubtitle=%q", got.Installer.BodyEmojiPrefSubtitle)
	}
	if got := i18n.ForCode("es-ES"); got.Installer.TitlePreflightCheck != i18n.SpanishSpain.Installer.TitlePreflightCheck {
		t.Errorf("ForCode(\"es-ES\"): expected SpanishSpain, got TitlePreflightCheck=%q", got.Installer.TitlePreflightCheck)
	}
	if got := i18n.ForCode("es-MX"); got.Installer.BodyEmojiPrefSubtitle != i18n.SpanishNeutral.Installer.BodyEmojiPrefSubtitle {
		t.Errorf("ForCode(\"es-MX\"): expected SpanishNeutral (es-419), got BodyEmojiPrefSubtitle=%q", got.Installer.BodyEmojiPrefSubtitle)
	}
	if got := i18n.ForCode("xx"); got.Installer.BtnContinue != i18n.English.Installer.BtnContinue {
		t.Errorf("ForCode(\"xx\"): expected English fallback, got BtnContinue=%q", got.Installer.BtnContinue)
	}
}

func TestEnglishCatalogComplete(t *testing.T) {
	walkCatalog(t, "English", i18n.English)
}

func TestES419CatalogComplete(t *testing.T) {
	walkCatalog(t, "SpanishNeutral", i18n.SpanishNeutral)
}

func TestESESCatalogComplete(t *testing.T) {
	walkCatalog(t, "SpanishSpain", i18n.SpanishSpain)
}

// walkCatalog walks every top-level struct field of Catalog (Installer,
// Dashboard, Personas, and any future feature-area struct) and verifies that
// every string field has a non-empty value. A missing field in a new catalog
// is a compile-time zero-value that only surfaces at runtime — this test
// catches it early without needing one assertion helper per struct.
func walkCatalog(t *testing.T, locale string, c i18n.Catalog) {
	t.Helper()
	root := reflect.ValueOf(c)
	rootType := root.Type()
	for i := range root.NumField() {
		section := root.Field(i)
		if section.Kind() != reflect.Struct {
			continue
		}
		sectionName := rootType.Field(i).Name
		sectionType := section.Type()
		for j := range section.NumField() {
			if section.Field(j).Kind() == reflect.String && section.Field(j).String() == "" {
				t.Errorf("%s catalog: %s.%s is empty", locale, sectionName, sectionType.Field(j).Name)
			}
		}
	}
}

func TestPersonaDescription_KnownIDs(t *testing.T) {
	for _, id := range []string{"sky", "toffy", "atreus", "babi", "lee-palacios"} {
		if got := i18n.English.PersonaDescription(id); got == "" {
			t.Errorf("English.PersonaDescription(%q) = \"\", want non-empty", id)
		}
		if got := i18n.SpanishNeutral.PersonaDescription(id); got == "" {
			t.Errorf("SpanishNeutral.PersonaDescription(%q) = \"\", want non-empty", id)
		}
	}
	if en, es := i18n.English.PersonaDescription("sky"), i18n.SpanishNeutral.PersonaDescription("sky"); en == es {
		t.Errorf("Spanish persona description for sky should differ from English, both = %q", en)
	}
}

func TestPersonaDescription_EmptyFieldFallsBackToEnglish(t *testing.T) {
	incomplete := i18n.SpanishNeutral
	incomplete.Personas.Sky = "" // simulate a catalog missing one persona description
	got := incomplete.PersonaDescription("sky")
	want := i18n.English.Personas.Sky
	if got != want {
		t.Errorf("PersonaDescription with empty field = %q, want English fallback %q", got, want)
	}
}

func TestPersonaDescription_UnknownIDReturnsEmpty(t *testing.T) {
	if got := i18n.English.PersonaDescription("nyan-cat"); got != "" {
		t.Errorf("PersonaDescription(unknown) = %q, want \"\"", got)
	}
}
