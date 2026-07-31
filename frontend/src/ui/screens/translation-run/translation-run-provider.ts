import type { TranslationProvider } from "./translation-run-view"

const SYNC_DEFAULT_ENDPOINT = "http://127.0.0.1:1234"
const OPENAI_DEFAULT_ENDPOINT = "https://api.openai.com/v1"
const XAI_DEFAULT_ENDPOINT = "https://api.x.ai"
const OPENAI_DEFAULT_MODEL = "gpt-5.6-luna"

interface ProviderDefaults {
  endpoint: string
  models: string[]
  model: string
}

// providerDefaults は提供元の切替時に、前の接続先、モデル一覧、進行表示を持ち越さないための初期値を返す。
export function providerDefaults(
  provider: TranslationProvider
): ProviderDefaults {
  switch (provider) {
    case "openai":
      return {
        endpoint: OPENAI_DEFAULT_ENDPOINT,
        models: [OPENAI_DEFAULT_MODEL],
        model: OPENAI_DEFAULT_MODEL
      }
    case "xai":
      return { endpoint: XAI_DEFAULT_ENDPOINT, models: [], model: "" }
    default:
      return { endpoint: SYNC_DEFAULT_ENDPOINT, models: [], model: "" }
  }
}

// orderProviderModels は OpenAI の取得結果に Luna があれば先頭へ置き、他モデルの順序と内容は残す。
export function orderProviderModels(
  provider: TranslationProvider,
  fetched: string[]
): string[] {
  if (provider !== "openai" || !fetched.includes(OPENAI_DEFAULT_MODEL)) {
    return [...fetched]
  }
  return [
    OPENAI_DEFAULT_MODEL,
    ...fetched.filter((item) => item !== OPENAI_DEFAULT_MODEL)
  ]
}

// canOperateBatch は OpenAI だけ API キーを batch 全操作の必須条件にする。
export function canOperateBatch(
  provider: TranslationProvider,
  apiKey: string
): boolean {
  return (
    provider !== "sync" && (provider !== "openai" || apiKey.trim().length > 0)
  )
}
