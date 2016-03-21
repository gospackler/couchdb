// File stores the contents in the template file so that its part of the package.
package couchdb

import (
	"strings"
	"text/template"
)

var DESIGNTMPL *template.Template = template.Must(template.New("design").Parse(strings.Replace(strings.Replace(`
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
			"map": "function({{.VariableName}}) { 
				{{if .CondStatus}}
					if({{.Condition}}) {
						emit({{.EmitStr}});
						}
				{{else}}
						emit{{.EmitStr}});
				{{end}}
			}"
		   {{end}}
		   },

		{{end}}
		   "{{.LastView.Name}}": {
	            {{if .LastView.RawStatus}}
			"map": "{{.LastView.RawJson}}"
		    {{else}}
			"map": "function({{.LastView.VariableName}}) { \
				{{if .LastView.CondStatus}}
					if({{.LastView.Condition}}) { \
						emit({{.LastView.EmitStr}}); \
						} \
				{{else}} \
						emit({{.LastView.EmitStr}});\
				{{end}} \
			}"
		    {{end}}
		   }
	},
	"language": "javascript"
}
`, "\\\n", "", -1), "\t", "", -1)))
