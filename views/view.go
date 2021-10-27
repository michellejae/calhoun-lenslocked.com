package views

import "html/template"

//creates new view object, parses all the template files given, and returns what is necessary to us as View type
func NewView(layout string, files ...string) *View {
	// appends all files we provide to this func to the files arg in this case it's our template for footer, may also need one for header, etc. may change for each section
	files = append(files,
		"views/layouts/footer.gohtml",
		"views/layouts/bootstrap.gohtml",
		"views/layouts/navbar.gohtml")

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
