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
var startUser bool
var userNew bool
var username string
var topic []string
var topic2 []byte
var sh *shell.Shell
var pubsubHash string
var count bool
var recordData string

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
			startUser = true
			userNew = true
			count = false

			//http.Redirect(w, r, "/sub/"+title, http.StatusFound)
}

/**
	This will execute when it ends in /edit. This
	will load the page with what is after /edit.
**/
func editHandler(w http.ResponseWriter, r *http.Request) {

		title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)

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
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save() //where it saves the new text


		if userNew == true {
				fmt.Println("user new is true")
				username = string(p.Body)
				message = username + " has joined" + "--\r\n"
				//fmt.Println("count " + count)
				//count++
				userNew = false
		} else {
				fmt.Println("user new is false")
				message = username + ": " + string(p.Body) + "--\r\n"
				fmt.Println("message " + message)
		}

    WriteStringToFile("userInput.txt", message)

		sh = shell.NewShell("localhost:5001")

    buf, err4 := os.Open("userInput.txt")

    hash, err4 := sh.Add(buf)
    if err4 != nil {
       os.Stderr.WriteString(err4.Error())
    }
    //fmt.Println(mhash)

		if count == true {
				getMessagingFile()
				AppendFile(string(hash))
		}


    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

		title = "messaging"
		if startUser == true {
				http.Redirect(w, r, "/sub/"+title, http.StatusFound)
		} else {
				http.Redirect(w, r, "/messenger/"+title, http.StatusFound)
		}

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

		getMessagingFile()

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
		fmt.Fprintf(w, "<h1>Here are all your messages</h1>" +
				"<p> %s </p>" +
				"<form action=\"/save/\" method=\"POST\">" +
				"<textarea name=\"body\"> </textarea>" +
				"<input type=\"submit\" value=\"Save\">" +
				"</form>",
				sli)

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

		output1, err1 := exec.Command("ipfs", "pubsub", topic[0], "peers").CombinedOutput()
		if err1 != nil {
			os.Stderr.WriteString(err1.Error())
		}
		swarmPeersOutput := string(output1)

		peers := strings.Split(swarmPeersOutput, "\n")
		fmt.Println("peers")
		fmt.Println(peers)

		counter := len(peers)
		fmt.Println( "counter ")
		fmt.Println(counter)

		listPeers := []string{}
		peersLength := 0

		for i := range peers {

				listPeers = append(listPeers, peers[i])
				fmt.Println(listPeers[i])

				peersLength = len(listPeers) - 1
				fmt.Println("peers length ")
				fmt.Println(peersLength)
		}

		if start == true {
				fmt.Println("start is true")

				if peersLength == 0 {
						WriteStringToFile("messaging.txt", message)
				} else {
						pubsubGetMessaging()
				}
				start = false
		} else {
				fmt.Println("start is false")
		}

}

/**
	This will use IPFS through the command line to take the userInput
	and add it to the messaging.txt file to update it.
**/
func AppendFile(hash string) {

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

		catHash, err3 := exec.Command("cmd","/C", "ipfs", "cat", recordData, ">messaging.txt").Output()
					if err3 != nil {
						os.Stderr.WriteString(err3.Error())
						}
		fmt.Println("catHash " + string(catHash))

		//this will execute if there are peers, therefore there will be a

		// startPubsub, err := exec.Command("cmd","/C", "ipfs", "pubsub", "NAVY_messaging").Output()
		// 			if err != nil {
		// 				os.Stderr.WriteString(err.Error())
		// 				}
		// fmt.Println("start pubsub " + string(startPubsub))
		//
		// getHash, err2 := exec.Command("cmd","/C", "ipfs", "get", "QmbDt7wdRVtAnJ9jRXwerxddA1aTLLFAucEQVieKEk99xd").Output()
		// 			if err2 != nil {
		// 				os.Stderr.WriteString(err2.Error())
		// 				}
		// fmt.Println("getHash " +string(getHash))
		//
		// catHash, err3 := exec.Command("cmd","/C", "ipfs", "cat", "QmbDt7wdRVtAnJ9jRXwerxddA1aTLLFAucEQVieKEk99xd", ">messaging.txt").Output()
		// 			if err3 != nil {
		// 				os.Stderr.WriteString(err3.Error())
		// 				}
		// fmt.Println("catHash " + string(catHash))



}

