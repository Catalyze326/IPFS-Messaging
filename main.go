package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"github.com/ipfs/go-ipfs-api"
)

var topic []string
var sh *shell.Shell //to be able to use pubsub
//var pubsubHash string
var recordData string //what is published in pubsub
var recordDataHolder string
var output *shell.PubSubSubscription  //to find out what is in pubsub
var err1 error                        //if there is error in getting message from pubsub
var output2 *shell.PubSubSubscription //to find out what is in pubsub
var err2 error                        //if there is error in getting message from pubsub
var message string                    //updated every time save is called with the new message
var username string                   //to save the username
var err3 error                        //if there is an error publishing
var uploadFileName string
var uploadAudioName string
var fileHash []byte
var quit chan bool
var first bool

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
	Start the daemon.
**/
func daemon() {
	output, err1 := exec.Command("ipfs", "daemon", "--enable-pubsub-experiment").Output()
	if err1 != nil {
		os.Stderr.WriteString(err1.Error())
	}
	fmt.Println(string(output))
}

//ALL HANDLERS ************************************************

/**
	The first page you see; it iwll start a daemon, and initialize global
	variables. Subscribes you to a pubsub, but you unsubscribe when you
	go to a different handler.
**/
func viewHandler(w http.ResponseWriter, r *http.Request) {

	//make a file so that when you load the page, it will call a txt file
	//instead of giving you an error
	WriteStringToFile("Lattice.txt", "Enter your username")

	title := "Lattice"
	p, _ := loadPage(title)
	//renderTemplate(w, "newViewSub", p)
	renderTemplate(w, "cViewNew", p)

	//initialize the shell for pubsub
	sh = shell.NewShell("localhost:5001")

	//when the page loads, it wants to load the messaging.txt file
	//so set title to messaging so load page will call messaging.txt
	title = "messaging"
	first = true

}

func aboutUsAfterHandler(w http.ResponseWriter, r *http.Request) {

	WriteStringToFile("AboutUs.txt", "Here's us!")

	title := "AboutUs"
	p, _ := loadPage(title)

	renderTemplate(w, "aboutUs", p)

}

func aboutUsBeforeHandler(w http.ResponseWriter, r *http.Request) {

	WriteStringToFile("AboutUs.txt", "Here's us!")

	title := "AboutUs"
	p, _ := loadPage(title)

	renderTemplate(w, "aboutUsBefore", p)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("cLoginNew.html")
		t.Execute(w, nil)
	} else {
		username = r.FormValue("username") // stores the input username and password as variables
		password := r.FormValue("password")
		fmt.Println(username)
		fmt.Println(password)
		attempt := login(username, password) //runs the login function with inputs username and password
		if attempt == true {                 //if the user login info is correct, page redirects to chat page
			fmt.Println("redirecting")
			http.Redirect(w, r, "/sub/", 301)
		} else {
			t, _ := template.ParseFiles("cLoginNew.html") // if the login info is incorrect, resets
			t.Execute(w, nil)
			fmt.Fprintf(w, "incorrect username and password")
		}
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("cAdminNew.html")
		t.Execute(w, nil)
	} else {
		adminUsername := r.FormValue("adminUsername")
		adminPassword := r.FormValue("adminPassword")
		fmt.Println(adminUsername)
		fmt.Println(adminPassword)
		admin := adminLogin(adminUsername, adminPassword)
		if admin == true {
			http.Redirect(w, r, "/admin/addaccount/", 301)
		} else {
			t, _ := template.ParseFiles("cAdminNew.html")
			t.Execute(w, nil)
			fmt.Fprintf(w, "incorrect admin username and password")
		}
	}
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("cSignUpNew.html")
		t.Execute(w, nil)
	} else {
		t, _ := template.ParseFiles("cSignUpNew.html")
		t.Execute(w, nil)
		username := r.FormValue("newUsername")
		password := r.FormValue("newPassword")
		password2 := r.FormValue("newPassword2")
		if password == password2 {
			addAccount(username, password)
			fmt.Fprintf(w, "account successfully added")
		} else {
			fmt.Fprintf(w, "please ensure you re-entered password correctly")
		}
	}
}

