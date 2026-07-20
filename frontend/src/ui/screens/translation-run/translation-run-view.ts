// 翻訳実行画面の表示用の型。値（表示定数・純関数）は translation-run-presentation.ts に置く。
// 型と値を分けるのは、import 規約（値と型を同一 import で混在させない）に沿うため。

// 実行の進行段階。表示の出し分けに使う。
export type RunPhase = "idle" | "running" | "done" | "error"

// 実行中の処理段階。抽出は件数が出ないため不定、本文翻訳は done/total で進捗を出す。
export type ProgressStage = "extract" | "translate"

// extract 段（翻訳前区間）の内訳サブ段。translate 段では持たない。
// 値は backend の処理段に対応する。extract=C# 抽出子、derive=固有名派生（DeriveMasterTerms）、
// reference=既存訳取込（LoadReferenceTranslations）、ingest=取込段（Ingest）。
// extract 段は件数が出ないため、進んでいるかを示すためにサブ段名を出す。
export type ExtractStep = "extract" | "derive" | "reference" | "ingest"

// 実行中の進捗。表示専用。phase==="running" のときだけ意味を持つ。
export interface RunProgress {
  // どの段階か。extract=台詞抽出（不定）、translate=本文翻訳（done/total）。
  stage: ProgressStage
  // 処理済み件数。translate で意味を持つ。extract では 0。
  done: number
  // 総件数。translate で意味を持つ。extract では 0。
  total: number
  // extract 段の内訳サブ段。指定時はサブ段ラベルを見出しに出す。
  // 未指定または translate 段では従来の段ラベルを出す。
  step?: ExtractStep
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
  // トークン予算。1 リクエストにまとめる原文の最大トークン数の目安。表示は文字列で持ち、数値変換はここでは行わない。
  tokenBudget: string
}

export type TranslationRunFormField = keyof TranslationRunForm

// 結果行の元レコード種別（箱と REC:FIELD）。会話・本以外の種別が増えたため、どの種別の結果かを一覧で示す。
export interface ResultRecordType {
  // 概念モデルの箱（叙述文・固有名・定型句・台詞）。
  box: string
  // REC:FIELD（例: WEAP:DESC）。どの種別から抽出した行かを示す。
  recField: string
}

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

// 決定経路。対人段階をどの信号で決めたか。
// 本文 2 軸 / 声質 prior / 保留 は話者を解決できた台詞の集計経路。
// 汎用 / PC は話者を解決できない台詞へ、利用者が指定した既定の対人段階を固定で当てた経路。
export type DecisionPath = "本文" | "voice" | "保留" | "汎用" | "PC"

// 結果行を展開したときにメタデータとして出す口調情報。
// 名指し話者は基底口調（persona_character 由来。対人段階×感情段階のセル）を持つ。
// 汎用台詞・PC 発話は対人段階・セル・印を持たず、決定経路・感情段階・性別だけを持つ。
// 叙述文や口調が付かない行は未設定。
export interface PersonaMeta {
  // 基底口調セル名（対人段階×感情段階。例: 冷然・見下し）。名指し話者だけ持つ。汎用・PC は未設定。
  cell?: string
  // 基底口調の性質文（口調をふつうの言葉で説明した一文）。名指し話者だけ持つ。汎用・PC は未設定。
  trait?: string
  // 対人段階。名指し話者だけ持つ。汎用・PC は対人を自由記述で書くため未設定。
  attitudeBand?: AttitudeBand
  // 感情段階。台詞 1 行から決まり、名指し・汎用・PC すべてが持つ。
  emotionBand: EmotionBand
  // 印。対人マーカーを含む台詞数で信頼度の目安。話者集計経路（本文 / voice / 保留）だけ持つ。
  marked?: number
  // 決定経路。
  decisionPath: DecisionPath
  // 性別ラベル（女性 / 男性）。口調の一人称・語尾の根拠。
  // 汎用台詞は INFO の条件由来、PC 発話は利用者選択の PC 性別。性別を定められない台詞では未設定。
  sex?: string
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
  // 元レコード種別（箱と REC:FIELD）。会話・本以外の種別が増えたため、一覧でどの種別の結果かを示す。
  // 既存の本文（叙述文・台詞）にも順次付く。未設定の行は種別バッジを出さない。
  recordType?: ResultRecordType
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
  field: "endpoint" | "apiKey" | "tokenBudget"
  label: string
  hint: string
  placeholder: string
  secret: boolean
  // 入力の右に出す静的な単位表示（例: "k"）。無い欄では省く。
  suffix?: string
}
