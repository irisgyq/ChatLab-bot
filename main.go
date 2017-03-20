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
	"math/rand"
	"time"
	"bufio"
	"strconv"
)

const FacebookEndPoint = "https://graph.facebook.com/v2.6/me/messages"
var (
	accessToken = os.Getenv("FB_PAGE_ACCESS_TOKEN")
	validationToken = os.Getenv("VALIDATION_TOKEN")
	port = os.Getenv("PORT")
	random = rand.New(rand.NewSource(time.Now().Unix()))
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
	Template_type string `json:"type"`
	Text string `json:"text"`
	Buttons *[]Button `json:"buttons"`
	//Elements      *[]Elements `json:"elements"`
}

/*type Elements struct {
	Title string `json:"title"`
	Subtitle string `json:"subtitle, omitempty"`
	Image_Url   string `json:"imageurl,omitempty"`
	Buttons     *[]Button `json:"buttons"`
}*/

type Button struct {
	Type string `json:"type"`
	Title string `json:"title"`
	Payload string `json:"payload"`
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
							msg = "OH"
							go sendGenericMessage(message.Sender.ID)
						} else if strings.Contains(mes, "PROGRAMMING LANGUAGES") {
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
						} else if strings.Contains(mes, "BLACKJACK") {
							go PlayBlackjack(message.Sender.ID)
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
					Template_type: "button",
					Text: "These are some tools.",
					Buttons: &[]Button{{
							Type: "postback",
							Title: "Study",
							Payload: "What kind of programming languages do you want to learn",
						},{
							Type: "postback",
							Title: "entertainment",
							Payload: "Let's play blackjack! If you are ready, please input 'blackjack'",
						}},

					}},

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