/**
	This will subscribe you to the topic you type in
**/

func subHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("sub handler")

	if r.Method == "GET" {
		t, _ := template.ParseFiles("cNewSubNew.html")
		t.Execute(w, nil)
	} else {
		sub(w, r)

		fmt.Fprintf(w, "subscribed!")

		colors := []string{"aqua", "aquamarine", "azure", "bisque",
			"burlywood", "cadetblue", "chartreuse", "chocolate",
			"coral", "cornflowerblue", "cyan", "darkkhaki",
			"darkorange", "darksalmon", "fuchsia", "gold",
			"greenyellow", "green", "hotpink", "lawngreen",
			"lightcyan", "lemonchiffon", "lightgreen", "lime",
			"mediumorchid", "navajowhite", "orange", "palegreen",
			"pink", "red", "seagreen", "royalblue", "springgreen",
			"thistle", "yellow", "tomato"}

		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s) // initialize local pseudorandom generator
		r.Intn(len(colors))

		theColor := colors[r.Intn(len(colors))]

		currTime := time.Now().Format("15:04:05")
		fmt.Println(currTime)

		currDate := time.Now().Format("Mon Jan _2")
		fmt.Println(currDate)

		addDate := "<font color='lightgray'>" + string(currDate) + "</font><br>"

		username = "<font color=" + theColor + ">" + username + "</font>"

		WriteStringToFile("messaging.txt", string(addDate)+"<font color='lightgray'>"+string(currTime)+" | </font>"+username+"<font color='white'>"+" has joined </font><br>")

		http.HandleFunc("/java/", java)

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

	currTime := time.Now().Format("15:04:05")
	fmt.Println(currTime)

	message = "<font color='lightgray'>" + string(currTime) + " | </font>" + username + "<font color='white'>" + ": " + string(p.Body) + "</font><br>"
	//message = username + "<font color='black'>" + ": " + string(p.Body) + "</font><br>"
	fmt.Println("message " + message)

	sendFile(message) //publish the message in your pubsub chat room
	fmt.Println("after send file")

	AppendFile(message) //append the message to your messaging.txt file
	fmt.Println("after append")

	recordDataHolder = recordData

	if err7 != nil {
		//http.Error(w, err7.Error(), http.StatusInternalServerError)
		//return
		os.Stderr.WriteString(err7.Error())
	}

	//title = "messaging"
	//title = "index/"

	http.Redirect(w, r, "/messenger/", http.StatusFound)

	first = true

}

/**
	This handler will display the messaging.txt file that you have.
**/
func messengerHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hello")
	WriteStringToFile("Lattice.txt", "Enter your username")

	//title := "Lattice"
	title := "messaging"
	p, _ := loadPage(title)

	fmt.Println("jquery handler")

	s := strings.Split(string(p.Body), "--")
	fmt.Println(s)

	sLen := len(s)

	sli := ""
	for i := 0; i < sLen; i++ {
		sli = sli + s[i] + " <br>"
	}

	fmt.Println("sli is " + sli)

	renderTemplate(w, "index", p)

	fmt.Println("r method is " + r.Method)

	if first == false {

		if r.Method == "GET" {

			urlName := r.URL.Path

			fmt.Println("url name is " + urlName)

			fileName := r.FormValue("file")

			fmt.Println("file name is " + fileName)

			uploadFileName = fileName

			buf, err4 := os.Open(uploadFileName)

			otherFileHash, err4 := sh.Add(buf)
			if err4 != nil {
				os.Stderr.WriteString(err4.Error())
			}
			fmt.Println("other file hash is " + otherFileHash)

			// sendFile(username + ":***" + uploadFileName + "***" + otherFileHash)

			//otherFileMessage := "<a href=\"/getOtherFile/\">" + uploadFileName + "</a>"
			theLink := "\"http://localhost/ipfs/" + string(otherFileHash) + "\""

			currTime := time.Now().Format("15:04:05")
			fmt.Println(currTime)

			//message = "<font color='lightgray'>" + string(currTime) + " | </font>" + username + "<font color='white'>" + ": " + string(p.Body) + "</font><br>"

			otherFileMessage := "<font color='lightgray'>" + string(currTime) + " | </font>" + username + "<font color='white'>: <a href=" + theLink + " target=\"_blank\"></font>" + uploadFileName + "</a><br>"

			sendFile(otherFileMessage)
			AppendFile(otherFileMessage)

			first = true

		}

	}

	fmt.Println("save hand is false")

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				getPeersMessageConst(w, r)
			case <-quit:
				fmt.Println("QUIT")
				ticker.Stop()
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}
	}()

	fmt.Println("it works!")

	first = false

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//var fileUp string
	//var err error

	fmt.Println("upload handler")

	//_, handler, _ := r.FormFile("uploadfile")

	if r.Method == "GET" {
		// GET
		t, _ := template.ParseFiles("upload.html")
		t.Execute(w, nil)

	} else if r.Method == "POST" {
		// Post
		r.ParseForm()

		_, handler, _ := r.FormFile("uploadfile")
		//fileUp = handler.Filename

		uploadFileName = handler.Filename
		fmt.Println("upload file name is " + uploadFileName)
		periodIndex := strings.Index(uploadFileName, ".")
		fmt.Println("period index is " + string(periodIndex))
		fileType := uploadFileName[periodIndex:]
		fmt.Println("file type is " + string(fileType))

		http.Redirect(w, r, "/uploadFile/"+uploadFileName, http.StatusFound)

	} else {
		fmt.Println("Unknown HTTP " + r.Method + "  Method")
	}

}

