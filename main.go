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
	Message   Message   `json:"message,omitempty"`
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
	Title  string `json:"title, omitempty"`
	Subtitle string `json:"subtitle, omitempty"`
	Image_Url    string `json:"imageurl,omitempty"`
	Buttons      []Button `json:"buttons"`
}

type Button struct {
	Title string `json:"title"`
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
			if message.Message.Text != "" {
				/*if message.Message.Text == "typing_on" {
					go sendActionMessage(message.Sender.ID, "typing_on")
				} else if message.Message.Text == "typing_off" {
					go sendActionMessage(message.Sender.ID, "typing_off")
				} else if message.Message.Text == "image" {
					go sendAttachmentMessage(message.Sender.ID, "image", "http://cdn.morguefile.com/imageData/public/files/b/Baydog64/08/p/8b2facd9fffc84af6cbdc5e7e24ede70.jpg")
				}*/
				mes := strings.ToUpper(message.Message.Text)
					info, errr := getSenderInfo(message.Sender.ID)
					msg := "There is something wrong"
					if errr == nil {
						if mes == "HI" {
							msg = "Hello " + info.FirstName + " " + info.LastName + ", this is a lovely chat bot. How are you today? Good or bad?"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "GOOD" {
							msg = "That's great! Do you want to learn something? You can input 'tools' to learn programming languages."
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "BAD" {
							msg = "I'm sorry. Maybe you can play some games. You can input 'tools' to play blackjack!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "THANK YOU" || mes == "THANKS" || strings.Contains(mes, "APPRECIATE") {
							msg = info.FirstName + " " + info.LastName + "You are welcome!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "BYE" || mes == "SEE YOU" || mes == "GOODBYE" {
							msg = "Bye " + info.FirstName + " " + info.LastName + "Have a nice day! See you next time."
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "TOOLS" {
							go sendGenericMessage(message.Sender.ID)
						} else {
							msg = "Hello " + info.FirstName + " " + info.LastName + ", this is a lovely chat bot. I like repeat your words, so " + message.Message.Text
							go sendTextMessage(message.Sender.ID, msg)
						}

					}

				}

			}

		}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))
}

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

func sendGenericMessage(sender string) {
	sendMessage(MessageToSend{
		Recipient: Recipient{
			ID: sender,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: &AttachmentPayload{
					Title: "Chat Lab tools",
					Image_Url: "./tools.png",
					Buttons: {
						Title: "Study",
						Title: "Entertainment",
						Title: "Calculator",
					},
				},
			},
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

