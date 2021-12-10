package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// created variables to use in func layoutFiles so we don't have to hardcode into Glob
var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

//creates new view object, parses all the template files given, and returns what is necessary to us as View type
func NewView(layout string, files ...string) *View {

	addTemplatePath(files)

	addTemplateExt(files)

	// appends all files we provide to this func to the files arg in this case it's our template for footer, may also need one for header, etc. may change for each section
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
		// going to panic here cause this function shouly only be ran during set up or initialization of entire app,
		//should not be used every time a user hits the site or server serves up templates
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// render is used to render the view with pre defined layout
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	switch data.(type) {
	case Data:
		//
	default:
		data = Data{
			Yield: data,
		}
	}
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// layout files returns a slice of stirngs representing
// the layout files used in our application
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// takes in a slice of strings, representing file pathss for templates and it prepends
// the templateDir director to each string
// Eg the input {"home"} would resullt in the output
//{"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		// don't have to return anyting cause we are modifying the slice that was provided sine we
		// sitting it to the index (OR SOMETHING)
		files[i] = TemplateDir + f
	}
}

// takes in slice of strings, representing file paths for templates and it appends
// the TemplateExt to each string in the slice

// Eg: the input is {"home"} would result in the output {"home.gohtml"}
// if the TemplateExt = ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
