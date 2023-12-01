package main

import (
	chatv1 "backend/api/chat/v1"
	"backend/api/chat/v1/chatv1connect"
	"connectrpc.com/connect"
	"context"
	"fmt"
	googleuuid "github.com/google/uuid"
	"github.com/rs/cors"
	"golang.org/x/exp/slog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"strconv"
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
	id   string
	uuid string
	ch   chan Message
}

type Message struct {
	name      string
	id        string
	message   string
	timestamp string
}

type ChatServer struct {
	users []User
}

func (server *ChatServer) GetUser(uuid string) (*User, bool) {
	for _, u := range server.users {
		if u.uuid == uuid {
			return &u, true
		}
	}
	return nil, false
}

func (server *ChatServer) AddUser(user *User) (*User, error) {
	if _, exist := server.GetUser(user.uuid); !exist {
		// TODO: 途中でユーザ削除が発生することを考慮し、idをindexで振りたいが、並列でリクエスト処理するとidが衝突しないか？
		u := User{user.name, strconv.Itoa(len(server.users)), user.uuid, make(chan Message)}
		server.users = append(server.users, u)
		return &u, nil
	} else {
		return nil, fmt.Errorf("failed to add user: uuid duplicated")
	}
}

func (server *ChatServer) DeleteUser(uuid string) (*User, error) {
	for i, u := range server.users {
		if u.uuid == uuid {
			server.users = append(server.users[:i], server.users[i+1:]...)
			return &u, nil
		}
	}
	return nil, fmt.Errorf("failed to delete user: uuid not found")
}

func (server *ChatServer) Broadcast(msg Message) error {
	for _, u := range server.users {
		u.ch <- msg
	}
	return nil
}

func (server *ChatServer) Connect(ctx context.Context, req *connect.Request[chatv1.ConnectRequest]) (*connect.Response[chatv1.ConnectResponse], error) {
	// Generate UUID
	uuid, err := googleuuid.NewRandom()
	if err != nil {
		slog.Error("[Connect]", "err", err)
		return nil, err
	}
	uuidStr := uuid.String()

	// Add User
	u := User{name: req.Msg.Name, uuid: uuidStr}
	if user, err := server.AddUser(&u); err == nil {
		slog.Info("[Connect]", "name", user.name, "uuid", user.uuid)
		res := connect.NewResponse(&chatv1.ConnectResponse{
			Id:   user.id,
			Uuid: user.uuid,
		})
		return res, nil
	} else {
		slog.Error("[Connect]", "err", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
}

func (server *ChatServer) Talk(ctx context.Context, req *connect.Request[chatv1.TalkRequest]) (*connect.Response[chatv1.TalkResponse], error) {
	uuid := req.Msg.Uuid
	msg := req.Msg.Message

	// uuid check
	user, ok := server.GetUser(uuid)

	if ok {
		slog.Info("[Talk]", "user", user.name, "msg", msg)
		if err := server.Broadcast(Message{
			name:      user.name,
			id:        user.id,
			message:   msg,
			timestamp: "",
		}); err != nil {
			return nil, err
		}

		res := connect.NewResponse(&chatv1.TalkResponse{
			Message: msg,
		})
		return res, nil
	} else {
		slog.Error("[Talk]", "error", "uuid not exist")
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("uuid not exist: %v", uuid))
	}
}

func (server *ChatServer) Disconnect(ctx context.Context, req *connect.Request[chatv1.DisconnectRequest]) (*connect.Response[chatv1.DisconnectResponse], error) {
	uuid := req.Msg.Uuid

	user, err := server.DeleteUser(uuid)

	if err == nil {
		slog.Info("[Disconnect]", "name", user.name, "uuid", user.uuid)
		res := connect.NewResponse(&chatv1.DisconnectResponse{})
		return res, nil
	} else {
		slog.Error("[Disconnect]", "error", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
}

func (server *ChatServer) Subscribe(ctx context.Context, req *connect.Request[chatv1.SubscribeRequest], stream *connect.ServerStream[chatv1.SubscribeStreamResponse]) error {
	uuid := req.Msg.Uuid
	user, ok := server.GetUser(uuid)

	if !ok {
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("uuid not exist"))
	}

	for msg := range user.ch {
		res := &chatv1.SubscribeStreamResponse{
			Name:    msg.name,
			Id:      msg.id,
			Message: msg.message,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		slog.Info("[Subscribe]", "send_to", user.name, "uuid", user.uuid)
	}

	return nil
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
