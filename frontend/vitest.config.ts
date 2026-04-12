import { defineConfig, mergeConfig } from "vitest/config"
import viteConfig from "./vite.config"

export default mergeConfig(viteConfig, defineConfig({
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./src/test/setup.ts"],
    include: ["src/**/*.test.ts"],
    coverage: {
      provider: "v8",
      reportsDirectory: "../test-results/frontend-coverage",
      reporter: ["text", "json-summary", "lcov"],
      include: ["src/**/*.ts"],
      exclude: ["src/**/*.test.ts", "src/test/**"]
    }
  }
}))
