package localize

import (
	"github.com/ugent-library/biblio-backend/internal/locale"
	"github.com/ugent-library/biblio-backend/internal/render/form"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

func ValidationErrors(loc *locale.Locale, errs validation.Errors) []string {
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = loc.TranslateScope("validation", e.Code)
	}
	return msgs
}

func ValidationErrorAt(loc *locale.Locale, errs validation.Errors, ptr string) string {
	err := errs.At(ptr)
	if err == nil {
		return ""
	}
	return loc.TranslateScope("validation", err.Code)
}

func LanguageNames(loc *locale.Locale, codes []string) []string {
	names := make([]string, len(codes))
	for i, code := range codes {
		names[i] = loc.LanguageName(code)
	}
	return names
}

func LanguageSelectOptions(locale *locale.Locale) []form.SelectOption {
	vals, ok := vocabularies.Map["language_codes"]
	if !ok {
		return nil
	}

	opts := make([]form.SelectOption, len(vals))

	for i, v := range vals {
		opts[i] = form.SelectOption{
			Value: v,
			Label: locale.LanguageName(v),
		}
	}

	return opts
}

func VocabularySelectOptions(locale *locale.Locale, key string) []form.SelectOption {
	vals, ok := vocabularies.Map[key]
	if !ok {
		return nil
	}

	opts := make([]form.SelectOption, len(vals))

	for i, v := range vals {
		opts[i] = form.SelectOption{
			Value: v,
			Label: locale.TS(key, v),
		}
	}

	return opts
}
