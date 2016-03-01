{
	"_id":"_design/by_age",
	"language": "javascript",
	"views":
	{
		"{{.Name}}": {
			"map": "function({{.VariableName}}) { if ({{.Condition}})  emit({{.EmitStr}}) }"
		},
	}
}
