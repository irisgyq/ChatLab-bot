package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/webhook", Handle)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func Handle(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		query := req.URL.Query()
		if query.Get("hub.verify_token") != os.Getenv("VALIDATION_TOKEN") {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(query.Get("hub.challenge")))
	} else if req.Method == "POST" {
		handlePOST(rw, req)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlePOST(rw http.ResponseWriter, req *http.Request) {
	read, err := ioutil.ReadAll(req.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	if len(req.Header.Get("x-hub-signature")) < 6 || !checkSignature(read, req.Header.Get("x-hub-signature")[5:]) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	payload := &Payload{}
	err = json.Unmarshal(read, payload)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Mensagem recebida!")
	pj, _ := json.Marshal(payload)
	fmt.Println(string(pj))

	for _, entry := range payload.Entries {
		for _, message := range entry.Messaging {
			if message.Message != nil {
				if message.Message.Text == "typing_on" {
					go sendActionMessage(message.Sender.ID, "typing_on")
				} else if message.Message.Text == "typing_off" {
					go sendActionMessage(message.Sender.ID, "typing_off")
				} else if message.Message.Text == "image" {
					go sendAttachmentMessage(message.Sender.ID, "image", "http://cdn.morguefile.com/imageData/public/files/b/Baydog64/08/p/8b2facd9fffc84af6cbdc5e7e24ede70.jpg")
				} else {
					info, errr := getSenderInfo(message.Sender.ID)
					msg := "Desculpe. Que é você?"
					if errr == nil {
						msg = "Olá " + info.FirstName + " " + info.LastName + ". Sou um papagaio. " + message.Message.Text
					}
					go sendTextMessage(message.Sender.ID, msg)
				}

			} else if message.Postback != nil {
				msg := "Não entendi!"
				if message.Postback.Payload == "USER_DEFINED_PAYLOAD" {
					msg = "Obrigado. Começou."
				} else if message.Postback.Payload == "DEVELOPER_DEFINED_PAYLOAD_FOR_HELP" {
					msg = "Mande a mensagem 'typing_on' e vc verá que estou digitando, mande 'typing_off' que eu paro, mande 'image' que eu mando uma imagem. Digite outra coisa e verá."
				} else if message.Postback.Payload == "DEVELOPER_DEFINED_PAYLOAD_FOR_START_ORDER" {
					msg = "Pode pedir."
				}
				go sendTextMessage(message.Sender.ID, msg)
			}
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))
}

// https://developers.facebook.com/docs/messenger-platform/webhook-reference
func checkSignature(bytes []byte, expectedSignature string) bool {
	appSecret := os.Getenv("APP_SECRET")

	mac := hmac.New(sha1.New, []byte(appSecret))
	mac.Write(bytes)
	if fmt.Sprintf("%x", mac.Sum(nil)) != expectedSignature {
		return false
	}
	return true
}

func sendAttachmentMessage(sender string, attachmentType string, url string) {
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

func sendActionMessage(sender string, action string) {
	sendMessage(ActionToSend{
		Recipient: Recipient{
			ID: sender,
		},
		SenderAction: action,
	})
}

func sendMessage(m interface{}) {

	msg, err := json.Marshal(m)

	if err != nil {
		fmt.Println("Erro criando mensagem de retorno!")
		fmt.Println(err)
	}

	fmt.Println("Enviando mensagem!")
	fmt.Println(string(msg))

	graphAPI := "https://graph.facebook.com"
	resp, err := doRequest("POST", graphAPI+"/v2.6/me/messages", bytes.NewReader(msg))

	if err != nil {
		fmt.Println("Erro na mensagem enviada para o Messenger!")
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
	query.Set("access_token", os.Getenv("FB_PAGE_ACCESS_TOKEN"))
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
		fmt.Println("Erro pegando informações no usuário!")
		fmt.Println(string(read))
		return nil, errors.New(string(read))
	}
	Info := new(Info)
	return Info, json.Unmarshal(read, Info)
}

type Payload struct {
	Object  string   `json:"object"`
	Entries []*Entry `json:"entry"`
}

// https://developers.facebook.com/docs/messenger-platform/webhook-reference
// https://developers.facebook.com/docs/messenger-platform/webhook-reference/message-received
type Entry struct {
	ID        json.Number `json:"id"`
	Time      int64       `json:"time"`
	Messaging []struct {
		Sender *struct {
			ID string `json:"id"`
		} `json:"sender"`
		Recipient *struct {
			ID string `json:"id"`
		} `json:"recipient"`
		Message *struct {
			Text string `json:"text"`
		} `json:"message,omitempty"`
		Postback *struct {
			Payload string `json:"payload"`
		} `json:"postback,omitempty"`
		Timestamp int64 `json:"timestamp"`
	} `json:"messaging"`
}

type Recipient struct {
	ID string `json:"id"`
}

type AttachmentPayload struct {
	Url string `json:"url,omitempty"`
}

type Attachment struct {
	Type    string             `json:"type,omitempty"`
	Payload *AttachmentPayload `json:"payload,omitempty"`
}

type Message struct {
	Text       string      `json:"text,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

// https://developers.facebook.com/docs/messenger-platform/send-api-reference/text-message
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