/**
	Handler to exit you from the chat room.
**/
func exitHandler(w http.ResponseWriter, r *http.Request) {

	exitMessage := username + " has left --"
	sendFile(exitMessage)

	fmt.Println("after true")

	output.Cancel()

	fmt.Println("after cancel")

	//quit = make(true)
	//quit <- true

	//title := "messaging"

	//PROBLEM HERE IS THE THREADING IN THE BACKGROUND DOESN'T STOP

	//http.Redirect(w, r, "/login/"+title, http.StatusFound)
	http.Redirect(w, r, "/", http.StatusFound)

}

//****fileName**hash

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	buf, err4 := os.Open(uploadFileName)

	otherFileHash, err4 := sh.Add(buf)
	if err4 != nil {
		os.Stderr.WriteString(err4.Error())
	}
	fmt.Println("other file hash is " + otherFileHash)

	// sendFile(username + ":***" + uploadFileName + "***" + otherFileHash)

	//otherFileMessage := "<a href=\"/getOtherFile/\">" + uploadFileName + "</a>"
	theLink := "\"http://localhost/ipfs/" + string(otherFileHash) + "\""
	otherFileMessage := username + ": <a href=" + theLink + " target=\"_blank\">" + uploadFileName + "</a><br>"

	sendFile(otherFileMessage)
	AppendFile(otherFileMessage)

	title := "messaging"

	http.Redirect(w, r, "/messenger/"+title, http.StatusFound)

}

func uploadAudioHandler(w http.ResponseWriter, r *http.Request) {

	//uploadFileName = "adidas.png"
	//periodIndex := strings.Index(uploadFileName, ".")
	//fileType := uploadFileName[periodIndex:]

	buf, err4 := os.Open(uploadAudioName)

	otherAudioHash, err4 := sh.Add(buf)
	if err4 != nil {
		os.Stderr.WriteString(err4.Error())
	}
	fmt.Println("other audio hash is " + otherAudioHash)

	//otherFileMessage := "<a href=\"/getOtherFile/\">" + uploadFileName + "</a>"
	theLink := "\"http://localhost/ipfs/" + string(otherAudioHash) + "\""

	// <audio controls>
	// <source src="https://www.computerhope.com/jargon/m/example.mp3" />
	// </audio>

	//otherAudioMessage := username + ": <a href=" + theLink + " target=\"_blank\">" + uploadAudioName + "</a><br>"
	otherAudioMessage := username + ": <audio controls> <source src=" + theLink + " />" + uploadAudioName + "</audio><br>"

	sendFile(otherAudioMessage)
	AppendFile(otherAudioMessage)

	title := "messaging"

	http.Redirect(w, r, "/messenger/"+title, http.StatusFound)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("audio.html")
		t.Execute(w, nil)
	} else {
		// StartRecord()
		//PlaySound()

	}
}

