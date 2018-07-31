package main

import (

    "fmt"
    "net/http"
)

var saveHand bool
var firstMessage bool
var firstAppend bool
var alreadySubscribed bool
var saveHand bool
var newSaveHand bool
var count int
var countAppend int


/**
	The page struct that includes a title and a body.
**/
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
	Start the daemon.
**/
func daemon() {
		output, err1 := exec.Command("ipfs", "daemon", "--enable-pubsub-experiment").Output()
			if err1 != nil {
				os.Stderr.WriteString(err1.Error())
		}
		fmt.Println(output)
}

/**
	The first page you see; it iwll start a daemon, and initialize global
	variables. Subscribes you to a pubsub, but you unsubscribe when you
	go to a different handler.
**/
func viewHandler(w http.ResponseWriter, r *http.Request) {

		//make a file so that when you load the page, it will call a txt file
		//instead of giving you an error
		WriteStringToFile("newSub.txt", "Enter your username")

		title := "newSub"
		p, _ := loadPage(title)
		renderTemplate(w, "newViewSub", p)

		//initialize the shell for pubsub
		sh = shell.NewShell("localhost:5001")

		//when the page loads, it wants to load the messaging.txt file
		//so set title to messaging so load page will call messaging.txt
		title = "messaging"
		firstMessage = true
		firstAppend = true
		alreadySubscribed = false
		saveHand = false
		newSaveHand = true
		count = 0
		countAppend = 0

		//check if already have a messaging.txt file
		_, err := os.Open("messaging.txt")
		if err != nil {
				//dont have the file, so create one
				os.Stderr.WriteString(err.Error())
				WriteStringToFile("messaging.txt", "Start chatting!")
		}

}

/**
	This handler will save the message you typed into the message page, and
	it will append the file. It will also send the message to the pubsub topic
	you are subscribed to.
**/
func saveHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("save handler")
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err7 := p.save() //where it saves the new text as p.Body

		//this will display username "has joined" when you first enter the chat
		if firstMessage == true {
				fmt.Println("first message is true")
				message = username + " has joined" + "--\r\n"
				fmt.Println("message " + message)
				firstMessage = false

		} else { //or else it will displace your message that you type in
				fmt.Println("first message is true")
				message = username + ": " + string(p.Body) + "--\r\n"
				fmt.Println("message " + message)
		}

		//it changes the recordData, but still appends the username function above

		// userInput = message
		//
    // WriteStringToFile("userInput.txt", message)

		// sendFile(message) //publish the message in your pubsub chat room
		// fmt.Println("after send file")
		// AppendFile(message) //append the message to your messaging.txt file
		// fmt.Println("after append")

		//recordDataHolder = recordData

    if err7 != nil {
        //http.Error(w, err7.Error(), http.StatusInternalServerError)
        //return
				os.Stderr.WriteString(err7.Error())
    }

		// fmt.Println("count is " + string(count))
		// if count > 1 {
		// 	 saveHand = true
		// 	 newSaveHand = true
		// }
		// count ++

		saveHand = true

		title = "messaging"

		http.Redirect(w, r, "/messenger/"+title, http.StatusFound)

}

