import "@testing-library/jest-dom/vitest"

import type { MasterDictionaryGatewayContract } from "@application/gateway-contract/master-dictionary"
import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
import { createMasterDictionaryScreenControllerFactory } from "@controller/master-dictionary"

export function createTestMasterDictionaryScreenControllerFactory(
  gateway: MasterDictionaryGatewayContract | null
): CreateMasterDictionaryScreenController {
  return createMasterDictionaryScreenControllerFactory(gateway)
}
