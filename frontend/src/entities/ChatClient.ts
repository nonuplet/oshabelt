import type {User} from "./User";
import type {PromiseClient} from "@connectrpc/connect";
import {createConnectTransport} from "@connectrpc/connect-web";
import {createPromiseClient} from "@connectrpc/connect";
import {ChatService} from "../api/chat/v1/chat_connect";
import type {MessageResponse} from "../api/chat/v1/chat_pb";

export class ChatClient {
    user: User
    client: PromiseClient<typeof ChatService>

    constructor(url: string) {
        this.user = {
            name: "",
            id: 0,
            uuid: "",
            connecting: false
        }
        const transport = createConnectTransport({
            baseUrl: url
        })
        this.client = createPromiseClient(ChatService, transport)
    }

    async connect(name: string) {
        if (this.user.connecting) {
            throw new Error("[Connect] Already connected.")
        }
        try {
            const res = await this.client.connect({
                name
            })
            this.user = {
                name,
                id: res.id,
                uuid: res.uuid,
                connecting: true
            }
        } catch (e) {
            console.error(e)
        }
    }

    async disconnect() {
        if (!this.user.connecting) {
            throw new Error("[Disconnect] No connection.")
        }
        try {
            this.user.connecting = false
            const res = await this.client.disconnect({
                uuid: this.user.uuid
            })
        } catch (e) {
            this.user.connecting = true
            console.error(e)
        }
    }

    async subscribe(callback: (message: MessageResponse) => void) {
        const stream = this.client.subscribe({
            uuid: this.user.uuid
        })
        try {
            for await (const message of stream) {
                callback(message)
            }
        } catch (error) {
            console.error(error)
        }
    }

    async talk(message: string) {
        if (!this.user.connecting) {
            console.error("[Talk]No connection.")
            return
        }

        try {
            const res = this.client.talk({
                uuid: this.user.uuid,
                message
            })
            return res
        } catch (e) {
            console.log(e)
        }
    }
}