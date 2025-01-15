package LIB

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var upgrader websocket.Upgrader
var clients = make(map[*websocket.Conn]bool)

var logger = zap.NewExample()

type WS_NEW_USER struct {
	Action string
	DB_INSERT_RESPONSE
}

type WS_DELETE_USER struct {
	Action string
	DB_DELETE_RESPONSE
}

type WS_UPDATE_USER struct {
	Action string
	DB_UPDATE_RESPONSE
}

func Upgrader() {
	upgrader = websocket.Upgrader{}
}

var c *websocket.Conn

func Echo(w http.ResponseWriter, r *http.Request) {
	logger.Info("ws call")
	var err error
	c, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("ws call error", zap.Error(err))
		return
	}
	defer func() {
		delete(clients, c)
		c.Close()
	}()

	clients[c] = true

	for {
		mt, message, err := c.ReadMessage()

		var data WS_NEW_USER
		json.Unmarshal(message, &data)

		if err != nil {
			logger.Error("ws call error", zap.Error(err))
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			logger.Error("ws call error", zap.Error(err))
			break
		}
	}
}

func NewUser(data DB_INSERT_RESPONSE) {
	logger.Info("ws new user")

	if c == nil {
		return
	}

	for client := range clients {
		err := client.WriteJSON(WS_NEW_USER{
			Action: "new_user",
			DB_INSERT_RESPONSE: DB_INSERT_RESPONSE{
				Id:   data.Id,
				Name: data.Name,
			},
		})

		if err != nil {
			logger.Error("ws call error", zap.Error(err))
			client.Close()
			delete(clients, client)
		}
	}
}

func UpdateUser(data DB_UPDATE_RESPONSE) {
	logger.Info("ws update user")

	if c == nil {
		return
	}

	for client := range clients {
		err := client.WriteJSON(WS_UPDATE_USER{
			Action: "update_user",
			DB_UPDATE_RESPONSE: DB_UPDATE_RESPONSE{
				Id:   data.Id,
				Name: data.Name,
			},
		})

		if err != nil {
			logger.Error("ws call error", zap.Error(err))
			client.Close()
			delete(clients, client)
		}
	}
}

func DeleteUser(data DB_DELETE_RESPONSE) {
	logger.Info("ws delete user")

	if c == nil {
		return
	}

	for client := range clients {
		err := client.WriteJSON(WS_DELETE_USER{
			Action: "delete_user",
			DB_DELETE_RESPONSE: DB_DELETE_RESPONSE{
				Id: data.Id,
			},
		})

		if err != nil {
			logger.Error("ws call error", zap.Error(err))
			client.Close()
			delete(clients, client)
		}
	}
}
