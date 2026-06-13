import type { Preview } from "@storybook/svelte-vite"
import "../src/ui/styles/app.css"

const preview: Preview = {
  parameters: {
    layout: "fullscreen"
  }
}

export default preview
