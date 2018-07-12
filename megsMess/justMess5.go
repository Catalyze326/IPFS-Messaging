package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
	"os/exec"
    "os"
    "io"
    "strings"
		"github.com/ipfs/go-ipfs-api"
		//"bufio"
		"bytes"
)

var chatRoom string //topic[0]
var topic []string
var sh *shell.Shell //to be able to use pubsub
//var pubsubHash string
var recordData string //what is published in pubsub
var alreadySubscribed bool
var firstMessage bool //to help with appending the file with "has joined"
var output *shell.PubSubSubscription //to find out what is in pubsub
var err1 error //if there is error in getting message from pubsub
var output2 *shell.PubSubSubscription //to find out what is in pubsub
var err2 error //if there is error in getting message from pubsub
var message string //updated every time save is called with the new message
var username string //to save the username
//var userInput string //
var firstAppend bool //to help with appending the file correctly
var err3 error //if there is an error publishing
var saveHand bool //to help with saving messages

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
	To render the page with an html file.
**/
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
		t, _ := template.ParseFiles(tmpl + ".html")
		t.Execute(w, p)
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

		//check if already have a messaging.txt file
		_, err := os.Open("messaging.txt")
		if err != nil {
				//dont have the file, so create one
				os.Stderr.WriteString(err.Error())
				WriteStringToFile("messaging.txt", "Start chatting!")
		}

}

/**
	This handler will give you a textbox to enter your username.
**/
func editHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("edit handler")

		title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "newEdit", p) //open the html page

}

/**
	This will save the username that you typed into the edit handler
	as the global variable username.
**/
func userHandler(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/user/"):]
		body := r.FormValue("body")
		p := &Page{Title: title, Body: []byte(body)}
		err := p.save() //where it saves the new text

		username = string(p.Body) //set the username as a global variable

		fmt.Println(username)

		title = "messaging"

		//message = username + " has joined" + "--\r\n"
		http.Redirect(w, r, "/sub/"+title, http.StatusFound)

		if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func subHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("sub handler")

		if r.Method == "GET" {
				t, _ := template.ParseFiles("newSub.html")
				t.Execute(w, nil)
		} else {
				sub(w, r)
				// t, _ := template.ParseFiles("newSub.html")
				// t.Execute(w, nil)
	      fmt.Fprintf(w, "subscribed!")
		}

}

/**
	Function to subscribe you to a pubsub topic.
**/
func sub(w http.ResponseWriter, r *http.Request){

		fmt.Println("sub function")
	  r.ParseForm()
	  t, _ := template.ParseFiles("newSub.html")
	  t.Execute(w, nil)
	  topic = r.Form["subscribe"]
	  fmt.Println(topic)

		fmt.Println("about to subscribe")

		go subscribe()

		fmt.Println("subscribed")

}

func subscribe() {
		output, err1 = sh.PubSubSubscribe(topic[0])
		if err1 != nil {
			os.Stderr.WriteString(err1.Error())
		}

		output2, err2 = sh.PubSubSubscribe("fileChannel" + topic[0])
		if err2 != nil {
			os.Stderr.WriteString(err2.Error())
		}
}