func func1Handler(w http.ResponseWriter, r *http.Request) {

	// StartRecord()
	fmt.Println("after recording")
	uploadAudioName = "navy.wav"
	//http.Redirect(w, r, "/uploadAudio/", http.StatusFound)

	buf, err4 := os.Open(uploadAudioName)

	otherAudioHash, err4 := sh.Add(buf)
	if err4 != nil {
		os.Stderr.WriteString(err4.Error())
	}
	fmt.Println("other audio hash is " + otherAudioHash)

	//otherFileMessage := "<a href=\"/getOtherFile/\">" + uploadFileName + "</a>"
	//theLink := "\"http://localhost/ipfs/" + string(otherAudioHash) + "\""

	// <audio controls>
	// <source src="https://www.computerhope.com/jargon/m/example.mp3" />
	// </audio>

	//otherAudioMessage := username + ": <a href=" + theLink + " target=\"_blank\">" + uploadAudioName + "</a><br>"
	//otherAudioMessage := username + ": <audio controls> <source src=" + theLink + " />" + uploadAudioName + "</audio><br>"

	//otherFileMessage := "<a href=\"/getOtherFile/\">" + uploadFileName + "</a>"
	theLink := "\"http://localhost/ipfs/" + string(otherAudioHash) + "\""

	currTime := time.Now().Format("15:04:05")
	fmt.Println(currTime)

	//message = "<font color='lightgray'>" + string(currTime) + " | </font>" + username + "<font color='white'>" + ": " + string(p.Body) + "</font><br>"

	//otherAudioMessage := "<font color='lightgray'>" + string(currTime) + " | </font>" + username + "<font color='white'>: <a href=" + theLink + " target=\"_blank\"></font>" + uploadAudioName + "</a><br>"
	otherAudioMessage := "<font color='lightgray'>" + string(currTime) + " | </font>" + username + ": <audio controls> <source src=" + theLink + " />" + uploadAudioName + "</audio><br>"

	sendFile(otherAudioMessage)
	AppendFile(otherAudioMessage)

	first = true

}

func func2Handler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		fmt.Println("yuhhhhhhhhhhhhhhhhh")
		r.ParseForm()

		file, handler, err := r.FormFile("uploadfile")

		fmt.Println(handler)
		fileUp := handler.Filename

		uploadFileName = handler.Filename

		fmt.Println("upload file name is " + uploadFileName)
		periodIndex := strings.Index(uploadFileName, ".")
		fmt.Println("period index is " + string(periodIndex))
		fileType := uploadFileName[periodIndex:]
		fmt.Println("file type is " + string(fileType))

		fmt.Println(fileUp)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)
		file.Close()

		fullPath, _ := filepath.Abs("./uploads/" + uploadFileName)
		fmt.Println("full path is " + fullPath)

		buf, err4 := os.Open(fullPath)

		otherFileHash, err4 := sh.Add(buf)
		if err4 != nil {
			os.Stderr.WriteString(err4.Error())
		}
		fmt.Println("other file hash is " + otherFileHash)

		// sendFile(username + ":***" + uploadFileName + "***" + otherFileHash)

		//otherFileMessage := "<a href=\"/getOtherFile/\">" + uploadFileName + "</a>"
		theLink := "\"http://localhost/ipfs/" + string(otherFileHash) + "\""
		otherFileMessage := username + ": <a href=" + theLink + " target=\"_blank\">" + uploadFileName + "</a><br>"

		sendFile(otherFileMessage)
		AppendFile(otherFileMessage)
	}
}

func func3Handler(w http.ResponseWriter, r *http.Request) {

	_, handler, _ := r.FormFile("uploadfile")
	//fileUp = handler.Filename

	uploadFileName = handler.Filename
	fmt.Println("upload file name is " + uploadFileName)
	periodIndex := strings.Index(uploadFileName, ".")
	fmt.Println("period index is " + string(periodIndex))
	fileType := uploadFileName[periodIndex:]
	fmt.Println("file type is " + string(fileType))

}

