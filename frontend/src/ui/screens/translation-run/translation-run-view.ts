// 翻訳実行画面の表示用の型。値（表示定数・純関数）は translation-run-presentation.ts に置く。
// 型と値を分けるのは、import 規約（値と型を同一 import で混在させない）に沿うため。

// 実行の進行段階。表示の出し分けに使う。
export type RunPhase = "idle" | "running" | "done" | "error"

// 実行中の処理段階。抽出は件数が出ないため不定、本文翻訳は done/total で進捗を出す。
export type ProgressStage = "extract" | "translate"

// 実行中の進捗。表示専用。phase==="running" のときだけ意味を持つ。
export interface RunProgress {
  // どの段階か。extract=台詞抽出（不定）、translate=本文翻訳（done/total）。
  stage: ProgressStage
  // 処理済み件数。translate で意味を持つ。extract では 0。
  done: number
  // 総件数。translate で意味を持つ。extract では 0。
  total: number
}

// 入力フォームの表示値。
export interface TranslationRunForm {
  // 翻訳対象 plugin のフルパス。ファイル選択ダイアログで設定する。表示専用。
  // Data フォルダはこのパスの親ディレクトリとして派生するので、別入力は持たない。
  pluginPath: string
  // OpenAI 互換 API のエンドポイント URL。
  endpoint: string
  // API キー。表示はマスクする。
  apiKey: string
  // 利用するモデル名。getModels で取得した一覧から選ぶ。
  model: string
}

export type TranslationRunFormField = keyof TranslationRunForm

// 結果一覧の 1 行。叙述文（narration）と台詞（line）の原文・訳文・状態を表示用に整形した値。
export interface NarrationResultRow {
  // レコードの EditorID。どの書物・台詞かを利用者が識別するための表示。
  edid: string
  // 原文（抽出した叙述文または台詞）。
  source: string
  // 訳文。未訳のときは空文字。
  dest: string
  // 訳状態の表示ラベル（未訳・仮訳・訳済・承認）。status コードからの変換は Presenter が行い、ここには表示文字列だけ入る。
  statusLabel: string
  // 注入したペルソナ口調指示文の全文。話者を解決できた台詞だけ入る。叙述文や話者なしの台詞は空または未設定。
  // 既定では畳み、行を展開したときに見せる。
  directive?: string
  // 口調チップ用の短い要約。話者の最も効く特徴 1 つ（声質など）。一覧のまま口調差を観測するための表示。
  // directive を持つ行だけ持つ。
  personaLabel?: string
}

// テキスト入力欄の表示メタ。
export interface FormFieldDescriptor {
  field: "endpoint" | "apiKey"
  label: string
  hint: string
  placeholder: string
  secret: boolean
}
