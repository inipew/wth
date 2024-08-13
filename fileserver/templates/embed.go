package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html
var templates embed.FS

// GetTemplate returns the parsed template from the specified template file, including a funcMap.
func GetTemplate(name string, funcMap template.FuncMap) (*template.Template, error) {
	content, err := templates.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return template.New(name).Funcs(funcMap).Parse(string(content))
}