//ALL FUNCTIONS********************************************

func java(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "messaging.txt")

}

/**
	Function to subscribe you to a pubsub topic.
**/
func sub(w http.ResponseWriter, r *http.Request) {

	fmt.Println("sub function")
	r.ParseForm()
	t, _ := template.ParseFiles("cNewSubNew.html")
	t.Execute(w, nil)
	topic = r.Form["subscribe"]
	fmt.Println(topic)

	fmt.Println("about to subscribe")

	go subscribe()

	fmt.Println("subscribed to " + topic[0])

}

/**
	Second channel you are subscribed to for files
**/
func subscribe() {
	output, err1 = sh.PubSubSubscribe(topic[0])
	if err1 != nil {
		os.Stderr.WriteString(err1.Error())
	}

	// output2, err2 = sh.PubSubSubscribe("fileChannel" + topic[0])
	// if err2 != nil {
	// 	os.Stderr.WriteString(err2.Error())
	// }
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
	To log you in with your username and password
**/
func login(username string, password string) (login bool) {
	userExist := false           //initialize boolean status of user login validation
	f, _ := os.Open("login.csv") //open the precreated CSV file
	defer f.Close()              //defer the closing function to allow for appending
	r := csv.NewReader(f)        //initialize new reader
	rows, _ := r.ReadAll()       //turn CSV data into a 2 dimensional slice [row][column]
	fmt.Println(rows)
	fmt.Println(rows[0][0])
	fmt.Println(rows[0][1])
	for i := range rows {
		fmt.Println(rows[i][0], rows[i][1])
		if (rows[i][0] == username) && (rows[i][1] == password) { //if the user login info matches stored data
			fmt.Println("username: ", username)
			fmt.Println("password: ", password)
			fmt.Println("known user detected")
			fmt.Println("login successful")
			userExist = true //set validation to true
			login = true
			break
		}
	}
	if userExist == false { //if info doesn't match, set validation to false
		fmt.Println("incorrect username and password")
		login = false
	}
	return login // return validation boolean
}

/**
	To log in as an administrator
**/
func adminLogin(adminUsername string, adminPassword string) (admin bool) {
	f, _ := os.Open("admin.csv") //open admin CSV file
	defer f.Close()
	r := csv.NewReader(f)
	rows, _ := r.ReadAll()
	fmt.Println(rows[0][0])
	fmt.Println(rows[0][1])
	if (adminUsername == rows[0][0]) && (adminPassword == rows[0][1]) {
		admin = true
	} else {
		fmt.Println("incorrect admin username and password")
		admin = false
	}
	return admin
}

/**
	To add an account when you are the admin
**/
func addAccount(username string, password string) {
	f, _ := os.Open("login.csv")
	defer f.Close()
	r := csv.NewReader(f)
	rows, _ := r.ReadAll()
	rows = append(rows, []string{username, password}) //add new user login info to login info slice
	file, err := os.Create("login.csv")               //create new CSV file
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", "keyMap.csv", err.Error())
	}
	defer func() {
		e := f.Close()
		if e != nil {
			log.Fatalf("Cannot close '%s': %s\n", "keyMap.csv", e.Error())
		}
	}()
	w := csv.NewWriter(file) //initialize a data writer
	err = w.WriteAll(rows)   //write slice data to CSV file
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
	fmt.Println("numOfPeers ")
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
// func getPBFile(w http.ResponseWriter, r *http.Request) {
//
// 	output2, err3 := output.Next()
// 	if err3 != nil {
// 		os.Stderr.WriteString(err3.Error())
// 	}
// 	//continue;
// 	pubsubHash := string(output2.Data())
// 	fmt.Println("the pubsub hash is " + pubsubHash)
//
// 	recordPeer := output2.From()
// 	recordData = string(output2.Data())
// 	recordNo := output2.SeqNo()
// 	recordTopics := output2.TopicIDs()
// 	fmt.Println("record peer" + string(recordPeer))
// 	fmt.Println("record data " + string(recordData))
// 	fmt.Println("record number " + string(recordNo))
// 	fmt.Println("record topics ")
// 	for i := range recordTopics {
// 		fmt.Println(recordTopics[i])
// 	}
//
// 	fmt.Println(recordData)
//
// 	AppendFile(recordData)
//
// 	fmt.Println("hasoifaoiwejfoiawjeoifjaowiejfoijwaea;fioj")
// 	//http.Redirect(w, r, "/messenger/", http.StatusFound)
//
// 	fmt.Fprintf(w, "<head>"+
// 		"<meta http-equiv=%s content=%s />"+
// 		"</head>", "refresh", ".1")
//
// }

/**
	This will add the new message string into your messaging.txt file.
**/
func AppendFile(message string) {

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

	// buf2, _ := os.Open("messaging.txt")
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

/**
	Function to run getPBFile in the background
**/
func getPeersMessageConst(w http.ResponseWriter, r *http.Request) {

	getPeers()

	//idea to run go getPBFile() in background

	// go getPBFile(w, r)

}

func parseFileChannel(input string) (string, string) {
	first := strings.Index(input, "***")
	last := strings.LastIndex(input, "***")
	fileName := string(input[first+3 : last])
	hash := string(input[last+3:])

	return fileName, hash
}

/**
	What jens wrote to record an audio file
**/
// var (
// 	winmm         = syscall.MustLoadDLL("winmm.dll")
// 	mciSendString = winmm.MustFindProc("mciSendStringW")
// 	stop          bool
// )
//
// func MCIWorker(lpstrCommand string, lpstrReturnString string, uReturnLength int, hwndCallback int) uintptr {
// 	i, _, _ := mciSendString.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpstrCommand))),
// 		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpstrReturnString))),
// 		uintptr(uReturnLength), uintptr(hwndCallback))
// 	return i
// }
//
// func StartRecord() {
// 	fmt.Println("windmm.dll Record Audio to .wav file.")
//
// 	i := MCIWorker("open new type waveaudio alias capture", "", 0, 0)
// 	if i != 0 {
// 		log.Fatal("Error Code A: ", i)
// 	}
//
// 	i = MCIWorker("record capture", "", 0, 0)
// 	if i != 0 {
// 		log.Fatal("Error Code B: ", i)
// 	}
//
// 	fmt.Println("Listening...")
// 	time.Sleep(10 * time.Second)
//
// 	i = MCIWorker("save capture navy.wav", "", 0, 0)
// 	if i != 0 {
// 		log.Fatal("Error Code C: ", i)
// 	}
//
// 	i = MCIWorker("close capture", "", 0, 0)
// 	if i != 0 {
// 		log.Fatal("Error Code D: ", i)
// 	}
// 	fmt.Println("saved to navy.wav")
//
// }
//
// func PlaySound() {
// 	f, _ := os.Open("navy.wav")
// 	s, format, _ := wav.Decode(f)
// 	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
// 	speaker.Play(s)
// 	select {}
// 	done := make(chan struct{})
// 	speaker.Play(beep.Seq(s, beep.Callback(func() {
// 		close(done)
// 	})))
// 	<-done
// }

/**
	This will publish your message to the chat room, and also update your
	record data to that message.
**/
func sendFile(theMessage string) {

	err3 = sh.PubSubPublish(topic[0], theMessage)
	if err3 != nil {
		os.Stderr.WriteString(err3.Error())
	}

	fmt.Println("testing point 1")

}

func main() {

	go daemon()
	var err error
	err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:8080").Start()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/func1/", func1Handler)
	http.HandleFunc("/func2/", func2Handler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/aboutUsAfter/", aboutUsAfterHandler)
	http.HandleFunc("/aboutUsBefore/", aboutUsBeforeHandler)
	http.HandleFunc("/exit/", exitHandler)
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/uploadFile/", uploadFileHandler)
	http.HandleFunc("/uploadAudio/", uploadAudioHandler)
	http.HandleFunc("/recordFile/", homeHandler)
	http.HandleFunc("/admin/addaccount/", accountHandler)
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/upload/", uploadHandler)
	http.HandleFunc("/sub/", subHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/messenger/", messengerHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
