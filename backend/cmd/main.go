package main

import (
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"oshabelt/backend/api/chat/v1/chatv1connect"
	"oshabelt/backend/grpc"
	"sync"
)

func main() {
	chat := &grpc.ChatServer{
		Users:     []grpc.User{},
		UserMutex: sync.RWMutex{},
		UserIndex: 0,
	}
	mux := http.NewServeMux()
	path, handler := chatv1connect.NewChatServiceHandler(chat)
	mux.Handle(path, handler)

	corsHandler := cors.AllowAll().Handler(h2c.NewHandler(mux, &http2.Server{}))
	err := http.ListenAndServe(
		"localhost:8085",
		corsHandler,
	)
	if err != nil {
		panic(err)
	}
}
