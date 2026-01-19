package asynqcodegen

import (
	"bytes"
	"errors"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var (
	// ErrRender is returned when the template cannot be rendered.
	ErrRender = errors.New("cannot render template")
	// ErrTemplate is returned when a template file cannot be parsed.
	ErrTemplate = errors.New("invalid template")
)

// Render renders a template, provided as a string along its render context and functions.
func Render(tpl string, context any, functions template.FuncMap) ([]byte, error) {
	t, err := template.New("text").Funcs(functions).Parse(tpl)
	if err != nil {
		return nil, errors.Join(ErrTemplate, err)
	}

	out := bytes.Buffer{}

	if err := t.Execute(&out, context); err != nil {
		return nil, errors.Join(ErrRender, err)
	}

	return out.Bytes(), nil
}

// snakeCase returns the input string, but in snake case. It supports utf8.
//
// Found on: https://go.dev/play/p/tvC-pjBM1S4
func snakeCase(input string) string {
	var (
		out  strings.Builder
		prev rune
	)

	for i, v := range input {
		if unicode.IsLower(v) {
			out.WriteRune(v)
			prev = v

			continue
		}

		if i > 0 && (unicode.IsLower(prev) ||
			unicode.IsLower(nextRune(input[i+utf8.RuneLen(v):]))) {
			out.WriteByte('_')
		}

		out.WriteRune(unicode.ToLower(v))
		prev = v
	}

	return out.String()
}

func nextRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)

	return r
}
