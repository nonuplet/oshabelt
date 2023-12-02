<script lang="ts">
    let name = ""
    export let connect: (name: string) => void

    const onConnect = () => {
        connect(name)
    }

    const keydownHandler = (e: KeyboardEvent) => {
        if (e.key === "Enter") connect(name)
    }
</script>

<div id="register">
    <div class="dialog">
        <div class="description">
            <h2 class="title">ルーム入室</h2>
            <div class="desc">
                <span class="main">名前を入力してください</span>
                <span class="sub">(他の人に公開されます)</span>
            </div>
        </div>
        <div class="input">
            <input class="text" bind:value={name} on:keydown={keydownHandler}/>
            <button class="join" on:click={onConnect}>入室</button>
        </div>
    </div>
</div>

<style lang="sass">
    #register
        position: absolute
        z-index: 10
        width: 100%
        height: 100%
        background-color: rgba(0, 0, 0, 0.8)

        .dialog
            position: absolute
            top: 0
            bottom: 0
            left: 0
            right: 0
            margin: auto
            width: 100%
            max-width: 550px
            height: 100%
            max-height: 200px
            padding: 1rem 2rem
            background-color: #202020
            border-radius: 10px

            .description
                container-type: inline-size
                margin-bottom: 2rem

                .title
                    border-left: white 4px solid
                    padding: 0.3rem 0 0.3rem 1rem
                    margin-bottom: 1rem

                .desc
                    padding: 0 1rem


                    .sub
                        color: #a0a0a0
                        font-size: small


                    @media (max-width: 500px)
                        .main
                            display: block
                            width: 100%

                        .sub
                            display: block
                            width: 100%


            .input
                display: flex
                margin-top: auto
                padding: 0 1.5rem
                column-gap: 1rem
                align-items: center

                .text
                    flex-grow: 1
                    background-color: #202020
                    border: #606060 1px solid
                    padding: 0.3rem 1rem
                    border-radius: 5px
                    color: white
                    outline: none

                    &:focus
                        border: #a0a0a0 1px solid

                .join
                    color: white
                    min-width: 50px
                    padding: 4px 10px
                    background-color: #198829
                    border: #0e6c1c 1px solid
                    border-radius: 5px

                    &:hover
                        background-color: #0a5b16


</style>