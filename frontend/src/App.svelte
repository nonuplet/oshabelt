<script lang="ts">
    import {ChatClient} from "./entities/ChatClient";
    import type {MessageResponse} from "./api/chat/v1/chat_pb";
    import Header from "./components/Header.svelte";

    const client = new ChatClient("http://localhost:8085")

    let name = ""
    let message = ""

    let talks: MessageResponse[] = []

    const connect = async () => {
        if (name === "") return
        try {
            await client.connect(name)
            client.subscribe(onMessage)
        } catch (e) {
            console.error(e)
        }
    }
    const disconnect = async () => {
        await client.disconnect()
    }

    const talk = async () => {
        if (message === "") return
        await client.talk(message)
    }

    const onMessage = (message: MessageResponse) => {
        console.log(message)
        talks = [...talks, message]
    }
</script>

<Header/>

<main>
    <div class="control">
        <div class="connect">
            <span>name:</span>
            <input bind:value={name} placeholder="名前を入力してください"/>
            <button class="connect-button" on:click={connect}>Connect</button>
            <button on:click={disconnect}>Disconnect</button>
        </div>

        <div class="talk">
            <input bind:value={message}/>
            <button on:click={talk}>Send</button>
        </div>
    </div>
    <div class="message-container">
        {#each talks as talk}
            <div class="message-box">
                <p class="message-name">{talk.name}</p>
                <p class="message-text">{talk.message}</p>
            </div>
        {/each}
    </div>
</main>

<style lang="sass">
    main
        width: 100vw
        margin: 0
        height: 0

    .control
        border: #808080 2px solid
        padding: 1rem

        .connect
            margin-bottom: 2rem

            .connect-button
                margin-right: 4rem


    .message-container
        width: 100%
        min-height: 30rem
        margin: 0
        padding: 0.5rem 1rem

        .message-box
            background-color: #eaeaea
            width: 100%
            margin-bottom: 2rem
            color: black

            .message-name
                width: 100%

            .message-text
                width: 100%


</style>
