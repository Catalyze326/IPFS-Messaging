package main
import(
"io/ioutil"
	"net/http"
		"log"
		"text/template"
		"fmt"
		"os"
		"io"
		"strings"
)



func java(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "Ag.txt")

}


func SendJqueryJs(w http.ResponseWriter, r *http.Request) {
    // data, err := ioutil.ReadFile("lead.js")
    // if err != nil {
    //     http.Error(w, "Couldn't read file", http.StatusInternalServerError)
    //     return
    // }
    // w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
    // w.Write(data)
		http.HandleFunc("/java/", java)
		fmt.Println("hello")
		WriteStringToFile("Lattice.txt", "Enter your username")

		title := "Lattice"
		p, _ := loadPage(title)
		renderTemplate(w, "index", p)
}


func main() {

    // http.HandleFunc("/lead.js", SendJqueryJs)
    // log.Fatal(http.ListenAndServe(":8080", nil))
// http.HandleFunc("/java/", java)
		http.HandleFunc("/index/", SendJqueryJs)
		log.Fatal(http.ListenAndServe(":8080", nil))


}


type Page struct {
		Title string
		Body  []byte
}

/**
	This will save a text file with the name of the
	page and the text of the body.
**/
func (p *Page) save() error {

		filename := p.Title + ".txt"
		return ioutil.WriteFile(filename, p.Body, 0600)
}

/**
	This will load the page by taking the title string,
	reading the textfile that has that name, and then
	returning the page title and body.
**/
func loadPage(title string) (*Page, error) {
		filename := title + ".txt"
		body, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		return &Page{Title: title, Body: body}, nil
}

/**
	To render the page with an html file.
**/
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
		t, _ := template.ParseFiles(tmpl + ".html")
		t.Execute(w, p)
}



func WriteStringToFile(filepath, s string) error {
		fo, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer fo.Close()

		_, err = io.Copy(fo, strings.NewReader(s))
		if err != nil {
			return err
		}

		return nil
}
