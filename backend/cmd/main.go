package main

import (
	chatv1 "backend/api/chat/v1"
	"backend/api/chat/v1/chatv1connect"
	"connectrpc.com/connect"
	"context"
	"fmt"
	googleuuid "github.com/google/uuid"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
)

func main() {
	chat := &ChatServer{}
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

type User struct {
	name string
	uuid string
}

type ChatServer struct {
	users []User
}

func (server *ChatServer) Connect(ctx context.Context, req *connect.Request[chatv1.ConnectRequest]) (*connect.Response[chatv1.ConnectResponse], error) {
	// Generate UUID
	var user User
	for user == (User{}) {
		uuid, err := googleuuid.NewRandom()
		if err != nil {
			log.Printf("uuid generate error : \n %v", err)
			continue
		}
		uuidStr := uuid.String()

		ok := true
		for _, u := range server.users {
			if u.uuid == uuidStr {
				ok = false
				break
			}
		}
		if ok {
			user = User{
				req.Msg.Name,
				uuidStr,
			}
		}
	}
	server.users = append(server.users, user)

	res := connect.NewResponse(&chatv1.ConnectResponse{
		Uuid: user.uuid,
	})
	//res.Header().Set("Chat-Version", "v1")
	return res, nil
}

func (server *ChatServer) Talk(ctx context.Context, req *connect.Request[chatv1.TalkRequest]) (*connect.Response[chatv1.TalkResponse], error) {
	uuid := req.Msg.Uuid
	msg := req.Msg.Message

	// uuid check
	ok := false

	fmt.Printf("%v", server.users)
	for _, u := range server.users {
		if u.uuid == uuid {
			ok = true
			break
		}
	}

	if ok {
		res := connect.NewResponse(&chatv1.TalkResponse{
			Message: msg,
		})
		return res, nil
	} else {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("not exist uuid: %v", uuid))
	}
}

//type chatServer struct {
//	chat.UnimplementedChatServiceServer
//}
//
//func NewChatServer() *chatServer {
//	return &chatServer{}
//}
//
//func (server *chatServer) Hello(ctx context.Context, req *chat.HelloRequest) (*chat.HelloResponse, error) {
//	return &chat.HelloResponse{
//		Message: fmt.Sprintf("Hello %s", req.GetName()),
//	}, nil
//}
//
//func (server *chatServer) HelloStream(req *chat.HelloRequest, stream chat.ChatService_HelloStreamServer) error {
//	count := 5
//	for i := 0; i < count; i++ {
//		if err := stream.Send(&chat.HelloResponse{
//			Message: fmt.Sprintf("[%d] Hello, %s!", i, req.GetName()),
//		}); err != nil {
//			return err
//		}
//		time.Sleep(time.Second * 1)
//	}
//	return nil
//}
