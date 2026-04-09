import js from "@eslint/js";
import { defineConfig } from "eslint/config";
import { createTypeScriptImportResolver } from "eslint-import-resolver-typescript";
import { importX } from "eslint-plugin-import-x";
import svelte from "eslint-plugin-svelte";
import globals from "globals";
import { dirname } from "node:path";
import { fileURLToPath } from "node:url";
import svelteConfig from "./svelte.config.js";
import tseslint from "typescript-eslint";
import { repositoryBoundaryPlugin } from "../scripts/eslint/repository-boundary-plugin.mjs";

const rootDir = dirname(fileURLToPath(import.meta.url));
const extraFileExtensions = [".svelte"];

export default defineConfig([
  {
    ignores: ["dist/**", "node_modules/**", "wailsjs/**"]
  },
  {
    linterOptions: {
      reportUnusedDisableDirectives: "error"
    }
  },
  js.configs.recommended,
  ...tseslint.configs.recommendedTypeChecked,
  ...svelte.configs["flat/recommended"],
  {
    settings: {
      "import-x/extensions": [".js", ".ts", ".svelte"],
      "import-x/resolver-next": [
        createTypeScriptImportResolver({
          project: "./tsconfig.json"
        })
      ]
    },
    languageOptions: {
      globals: {
        ...globals.browser
      },
      parserOptions: {
        extraFileExtensions,
        projectService: true,
        tsconfigRootDir: rootDir
      }
    }
  },
  {
    files: ["**/*.svelte", "**/*.svelte.ts", "**/*.svelte.js"],
    languageOptions: {
      parserOptions: {
        extraFileExtensions,
        parser: tseslint.parser,
        projectService: true,
        svelteConfig,
        tsconfigRootDir: rootDir
      }
    }
  },
  {
    files: ["src/**/*.{ts,svelte}"],
    plugins: {
      "import-x": importX,
      local: repositoryBoundaryPlugin
    },
    rules: {
      "local/no-commented-out-code": "error",
      "no-unreachable": "error",
      "no-unreachable-loop": "error",
      "import-x/consistent-type-specifier-style": ["error", "prefer-top-level"],
      "import-x/no-duplicates": ["error", { "prefer-inline": true }],
      "local/enforce-layer-boundaries": "error"
    }
  },
  {
    files: ["*.js"],
    extends: [tseslint.configs.disableTypeChecked],
    languageOptions: {
      globals: {
        ...globals.node
      }
    }
  }
]);
