package main

import "github.com/dapetoo/snippetbox/pkg/models"

// This will hold the structure for any dynamic data that we want to pass to HTML templates
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
