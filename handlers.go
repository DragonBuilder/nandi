package main

import (
	"bytes"
	"net/http"
	"text/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		w.Write([]byte("error fetching file man"))
	}

	var op bytes.Buffer
	if err := tmpl.Execute(&op, nil); err != nil {
		w.Write([]byte("error executing file man"))
	}

	w.Write(op.Bytes())
}
