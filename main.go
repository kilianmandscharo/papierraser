package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"github.com/kilianmandscharo/papierraser/components"
	"github.com/kilianmandscharo/papierraser/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return r.Header.Get("Origin") == "http://localhost:8080"
}

var track = types.Track{
	Width:  20,
	Height: 10,
	Outer: []types.Point{
		{X: 0, Y: 0}, {X: 20, Y: 0}, {X: 20, Y: 10}, {X: 0, Y: 10},
	},
	Inner: []types.Point{
		{X: 4, Y: 3}, {X: 16, Y: 3}, {X: 16, Y: 7}, {X: 4, Y: 7},
	},
	Finish: [2]types.Point{
		{X: 0, Y: 5}, {X: 4, Y: 5},
	},
}

type ConnectionRequestType = string

const (
	ConnectionRequestConnect    ConnectionRequestType = "connect"
	ConnectionRequestDisconnect ConnectionRequestType = "disconnect"
)

type ConnectionRequest struct {
	Type    ConnectionRequestType
	Request *http.Request
	Conn    *websocket.Conn
	GameId  string
}

func main() {
	ch := make(chan ConnectionRequest)

	go connectionHandler(ch)

	http.Handle("/", templ.Handler(components.Index()))
	http.HandleFunc("/ws", websocketHandler(ch))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func broadcastConnections(conns types.Connections) {
	for _, conn := range conns {
		if conn == nil {
			continue
		}

		var buf bytes.Buffer
		err := components.Lobby(conns).Render(
			context.Background(),
			&buf,
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := conn.WriteMessage(1, buf.Bytes()); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func connectionHandler(ch <-chan ConnectionRequest) {
	conns := make(map[string]types.Connections)

	for message := range ch {
		remoteAddr := message.Request.RemoteAddr
		gameId := message.GameId

		switch message.Type {
		case ConnectionRequestConnect:
			if conns[gameId] == nil {
				conns[gameId] = make(types.Connections)
			}
			conns[gameId][remoteAddr] = message.Conn
			fmt.Printf("connected %s to %s\n", remoteAddr, gameId)
			broadcastConnections(conns[gameId])
		case ConnectionRequestDisconnect:
			conns[gameId][remoteAddr] = nil
			fmt.Printf("disconnected %s from %s\n", remoteAddr, gameId)
			broadcastConnections(conns[gameId])
		default:
			fmt.Println("unknown ConnectionRequestType:", message.Type)
		}
	}
}

func getGameId(r *http.Request) (string, bool) {
	queryStrings := r.URL.Query()["id"]

	if len(queryStrings) == 0 {
		return "", false
	}

	return queryStrings[0], true
}

func websocketHandler(ch chan<- ConnectionRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		gameId, ok := getGameId(r)
		if !ok {
			fmt.Println("no game id provided")
			return
		}

		defer func() {
			ch <- ConnectionRequest{
				Type:    ConnectionRequestDisconnect,
				Request: r,
				Conn:    conn,
				GameId:  gameId,
			}
			conn.Close()
		}()

		ch <- ConnectionRequest{
			Type:    ConnectionRequestConnect,
			Request: r,
			Conn:    conn,
			GameId:  gameId,
		}

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(messageType, p)
		}

		// race := types.NewRace(track)
		// race.AddPlayer("Mason")
		//
		// for {
		// 	messageType, p, err := conn.ReadMessage()
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return
		// 	}
		//
		// 	var point types.Point
		// 	err = json.Unmarshal(p, &point)
		// 	if err != nil {
		// 		fmt.Println("ERROR:", err)
		// 		return
		// 	}
		//
		// 	race.Move(point)
		//
		// 	var buf bytes.Buffer
		// 	err = components.Track(race).Render(context.Background(), &buf)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return
		// 	}
		//
		// 	if err := conn.WriteMessage(messageType, buf.Bytes()); err != nil {
		// 		fmt.Println(err)
		// 		return
		// 	}
		// }
	}
}
