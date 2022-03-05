package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

// Channel to handle websocket payloads
var wsChannel = make(chan WebsocketPayload)

var clients = make(map[WebsocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

type WebsocketConnection struct {
	*websocket.Conn
}

// WebsocketJsonResponse - response sent back from webSocketEndpoint
type WebsocketJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WebsocketPayload struct {
	Action     string              `json:"action"`
	Message    string              `json:"message"`
	User       string              `json:"user"`
	Connection WebsocketConnection `json:"-"`
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)

	}
}

// Upgrades http connection to websocket
func WebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Connection Upgraded: Client connected to endpoint")

	var response WebsocketJsonResponse

	response.Message = `<em><small>Connected To Server</small></em>`

	conn := WebsocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)

	if err != nil {
		log.Println(err)
	}

	go ListenForWS(&conn)

}

// Runs continuously and listens for payload, if payload is available it sends it to the channel
func ListenForWS(conn *WebsocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WebsocketPayload

	for {
		err := conn.ReadJSON(&payload)

		if err != nil {
			// do nothing
		} else {
			payload.Connection = *conn
			wsChannel <- payload
		}
	}
}

// Listens to Websocket channel for any websockets and broadcasts to all connected users
func ListenToWebsocketChannel() {
	var response WebsocketJsonResponse

	for {
		event := <-wsChannel

		response.Action = "We Here"
		response.Message = fmt.Sprintf("some Message, and Action is %s\n", event.Action)
		BroadcastToAllUsers(response)
	}

}

func BroadcastToAllUsers(r WebsocketJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(r)
		if err != nil {
			log.Println("Websocket Error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func renderPage(w http.ResponseWriter, tmp string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmp)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
