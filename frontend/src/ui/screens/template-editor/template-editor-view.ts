// テンプレート編集画面の表示用の型。値（表示定数・純関数）が要る場合は別ファイルに置く。
// state・保存・validation は持たず、表示と入力 event の中継だけを型で表す。

// プロンプトテンプレート画面のサブタブ。base のプロンプトと、レコード別の変数（口調・文体）を分ける。
export type TemplateTab = "base" | "record"

// PC（プレイヤーキャラクター）の性別。一人称・語尾の根拠。空は未指定（一人称・語尾を付けない）。
export type PcSex = "" | "Male" | "Female"

// 編集対象のプロンプトテンプレート。翻訳 AI へ送る指示文の雛形と、話者なし台詞の口調設定。
export interface PromptTemplateForm {
  // base 翻訳指示文。叙述文・台詞のどちらにも付く冒頭の system 指示。ベースタブで編集する。
  baseDirective: string
  // 口調指示テンプレート。話者のいる台詞だけに付く。話者の性質列を差し込み口で展開する。レコード別タブで編集する。
  personaTemplate: string
  // 汎用台詞（話者なし）の口調指示文。自由記述。感情と性別はシステムが自動で重ねる。レコード別タブで編集する。
  // 配線前（Container 未更新）でも壊れないよう任意にする。
  genericToneText?: string
  // PC 発話（選択肢）の口調指示文。自由記述。感情と PC 性別はシステムが自動で重ねる。レコード別タブで編集する。
  pcToneText?: string
  // PC の性別。一人称・語尾の根拠。レコード別タブで選ぶ。
  pcSex?: PcSex
}

export type PromptTemplateField = keyof PromptTemplateForm
