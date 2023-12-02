import './app.css'
import App from './App.svelte'

const app = new App({
  target: document.body
})

export default app


import {createPromiseClient} from "@connectrpc/connect";
import {createConnectTransport} from "@connectrpc/connect-web";
import {ChatService} from "./api/chat/v1/chat_connect";
import {onMount} from "svelte";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8085"
})

const client = createPromiseClient(ChatService, transport)
let uuid: string;

async function getStream(uuid: string) {
  const stream = client.subscribe({
    uuid
  })
  try {
    for await (const message of stream) {
      console.log(message)
    }
  } catch (error) {
    console.error(error)
  }
}

(async () => {
  const res = await client.connect({
    name: "hoge fuga"
  })
  console.log(res)
  uuid = res.uuid

  getStream(res.uuid)


  setTimeout(async () => {
    const talk = await client.talk({
      uuid: res.uuid,
      message: "test message"
    })
  }, 1000)

})()

setTimeout(async () => {
  console.log("disconnecting...");
  const dlt = await client.disconnect({
    uuid
  })
  console.log(dlt)
}, 4000)