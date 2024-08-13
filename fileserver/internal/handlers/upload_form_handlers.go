package handlers

import (
	"fileserver/templates"
	"html/template"
	"log"
	"net/http"
)

// Template variable
var uploadTemplate *template.Template

// init function to load the template
func init() {
	var err error
	funcMap := template.FuncMap{}
	uploadTemplate, err = templates.GetTemplate("upload_form.html",funcMap)
	if err != nil {
		log.Fatalf("Failed to load template: %v", err)
	}
}
func UploadFormHandler(w http.ResponseWriter, r *http.Request) {
    dir := r.URL.Query().Get("dir")
    if dir == "" {
        dir = ""
    }

    // Prepare data to be passed to the template
    data := struct {
        CurrentPath string
    }{
        CurrentPath: dir,
    }

    // Execute the template with the provided data
    if err := uploadTemplate.Execute(w, data); err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
        return
    }

}
