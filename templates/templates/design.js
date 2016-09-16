{
	"_id":"_design/by_age",
	"language": "javascript",
	"views":
	{
		"{{.Name}}": {
			{{if .RawJson}}
			    {{.RawJson}}
			 {{else}}
			"map": "function({{.VariableName}}) { if ({{.Condition}})  emit({{.EmitStr}}) }"
			,
			"reduce":"function(keys, values) { return { _id : keys[0][0], users : values.length};}"
			{{end}}
		}
	}
}