/**
	This handler will display the messaging.txt file that you have.
**/
func messengerHandler(w http.ResponseWriter, r *http.Request) {

			fmt.Println("messenger handler")

			//title := r.URL.Path[len("/messenger/"):]
			title := "messaging"
			p, err := loadPage(title)
			if err != nil {
					p = &Page{Title: title}
			}

			fmt.Println("hefeowafaosief")



			//messageFile := readFile("messaging.txt")
			s := strings.Split(string(p.Body), "--")
			fmt.Println(s)

			sLen := len(s)

			sli := ""
			for i := 0; i < sLen; i++ {
					sli = sli + s[i] + " <br>"
			}

			fmt.Println("sli is " + sli)

			renderTemplate(w, "newTryMessenger", p)

			fmt.Fprintf(w, "<body>" + "<p> %s </p>" + "</body>", "hi there")

				// fmt.Fprintf(w,
				// 	  "<body>" +
				// 		"<h1>Chat room: %s</h1>" +
				// 		"<p> %s </p>" +
				// 		"<form action=\"/save/\" method=\"POST\">" +
				// 		"<textarea name=\"body\"> </textarea>" +
				// 		"<input type=\"submit\" value=\"Save\">" +
				// 		"</form>" +
				// 		"<form action=\"/upload/\">" +
	      //     "<input type=\"submit\" value='Upload a File'>"+
	      //     "</form>"+
				// 		"<form action=\"/recordFile/\">" +
	      //     "<input type=\"submit\" value='Record Audio'>"+
	      //     "</form>"+
	      //     "<form action=\"/exit/{{.Title}}\">" +
	      //     "<input type=\"submit\" value='EXIT'>"+
	      //     "</form>"+
				// 		"</body>",
				// 		topic[0], sli)

			//r.FormValue


			title = "messaging"

			//recordDataHolder := recordData

			if saveHand == false {
					fmt.Println("save hand is false")
					// if getPeers() > 0 {
					// 		getPeersMessageConst()
					// 		if recordDataHolder != recordDataHolder {
					// 				http.Redirect(w, r, "/refresh/", http.StatusFound)
					// 		}
					// }


					ticker := time.NewTicker(5 * time.Second)
					quit := make(chan struct{})
					go func() {
					    for {
					       select {
					        case <- ticker.C:
					            // do stuff
											getPeersMessageConst(w, r)
					        case <- quit:
					            ticker.Stop()
					            return
					        }
					    }
					 }()

					saveHand = true
					count = 0
					newSaveHand = true

					fmt.Println("it works!")

					fmt.Fprintf(w, "<head>" +
							"<meta http-equiv=%s content=%s />" +
							"</head>", "refresh", "15")

					//window.location.reload(true);

					// fmt.Fprintf(setInterval(function() {
					// 	$("#reloadContent").load(location.href+" #reloadContent>*","")
					// }, 200000))


					//http.Redirect(w, r, "/save/", http.StatusFound)

					// for {
					// 		getPeersMessageConst(w, r)
					// 		fmt.Println("record data holder is " + recordDataHolder)
					// 		fmt.Println("record data is " + recordData)
					// 		if recordDataHolder != recordData {
					// 			fmt.Println("NOT THE SAME***********")
					// 			saveHand = true
					// 			count = 0
					// 			newSaveHand = true
					// 			fmt.Println("what about here")
					// 			// fmt.Fprintf(w, "<head>" +
					// 			// 		"<meta http-equiv=%s content=%s />" +
					// 			// 		"</head>", "refresh", "2")
					// 			//http.Redirect(w, r, "/save/", http.StatusFound)
					//
					//
					//
					//
					// 				break
					// 		}
					// 		fmt.Println("EHASEOFHAEOS")
					// 		fmt.Fprintf(w, "<head>" +
					// 				"<meta http-equiv=%s content=%s />" +
					// 				"</head>", "refresh", "2")
					// 		break
					// 		recordDataHolder = recordData
					//}

			} else {

					fmt.Println("save hand is true")
						saveHand = false
						fmt.Println("save hand is false")
						if newSaveHand == true {
							fmt.Fprintf(w, "<head>" +
									"<meta http-equiv=%s content=%s />" +
									"</head>", "refresh", "1")
								newSaveHand = false
						} else {
							fmt.Fprintf(w, "<head>" +
									"<meta http-equiv=%s content=%s />" +
									"</head>", "refresh", "50")
						}

						//http.Redirect(w, r, "/save/", http.StatusFound)
						//newSaveHand = false
			}


			//getPeersMessageConst(w, r)

}

func main(){
  go daemon()
  http.HandleFunc("/login/", loginHandler)
}
