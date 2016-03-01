package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type TestTempl struct {
	Name         string
	VariableName string
	Condition    string
	EmitStr      string
}

func NewTempl(name string, varName string, condition string, emitStr string) (templ *TestTempl) {
	templ = &TestTempl{
		Name:         name,
		VariableName: varName,
		Condition:    condition,
		EmitStr:      emitStr,
	}
	return templ
}

func main() {

	tmpl, err := template.ParseFiles("templates/design.js")
	if err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	testTempl := NewTempl("by_age", "doc", "doc.age < 25", "doc.age, doc.name")
	tmpl.Execute(buffer, testTempl)
	fmt.Println(buffer)
}
