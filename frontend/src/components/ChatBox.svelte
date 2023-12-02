<script lang="ts">
    let text = ""

    export let send: (message: string) => void

    const onSend = () => {
        send(text)
        text = ""
    }

    const keydownHandler = (e: KeyboardEvent) => {
        if (e.ctrlKey && e.key === "Enter"){
            onSend()
        }
    }
</script>

<div id="chatbox">
    <div class="chatbox-container">
        <div class="chatbox-textarea-container">
            <textarea class="chatbox-textarea" bind:value={text} on:keydown={keydownHandler}/>
            <p class="chatbox-shortcut">Ctrl + Enterで送信</p>
        </div>

        <button class="send-button" on:click={onSend}>
            <svg class="send-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
                <!--!Font Awesome Free 6.5.1 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2023 Fonticons, Inc.-->
                <path d="M498.1 5.6c10.1 7 15.4 19.1 13.5 31.2l-64 416c-1.5 9.7-7.4 18.2-16 23s-18.9 5.4-28 1.6L284 427.7l-68.5 74.1c-8.9 9.7-22.9 12.9-35.2 8.1S160 493.2 160 480V396.4c0-4 1.5-7.8 4.2-10.7L331.8 202.8c5.8-6.3 5.6-16-.4-22s-15.7-6.4-22-.7L106 360.8 17.7 316.6C7.1 311.3 .3 300.7 0 288.9s5.9-22.8 16.1-28.7l448-256c10.7-6.1 23.9-5.5 34 1.4z"/>
            </svg>
        </button>
    </div>
</div>

<style lang="sass">
    #chatbox
        position: sticky
        bottom: 0
        padding-bottom: 0.8rem
        width: 100%
        height: fit-content

        .chatbox-container
            display: flex
            column-gap: 12px
            width: 95%
            max-width: 800px
            padding: 0.7rem 1rem
            margin: 0 auto
            align-items: center

            .chatbox-textarea-container
                flex-grow: 1
                height: fit-content
                position: relative

                .chatbox-textarea
                    color: white
                    width: 100%
                    resize: none
                    padding: 0.2rem
                    background-color: #353535
                    border-radius: 4px
                    overflow: hidden

                    &:focus
                        border-color: #b0b0b0
                        outline: none

                .chatbox-shortcut
                    position: absolute
                    bottom: -10px
                    right: 0
                    color: gray
                    font-size: xx-small
                    text-align: right



            .send-button
                display: flex
                align-items: center
                justify-content: center
                width: 32px
                height: 32px
                border: none
                background-color: #32cb66
                border-radius: 8px

                &:hover
                    background-color: #21a44d

                .send-icon
                    width: 16px
                    height: 16px
                    fill: white
</style>