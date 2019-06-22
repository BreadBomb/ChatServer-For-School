package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"goehlerChatServer/model"
	"goehlerChatServer/util"
	"gopkg.in/square/go-jose.v1/json"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	Sessions map[uuid.UUID]model.Session
)

func main() {

	Sessions = make(map[uuid.UUID]model.Session)

	http.HandleFunc("/chat", HandleChat)

	http.HandleFunc("/room/create", CreateRoom)
	http.HandleFunc("/room/join", JoinRoom)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":8080", nil)
}


func CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	motd := r.URL.Query().Get("motd")

	if motd == "" {
		w.Write([]byte("motd needed in order to create room"))
		return
	}

	session := model.Session{}
	session.Id = uuid.New()
	session.CreateTime = util.GetCurrentTimestamp()
	session.Motd = motd
	session.Users = make(map[uuid.UUID]model.User)

	w.WriteHeader(200)

	response := model.Response{}

	response.Code = 200
	response.Time = util.GetCurrentTimestamp()
	response.Message = "Room successfull created"
	response.Data = session

	json, err := json.Marshal(response)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	Sessions[session.Id] = session

	w.Write(json)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	sessionId, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	username := r.URL.Query().Get("username")

	session := Sessions[sessionId]

	user := model.NewUser(username)

	user.SessionMotd = session.Motd

	session.Users[user.Id] = user

	fmt.Println(Sessions[sessionId].Users)

	response := model.Response{}

	response.Code = 200
	response.Time = util.GetCurrentTimestamp()
	response.Message = "Successfully joined room"
	response.Data = user

	json, err := json.Marshal(response)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(json)
}

func HandleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	sessionId, err := uuid.Parse(r.URL.Query().Get("session"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	userId, err := uuid.Parse(r.URL.Query().Get("user"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var user model.User

	user = Sessions[sessionId].Users[userId]

	user.Connection = conn

	Sessions[sessionId].Users[userId] = user

	//Sessions[sessionId].Users[userId] = user

	fmt.Println(Sessions[sessionId].Users)

	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()

		if err != nil {
			return
		}

		var message model.Message

		fmt.Println(msg)

		err = json.Unmarshal(msg, message)
		if err != nil {
			fmt.Println("Message cant be unmarshaled")
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		for _, v := range Sessions[sessionId].Users {
			fmt.Println(v)
			if v.Connection != nil {
				v.Connection.WriteJSON(message)
			}
		}
	}
}