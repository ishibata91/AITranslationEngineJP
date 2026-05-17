import type { StorybookConfig } from "@storybook/svelte-vite"
import { fileURLToPath, URL } from "node:url"

const config: StorybookConfig = {
  stories: ["../src/ui/**/*.stories.@(ts|svelte)"],
  addons: [],
  framework: {
    name: "@storybook/svelte-vite",
    options: {}
  },
  viteFinal: (config) => ({
    ...config,
    resolve: {
      ...config.resolve,
      alias: {
        ...config.resolve?.alias,
        "@ui": fileURLToPath(new URL("../src/ui", import.meta.url)),
        "@application": fileURLToPath(
          new URL("../src/application", import.meta.url)
        ),
        "@controller": fileURLToPath(
          new URL("../src/controller", import.meta.url)
        )
      }
    }
  })
}

export default config
