package main

import (
	"fmt"
	"go-chat/front_api"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

type Server struct {
	conns map[*websocket.Conn]bool
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)

}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 2048)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
		bString := buf[:n]
		msg := `<ul id="chat_room" hx-swap-oob="beforeend"><li>` + string(bString[:]) + `</li></ul>`
		s.broadcast([]byte(msg))
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}

func main() {
	s := NewServer()
	front_api.HandleEndPoints()
	http.Handle("/ws", websocket.Handler(s.handleWS))
	http.ListenAndServe(":3000", nil)
}
