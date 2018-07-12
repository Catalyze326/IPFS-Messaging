package main

import (
	"fmt"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<!DOCTYPE html>
	<html>

	<head>
	  <link rel="stylesheet" href="files/bootstrap.min.css">
	  <link rel="stylesheet" href="mess.css">
	  <!-- THis is the header that includes the two buttons for audio and file transfer -->
	  <link rel="icon" href="pictures/favicon.ico">

	  <body>
	    <div class="header">
	      <font color="white"><span style="font-size:30px;cursor:pointer" onclick="openNav()">&#9776; IPFS Navy</span></font>
	      <div class="header-right">
	        <button type="button" class="btn btn-default btn-circle btn-lg"><img src = "pictures/upload2.png" class="img-responsive" alt="Upload Icon"></i></button>
	        <button type="button" class="btn btn-default btn-circle btn-lg"><img src = "pictures/mic.png" class="img-responsive" alt="Mic Icon"></i></button>
	        <a href="bootstrapNavy2.html">Home</a>
	        <a href="messaging.html">Messaging</a>
	      </div>
	    </div>

	    <div id="mySidenav" class="sidenav">
	      <a href="javascript:void(0)" class="closebtn" onclick="closeNav()">&times;</a>
	      <a href="https://en.wikipedia.org/wiki/InterPlanetary_File_System">About IPFS</a>
	      <a href="messaging.html">Messenger</a>
	      <a href="ourTeam.html">Our Team</a>
	      </font>
	    </div>

	    <script>
	      function openNav() {
	        document.getElementById("mySidenav").style.width = "250px";
	      }

	      function closeNav() {
	        document.getElementById("mySidenav").style.width = "0";
	      }

	    </script>
	    <div class="footer ">
	      <div class="col-lg-1">
	      </div>
	      <div class="col-lg-10">
	        <ul style="overflow: auto; max-height: 80vh;"><strong>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	          <li class="list-group-item text-right list-group-item-info ">This is a message I posted</li>
	          <li class="list-group-item text-left list-group-item-success ">This is a message they posted</li>
	    </ul></strong>
	      </div>
	      <div class="col-lg-1 ">
	      </div>

	      <div class="form-group " class="text-center ">
	        <label class="control-label col-lg-12 " for="pwd ">Message:</label>
	      </div>
	      <div class="col-lg-12 ">
	        <input type="Message " class="form-control input-lg " id="pwd " placeholder="Enter Message " name="pwd ">
	        <br>
	      </div>
	    </div>

	</html>
")
}

func main() {
	http.HandleFunc("/", index_handler)
	http.ListenAndServe(":8080", nil)
}

//
// func (p *Page) save() error {
// 	f := p.Title + ".txt"
// 	return ioutil.WriteFile(f, p.Body, 0600)
// }
//
// func load(title string) (*Page, error) {
// 	f := title + ".txt"
// 	body, err := ioutil.ReadFile(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Page{Title: title, Body: body}, nil
// }
//
// func view(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/test/"):]
// 	p, _ := load(title)
// 	t, _ := template.ParseFiles("bootstrapNavy2.html")
// 	t.Execute(w, p)
// }
//
// func edit(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/edit/"):]
// 	p, _ := load(title)
// 	t, _ := template.ParseFiles("messaging.html")
// 	t.Execute(w, p)
// }
//
// func save(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/save/"):]
// 	body := r.FormValue("body")
// 	p := &Page{Title: title, Body: []byte(body)}
// 	p.save()
// 	http.Redirect(w, r, "/test/"+title, http.StatusFound)
// }
//
// func main() {
// 	http.HandleFunc("/test/", view)
// 	http.HandleFunc("/edit/", edit)
// 	http.HandleFunc("/save/", save)
// 	http.ListenAndServe(":8080", nil)
// }