/**
	This will write a string to a new text file of any name that you want.
	It will overwrite a text file if you already have one with that name.
**/
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

			//messageFile := readFile("messaging.txt")
			s := strings.Split(string(p.Body), "--")
			fmt.Println(s)

			sLen := len(s)

			sli := ""
			for i := 0; i < sLen; i++ {
					sli = sli + s[i] + " <br>"
			}

			fmt.Println("sli is " + sli)

			//renderTemplate(w, "messenger", p)
			fmt.Fprintf(w, "<head>" +
					"<meta http-equiv=%s content=%s />" +
					"</head>" +
				  "<body>" +
					"<h1>Chat room: %s</h1>" +
					"<p> %s </p>" +
					"<form action=\"/save/\" method=\"POST\">" +
					"<textarea name=\"body\"> </textarea>" +
					"<input type=\"submit\" value=\"Save\">" +
					"</form>" +
					"<form action=\"/refresh/\" method=\"POST\">" +
					"<input type=\"submit\" value=\"Refresh Page\">" +
					"</form>" +
					"</body>",

					"refresh", "5", topic[0], sli)

			// fmt.Fprintf(w, "<h1>Chat room: %s</h1>" +
			// 		"<p> %s </p>" +
			// 		"<form action=\"/save/\" method=\"POST\">" +
			// 		"<textarea name=\"body\"> </textarea>" +
			// 		"<input type=\"submit\" value=\"Save\">" +
			// 		"</form>" +
			// 		"<form action=\"/refresh/\" method=\"POST\">" +
			// 		"<input type=\"submit\" value=\"Refresh Page\">" +
			// 		"</form>",
			// 		topic[0], sli)

			// getPeers()
			//
			// getPBFile()
			//
			// AppendFile(recordData + "--")
			//
			// title := "messaging"

			// log.Fatal(http.ListenAndServe(":5001", nil))
			//
			// http.Redirect(w, r, "/messenger/"+title, http.StatusFound)

			title = "messaging"

			if saveHand == false {
					fmt.Println("save hand is false")
					if getPeers() > 0 {
							getPeersMessageConst()
							http.Redirect(w, r, "/refresh/"+title, http.StatusFound)
					}
			} else {
					fmt.Println("save hand is true")
			}
			saveHand = false

			//getPeersMessageConst(w, r)

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
				message = username + " has joined" + "--\r\n"
				fmt.Println("message " + message)
				firstMessage = false

		} else { //or else it will displace your message that you type in
				message = username + ": " + string(p.Body) + "--\r\n"
				fmt.Println("message " + message)
		}

		//it changes the recordData, but still appends the username function above

		// userInput = message
		//
    // WriteStringToFile("userInput.txt", message)

		sendFile(message) //publish the message in your pubsub chat room
		fmt.Println("after send file")
		AppendFile(message) //append the message to your messaging.txt file
		fmt.Println("after append")

    if err7 != nil {
        //http.Error(w, err7.Error(), http.StatusInternalServerError)
        //return
				os.Stderr.WriteString(err7.Error())
    }

		saveHand = true

		title = "messaging"

		http.Redirect(w, r, "/messenger/"+title, http.StatusFound)

}

/**
	Return the amount of peers in your chat room.
**/
func getPeers() int {
		fmt.Println("get peers")
		output1, err1 := exec.Command("ipfs", "pubsub", topic[0], "peers").CombinedOutput()
		if err1 != nil {
			os.Stderr.WriteString(err1.Error())
		}
		swarmPeersOutput := string(output1)

		peers := strings.Split(swarmPeersOutput, "\n")
		fmt.Println("peers")
		fmt.Println(peers)

 		numOfPeers := len(peers)
		fmt.Println( "numOfPeers ")
		fmt.Println(numOfPeers)

		listPeers := []string{}
		peersLength := 0

		for i := range peers {

				listPeers = append(listPeers, peers[i])
				fmt.Println(listPeers[i])

				peersLength = len(listPeers) - 1
				fmt.Println("peers length ")
				fmt.Println(peersLength)
		}

		return peersLength
}

/**
	This will get the message that was published to your chat room and
	update the record data to that message.
**/
func getPBFile(){

			output2, err3:= output.Next()
			if err3 != nil {
				os.Stderr.WriteString(err3.Error())
			}
			//continue;
			pubsubHash := string(output2.Data())
			fmt.Println(pubsubHash)

			recordPeer := output2.From()
			recordData = string(output2.Data())
			recordNo := output2.SeqNo()
			recordTopics := output2.TopicIDs()
			fmt.Println("record peer" + string(recordPeer))
			fmt.Println("record data " + string(recordData))
			fmt.Println("record number " + string(recordNo))
			fmt.Println("record topics ")
			for i := range recordTopics {
					fmt.Println(recordTopics[i])
			}

			fmt.Println(recordData)

			AppendFile(recordData + "--")

}

