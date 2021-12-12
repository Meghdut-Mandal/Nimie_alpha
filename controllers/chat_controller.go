package controllers

import (
	"Nimie_alpha/models"
	"Nimie_alpha/utils"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Map of connection id to list of active clients
var clients = make(map[int64][]*websocket.Conn) // connected clients
var broadcast = make(chan models.ChatMessage)   // broadcast channel
// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleChatConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	conversationId := utils.GetConversationId(r)
	// add the connection to the map
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(ws)

	clients[conversationId] = append(clients[conversationId], ws)

	for {
		var msg models.ChatMessage
		// Read in a new message as JSON and map it to a ChatMessage object
		err := ws.ReadJSON(&msg)
		msg.ConversationId = conversationId
		models.AddMessage(&msg)
		// print the message to the console
		if err != nil {
			log.Printf("error: %v", err)
			deleteClient(clients[conversationId], ws)
			break
		}
		println("message sent by ", msg.UserId)
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		println("Broadcast message ", msg.Message)
		// Send it out to every client that is currently connected
		// Get the clint list having the same Conversation id
		clientList := clients[msg.ConversationId]

		for _, client := range clientList {
			// Send the message
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				// Remove the client from the list
				deleteClient(clientList, client)
			}
		}
	}
}

func deleteClient(clientList []*websocket.Conn, client *websocket.Conn) {
	for i, c := range clientList {
		if c == client {
			clientList = append(clientList[:i], clientList[i+1:]...)
			break
		}
	}
}