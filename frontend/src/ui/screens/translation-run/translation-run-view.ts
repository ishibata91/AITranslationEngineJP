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

// 本文で辞書から確定訳語へ置換した固有名 1 件（原語 → 確定訳語）。
// 同一原語へ常に同一訳語が当たることを行ごとに確かめるための表示値。
export interface ReplacedTerm {
  // 原語（英語 FULL）。本文中に出ていた固有名。
  source: string
  // 確定訳語（日本語の公式既訳）。置換後に本文へ入った訳。
  dest: string
}

// 基底口調の対人段階。0 尊大 / 1 中立 / 2 丁寧。
export type AttitudeBand = 0 | 1 | 2

// 基底口調の感情段階。0 抑制 / 1 中 / 2 激情。
export type EmotionBand = 0 | 1 | 2

// 決定経路。対人段階をどの信号で決めたか（本文 2 軸 / 声質 prior / 保留）。
export type DecisionPath = "本文" | "voice" | "保留"

// 結果行を展開したときにメタデータとして出す口調情報。台詞の話者の生成済み基底口調（persona_character 由来）。
// 話者を解決できた台詞だけ持つ。叙述文や話者なしの台詞は未設定。
export interface PersonaMeta {
  // 基底口調セル名（対人段階×感情段階。例: 冷然・見下し）。判定結果として強調して出す。
  cell: string
  // 基底口調の性質文（口調をふつうの言葉で説明した一文）。判定結果の補足として出す。
  trait: string
  // 対人段階。
  attitudeBand: AttitudeBand
  // 感情段階。
  emotionBand: EmotionBand
  // 印。対人マーカーを含む台詞数で信頼度の目安。
  marked: number
  // 決定経路。
  decisionPath: DecisionPath
}

// 結果行を展開したときに出す話者情報。誰の台詞かと、口調指示の根拠（性別・年齢・声型）。
// 話者を解決できた台詞だけ持つ。叙述文や話者なしの台詞は未設定。
export interface SpeakerMeta {
  // 話者の EditorID（例: AventusAretino）。誰の台詞かを示す。
  edid: string
  // 性別ラベル（女性 / 男性 / 空）。役割語（一人称・語尾）の根拠。
  sex: string
  // 年齢区分ラベル（老人 / 子供 / 成人 / 空）。役割語の根拠（種族 EditorID 由来）。
  age: string
  // 声型 EditorID（例: FemaleOldGrumpy）。印不足時の対人 prior の根拠。
  voice: string
}

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
  // 行を展開したときに出す話者情報（誰の台詞か＋属性）。話者を解決できた台詞だけ持つ。
  speaker?: SpeakerMeta
  // 注入したペルソナ口調指示文の全文。話者を解決できた台詞だけ入る。叙述文や話者なしの台詞は空または未設定。
  // 既定では畳み、行を展開したときに見せる。
  directive?: string
  // 口調チップ用の短い要約。基底口調セル名。一覧のまま口調差を観測するための表示。
  // directive を持つ行だけ持つ。
  personaLabel?: string
  // 行を展開したときにメタデータとして出す口調情報（基底口調・対人/感情段階・印・決定経路）。
  // 話者を解決できた台詞だけ持つ。
  persona?: PersonaMeta
  // この本文で辞書から置換した固有名（原語 → 確定訳語）。置換が無い行は空または未設定。
  // 畳んだ行には件数チップ、展開で一覧を出し、同一原語への同一訳語を行ごとに確かめる。
  terms?: ReplacedTerm[]
  // この行の翻訳で実際に翻訳 AI へ投げた完全プロンプト（base 指示 ＋ 口調指示 ＋ 置換済み原文）。
  // 取得時に再構成した文字列。口調指示が実プロンプトへ合成されたかを目視で確かめるため、展開時に全文を出す。
  // 口調指示（directive）は実プロンプトの一部の抜粋で、こちらは送信した全文を含む。再構成できない行は未設定。
  prompt?: string
}

// 結果一覧の keyset ページャの表示値。results はページ分のみを持つため、総件数と現在ページ番号は別に渡す。
export interface ResultsPaging {
  // 結果全体の総件数（叙述文＋台詞）。ページ送りに依らない合計。件数バッジに出す。
  total: number
  // 現在ページ番号（1 始まり）。
  pageNumber: number
  // 前ページがあるか。先頭ページでは false。
  canPrev: boolean
  // 次ページがあるか。末尾ページでは false。
  canNext: boolean
}

// テキスト入力欄の表示メタ。
export interface FormFieldDescriptor {
  field: "endpoint" | "apiKey"
  label: string
  hint: string
  placeholder: string
  secret: boolean
}
