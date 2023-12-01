// @generated by protoc-gen-connect-es v1.1.3 with parameter "target=ts"
// @generated from file chat/v1/chat.proto (package chat.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { ConnectRequest, ConnectResponse, DisconnectRequest, DisconnectResponse, SubscribeRequest, SubscribeStreamResponse, TalkRequest, TalkResponse } from "./chat_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service chat.v1.ChatService
 */
export const ChatService = {
  typeName: "chat.v1.ChatService",
  methods: {
    /**
     * @generated from rpc chat.v1.ChatService.Connect
     */
    connect: {
      name: "Connect",
      I: ConnectRequest,
      O: ConnectResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc chat.v1.ChatService.Talk
     */
    talk: {
      name: "Talk",
      I: TalkRequest,
      O: TalkResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc chat.v1.ChatService.Disconnect
     */
    disconnect: {
      name: "Disconnect",
      I: DisconnectRequest,
      O: DisconnectResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc chat.v1.ChatService.Subscribe
     */
    subscribe: {
      name: "Subscribe",
      I: SubscribeRequest,
      O: SubscribeStreamResponse,
      kind: MethodKind.ServerStreaming,
    },
  }
} as const;

