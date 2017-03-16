package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var accessToken = os.Getenv("FB_PAGE_ACCESS_TOKEN")
var validationToken = os.Getenv("VALIDATION_TOKEN")
var port = os.Getenv("PORT")
const FacebookEndPoint = "https://graph.facebook.com/v2.6/me/messages"

type Payload struct {
	Object  string   `json:"object"`
	Entries []*Entry `json:"entry"`
}

type Entry struct {
	ID        string  `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Postback  Postback  `json:"postback"`
	Message   Message   `json:"message"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Text       string      `json:"text,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

type Postback struct {
	Payload string `json:"payload"`
}

type AttachmentPayload struct {
	Url string `json:"url,omitempty"`
}

type Attachment struct {
	Type    string             `json:"type,omitempty"`
	Payload *AttachmentPayload `json:"payload,omitempty"`
}

type MessageToSend struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
}

type ActionToSend struct {
	Recipient    Recipient `json:"recipient"`
	SenderAction string    `json:"sender_action"`
}

type Info struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}


func main() {

	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", SiteHandle)
	http.HandleFunc("/webhook", Handle)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func SiteHandle(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Welcome to my world, I'm a golang chat-bot.")
}

func Handle(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		handleGet(rw, req)
	} else if req.Method == "POST" {
		handlePost(rw, req)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleGet(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("hub.verify_token") == validationToken {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, req.URL.Query().Get("hub.challenge"))

	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(rw, "Error, wrong validation token")
	}
}

func handlePost(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	payload := &Payload{}
	if err = json.Unmarshal(body, payload); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, entry := range payload.Entries {
		for _, message := range entry.Messaging {
			if message.Message != nil {
				/*if message.Message.Text == "typing_on" {
					go sendActionMessage(message.Sender.ID, "typing_on")
				} else if message.Message.Text == "typing_off" {
					go sendActionMessage(message.Sender.ID, "typing_off")
				} else if message.Message.Text == "image" {
					go sendAttachmentMessage(message.Sender.ID, "image", "http://cdn.morguefile.com/imageData/public/files/b/Baydog64/08/p/8b2facd9fffc84af6cbdc5e7e24ede70.jpg")
				}*/
				mes := strings.ToUpper(message.Message)
				if{
					info, errr := getSenderInfo(message.Sender.ID)
					msg := "There is something wrong"
					if errr == nil {
						if mes == "HI" {
							msg = "Hello " + info.FirstName + " " + info.LastName + ", this is a lovely chat bot. How are you today? Good or bad?"
						} else if mes == "GOOD" {
							msg = "That's great! What can I do for you? Study or entertainment?"
						} else if mes == "BAD" {
							msg = "If you talk to me, you will be happy! What can I do for you? Study or entertainment?"
						} else if mes == "STUDY" {
							msg = "What kind of language do you want to learn?"
						} else if mes == "GO" {
							msg = "Go is a free and open source programming language created at Google in 2007 by Robert Griesemer, Rob Pike, and Ken Thompson. It is a compiled, statically typed language in the tradition of Algol and C, with garbage collection, limited structural typing,[3] memory safety features and CSP-style concurrent programming features added."
						} else if mes == "JAVA" {
							msg = "Java is a general-purpose computer programming language that is concurrent, class-based, object-oriented, and specifically designed to have as few implementation dependencies as possible."
						} else if mes == "SCALA" {
							msg = "Scala is a general-purpose programming language providing support for functional programming and a strong static type system. Designed to be concise, many of Scala's design decisions were designed to build from criticisms of Java."
						} else if mes == "PROLOG" {
							msg = "Prolog is a general-purpose logic programming language associated with artificial intelligence and computational linguistics."
						} else if mes == "C" {
							msg = "C is a general-purpose, imperative computer programming language, supporting structured programming, lexical variable scope and recursion, while a static type system prevents many unintended operations. "
						} else if mes == "JAVASCRIPT" {
							msg = "JavaScript is a high-level, dynamic, untyped, and interpreted programming language. It has been standardized in the ECMAScript language specification."
						} else if mes == "RUBY" {
							msg = "Ruby is a dynamic, reflective, object-oriented, general-purpose programming language. It was designed and developed in the mid-1990s by Yukihiro Matz Matsumoto in Japan."
						} else if mes == "PYTHON" {
							msg = "Python is a widely used high-level programming language for general-purpose programming, created by Guido van Rossum and first released in 1991. "
						} else if mes == "LISP" {
							msg = "Lisp is a family of computer programming languages with a long history and a distinctive, fully parenthesized prefix notation."
						} else if mes == "ENTERTAINMENT" {
							msg = "Let's play blackjack! You are the player, and I am the dealer."
						} else if mes == "BYE" || mes == "SEE YOU" || mes == "GOODBYE" {
							msg = "Bye "+ info.FirstName + " " + info.LastName +"Have a nice day! See you next time."
						} else {
							msg = "Hello " + info.FirstName + " " + info.LastName + ", this is a lovely chat bot. I like repeat your words, so " + message.Message.Text
						}

					}
					go sendTextMessage(message.Sender.ID, msg)
				}

			}

		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))
}

/*func sendAttachmentMessage(sender string, attachmentType string, url string) {
	sendMessage(MessageToSend{
		Recipient: Recipient{
			ID: sender,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: attachmentType,
				Payload: &AttachmentPayload{
					Url: url,
				},
			},
		},
	})
}*/

func sendTextMessage(sender string, text string) {
	sendMessage(MessageToSend{
		Recipient: Recipient{
			ID: sender,
		},
		Message: Message{
			Text: text,
		},
	})
}

/*func sendActionMessage(sender string, action string) {
	sendMessage(ActionToSend{
		Recipient: Recipient{
			ID: sender,
		},
		SenderAction: action,
	})
}*/

func sendMessage(m interface{}) {

	msg, err := json.Marshal(m)

	if err != nil {
		fmt.Println("There is something wrong!")
		fmt.Println(err)
	}

	fmt.Println("Send message!")
	fmt.Println(string(msg))

	resp, err := doRequest("POST", FacebookEndPoint, bytes.NewReader(msg))

	if err != nil {
		fmt.Println("There is something wrong!")
		fmt.Println(err)
	}
	defer resp.Body.Close()

	read, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(read))
	}

}

func doRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	query := req.URL.Query()
	query.Set("access_token", accessToken)
	req.URL.RawQuery = query.Encode()
	return http.DefaultClient.Do(req)
}

func getSenderInfo(userID string) (*Info, error) {
	graphAPI := "https://graph.facebook.com"
	resp, err := doRequest("GET", fmt.Sprintf(graphAPI+"/v2.6/%s?fields=first_name,last_name", userID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	read, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("There is something wrong!")
		fmt.Println(string(read))
		return nil, errors.New(string(read))
	}
	Info := new(Info)
	return Info, json.Unmarshal(read, Info)
}

