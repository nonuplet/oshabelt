<script lang="ts">
    import type {Message} from "../entities/Message"
    import {MessageType} from "../api/chat/v1/chat_pb";
    import {client} from "../store";

    export let message: Message


    $: time = message.timestamp.toLocaleTimeString("ja-JP")
    $: messageStyle = message.id === $client.user.id ? "message self" : "message"
</script>

{#if message.type === MessageType.MSG_TALK}
    <div class={messageStyle}>
        <h3 class="name">{message.name}</h3>
        <p class="text">{message.message}</p>
        <p class="timestamp">{time}</p>
    </div>
{/if}

{#if message.type === MessageType.MSG_CONNECT}
    <div class="connected">
        {#if message.id === $client.user.id}
            <p class="text">ルームに接続しました</p>
        {:else}
            <p class="text">{message.name} さんが入室しました</p>
        {/if}
    </div>
{/if}

<style lang="sass">
  .message
    background-color: #3e3e3e
    border: rgba(0, 0, 0, 0.2) 1px solid
    border-radius: 10px
    padding: 0.3rem 1.5rem
    margin-bottom: 10px

    &.self
      background-color: #4c5b4a

    .name
      font-size: large
      margin-bottom: 3px

    .text
      padding: 0 1rem

    .timestamp
      text-align: right
      color: #909090
      font-size: small

  .connected
    width: 100%
    margin: 20px 0
    padding: 0 1rem

    .text
      color: #b0b0b0
      text-align: center

</style>