// File stores the contents in the template file so that its part of the package.
package couchdb

import (
	"text/template"
)

var DESIGNTMPL *template.Template = template.Must(template.New("design").Parse(`
{
	"_id": "{{.Id}}",
	{{if .RevStatus}}
	"_rev":"{{.Rev}}",
	{{end}}
	"views": {
		{{range .Views}}

		   "{{.Name}}": {
		   {{if .RawStatus}} 
			"map": "{{.RawJson}}"
		   {{else}}
			"map": "function({{.VariableName}}) { if({{.Condition}})  emit({{.EmitStr}});}"
		   {{end}}
		   },

		{{end}}
		   "{{.LastView.Name}}": {
	            {{if .LastView.RawStatus}}
			"map": "{{.LastView.RawJson}}"
		    {{else}}
			"map": "function({{.LastView.VariableName}}) { if({{.LastView.Condition}})  emit({{.LastView.EmitStr}});}"
		    {{end}}
		   }
	},
	"language": "javascript"
}
`))
