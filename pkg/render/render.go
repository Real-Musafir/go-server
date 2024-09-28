package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/alshahadath/go-web/pkg/config"
	"github.com/alshahadath/go-web/pkg/models"
)


var functions = template.FuncMap {

}

var app *config.AppConfig
// NewTemplate sets teh config for the template package
func NewTemplate(a *config.AppConfig){
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData{
	return td
}

//RenderTemplate renders templates using html/temp
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData){
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	}else {
		tc, _ = CreateTemplateCache()
	}

	// get the template cache from the app config
	

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_= t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error){
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
	
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if(err != nil){
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")

		if(err != nil){
			return myCache, err
		}

		if len(matches) >0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
				if(err != nil){
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}