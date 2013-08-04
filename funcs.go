package main

import (
	"html/template"
)

func htmlattr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

var funcs = template.FuncMap{
	"htmlattr": htmlattr,
}
