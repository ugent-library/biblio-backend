package locale

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"
)

type Locale struct {
	Tag           language.Tag
	languageNamer display.Namer
	printer       *message.Printer
}

func (l *Locale) Translate(key string, args ...interface{}) string {
	return l.printer.Sprintf(key, args...)
}

func (l *Locale) T(key string, args ...interface{}) string {
	return l.printer.Sprintf(key, args...)
}

func (l *Locale) TranslateScope(scope, key string, args ...interface{}) string {
	if scope != "" {
		key = scope + "." + key
	}
	return l.printer.Sprintf(key, args...)
}

func (l *Locale) TS(scope, key string, args ...interface{}) string {
	return l.TranslateScope(scope, key, args...)
}

func (l *Locale) LanguageName(code string) string {
	tag := language.Make(code)
	return l.languageNamer.Name(tag)
}
