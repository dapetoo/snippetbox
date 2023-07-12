package main

import (
	"github.com/dapetoo/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
)

// This will hold the structure for any dynamic data that we want to pass to HTML templates
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	//Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	//Filepath.Glob() to get a slice of all filepaths with the extensions "page.tmpl"
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	//Loop through the pages
	for _, page := range pages {
		//Extract file name
		name := filepath.Base(page)
		//Parse the page template in to a template set
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		//Use the ParseGlob() method to add any layout templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		//Use the ParseGlob() method to add any partial templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		//Add the template set to the cache using the name of the page
		cache[name] = ts
	}
	return cache, nil
}
