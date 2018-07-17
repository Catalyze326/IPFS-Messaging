package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

// func index_handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "")
// }
//
// func main() {
// 	http.HandleFunc("/", index_handler)
// 	http.ListenAndServe(":8080", nil)
// }

func (p *Page) save() error {
	f := p.Title + ".txt"
	return ioutil.WriteFile(f, p.Body, 0600)
}

func load(title string) (*Page, error) {
	f := title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func view(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	p, _ := load(title)
	t, _ := template.ParseFiles("bootstrapNavy2.html")
	t.Execute(w, p)
}

func edit(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/test/"):]
	p, _ := load(title)
	t, _ := template.ParseFiles("messaging.html")
	t.Execute(w, p)
}

func save(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/test/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/test/", view)
	http.HandleFunc("/edit/", edit)
	http.HandleFunc("/save/", save)
	http.ListenAndServe(":8080", nil)
}
