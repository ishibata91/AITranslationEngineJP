// 翻訳実行画面の表示用の型。値（表示定数・純関数）は translation-run-presentation.ts に置く。
// 型と値を分けるのは、import 規約（値と型を同一 import で混在させない）に沿うため。

// 実行の進行段階。表示の出し分けに使う。
export type RunPhase = "idle" | "running" | "done" | "error"

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

// 結果一覧の 1 行。叙述文（narration）の原文・訳文・状態を表示用に整形した値。
export interface NarrationResultRow {
  // レコードの EditorID。どの書物かを利用者が識別するための表示。
  edid: string
  // 原文（抽出した叙述文）。
  source: string
  // 訳文。未訳のときは空文字。
  dest: string
  // 訳状態の表示ラベル（未訳・仮訳・訳済・承認）。status コードからの変換は Presenter が行い、ここには表示文字列だけ入る。
  statusLabel: string
}

// テキスト入力欄の表示メタ。
export interface FormFieldDescriptor {
  field: "endpoint" | "apiKey"
  label: string
  hint: string
  placeholder: string
  secret: boolean
}
