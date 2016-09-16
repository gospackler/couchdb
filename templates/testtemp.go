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
	RawStatus    bool
	RawJson      string
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

func case1() {
	tmpl, err := template.ParseFiles("templates/design.js")
	if err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	testTempl := NewTempl("by_age", "doc", "doc.age < 25", "doc.age, doc.name")
	tmpl.Execute(buffer, testTempl)
	fmt.Println(buffer)
}

func case2() {
	tmpl, err := template.ParseFiles("templates/design.js")
	if err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	testTempl := &TestTempl{
		RawStatus: true,
		RawJson:   "\"map\" : \"function(doc) { cosole.log(1234)}\", \"reduce\":\"function(keys,value) {console.log(1234);}\"",
	}
	tmpl.Execute(buffer, testTempl)
	fmt.Println(buffer)
}

func main() {

	case1()
	case2()
}
