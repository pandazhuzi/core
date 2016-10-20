package core

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

func camelName(source []string) string {

	var value string

	for _, slice := range source {

		if len(slice) == 0 {
			continue
		}

		slice = strings.ToUpper(slice[0:1]) + slice[1:]

		value += slice

	}

	return value

}

func unixName(source []string) string {

	return strings.Join(source, "_")
}

func formatName(source string) (string, string) {

	var camel string
	var unix string
	var runs []rune

	buf := bytes.NewBuffer(nil)

	for len(source) > 0 {
		r, size := utf8.DecodeRuneInString(source)

		runs = append(runs, r)

		source = source[size:]
	}

	length := len(runs)
	for index, r := range runs {

		if index == 0 {
			buf.WriteRune(r)
		} else {

			if unicode.IsUpper(r) {

				if index+1 == length {
					if unicode.IsLower(runs[index-1]) {
						buf.WriteString("_")
						buf.WriteRune(r)
					} else {
						buf.WriteRune(r)
					}
				} else if unicode.IsLower(runs[index+1]) {
					buf.WriteString("_")
					buf.WriteRune(r)
				} else if unicode.IsLower(runs[index-1]) {
					buf.WriteString("_")
					buf.WriteRune(r)
				} else {
					buf.WriteRune(r)
				}

			} else {
				buf.WriteRune(r)
			}

		}
	}

	names := strings.Split(
		strings.ToLower(buf.String()), "_")

	camel = camelName(names)
	unix = unixName(names)

	return camel, unix

}
