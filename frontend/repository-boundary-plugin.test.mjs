// @vitest-environment node

import { RuleTester } from "eslint"
import tseslint from "typescript-eslint"
import { repositoryBoundaryPlugin } from "../scripts/eslint/repository-boundary-plugin.mjs"

const rule = repositoryBoundaryPlugin.rules["enforce-layer-boundaries"]

const ruleTester = new RuleTester({
  languageOptions: {
    ecmaVersion: 2022,
    sourceType: "module",
    parser: tseslint.parser
  }
})

ruleTester.run("enforce-layer-boundaries", rule, {
  valid: [
    {
      filename: "F:/AITranslationEngineJp/frontend/src/ui/App.svelte",
      code: "import { createShellState } from '@ui/stores/shell-state'"
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/ui/screens/bootstrap/BootstrapScreen.svelte",
      code: "import type { BootstrapGatewayContract } from '@application/gateway-contract'"
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/controller/wails/bootstrap.gateway.ts",
      code: "import type { BootstrapGatewayContract } from '@application/gateway-contract'"
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/controller/wails/bootstrap.gateway.ts",
      code: "import { AppController } from '../../wailsjs/go/wails/AppController'"
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.test.ts",
      code: "import { fixture } from '../../fixtures/bootstrap'"
    },
    {
      filename: "F:/AITranslationEngineJp/frontend/src/ui/App.test.ts",
      code: "import { createTestMasterDictionaryScreenControllerFactory } from '../test/setup'"
    },
    {
      filename: "F:/AITranslationEngineJp/frontend/src/ui/App.svelte",
      code: "import AppShell from '@ui/views/AppShell.svelte'"
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte",
      code: "import type { CreateMasterDictionaryScreenController } from '@application/contract/master-dictionary'"
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/application/usecase/master-dictionary/master-dictionary.usecase.ts",
      code: "import type { MasterDictionaryGatewayContract } from '@application/gateway-contract/master-dictionary'"
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/application/usecase/master-dictionary/master-dictionary.usecase.ts",
      code: "import { MasterDictionaryStore } from '@application/store/master-dictionary'"
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/application/store/master-dictionary/master-dictionary.store.ts",
      code: "import { DEFAULT_CATEGORY } from '@application/contract/master-dictionary'"
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.ts",
      code: "import { MasterDictionaryRuntimeEventAdapter } from '@controller/runtime/master-dictionary'"
    }
  ],
  invalid: [
    {
      filename: "F:/AITranslationEngineJp/frontend/src/ui/App.svelte",
      code: "import { invokeBootstrap } from '@controller/wails/bootstrap.gateway'",
      errors: [{ message: "ui code must not import controller code directly." }]
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.ts",
      code: "import { invokeBootstrap } from '@controller/wails/bootstrap.gateway'",
      errors: [
        {
          message: "application code must not import controller code directly."
        }
      ]
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.ts",
      code: "import { AppController } from '../../wailsjs/go/wails/AppController'",
      errors: [
        {
          message:
            "application code must not import Wails bindings directly. Go through gateway ports or gateway adapters instead."
        }
      ]
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/controller/wails/bootstrap.gateway.ts",
      code: "import AppShell from '@ui/views/AppShell.svelte'",
      errors: [{ message: "controller code must not import ui code directly." }]
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/application/gateway-contract/index.ts",
      code: "import { fixture } from '../../fixtures/bootstrap'",
      errors: [
        {
          message:
            "application production code must not import test, fixture, or generated support files."
        }
      ]
    },
    {
      filename: "F:/AITranslationEngineJp/frontend/src/ui/App.svelte",
      code: "import { buildScreenState } from '@ui/screens/bootstrap/internal/build-screen-state'",
      errors: [{ messageId: "forbiddenImport" }]
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/application/usecase/master-dictionary/master-dictionary.usecase.ts",
      code: "import { createMasterDictionaryScreenControllerFactory } from '@controller/master-dictionary'",
      errors: [
        {
          message: "application code must not import controller code directly."
        }
      ]
    },
    {
      filename: "F:/AITranslationEngineJp/frontend/src/ui/App.test.ts",
      code: "import { createMasterDictionaryScreenControllerFactory } from '@controller/master-dictionary'",
      errors: [{ message: "ui code must not import controller code directly." }]
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/application/usecase/master-dictionary/master-dictionary.usecase.ts",
      code: "import { OtherFeaturePresenter } from '@application/presenter/other-feature'",
      errors: [
        {
          message:
            "application code must not import other application roots directly."
        }
      ]
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/controller/master-dictionary/master-dictionary-screen-controller.ts",
      code: "import App from '@ui/App.svelte'",
      errors: [{ message: "controller code must not import ui code directly." }]
    },
    {
      filename:
        "F:/AITranslationEngineJp/frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.ts",
      code: "import { helper } from '@controller/wails/internal/helper'",
      errors: [
        {
          message:
            "controller code must not import other controller roots directly."
        }
      ]
    }
  ]
})
