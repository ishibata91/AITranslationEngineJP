// テンプレート編集画面の表示用の型。値（表示定数・純関数）が要る場合は別ファイルに置く。
// state・保存・validation は持たず、表示と入力 event の中継だけを型で表す。

// 編集対象のプロンプトテンプレート。翻訳 AI へ送る指示文の雛形。
export interface PromptTemplateForm {
  // base 翻訳指示文。叙述文・台詞のどちらにも付く冒頭の system 指示。
  baseDirective: string
  // 口調指示テンプレート。話者のいる台詞だけに付く。話者の性質列を差し込み口で展開する。
  personaTemplate: string
}

export type PromptTemplateField = keyof PromptTemplateForm

// 口調指示テンプレートに差し込めるプレースホルダ。編集者がキーと意味を見て使う。
export interface TemplatePlaceholder {
  // 差し込み口のキー（例: {traits}）。
  token: string
  // 差し込まれる内容の説明。
  description: string
}