/**
	This will add the new message string into your messaging.txt file.
**/
func AppendFile(message string) {

	if firstAppend == true {
			fmt.Println("first append is true")
	} else {
			fmt.Println("first append is false")
	}

		  if firstAppend == true {
			WriteStringToFile("newMessage.txt", message)

			buf, err4 := os.Open("newMessage.txt")

			hash, err4 := sh.Add(buf)
			if err4 != nil {
				 os.Stderr.WriteString(err4.Error())
			}
			fmt.Println("hash is " + hash)

		 theReader, err3 := sh.Cat("/ipfs/" + string(hash))
		 if err3 != nil {
				os.Stderr.WriteString(err3.Error())
		 }

		 buf3 := new(bytes.Buffer)
		 buf3.ReadFrom(theReader)
		 newOutput := buf3.String()

		 WriteStringToFile("messaging.txt", newOutput)

		 firstAppend = false

		} else {

			 WriteStringToFile("newMessage.txt", message)

			 buf, err4 := os.Open("newMessage.txt")

			 hash, err4 := sh.Add(buf)
			 if err4 != nil {
			 	 os.Stderr.WriteString(err4.Error())
			 }
			 fmt.Println("hash is " + hash)

			 theReader, err3 := sh.Cat("/ipfs/" + string(hash))
			 if err3 != nil {
					os.Stderr.WriteString(err3.Error())
			 }

			 buf3 := new(bytes.Buffer)
			 buf3.ReadFrom(theReader)
			 newOutput := buf3.String()

			 buf2, _ := os.Open("messaging.txt")
			 messagingHash, err2 := sh.Add(buf2)

				if err2 != nil {
						os.Stderr.WriteString(err2.Error())
				}

			 catReader, err5 := sh.Cat("/ipfs/" + string(messagingHash))
			 if err5 != nil {
					os.Stderr.WriteString(err5.Error())
			 }

			 buf4 := new(bytes.Buffer)
			 buf4.ReadFrom(catReader)
			 catMessage := buf4.String()
			 fmt.Println("cat message messaging.txt is " + catMessage)

			 var newMessagingText string = string(catMessage) + string(newOutput)
			 WriteStringToFile("messaging.txt", newMessagingText)
			 fmt.Println("new messaging txt is " + newMessagingText)

	 }

}

/**
	This will wait for someone to publish something in the chat room, and
	it will change the record data and append the other person's message into
	your messaging.txt file.
**/
func refreshHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("refresh handler")

		getPeersMessageConst()

		title := "messaging"

		http.Redirect(w, r, "/messenger/"+title, http.StatusFound)
}

func getPeersMessageConst() {

		saveHand = false

		getPeers()

		//idea to run go getPBFile() in background

		go getPBFile()

}

/**
	This will publish your message to the chat room, and also update your
	record data to that message.
**/
func sendFile(theMessage string){

		err3 = sh.PubSubPublish(topic[0], theMessage)
		if err3 != nil {
			os.Stderr.WriteString(err3.Error())
		}

		fmt.Println("testing point 1")

		// output2, err3:= output.Next()
		// if err3 != nil {
		// 	os.Stderr.WriteString(err3.Error())
		// }
		//
		// fmt.Println("testing point 2")
		// pubsubHash := string(output2.Data())
		// fmt.Println(pubsubHash)
		//
		// recordPeer := output2.From()
		// recordData = string(output2.Data())
		// recordNo := output2.SeqNo()
		// recordTopics := output2.TopicIDs()
		// fmt.Println("record peer" + string(recordPeer))
		// fmt.Println("record data " + string(recordData))
		// fmt.Println("record number " + string(recordNo))
		// fmt.Println("record topics ")
		// for i := range recordTopics {
		// 		fmt.Println(recordTopics[i])
		// }
		//
		// fmt.Println("testing point 3")
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



func main(){

		go daemon()
		//sh = shell.NewShell("localhost:5001")
		//go sub()

		http.HandleFunc("/view/", viewHandler)
		http.HandleFunc("/sub/", subHandler)
		http.HandleFunc("/save/", saveHandler)
		http.HandleFunc("/edit/", editHandler)
		http.HandleFunc("/messenger/", messengerHandler)
		http.HandleFunc("/user/", userHandler)
		http.HandleFunc("/refresh/", refreshHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
}
