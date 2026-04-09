// @vitest-environment node

import { RuleTester } from "eslint";
import { describe, it } from "vitest";
import tseslint from "typescript-eslint";
import { repositoryBoundaryPlugin } from "./repository-boundary-plugin.mjs";

const rule = repositoryBoundaryPlugin.rules["enforce-layer-boundaries"];

const ruleTester = new RuleTester({
  languageOptions: {
    ecmaVersion: 2022,
    sourceType: "module",
    parser: tseslint.parser
  }
});

describe("repository boundary plugin", () => {
  it("enforces current frontend layer boundaries", () => {
    ruleTester.run("enforce-layer-boundaries", rule, {
      valid: [
        {
          filename: "F:/AITranslationEngineJp/frontend/src/ui/App.svelte",
          code: "import { createShellState } from '@ui/stores/shell-state';"
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/ui/screens/bootstrap/BootstrapScreen.svelte",
          code: "import type { BootstrapGatewayContract } from '@application/gateway-contract';"
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/controller/wails/bootstrap.gateway.ts",
          code: "import type { BootstrapGatewayContract } from '@application/gateway-contract';"
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/controller/wails/bootstrap.gateway.ts",
          code: "import { AppController } from '../../wailsjs/go/wails/AppController';"
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.test.ts",
          code: "import { fixture } from '../../fixtures/bootstrap';"
        }
      ],
      invalid: [
        {
          filename: "F:/AITranslationEngineJp/frontend/src/ui/App.svelte",
          code: "import { invokeBootstrap } from '@controller/wails/bootstrap.gateway';",
          errors: [{ message: "ui code must not import controller code directly." }]
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.ts",
          code: "import { invokeBootstrap } from '@controller/wails/bootstrap.gateway';",
          errors: [{ message: "application code must not import controller code directly." }]
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.ts",
          code: "import { AppController } from '../../wailsjs/go/wails/AppController';",
          errors: [
            {
              message:
                "application code must not import Wails bindings directly. Go through gateway ports or gateway adapters instead."
            }
          ]
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/controller/wails/bootstrap.gateway.ts",
          code: "import AppShell from '@ui/views/AppShell.svelte';",
          errors: [{ message: "controller code must not import ui code directly." }]
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.ts",
          code: "import { fixture } from '../../fixtures/bootstrap';",
          errors: [
            {
              message:
                "application production code must not import test, fixture, or generated support files."
            }
          ]
        }
      ]
    });
  });

  it("keeps same-layer imports on public entrypoints", () => {
    ruleTester.run("enforce-layer-boundaries", rule, {
      valid: [
        {
          filename: "F:/AITranslationEngineJp/frontend/src/ui/App.svelte",
          code: "import AppShell from '@ui/views/AppShell.svelte';"
        }
      ],
      invalid: [
        {
          filename: "F:/AITranslationEngineJp/frontend/src/ui/App.svelte",
          code: "import { buildScreenState } from '@ui/screens/bootstrap/internal/build-screen-state';",
          errors: [{ messageId: "forbiddenImport" }]
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/controller/wails/bootstrap.gateway.ts",
          code: "import type { BootstrapGatewayDto } from '@controller/wails/gateway-dto';",
          errors: [{ message: "controller code must not import other controller roots directly." }]
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/controller/wails/bootstrap.gateway.ts",
          code: "import { createGatewayDto } from '@controller/wails/gateway-dto/internal/create-gateway-dto';",
          errors: [{ message: "controller code must not import other controller roots directly." }]
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.ts",
          code: "import type { AnotherContract } from '@application/other-contract';",
          errors: [{ message: "application code must not import other application roots directly." }]
        },
        {
          filename: "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.ts",
          code: "import { buildContract } from '@application/other-contract/internal/build-contract';",
          errors: [{ messageId: "forbiddenImport" }]
        }
      ]
    });
  });
});
