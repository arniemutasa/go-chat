package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

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

	err = ws.WriteJSON(response)

	if err != nil {
		log.Println(err)
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
