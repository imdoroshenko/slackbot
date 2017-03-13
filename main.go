package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	"log"
)

var dialHeader http.Header = map[string][]string{
	"Origin": []string{"http://localhost/",},
}

type SlackBot struct {
	ApiToken string
	EventHandler HandleFunc
	conn *websocket.Conn
	details *botDetails
}

type HandleFunc func(*Event)

type Event struct {
	Id int `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	User string `json:"user,omitempty"`
	Text string `json:"text,omitempty"`
	Channel string `json:"channel,omitempty"`
	Ts string `json:"ts,omitempty"`
}

type botDetails struct {
	URL string
}

func NewSlackBot() *SlackBot {
	return &SlackBot{}
}

// Start will start slack bot routines.
func (sb *SlackBot) Run() {
	log.Println("slack Run")
	sb.details = rtmStart(sb.ApiToken)
	sb.conn, _ = establishWSConnection(sb.details.URL)
	defer sb.conn.Close()
	handleConnection(sb.conn, sb.EventHandler)
	//handleEvents(sb.Messages)
}

func (sb *SlackBot) Send(event *Event) {
	err := sb.conn.WriteJSON(event)
	if err != nil {
		log.Println("Slackbot#Send:", err)
	}
}

func handleConnection(conn *websocket.Conn, handler HandleFunc) {
	for {
		data := &Event{}
		err := conn.ReadJSON(data)
		if err != nil {
			log.Println("Error:", err)
		}
		if (handler != nil) {
			handler(data)
		}
		//log.Printf("%+v\n", data)
	}
}

// establishWSConnection will establish WS connection to slack API.
func establishWSConnection(url string) (*websocket.Conn, error) {
	dialer := &websocket.Dialer{}
	conn, _, err := dialer.Dial(url, dialHeader)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	log.Println("WS connection established")
	return conn, nil
}

// rtmStart will request new rtm session from slack API.
func rtmStart(token string) *botDetails {
	resp, httpErr := http.Get("https://slack.com/api/rtm.start?token=" + token)
	defer resp.Body.Close()

	if httpErr != nil {
		log.Println("http get error")
	}

	body, readAllErr := ioutil.ReadAll(resp.Body)
	if readAllErr != nil {
		log.Println("ioutil error")
	}

	var data = new(botDetails)
	unmarshalErr := json.Unmarshal(body, data)

	if unmarshalErr != nil {
		log.Println("Unmarshal error")
	}
	return data
}
