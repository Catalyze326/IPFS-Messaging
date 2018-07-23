package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", dog)
	http.HandleFunc("/toby.jpg", dogPic)
	http.HandleFunc("/java", java)
	http.ListenAndServe(":8080", nil)
}

func dog(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<html>
	<body>

	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
	<script type = "text/javascript" src="https://code.jquery.com/jquery-3.3.1.js"></script>

	<script type="text/javascript" src="http://localhost:8080/java">	</script>

	<div class="button">
	    <input type="button" id="lesen" value="Lesen!" />
	</div>

	<div class="text">
	    Lorem Ipsum <br>
	</div>
	</body>

<img src="toby.jpg">

	</html>	`)
}

func dogPic(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "toby.jpg")
}

func java(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "java.js")

}
