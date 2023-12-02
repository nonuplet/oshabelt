<script lang="ts">
    import {ChatClient} from "./entities/ChatClient";
    import type {MessageResponse} from "./api/chat/v1/chat_pb";
    import Header from "./components/Header.svelte";
    import ChatBox from "./components/ChatBox.svelte";
    import {client} from "./store";
    import RegisterOverlay from "./components/RegisterOverlay.svelte";
    import MessageCard from "./components/MessageCard.svelte";
    import type {Message} from "./entities/Message";
    import {convertMessage} from "./entities/Message";
    import {onMount, tick} from "svelte";

    client.set(new ChatClient("http://localhost:8085"))

    let connecting = false
    let messageContainer: HTMLElement

    let talks: Message[] = []

    const connect = async (name: string) => {
        if (name === "") return
        try {
            await $client.connect(name)
            connecting = true
            $client.subscribe(onMessage)
        } catch (e) {
            console.error(e)
        }
    }
    const disconnect = async () => {
        await $client.disconnect()
    }

    const talk = async (text: string) => {
        if (text === "") return
        await $client.talk(text)
    }

    const onMessage = async (message: MessageResponse) => {
        console.log(message)
        const talk = convertMessage(message)
        talks = [...talks, talk]
        autoScroll()
    }

    const autoScroll = async () => {
        await tick()
        messageContainer.scrollTo({
            top: messageContainer.scrollHeight,
            left: 0,
            behavior: "smooth"
        })
    }

    onMount(() => {
        messageContainer = document.getElementById("message-container") as HTMLElement
    })
</script>


{#if !connecting}
    <RegisterOverlay connect={connect}/>
{/if}
<main>
    <Header/>

    <div id="message-container">
        {#each talks as talk}
            <MessageCard message={talk}/>
        {/each}
    </div>

    <ChatBox send={talk}/>
</main>

<style lang="sass">
    main
        height: 100%
        display: flex
        flex-flow: column

    #message-container
        width: 100%
        flex-grow: 1
        overflow-y: scroll
        margin: 0
        padding: 0.5rem 1rem

        $sb-track-color: #2b2b2b
        $sb-thumb-color: #b0b0b0
        $sb-size: 10px
        scrollbar-color: $sb-thumb-color $sb-track-color

        &::-webkit-scrollbar
            width: $sb-size

        &::-webkit-scrollbar-track
            background: $sb-track-color
            border-radius: $sb-size

        &::-webkit-scrollbar-thumb
            background: $sb-thumb-color
            border-radius: $sb-size
</style>
