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
		"time"
)

var chatRoom string //topic[0]
var sh *shell.Shell
//var pubsubHash string
var haveMessagingFile bool
var recordData string
var alreadySubscribed bool
var firstMessage bool
var topic []string
var output *shell.PubSubSubscription
var err1 error
var message string
var username string
var userInput string
var firstAppend bool
var entered bool
var err2 error
var firstPB bool
var firstTry bool

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

func subHandler(w http.ResponseWriter, r *http.Request) {

		WriteStringToFile("newSub.txt", "Enter your username")

		title := "newSub"
		p, _ := loadPage(title)
		renderTemplate(w, "newViewSub", p)

		output, err1 := exec.Command("ipfs", "daemon", "--enable-pubsub-experiment").Output()
			if err1 != nil {
				os.Stderr.WriteString(err1.Error())
		}
		fmt.Println(output)

		//initialize the shell for pubsub
		sh = shell.NewShell("localhost:5001")

		//when the page loads, it wants to load the messaging.txt file
		//so set title to messaging so load page will call messaging.txt
		title = "messaging"
		firstMessage = true
		firstPB = true
		firstTry = true
		firstAppend = true

		//check if already have a messaging.txt file
		_, err := os.Open("messaging.txt")
		if err != nil {
				//dont have the file, so create one
				os.Stderr.WriteString(err.Error())
				WriteStringToFile("messaging.txt", "Chat in dubsonly")
		}

		if alreadySubscribed == false {
				sub(w, r)
				alreadySubscribed = true
		}

		//fmt.Println(string(haveFile))
}

func sub(w http.ResponseWriter, r *http.Request){
		// fmt.Println("sub function")
	  // r.ParseForm()
	  // //t, _ := template.ParseFiles("newViewSub.html")
	  // //t.Execute(w, nil)
	  // topic = r.Form["subscribe"]
	  // fmt.Println(topic)
		//
		// fmt.Println("about to subscribe")

		output, err1 = sh.PubSubSubscribe("dubsonly")
		if err1 != nil {
			os.Stderr.WriteString(err1.Error())
		}
		fmt.Println("subscribed")

}

func editHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("edit handler")

		title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "newEdit", p)

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
		http.Redirect(w, r, "/save/"+title, http.StatusFound)

		if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
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

func messengerHandler(w http.ResponseWriter, r *http.Request) {

		//subscribe the first time you open the page
		if alreadySubscribed == false {
				sub(w, r)
				alreadySubscribed = true
		}

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
		fmt.Fprintf(w, "<h1>Chat room: Dubsonly</h1>" +
				"<p> %s </p>" +
				"<form action=\"/save/\" method=\"POST\">" +
				"<textarea name=\"body\"> </textarea>" +
				"<input type=\"submit\" value=\"Save\">" +
				"</form>" +
				"<form action=\"/refresh/\" method=\"POST\">" +
				"<input type=\"submit\" value=\"Refresh Page\">" +
				"</form>",
				sli)


		// if getPeers() > 0 {
		//
		// 	//sendFile()
		//
		// 	// t := time.Now()
		// 	// pos := time.Since(t)
		// 	// getPBFile()
		// 	// if pos == time.Second {
		// 	// 		break;
		// 	// }
		// 	//break;
		//
		// 	// done := countToTen()
    //   // fmt.Println("countToTen() exited")
    //   // <-done // block until countToTen()'s goroutine is done
		//
		// 	for start := time.Now(); time.Since(start) < time.Second; {
		// 	    getPBFile()
		// 	}
		//
		// 	http.Redirect(w, r, "/save/"+title, http.StatusFound)
		//
		// 	// for start := time.Now(); time.Since(start) < time.Second; {
		// 	// 	if firstPB == true{
		// 	// 		getPBFile()
		// 	// 		firstPB = false
		// 	// 		break
		// 	// 	}
		// 	// 	else{
		// 	// 		tempRD = recordData
		// 	// 		getPBFile()
		// 	// 			if(tempRD != recordData){
		// 	// 				break
		// 	// 			}
		// 	// 		}
		// 	// }
		// }
}

func countToTen() chan bool {
    done := make(chan bool)
    go func() {
        for i := 0; i < 10; i++ {
            time.Sleep(1 * time.Second)
            fmt.Println(i)
        }
        done <- true
    }()
    return done
}

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

		//it changes the recordData, but still appends the username function above

		userInput = message

    WriteStringToFile("userInput.txt", message)

		// if firstTry == true {
		// 		fmt.Println("first try is true")
		// 		fmt.Println("message is " + message)
		// 		sendFile(message)
		// 		AppendFile(message)
		// 		firstTry = false
		// } else {
		// 		// getPBFile()
		// 		getPeers()
		//
		// 		fmt.Println("first try is false")
		//
		// 		getPBFile()
		//
		// 		AppendFile(recordData)
		// }

		sendFile(message)
		AppendFile(message)

		//fmt.Println("record data is " + recordData)
		//WriteStringToFile("messaging.txt", recordData)


		entered = false

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

		title = "messaging"

		http.Redirect(w, r, "/messenger/"+title, http.StatusFound)

}

func getPeers() int {
		fmt.Println("get peers")
		output1, err1 := exec.Command("ipfs", "pubsub", "dubsonly", "peers").CombinedOutput()
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
	basically pubsubGetMessaging
**/
func getPBFile(){
// start :=  time.Now()
// 	for time.Since(start) < time.Second {
// 			//getPBFile()
// 			//break;

			// time.Sleep(100 * time.Millisecond)
			// fmt.Println(time.Since(start))

			//firstAppend = false

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
	// }

}


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

			// buf1, err5 := os.Open("messaging.txt")
			//
			// hash1, err5 := sh.Add(buf1)
			// if err5 != nil {
			// 	os.Stderr.WriteString(err5.Error())
			// }
			// fmt.Println("hash1 is " + hash1)
			//
			// theReader3, err6 := sh.Cat("/ipfs/" + string(hash1))
			// if err6 != nil {
			// 	 os.Stderr.WriteString(err6.Error())
			// }
			// buf7 := new(bytes.Buffer)
			// buf7.ReadFrom(theReader3)
			// newOutput3 := buf7.String()
			// fmt.Println("new output 3 is " + newOutput3)


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

	 }

}

func refreshHandler(w http.ResponseWriter, r *http.Request) {

		fmt.Println("refresh handler")
		getPeers()

		getPBFile()

		AppendFile(recordData + "--")

		title := "messaging"

		http.Redirect(w, r, "/messenger/"+title, http.StatusFound)
}

/**
	basically pubsubTheHash
**/
func sendFile(theMessage string){

		err2 = sh.PubSubPublish("dubsonly", theMessage)
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

func main(){
		alreadySubscribed = false
		http.HandleFunc("/sub/", subHandler)
		http.HandleFunc("/save/", saveHandler)
		http.HandleFunc("/edit/", editHandler)
		http.HandleFunc("/messenger/", messengerHandler)
		http.HandleFunc("/user/", userHandler)
		http.HandleFunc("/refresh/", refreshHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
}
