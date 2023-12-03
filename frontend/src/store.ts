import {writable} from "svelte/store";
import {ChatClient} from "./entities/ChatClient";
import type {MessageResponse} from "./api/chat/v1/chat_pb";

export const client = writable<ChatClient>()
export const messages = writable<MessageResponse[]>()