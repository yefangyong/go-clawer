package view

import (
	"go-clawer/frontend/model"
	"html/template"
	"io"
)

type SearchResultView struct {
	template *template.Template
}

func CreateSearchResultView(filename string) SearchResultView {
	return SearchResultView{
		template: template.Must(template.ParseFiles(filename)),
	}
}

func (s SearchResultView) Render(io io.Writer, data model.SearchResult) error {
	return s.template.Execute(io, data)
}
