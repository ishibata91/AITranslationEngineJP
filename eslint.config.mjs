import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import svelte from "eslint-plugin-svelte";
import svelteParser from "svelte-eslint-parser";
import { repositoryBoundaryPlugin } from "./scripts/eslint/repository-boundary-plugin.mjs";
import allowlists from "./config/lint/allowlists.json" with { type: "json" };

export default [
  {
    ignores: allowlists.pathIgnores
  },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  ...svelte.configs["flat/recommended"],
  {
    files: ["src/**/*.{ts,svelte}"],
    plugins: {
      repository: repositoryBoundaryPlugin
    },
    languageOptions: {
      globals: {
        ...globals.browser
      }
    },
    rules: {
      "repository/enforce-layer-boundaries": "error"
    }
  },
  {
    files: ["src/**/*.svelte"],
    languageOptions: {
      parser: svelteParser,
      parserOptions: {
        parser: tseslint.parser,
        extraFileExtensions: [".svelte"]
      }
    }
  },
  {
    files: ["src/ui/**/*.svelte", "src/ui/**/*.ts"],
    rules: {
      "no-restricted-imports": [
        "error",
        {
          "patterns": [
            {
              "group": ["@gateway/*", "@tauri-apps/api/*"],
              "message": "UI views and stores must go through application ports, not gateway or Tauri APIs."
            }
          ]
        }
      ]
    }
  },
  {
    files: ["src/application/**/*.ts"],
    rules: {
      "no-restricted-imports": [
        "error",
        {
          "patterns": [
            {
              "group": ["@tauri-apps/api/*"],
              "message": "Application code may depend on gateway ports, not Tauri APIs directly."
            }
          ]
        }
      ]
    }
  },
  {
    files: ["src/shared/**/*.ts"],
    rules: {
      "no-restricted-imports": [
        "error",
        {
          "patterns": [
            {
              "group": ["@tauri-apps/api/*"],
              "message": "Shared contracts must stay Tauri-agnostic."
            }
          ]
        }
      ]
    }
  }
];
