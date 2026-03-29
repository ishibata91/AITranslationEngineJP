import allowlists from "./config/lint/allowlists.json" with { type: "json" };

export default {
  $schema: "https://unpkg.com/knip@5/schema.json",
  entry: ["index.html"],
  project: ["src/**/*.{ts,svelte}", "*.config.{js,mjs,ts}"],
  ignoreBinaries: allowlists.knipIgnoreBinaries,
  ignoreDependencies: allowlists.knipIgnoreDependencies,
  vite: true,
  vitest: true
};
