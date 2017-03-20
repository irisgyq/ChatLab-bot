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

const FacebookEndPoint = "https://graph.facebook.com/v2.6/me/messages"
const WeatherPoint = "api.openweathermap.org/data/2.5/weather"

var (
	accessToken = os.Getenv("FB_PAGE_ACCESS_TOKEN")
	validationToken = os.Getenv("VALIDATION_TOKEN")
	port = os.Getenv("PORT")
)

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

type Postback struct {
	Payload string `json:"payload"`
}

type Sender struct {
	ID string `json:"id"`
	Location string `json:"location"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Text       string      `json:"text,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

type AttachmentPayload struct {
	Template_type string `json:"template_type"`
	//Text string `json:"text"`
	Buttons *[]Button `json:"buttons"`
	Elements      *[]Elements `json:"elements"`
}

type Elements struct {
	Title string `json:"title"`
	Subtitle string `json:"subtitle, omitempty"`
	Image_Url   string `json:"image_url,omitempty"`
	Buttons     *[]Button `json:"buttons"`
}

type Button struct {
	Type string `json:"type"`
	Title string `json:"title"`
	Payload string `json:"payload"`
	Url string `json:"url"`
	Webview_height_ratio string `json:"webview_height_ratio"`
	Messenger_extensions bool `json:"messenger_extensions"`
}

type Attachment struct {
	Type    string             `json:"type,omitempty"`
	Payload *AttachmentPayload `json:"payload,omitempty"`
}

type MessageToSend struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
}


type Info struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Location string `json:"locale"`
}

func main() {

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
				mes := strings.ToUpper(message.Message.Text)
					info, errr := getSenderInfo(message.Sender.ID)
					msg := "Tooo is something wrong"
					if errr == nil {
						if strings.Contains(mes,"HI") ||  strings.Contains(mes,"HELLO"){
							msg = "Hello " + info.FirstName + " " + info.LastName + ", this is a lovely chat bot. How are you today? Good or bad?"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "GOOD" || mes=="GREAT" {
							msg = "That's great! Do you want to learn something? You can input 'tools' to learn programming languages."
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "BAD" || mes == "FINE"{
							msg = "I'm sorry. Maybe you can play some games. You can input 'tools' to play blackjack!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "THANK YOU" || mes == "THANKS" || strings.Contains(mes, "APPRECIATE") {
							msg = info.FirstName + " " + info.LastName + "You are welcome!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "BYE" || mes == "SEE YOU" || mes == "GOODBYE" {
							msg = "Bye " + info.FirstName + " " + info.LastName + "Have a nice day! See you next time."
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "TOOLS" {
							go sendGenericMessage(message.Sender.ID, message.Sender.Location)
						} else if strings.Contains(mes, "PROGRAMMING LANGUAGES") || mes == "STUDY" {
							msg = "What kind of programming languages do you want to learn"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "GO" || mes == "GOLANG" {
							msg = "GO is really good!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "PYTHON" {
							msg = "Python is really good!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "RUBY" || mes == "RUBY ON RAILS" {
							msg = "Ruby is really good!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "Java" {
							msg = "Java is really good!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "JavaScript" || mes == "JS" {
							msg = "JavaScript is really good!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "SCALA"{
							msg = "SCALA is really good!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "PROLOG"  {
							msg = "Prolog is really good!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "LISP" {
							msg = "LISP is really good!"
							go sendTextMessage(message.Sender.ID, msg)
						}  else if strings.Contains(mes, "WEATHER") {
							msg = "Please input a city and plus a &, like WashingtonDC&."
							go sendTextMessage(message.Sender.ID, msg)
						} else if strings.HasSuffix(mes, "&") {
							go sendUrlMessage(message.Sender.ID, "http://api.openweathermap.org/data/2.5/weather?q="+mes+"mode=html&APPID=404cd230fcf7a79e7dcb4f9abbaca518")
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

func sendUrlMessage(sender string, url string) {
	sendMessage(MessageToSend{
		Recipient: Recipient{
			ID: sender,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: &AttachmentPayload{
					Template_type: "button",
					Buttons: &[]Button{{
						Type:"web_url",
						Url: url,
						Webview_height_ratio: "fULL",
						Messenger_extensions: true,
					}},
				},
			},
		},
	})
}

func sendGenericMessage(sender string, city string) {
	/*resp, err := http.Get("https://graph.facebook.com/v2.8/city/fields=city")
	if err!=nil {
		//
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)*/

	sendMessage(MessageToSend{
		Recipient: Recipient{
			ID: sender,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: &AttachmentPayload{
					Template_type: "generic",
					Elements: &[]Elements{{
						Title: "Tools",
						Subtitle: "You can choose one of them!",
						Image_Url: "http://chuantu.biz/t5/50/1490034411x2890154370.png",
						Buttons: &[]Button{{
							Type: "postback",
							Title: "Study",
							Payload: "PROGRAMMING LANGUAGES",
						}, {
							Type: "web_url",
							Title: "Entertainment",
							Url:"http://www.websudoku.com/",
							Webview_height_ratio: "fULL",
							Messenger_extensions: true,
						},{
							Type: "web_url",
							Title: "Weather",
							Url: "http://api.openweathermap.org/data/2.5/weather?q=WashingtonDC&mode=html&APPID=404cd230fcf7a79e7dcb4f9abbaca518",
							Webview_height_ratio: "fULL",
							Messenger_extensions: true,
						}},
					}},
				},
			},
		},
	},)
}

func sendMessage(m interface{}) {

	msg, err := json.Marshal(m)

	if err != nil {
		fmt.Println("Thrrr is something wrong!")
		fmt.Println(err)
	}

	fmt.Println("Send message!")
	fmt.Println(string(msg))

	resp, err := doRequest("POST", FacebookEndPoint, bytes.NewReader(msg))

	if err != nil {
		fmt.Println("The is something wrong!")
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
		fmt.Println("Tttt is something wrong!")
		fmt.Println(string(read))
		return nil, errors.New(string(read))
	}
	Info := new(Info)
	return Info, json.Unmarshal(read, Info)
}


