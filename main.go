package main

import (
	"bytes"
	"encoding/json"
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
				go sendMessage(message.Recipient.ID, message.Message.Text)
			}
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))
}

func sendMessage(sender string, text string) {

	msg, err := json.Marshal(MessageToSend{
		Recipient: Recipient{
			ID: sender,
		},
		Message: TextMessage{
			Text: text,
		},
	})

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
		} `json:"message"`
		Timestamp int64 `json:"timestamp"`
	} `json:"messaging"`
}

type Recipient struct {
	ID string `json:"id"`
}

type TextMessage struct {
	Text string `json:"text"`
}

// https://developers.facebook.com/docs/messenger-platform/send-api-reference/text-message
type MessageToSend struct {
	Recipient Recipient   `json:"recipient"`
	Message   TextMessage `json:"message"`
}
