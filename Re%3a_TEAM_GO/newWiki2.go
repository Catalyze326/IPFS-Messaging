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
var hash string
var userInput string
var message string
var list []string
var start bool
var username string
var topic []string
var topic2 []byte
var sh *shell.Shell
var pubsubHash string
var haveMessagingFile bool
var recordData string
var output *shell.PubSubSubscription
var err1 error
var firstMessage bool
var getMess bool
var firstAppend bool

/**
	 The struct for the page that includes
	 the title and the body of the page.
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
		t, _ := template.ParseFiles(tmpl + ".html")
		t.Execute(w, p)
}

/**
	This will execute when the file ends in /view.
	This will write a new string to file navymessenger, and
	then make the title that so it will load the page with
	that text file. This will also run the daemon.
**/
func viewHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("view handler")

		WriteStringToFile("NAVY Messenger.txt", "Enter your username")

		title := "NAVY Messenger"
    p, _ := loadPage(title)
    renderTemplate(w, "view", p)

		output, err1 := exec.Command("ipfs", "daemon", "--enable-pubsub-experiment").Output()
        if err1 != nil {
          os.Stderr.WriteString(err1.Error())
	    }
			fmt.Println(output)

		start = true
		haveMessagingFile = false
		firstMessage = true
		getMess = true
		firstAppend = true

		sh = shell.NewShell("localhost:5001")

			//http.Redirect(w, r, "/sub/"+title, http.StatusFound)
}

/**
	This will execute when it ends in /edit. This
	will load the page with what is after /edit.
**/
func editHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("edit handler")

		title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)

}

func userHandler(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/user/"):]
		body := r.FormValue("body")
		p := &Page{Title: title, Body: []byte(body)}
		err := p.save() //where it saves the new text

		username = string(p.Body)

		fmt.Println(username)

		title = "messaging"

		//message = username + " has joined" + "--\r\n"
		http.Redirect(w, r, "/sub/"+title, http.StatusFound)

		if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

/**
	This will execute when any save button is pressed.
	This will take the word after /save as the title.
	The title will give you the body, and save what you
	typed in. It will take what you typed in as a message
	and save that message in userInput.txt. It will call
	the command line to add userInput to IPFS, and then
	it will give you the hash back. AppendFile is called,
	to add it to messaging.txt, and then redirects you
	to /messenger/messaging to call the messaging.txt file.
**/
func saveHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("save handler")
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save() //where it saves the new text


		if firstMessage == true {
				message = username + " has joined" + "--\r\n"
				fmt.Println("message " + message)
				firstMessage = false

		} else {
				message = username + ": " + string(p.Body) + "--\r\n"
				fmt.Println("message " + message)
		}

		userInput = message

    WriteStringToFile("userInput.txt", message)

    buf, err4 := os.Open("userInput.txt")

    hash, err4 := sh.Add(buf)
    if err4 != nil {
       os.Stderr.WriteString(err4.Error())
    }
    fmt.Println("hash is " + hash)

		// if haveMessagingFile == true {
		// 		fmt.Println("have messaging file true")
		// 		getMessagingFile()
		// 		//AppendFile(string(hash))
		// } else {
		// 		fmt.Println("have messaging file false")
		// 		getMessagingFile()
		// 		AppendFile(string(hash))
		// }

		if haveMessagingFile == false {
			getMessagingFile()
			AppendFile(hash)
		}

		AppendFile(hash)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

		title = "messaging"

		http.Redirect(w, r, "/messenger/"+title, http.StatusFound)

}

/**
	This is called when it ends in /messenger. This will swarm
	peers and getHashes to get the hashes of everyone in the local
	directory. If there is someone else in the chat, i.e. there
	is a peer, then it will find their messaging.txt file. If not,
	then you create a new messaging.txt file.
	)
**/
func messengerHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("messenger handler")

		//getMessagingFile()

		fmt.Println("afer get messaging file")

		title := r.URL.Path[len("/messenger/"):]
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
		fmt.Fprintf(w, "<h1>Chat room: %s</h1>" +
				"<p> %s </p>" +
				"<form action=\"/save/\" method=\"POST\">" +
				"<textarea name=\"body\"> </textarea>" +
				"<input type=\"submit\" value=\"Save\">" +
				"</form>",
				topic[0], sli)

		if getPeers() > 0 {

			buf, _ := os.Open("messaging.txt")

			messHash, err2 := sh.Add(buf)
			if err2 != nil {
				 os.Stderr.WriteString(err2.Error())
			}

			fmt.Println(messHash)

			pubsubTheHash(messHash)

		}

}

