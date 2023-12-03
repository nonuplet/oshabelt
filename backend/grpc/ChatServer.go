package grpc

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	googleuuid "github.com/google/uuid"
	"golang.org/x/exp/slog"
	chatv1 "oshabelt/backend/api/chat/v1"
	"sync"
	"time"
)

type User struct {
	name string
	id   uint32
	uuid string
	ch   chan Message
}

type Message struct {
	msgType   chatv1.MessageType
	name      string
	id        uint32
	message   string
	timestamp string
}

type ChatServer struct {
	Users     []User
	UserMutex sync.RWMutex
	UserIndex uint32
}

func (server *ChatServer) CurrentTime() string {
	current := time.Now()
	return current.Format(time.RFC3339)
}

func (server *ChatServer) GetUser(uuid string) (*User, bool) {
	for _, u := range server.Users {
		if u.uuid == uuid {
			return &u, true
		}
	}
	return nil, false
}

func (server *ChatServer) AddUser(user *User) (*User, error) {
	server.UserMutex.Lock()
	defer server.UserMutex.Unlock()

	if _, exist := server.GetUser(user.uuid); !exist {
		u := User{user.name, server.UserIndex, user.uuid, make(chan Message)}
		server.Users = append(server.Users, u)
		server.UserIndex++
		return &u, nil
	} else {
		return nil, fmt.Errorf("failed to add user: uuid duplicated")
	}
}

func (server *ChatServer) DeleteUser(uuid string) (*User, error) {
	for i, u := range server.Users {
		if u.uuid == uuid {
			server.Users = append(server.Users[:i], server.Users[i+1:]...)
			return &u, nil
		}
	}
	return nil, fmt.Errorf("failed to delete user: uuid not found")
}

func (server *ChatServer) Broadcast(msg Message) {
	server.UserMutex.RLock()
	defer server.UserMutex.RUnlock()
	for _, u := range server.Users {
		if u.id == msg.id { // 自分のtalkメッセージはbroadcastしない
			continue
		}
		u.ch <- msg
	}
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
	user, err := server.AddUser(&u)
	if err == nil {
		slog.Info("[Connect]", "name", user.name, "uuid", user.uuid)
		res := connect.NewResponse(&chatv1.ConnectResponse{
			Id:   user.id,
			Uuid: user.uuid,
		})
		broad := Message{
			msgType:   chatv1.MessageType_MSG_CONNECT,
			name:      user.name,
			id:        user.id,
			message:   "",
			timestamp: server.CurrentTime(),
		}
		go server.Broadcast(broad)
		return res, nil
	} else {
		slog.Error("[Connect]", "err", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
}

func (server *ChatServer) Disconnect(ctx context.Context, req *connect.Request[chatv1.DisconnectRequest]) (*connect.Response[chatv1.DisconnectResponse], error) {
	server.UserMutex.Lock()
	server.UserMutex.Unlock()
	uuid := req.Msg.Uuid
	user, err := server.DeleteUser(uuid)

	if err == nil {
		slog.Info("[Disconnect]", "name", user.name, "uuid", user.uuid)
		close(user.ch)
		res := connect.NewResponse(&chatv1.DisconnectResponse{})
		return res, nil
	} else {
		slog.Error("[Disconnect]", "error", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
}

func (server *ChatServer) Talk(ctx context.Context, req *connect.Request[chatv1.TalkRequest]) (*connect.Response[chatv1.MessageResponse], error) {
	uuid := req.Msg.Uuid
	msg := req.Msg.Message

	// uuid check
	user, ok := server.GetUser(uuid)

	if ok {
		slog.Info("[Talk]", "user", user.name, "msg", msg)
		talk := Message{
			msgType:   chatv1.MessageType_MSG_TALK,
			name:      user.name,
			id:        user.id,
			message:   msg,
			timestamp: server.CurrentTime(),
		}
		go server.Broadcast(talk)

		res := connect.NewResponse(&chatv1.MessageResponse{
			Type:      talk.msgType,
			Name:      talk.name,
			Id:        talk.id,
			Message:   talk.message,
			Timestamp: talk.timestamp,
		})
		return res, nil
	} else {
		slog.Error("[Talk]", "error", "uuid not exist")
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("uuid not exist: %v", uuid))
	}
}

func (server *ChatServer) Subscribe(ctx context.Context, req *connect.Request[chatv1.SubscribeRequest], stream *connect.ServerStream[chatv1.MessageResponse]) error {
	uuid := req.Msg.Uuid
	user, ok := server.GetUser(uuid)

	if !ok {
		err := fmt.Errorf("uuid not exist")
		slog.Error("[Subscribe]", "error", err)
		return connect.NewError(connect.CodeNotFound, err)
	}

	slog.Info("[Subscribe]", "send_to", user.name, "uuid", user.uuid)
	for msg := range user.ch {
		res := &chatv1.MessageResponse{
			Type:      msg.msgType,
			Name:      msg.name,
			Id:        msg.id,
			Message:   msg.message,
			Timestamp: msg.timestamp,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}

	return nil
}