func pubsubTheHash () {

		output,err1 := sh.PubSubSubscribe(topic[0])
		if err1 != nil {
			os.Stderr.WriteString(err1.Error())
		}
		err2 := sh.PubSubPublish(topic[0], hash)
		if err2 != nil {
			os.Stderr.WriteString(err2.Error())
		}
		output2, err3:= output.Next()
		if err3 != nil {
			os.Stderr.WriteString(err3.Error())
		}
		pubsubHash = string(output2.Data())
}

func subHandler(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
				t, _ := template.ParseFiles("sub.html")
				t.Execute(w, nil)
		} else {
				sub(w, r)
	      fmt.Fprintf(w, "subscribed!")
		}
		startUser = false
		count = true
}

func sub(w http.ResponseWriter, r *http.Request){
  r.ParseForm()
  t, _ := template.ParseFiles("sub.html")
  t.Execute(w, nil)
  topic = r.Form["subscribe"]
  fmt.Println(topic)


	subVar, subErr := sh.PubSubSubscribe(topic[0])
	if subErr != nil {
		os.Stderr.WriteString(subErr.Error())
	}

	subRecord, nextErr := subVar.Next()
	if nextErr != nil {
		os.Stderr.WriteString(nextErr.Error())
	}
	recordPeer := subRecord.From()
	recordData = string(subRecord.Data())
	recordNo := subRecord.SeqNo()
	recordTopics := subRecord.TopicIDs()
	fmt.Println("record peer" + string(recordPeer))
	fmt.Println("record data " + string(recordData))
	fmt.Println("record number " + string(recordNo))
	fmt.Println("record topics ")
	for i := range recordTopics {
			fmt.Println(recordTopics[i])
	}

//  exec.Command("cmd","/C","ipfs", "pubsub", "sub", topic[0]).Output()
  //fmt.Println(string(topic2))
}

func pubHandler(w http.ResponseWriter, r *http.Request) {

	title := "NAVY Messenger"
	p, err := loadPage(title)
	if err != nil {
			p = &Page{Title: title}
	}
	renderTemplate(w, "pub", p)

	//r.ParseForm()
  //t, _ := template.ParseFiles("pub.html")
  //t.Execute(w, nil)

  fmt.Println(topic[0])
  topic2,_ =  exec.Command("ipfs","pubsub","pub","dubsonly","hey").Output()
  fmt.Fprintf(w,string(topic2))
  exec.Command("cmd","/C","echo",string(topic2),">","pubsub.txt")
  //fmt.Println(string(output))

}


/**
	This will get all of the local repo hashes.
**/
func getHashes() {

		counter := false
		for true {
			tempList := []string{}

		//	time.Sleep(50000000)
			output, err1 := exec.Command("ipfs", "refs", "local").Output()
			if err1 != nil {
				os.Stderr.WriteString(err1.Error())
			}
			s := string(output)

			hashes := strings.Split(s, "\n")
			//where list was initiated

			if counter == true {
				d := difference(list, tempList)
				//receiver = d
				fmt.Println("Nothing Found")
				for i := range d {
					output, err1 := exec.Command("ipfs", "get", d[i]).Output()
					fmt.Println(output)
					if err1 != nil {
						os.Stderr.WriteString(err1.Error())
					}
				}
			} else {
				for i := range hashes {
					fmt.Println(hashes[i])
					list = append(list, hashes[i])
					output, err1 := exec.Command("ipfs", "get", hashes[i]).Output() //where you download all the local directory
					fmt.Println(string(output))
					if err1 != nil {
						os.Stderr.WriteString(err1.Error())
					}
				}
			}
			counter = true
			tempList = list
		}
}

/**
	Finding the newest messaging.txt file.
**/
func difference(slice1 []string, slice2 []string) []string {
		diffStr := []string{}
		m := map[string]int{}

		for _, s1Val := range slice1 {
			m[s1Val] = 1
		}
		for _, s2Val := range slice2 {
			m[s2Val] = m[s2Val] + 1
		}

		for mKey, mVal := range m {
			if mVal == 1 {
				diffStr = append(diffStr, mKey)
			}
		}

		return diffStr
}

func main() {

	  http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
		http.HandleFunc("/messenger/", messengerHandler)
		http.HandleFunc("/sub/", subHandler)
		http.HandleFunc("/pub/", pubHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
