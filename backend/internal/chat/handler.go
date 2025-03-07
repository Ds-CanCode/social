package funcs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	pkg "funcs/pkg"

	"github.com/gorilla/websocket"
)

func SendRealTimeNotification(userId []int, notif interface{}) {
	mu.Lock()
	for _, user := range userId {
		for _, conn := range conns[user] {
			conn.WriteJSON(notif)
		}
	}
	mu.Unlock()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	conns = make(map[int][]*websocket.Conn)
	mu    = &sync.Mutex{}
)

func GetConns() (map[int][]*websocket.Conn, *sync.Mutex) {
	return conns, mu
}

func ChatService(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		mu.Lock()
		conn.Close()
		mu.Unlock()
	}()

	mu.Lock()
	conns[id] = append(conns[id], conn)
	mu.Unlock()

	for {
		var exchangeData WebsocketData
		err := conn.ReadJSON(&exchangeData)
		if err != nil {
			log.Println("Read Error:", err)
			return
		}

		fmt.Println(exchangeData.Type)
		mu.Lock()
		switch exchangeData.Type {
		case "Private":

			var privateMsg Private_Message
			if err := json.Unmarshal(exchangeData.Data, &privateMsg); err != nil {

				log.Println("JSON Decode Error:", err)
				break
			}
			err = PrivateMessages(id, privateMsg.ReceiverID, privateMsg.Content)
			if err != nil {
				break
			}

			privateMsg.SenderID = id

			// Re-marshal the updated message back into the exchangeData.Data
			updatedData, err := json.Marshal(privateMsg)
			if err != nil {
				log.Println("JSON Marshal Error:", err)
				break
			}
			exchangeData.Data = updatedData

			broadcastMessage(privateMsg.ReceiverID, exchangeData)

		case "Group":
			var groupMsg Group_Mesaage
			if err := json.Unmarshal(exchangeData.Data, &groupMsg); err != nil {
				log.Println("JSON Decode Error:", err)
				break
			}
			fmt.Println(groupMsg.Group)
			err = GroupMessages(id, groupMsg.Group, groupMsg.Content)
			if err != nil {
				fmt.Println(err)
				break
			}

			groupMsg.SenderID = id
			groupMsg.SenderName , _ = Getusernamebyid(id) 



			// Re-marshal the updated message back into the exchangeData.Data
			updatedData, err := json.Marshal(groupMsg)
			if err != nil {
				log.Println("JSON Marshal Error:", err)
				break
			}
			exchangeData.Data = updatedData
			broadcastToGroup(id , groupMsg.Group, exchangeData)
		}
		mu.Unlock()
	}
}

func broadcastMessage(receiverID int, message WebsocketData) {
	for _, conn := range conns[receiverID] {
		conn.WriteJSON(message)
	}
	// for _, conn := range conns[senderID] {
	// 	fmt.Println(message)
	// 	conn.WriteJSON(message)
	// }
}

func broadcastToGroup(id int , groupID int, message WebsocketData) {
	userIDs, err := GetGroupusers(groupID)
	if err != nil {
		log.Println("Error getting group users:", err)
		return
	}
	for _, userID := range userIDs {
		if userID != id {
			for _, conn := range conns[userID] {
				conn.WriteJSON(message)
			}

		}
	}
}

func ChatHistory(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	query := r.URL.Query()
	recivierId, _ := strconv.Atoi(query.Get("recivierID"))
	if recivierId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}

	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}

	var message []Private_Message
	message, err = GetHistoryPrv(id, recivierId, offset)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &message)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func ChatGroupHistory(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		fmt.Println(err)
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	
	query := r.URL.Query()
	groupId, _ := strconv.Atoi(query.Get("groupId"))
	if groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}
	
	///check if member
	_, err = CheckIfMember(id, groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusForbidden, err)
		return
	}

	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}

	var message []Group_Mesaage
	message, err = GetGroupHistoryPrv(id, groupId, offset)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}


	err = pkg.Encode(w, &message)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}
