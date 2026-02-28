package main

import (
	"html/template"
	"strings"
)

func nl2br(b []byte) template.HTML {
	s := string(b) // Convert []byte to string
	return template.HTML(strings.Replace(template.HTMLEscapeString(s), "\n", "<br>", -1))
}