func PlayBlackjack(sender string) {

		initcards := []int{1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 8, 9, 9, 9, 9,
			10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}

		playerCard := make([]int, 0)
		dealerCard := make([]int, 0)

		PisBust := false
		DisBust := false

		go sendTextMessage(sender,"Game begins...")

		cards := shuffle(initcards)

		playerCard = append(playerCard, pop(&cards))
		go sendTextMessage(sender,"The player's first card is: "+strconv.Itoa(playerCard[0]))
		dealerCard = append(dealerCard, pop(&cards))
		go sendTextMessage(sender, "The dealer's first card is: "+strconv.Itoa(dealerCard[0]))
		playerCard = append(playerCard, pop(&cards))
		go sendTextMessage(sender, "The player's second card is: "+strconv.Itoa(playerCard[1]))

		if blackjack(playerCard) {
			go sendTextMessage(sender, "The player has blackjack! Game is over, the player is the winner.")
		} else {
			dealerCard = append(dealerCard, pop(&cards))
			if blackjack(dealerCard) {
				go sendTextMessage(sender, "The dealer's second card is: "+strconv.Itoa(dealerCard[1]))
				go sendTextMessage(sender, "The dealer has blackjack! Game is over, the dealer is the winner.")
			} else {
				playerSum := playerCard[0] + playerCard[1]
				dealerSum := dealerCard[0] + dealerCard[1]

				go sendTextMessage(sender, "The sum of player's cards is:" +strconv.Itoa(playerSum))

				isPValid := true
				for (isPValid) {
					go sendTextMessage(sender,"Does the player want one more card? yes or hit")
					inputReader := bufio.NewReader(os.Stdin)
					input, err := inputReader.ReadString('\n')

					if err != nil {
						go sendTextMessage(sender, "Your input is wrong.")
						return
					}

					switch input {
					case "yes\n":{
						playerCard = append(playerCard, pop(&cards))
						go sendTextMessage(sender,"This card is:"+strconv.Itoa(playerCard[len(playerCard) - 1]))
						playerSum += playerCard[len(playerCard) - 1]
						go sendTextMessage(sender, "The sum of player's card is:"+strconv.Itoa(playerSum))

						if blackjack(playerCard) {
							go sendTextMessage(sender,"The player has 21 points!")
							isPValid = false
							break
						} else if playerSum > 21 {
							go sendTextMessage(sender, "Player's cards are busting.")
							PisBust = true
							isPValid = false
							break
						}
						break
					}
					case "hit\n" :{
						isPValid = false
						break
					}

					}
				}

				go sendTextMessage(sender,"The dealer's second card is: "+strconv.Itoa(dealerCard[1]))
				go sendTextMessage(sender, "The sum of dealer's cards is:"+strconv.Itoa(dealerSum))

				isDValid := true
				for (isDValid) {
					for dealerSum < 17 {
						go sendTextMessage(sender,"Because the sum of dealer's cards is less than 17, he must add one more card.")
						dealerCard = append(dealerCard, pop(&cards))
						go sendTextMessage(sender,"The new card is:"+strconv.Itoa(dealerCard[len(dealerCard) - 1]))
						dealerSum += dealerCard[len(dealerCard) - 1]
						go sendTextMessage(sender,"The sum of dealer's cards is:"+strconv.Itoa(dealerSum))

						if dealerSum == 21 {
							go sendTextMessage(sender,"The dealer has 21 points!")
							break
						}
						if dealerSum > 21 {
							go sendTextMessage(sender,"dealer's cards are busting.")
							DisBust = true
							break
						}
					}

					if dealerSum < 21 {
						go sendTextMessage(sender,"Dose the dealer want one more card?")
						inputReader := bufio.NewReader(os.Stdin)
						input, err := inputReader.ReadString('\n')

						if err != nil {
							go sendTextMessage(sender,"Your input is wrong.")
							return
						}

						switch input {
						case "yes\n":{
							dealerCard = append(dealerCard, pop(&cards))
							go sendTextMessage(sender,"This card is:")
							go sendTextMessage(sender,strconv.Itoa(dealerCard[len(dealerCard) - 1]))
							dealerSum += dealerCard[len(dealerCard) - 1]
							go sendTextMessage(sender,"The sum of dealer's card is:"+strconv.Itoa(dealerSum))

							if blackjack(dealerCard) {
								go sendTextMessage(sender,"The dealer has 21 points!")
								isDValid = false
								break
							} else if dealerSum > 21 {
								go sendTextMessage(sender,"dealer's cards are busting.")
								DisBust = true
								isDValid = false
								break
							}
						}
						case "hit\n" :{
							isDValid = false
							break
						}

						}

					} else {
						isDValid = false
					}
				}

				if (DisBust && PisBust) || (!DisBust && !PisBust && (dealerSum == playerSum)) {
					go sendTextMessage(sender,"Game is over, it's a push")
				} else if (DisBust && !PisBust) || (!DisBust && !PisBust && (dealerSum < playerSum)) {
					go sendTextMessage(sender,"Game is over, the player wins")
				} else if (!DisBust && PisBust) || (!DisBust && !PisBust && (dealerSum > playerSum)) {
					go sendTextMessage(sender,"Game is over, the dealer wins")
				}

			}
		}
}

func shuffle (cards []int) []int {
	temp := [52]int{}
	l := len(cards)
	for i := l-1; i>0; i-- {
		r := random.Intn(i+1)
		cards[r], cards[i] = cards[i], cards[r]
	}
	temp[cards[0]] += 1
	return cards
}

//deal cards randomly
func pop (cards *[]int) int  {
	pos := rand.Intn(len(*cards)-1)
	card := (*cards)[pos]
	*cards = append((*cards)[1:pos],(*cards)[pos:]...)
	return card
}

//judge if it is blackjack
func blackjack (a []int) bool {
	sum := 0
	hasOne := false
	for i :=0; i<len(a);i++ {
		sum += a[i]
		if(a[i]==1) {
			hasOne = true
		}
	}

	if sum == 21{
		return true;
	} else if hasOne && sum+10==21 {
		return true;
	}
	return false;
}

