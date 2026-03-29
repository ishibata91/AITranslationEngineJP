import { mount } from "svelte";
import App from "./App.svelte";
import "./app.css";
import { configureBootstrapStatusPort } from "@application/bootstrap/load-bootstrap-status";
import { tauriBootstrapStatusGateway } from "@gateway/tauri/bootstrap-status.gateway";

configureBootstrapStatusPort(tauriBootstrapStatusGateway);

mount(App, {
  target: document.getElementById("app")!
});
