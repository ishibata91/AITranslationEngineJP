/** @type {import("prettier").Config} */
const config = {
  semi: false,
  singleQuote: false,
  trailingComma: "none",
  plugins: ["prettier-plugin-svelte"],
  overrides: [
    {
      files: "*.svelte",
      options: {
        parser: "svelte"
      }
    }
  ]
}

export default config
