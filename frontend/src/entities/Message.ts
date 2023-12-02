import type {MessageType} from "../api/chat/v1/chat_pb";
import {MessageResponse} from "../api/chat/v1/chat_pb";

export interface Message {
    type: MessageType
    name: string
    id: number
    message: string
    timestamp: Date
}

export const convertMessage = (msg: MessageResponse): Message => {
    return {
        type: msg.type,
        name: msg.name,
        id: msg.id,
        message: msg.message,
        timestamp: new Date(msg.timestamp)
    }
}