func readFile(file string) string {
    b, err := ioutil.ReadFile(file) // just pass the file name
    if err != nil {
        fmt.Print(err)
    }
    return string(b)
}

/**
	Either write or get a messaging file.
**/
func getMessagingFile() {
		fmt.Println("get messaging file")
		// peersLength := getPeers()
		//
		// if start == true {
		// 		fmt.Println("start is true")
		//
		// 		if peersLength == 0 {
		// 				WriteStringToFile("messaging.txt", message)
		// 		} else {
		// 				pubsubGetMessaging()
		// 		}
		// 		start = false
		// } else {
		// 		fmt.Println("start is false")
		// 		//pubsubGetMessaging()
		// }

		// if getMess == true {
		// 		fmt.Println("message is " + message)
		// 		WriteStringToFile("messaging.txt", message)
		// 		getMess = false
		// } else {
		// 		//pubsubGetMessaging()
		// }

		/**
			When you first log in to the program, if no one else is in the
			chat then make a messaging file. Otherwise, if there is no one in the
			chat, you already have a messaging file. Constantly check for the
			new messaing file from someone else, i.e. if the recordData changes,
			then get that hash.
		**/

		if haveMessagingFile == false {
				if getPeers() == 0 {
						WriteStringToFile("messaging.txt", userInput)
				} else {
						pubsubGetMessaging()
				}
		} else {
				if getPeers() != 0 {
						pubsubGetMessaging()
				}
		}

		haveMessagingFile = true

}

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
	This will use IPFS through the command line to take the userInput
	and add it to the messaging.txt file to update it.
**/
func AppendFile(hash string) {

		if firstAppend == true {
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

				 var newMessagingText string = string(catMessage) + string(newOutput)
				 WriteStringToFile("messaging.txt", newMessagingText)
		}
}

/**
  This will create a text file with whatever the string
	is as the body.
  )
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
	Get the hash from someone else, by getting it from pubsub.
**/
func pubsubGetMessaging() {

		firstAppend = false

		output2, err3:= output.Next()
		if err3 != nil {
			os.Stderr.WriteString(err3.Error())
		}
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

		getHash, err2 := exec.Command("cmd","/C", "ipfs", "get", recordData).Output()
					if err2 != nil {
						os.Stderr.WriteString(err2.Error())
						}
		fmt.Println("getHash " +string(getHash))


		catHash, err3 := exec.Command("cmd","/C", "ipfs", "cat", recordData, ">messaging.txt").Output()
					if err3 != nil {
						os.Stderr.WriteString(err3.Error())
						}
		fmt.Println("catHash " + string(catHash))

}

func pubsubTheHash(thehash string) {

		err2 := sh.PubSubPublish(topic[0], thehash)
		if err2 != nil {
			os.Stderr.WriteString(err2.Error())
		}

		output2, err3:= output.Next()
		if err3 != nil {
			os.Stderr.WriteString(err3.Error())
		}
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
}

func subHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("sub handler")

		if r.Method == "GET" {
				t, _ := template.ParseFiles("sub.html")
				t.Execute(w, nil)
		} else {
				sub(w, r)
	      fmt.Fprintf(w, "subscribed!")
		}


}

func sub(w http.ResponseWriter, r *http.Request){
	fmt.Println("sub function")
  r.ParseForm()
  t, _ := template.ParseFiles("sub.html")
  t.Execute(w, nil)
  topic = r.Form["subscribe"]
  fmt.Println(topic)

	fmt.Println("about to subscribe")

	output, err1 = sh.PubSubSubscribe(topic[0])
	if err1 != nil {
		os.Stderr.WriteString(err1.Error())
	}
	fmt.Println("subscribed")



	// subVar, subErr := sh.PubSubSubscribe(topic[0])
	// if subErr != nil {
	// 	os.Stderr.WriteString(subErr.Error())
	// }
	//
	// subRecord, nextErr := subVar.Next()
	// if nextErr != nil {
	// 	os.Stderr.WriteString(nextErr.Error())
	// }
	// recordPeer := subRecord.From()
	// recordData = string(subRecord.Data())
	// recordNo := subRecord.SeqNo()
	// recordTopics := subRecord.TopicIDs()
	// fmt.Println("record peer" + string(recordPeer))
	// fmt.Println("record data " + string(recordData))
	// fmt.Println("record number " + string(recordNo))
	// fmt.Println("record topics ")
	// for i := range recordTopics {
	// 		fmt.Println(recordTopics[i])
	// }

}


func main() {

	  http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
		http.HandleFunc("/messenger/", messengerHandler)
		http.HandleFunc("/sub/", subHandler)
		http.HandleFunc("/user/", userHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
