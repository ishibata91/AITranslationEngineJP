import { defineConfig } from "vite"
import { svelte } from "@sveltejs/vite-plugin-svelte"
import { fileURLToPath, URL } from "node:url"

export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      "@ui": fileURLToPath(new URL("./src/ui", import.meta.url)),
      "@application": fileURLToPath(
        new URL("./src/application", import.meta.url)
      ),
      "@controller": fileURLToPath(new URL("./src/controller", import.meta.url))
    }
  },
  build: {
    outDir: "dist"
  }
})
