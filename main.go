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
	//Buttons *[]Button `json:"buttons"`
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
						if mes == "HI" || mes == "HELLO"{
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
							go sendGenericMessage(message.Sender.ID)
						} else if mes== "PROGRAMMING LANGUAGES" || mes == "STUDY" {
							msg = "What kind of programming languages do you want to learn"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "GO" || mes == "GOLANG" {
							msg = "Go is a free and open source programming language created at Google in 2007 by Robert Griesemer, Rob Pike, and Ken Thompson.It is a compiled, statically typed language in the tradition of Algol and C, with garbage collection, limited structural typing，memory safety features and CSP-style concurrent programming features added."
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "PYTHON" {
							msg = "Python is a widely used high-level programming language for general-purpose programming, created by Guido van Rossum and first released in 1991. An interpreted language. Python has a design philosophy which emphasizes code readabilityPython features a dynamic type system and automatic memory management and supports multiple programming paradigms, including object-oriented, imperative, functional programming, and procedural styles. It has a large and comprehensive standard library."
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "RUBY" || mes == "RUBY ON RAILS" {
							msg = "Ruby is a dynamic, reflective, object-oriented, general-purpose programming language. It was designed and developed in the mid-1990s by Yukihiro Matsumoto in Japan. It supports multiple programming paradigms, including functional, object-oriented, and imperative. It also has a dynamic type system and automatic memory management"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "JAVA" {
							msg = "Java is a general-purpose computer programming language that is concurrent, class-based, object-oriented,and specifically designed to have as few implementation dependencies as possible. Java is a general-purpose computer programming language that is concurrent, class-based, object-oriented,and specifically designed to have as few implementation dependencies as possible.  "
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "JAVASCRIPT" || mes == "JS" {
							msg = "JavaScript is a high-level, dynamic, untyped, and interpreted programming language.It has been standardized in the ECMAScript language specification.Alongside HTML and CSS, JavaScript is one of the three core technologies of World Wide Web content production; the majority of websites employ it, and all modern Web browsers support it without the need for plug-ins."
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "SCALA"{
							msg = "Scala is a general-purpose programming language providing support for functional programming and a strong static type system. Designed to be concise,many of Scala's design decisions were designed to build from criticisms of Java.!"
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "PROLOG"  {
							msg = "Prolog is a general-purpose logic programming language associated with artificial intelligence and computational linguistics.Prolog has its roots in first-order logic, a formal logic, and unlike many other programming languages, Prolog is declarative: the program logic is expressed in terms of relations, represented as facts and rules. A computation is initiated by running a query over these relations."
							go sendTextMessage(message.Sender.ID, msg)
						} else if mes == "LISP" {
							msg = "Lisp is a family of computer programming languages with a long history and a distinctive, fully parenthesized prefix notation.[3] Originally specified in 1958, Lisp is the second-oldest high-level programming language in widespread use today. Only Fortran is older, by one year.Lisp has changed since its early days, and many dialects have existed over its history. Today, the best known general-purpose Lisp dialects are Common Lisp and Scheme.!"
							go sendTextMessage(message.Sender.ID, msg)
						}  else if mes == "WEATHER" {
							msg = "Please input a city and plus a &, like WashingtonDC&."
							go sendTextMessage(message.Sender.ID, msg)
						} else if strings.HasSuffix(mes, "&") {
							go sendUrlMessage(message.Sender.ID, "http://api.openweathermap.org/data/2.5/weather?q="+message.Message.Text+"mode=html&APPID=404cd230fcf7a79e7dcb4f9abbaca518")
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

func sendGenericMessage(sender string) {
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
						},{
							Type: "web_url",
							Url:"http://www.websudoku.com/",
							Title: "Entertainment",
						},{
							Type: "web_url",
							Url: "http://api.openweathermap.org/data/2.5/weather?q=WashingtonDC&mode=html&APPID=404cd230fcf7a79e7dcb4f9abbaca518",
							Title: "Weather",
						}},
					}},
				},
			},
		},
	},)
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
					Template_type: "generic",
					Elements: &[]Elements{{
						Title:"Weather",
						Buttons: &[]Button{{
							Type:"web_url",
							Title: "weather",
							Url: url,
						}},
					}},

				},
			},
		},
	})
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


