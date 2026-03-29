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
  it("enforces frontend layer boundaries", () => {
    ruleTester.run("enforce-layer-boundaries", rule, {
      valid: [
        {
          filename: "F:/AITranslationEngineJp/src/ui/screens/bootstrap/view-model.ts",
          code: "import { loadBootstrapStatus } from '@application/bootstrap/load-bootstrap-status';"
        },
        {
          filename: "F:/AITranslationEngineJp/src/gateway/tauri/bootstrap-status.gateway.ts",
          code: "import { loadBootstrapStatus } from '@application/bootstrap/load-bootstrap-status';"
        },
        {
          filename: "F:/AITranslationEngineJp/src/shared/contracts/bootstrap-status.ts",
          code: "export type BootstrapStatus = { ready: boolean };"
        },
        {
          filename: "F:/AITranslationEngineJp/src/ui/screens/bootstrap/view-model.ts",
          code: "import type { BootstrapStatus } from '@shared/contracts/bootstrap-status';"
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/bootstrap/load-bootstrap-status.test.ts",
          code: "import { fixture } from '../../fixtures/bootstrap-status';"
        },
        {
          filename: "F:/AITranslationEngineJp/src/gateway/tauri/bootstrap-status.gateway.ts",
          code: "import { invoke } from '@tauri-apps/api/core';"
        }
      ],
      invalid: [
        {
          filename: "F:/AITranslationEngineJp/src/ui/screens/bootstrap/BootstrapScreen.svelte",
          code: "import { getBootstrapStatusGateway } from '@gateway/tauri/bootstrap-status.gateway';",
          errors: [{ message: "ui code must not import gateway code directly." }]
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/bootstrap/load-bootstrap-status.ts",
          code: "import { invoke } from '@tauri-apps/api/core';",
          errors: [
            {
              message:
                "application code must not import Tauri APIs directly. Go through gateway ports or gateway adapters instead."
            }
          ]
        },
        {
          filename: "F:/AITranslationEngineJp/src/shared/contracts/bootstrap-status.ts",
          code: "import { loadBootstrapStatus } from '@application/bootstrap/load-bootstrap-status';",
          errors: [{ message: "shared code must not import application code directly." }]
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/bootstrap/load-bootstrap-status.ts",
          code: "import { fixture } from '../../fixtures/bootstrap-status';",
          errors: [
            {
              message:
                "application production code must not import test, fixture, or generated support files."
            }
          ]
        },
        {
          filename: "F:/AITranslationEngineJp/src/gateway/tauri/bootstrap-status.gateway.ts",
          code: "import type { BootstrapStatus } from '@ui/screens/bootstrap/bootstrap-status';",
          errors: [{ message: "gateway code must not import ui code directly." }]
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/bootstrap/load-bootstrap-status.ts",
          code: "import { testData } from './load-bootstrap-status.test';",
          errors: [
            {
              message:
                "application production code must not import test, fixture, or generated support files."
            }
          ]
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/bootstrap/load-bootstrap-status.ts",
          code: "import { setup } from '../../test/setup';",
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

  it("given same-layer cross-root imports when target visibility differs then only public entrypoints stay importable", () => {
    ruleTester.run("enforce-layer-boundaries", rule, {
      valid: [
        {
          filename: "F:/AITranslationEngineJp/src/ui/app-shell/AppShell.svelte",
          code: "import BootstrapScreen from '@ui/screens/bootstrap/BootstrapScreen.svelte';"
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/bootstrap/load-bootstrap-status.ts",
          code: "import type { BootstrapInputPort } from '@application/ports/input';"
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/usecases/start-import.ts",
          code: "import { loadBootstrapStatus } from '@application/bootstrap/load-bootstrap-status';"
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/usecases/start-import.ts",
          code: "import type { BootstrapInputPort } from '@application/ports/input/index';"
        },
        {
          filename: "F:/AITranslationEngineJp/src/ui/screens/bootstrap/BootstrapScreen.svelte",
          code: "import { buildScreenState } from '@ui/screens/bootstrap/internal/build-screen-state';"
        }
      ],
      invalid: [
        {
          filename: "F:/AITranslationEngineJp/src/ui/app-shell/AppShell.svelte",
          code: "import { buildScreenState } from '@ui/screens/bootstrap/internal/build-screen-state';",
          errors: [{ messageId: "forbiddenImport" }]
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/usecases/start-import.ts",
          code: "import { createGatewayRequest } from '@application/ports/gateway/internal/create-gateway-request';",
          errors: [{ messageId: "forbiddenImport" }]
        },
        {
          filename: "F:/AITranslationEngineJp/src/ui/app-shell/AppShell.svelte",
          code: "import { triggerBootstrap } from '@ui/screens/bootstrap/internal/commands/trigger-bootstrap';",
          errors: [{ messageId: "forbiddenImport" }]
        },
        {
          filename: "F:/AITranslationEngineJp/src/application/usecases/start-import.ts",
          code: "import { invokeTauri } from '@gateway/tauri/invoke/foo';",
          errors: [{ message: "application code must not import gateway code directly." }]
        }
      ]
    });
  });
});
