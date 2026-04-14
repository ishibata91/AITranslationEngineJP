import type {
  MasterDictionaryFrontendRefresh,
  MasterDictionaryUpsertPayload
} from "@application/gateway-contract/master-dictionary"

import type { MasterDictionaryScreenState } from "./master-dictionary-screen-types"

export const PAGE_SIZE = 30
export const DEFAULT_CATEGORY = "固有名詞"
export const DEFAULT_ORIGIN = "手動登録"

export function buildRefreshPayload(
  query: string,
  category: string,
  page: number
): MasterDictionaryFrontendRefresh {
  return {
    query,
    category: category === "すべて" ? "" : category,
    page,
    pageSize: PAGE_SIZE
  }
}

export function buildUpsertPayload(
  state: MasterDictionaryScreenState
): MasterDictionaryUpsertPayload {
  return {
    source: state.formSource.trim(),
    translation: state.formTranslation.trim(),
    category: state.formCategory,
    origin: state.formOrigin
  }
}
