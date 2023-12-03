import './app.sass'
import App from './App.svelte'

const root = document.getElementById("app") as HTMLElement

// @ts-ignore
const app = new App({
  target: root
})

export default app