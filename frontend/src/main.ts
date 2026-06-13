// production の composition root（frontend 側）。App を #app へ mount する。
import { mount } from "svelte"
import "./ui/styles/app.css"
import App from "./App.svelte"

const target = document.getElementById("app")
if (target === null) {
  throw new Error("#app が見つかりません")
}

export default mount(App, { target })
