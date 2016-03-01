{
	"_id": "_design/{{.Id}}",
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
