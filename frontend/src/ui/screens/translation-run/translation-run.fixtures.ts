// 翻訳実行画面の表示用 fixture。
// 画面単位の代表状態（空・入力済み・実行中・完了・エラー）を固定する。
import type {
  TranslationRunForm,
  TranslationProvider,
  BatchProgressView,
  NarrationResultRow,
  RunPhase,
  RunProgress,
  ResultsPaging
} from "./translation-run-view"
import { SUBMIT_NOTICE } from "./translation-run-presentation"

// 画面全体の表示状態。story の固定状態を組むための型。
interface TranslationRunViewModel {
  form: TranslationRunForm
  phase: RunPhase
  canRun: boolean
  models: string[]
  modelsLoading: boolean
  results: NarrationResultRow[]
  errorMessage: string
  // 実行中の進捗。phase==="running" のときだけ意味を持つ。
  progress?: RunProgress
  // 結果一覧のページング表示値。完了状態でだけ意味を持つ。未指定なら単一ページ扱い。
  paging?: ResultsPaging
  // 配送方式。未指定なら同期で従来表示。
  provider?: TranslationProvider
  // 送信の可否。xAI 選択時に意味を持つ。
  canSubmit?: boolean
  // 反映の可否。xAI 選択時に意味を持つ。
  canRefresh?: boolean
  // 送信中フラグ。
  submitting?: boolean
  // 反映中フラグ（旧仕様。Container 互換のため残す）。
  refreshing?: boolean
  // xAI の送信・取り込みの結果として出す案内。空なら出さない。
  notice?: string
  // xAI batch の進行状況（状態確認で取得）。未確認は未指定。
  batchProgress?: BatchProgressView
  // 取り込みの可否。完了段があるとき true。
  canApply?: boolean
  // 状態確認中フラグ。
  checking?: boolean
  // 取り込み中フラグ。
  applying?: boolean
}

const EMPTY_FORM: TranslationRunForm = {
  pluginPath: "",
  endpoint: "",
  apiKey: "",
  model: ""
}

const FILLED_FORM: TranslationRunForm = {
  pluginPath: "/Users/me/Skyrim/Data/Dawnguard.esm",
  endpoint: "https://api.openai.com/v1",
  apiKey: "sk-demo-key",
  model: "gpt-4o-mini"
}

// xAI（batch）選択時の入力済みフォーム。エンドポイントとモデルを xAI 用にする。
const XAI_FORM: TranslationRunForm = {
  pluginPath: "/Users/me/Skyrim/Data/Dawnguard.esm",
  endpoint: "https://api.x.ai",
  apiKey: "xai-demo-key",
  model: "grok-4"
}

// getModels で取得した一覧の表示用サンプル。
const MODELS = ["gpt-4o-mini", "gpt-4o", "gpt-4.1-mini", "o4-mini"]

// GetXAIModels で取得した batch 対応モデルの表示用サンプル（batch 非対応は除外済み）。
const XAI_MODELS = ["grok-4", "grok-4-fast", "grok-3", "grok-3-mini"]

const RESULT_ROWS: NarrationResultRow[] = [
  {
    edid: "DLC1BookSerana",
    source:
      "I have walked these halls for centuries, and still the cold of Castle Volkihar finds me.",
    dest: "私は何世紀もこの広間を歩いてきたが、それでもヴォルキハル城の冷気は私を捉えて離さない。",
    statusLabel: "仮訳",
    terms: [{ source: "Castle Volkihar", dest: "ヴォルキハル城" }]
  },
  {
    edid: "DLC1NoteFromArvak",
    source:
      "Whoever finds this: the Soul Cairn takes more than it gives. Do not bargain with the Ideal Masters.",
    dest: "これを見つけた者へ。ソウル・ケルンは与える以上に奪う。理想の支配者と取引してはならない。",
    statusLabel: "仮訳"
  },
  {
    edid: "DLC1BookAetherium",
    source:
      "The Aetherium Forge lies beneath Blackreach, its fires long since dimmed by the dwarves who fled.",
    dest: "",
    statusLabel: "未訳"
  }
]

// 台詞（line）の結果行。話者を解決できた行はペルソナ口調指示文を持つ。
// Innocence Lost - Quest Expansion の話者で、子供と老女の口調差が観測できることを示す。
const LINE_RESULT_ROWS: NarrationResultRow[] = [
  {
    edid: "AventusAretinoInnocenceLost",
    source: "Mother? Are you there? It's so cold without you.",
    dest: "母さん？ そこにいるの？ あなたがいないと、こんなに寒いんだ。",
    statusLabel: "仮訳",
    personaLabel: "声質: 幼い少年の声",
    directive:
      "この台詞の話者の人物像:\n- 声質: 幼い少年の声\n- 種族の気質: ノルド（実直で粘り強い北方の気質）\nこの人物像に合う口調と人称で訳すこと。"
  },
  {
    edid: "GrelodTheKindScold",
    source:
      "You are all here to serve and honor your matron. Now off to bed, the lot of you!",
    dest: "お前たちは院母に仕え、敬うためにここにいるんだよ。さあ、とっとと寝な、お前たち全員！",
    statusLabel: "仮訳",
    personaLabel: "声質: しわがれた老女の声",
    directive:
      "この台詞の話者の人物像:\n- 声質: しわがれた老女の声\n- 種族の気質: ノルド（実直で粘り強い北方の気質）\nこの人物像に合う口調と人称で訳すこと。"
  },
  {
    edid: "ConstanceMichelGreeting",
    source: "Welcome to Honorhall. I do what I can for the children here.",
    dest: "ようこそ、オネルハル孤児院へ。ここの子供たちのために、できる限りのことはしています。",
    statusLabel: "仮訳",
    personaLabel: "声質: 落ち着いた女性の声",
    directive:
      "この台詞の話者の人物像:\n- 声質: 落ち着いた女性の声\n- 種族の気質: ブレトン（如才ない交渉上手の気質）\nこの人物像に合う口調と人称で訳すこと。"
  },
  {
    edid: "AventusAretinoRitual",
    source: "I knew the Dark Brotherhood would answer my prayer.",
    dest: "闇の一党が僕の祈りに応えてくれるって、分かってたんだ。",
    statusLabel: "仮訳",
    personaLabel: "声質: 幼い少年の声",
    directive:
      "この台詞の話者の人物像:\n- 声質: 幼い少年の声\n- 種族の気質: ノルド（実直で粘り強い北方の気質）\nこの人物像に合う口調と人称で訳すこと。",
    terms: [{ source: "Dark Brotherhood", dest: "闇の一党" }]
  },
  {
    edid: "HonorhallDoorActivate",
    source: "The door is locked.",
    dest: "扉には鍵がかかっている。",
    statusLabel: "仮訳"
  },
  {
    edid: "RiftenGuardGenericHalt",
    source: "Halt! State your business.",
    dest: "止まれ！ 用件を述べろ。",
    statusLabel: "仮訳",
    personaLabel: "口調: 汎用台詞",
    directive:
      "この台詞の話者の人物像:\n- 口調: 衛兵などの不特定多数が話す汎用的な台詞。職務的で簡潔な口調で訳す。\n- 人称と言い回し: 一人称は「俺」。\nこの人物像に合う口調と人称で訳すこと。"
  },
  {
    edid: "DialogueRiftenThievesGuildJoin",
    source: "I'm looking to join the Thieves Guild.",
    dest: "盗賊ギルドに入りたいんだが。",
    statusLabel: "仮訳",
    personaLabel: "口調: PC発話",
    directive:
      "この台詞の話者の人物像:\n- 口調: プレイヤーキャラクターの選択肢。自然な口語で訳す。\n- 人称と言い回し: 一人称は「俺」。\nこの人物像に合う口調と人称で訳すこと。"
  }
]

// 初期表示。何も入力しておらず、モデルも未取得で実行できない。
export const EMPTY_STATE: TranslationRunViewModel = {
  form: EMPTY_FORM,
  phase: "idle",
  canRun: false,
  models: [],
  modelsLoading: false,
  results: [],
  errorMessage: ""
}

// エンドポイントと API キーを入れて、モデルを取得中。
export const MODELS_LOADING_STATE: TranslationRunViewModel = {
  form: { ...FILLED_FORM, model: "" },
  phase: "idle",
  canRun: false,
  models: [],
  modelsLoading: true,
  results: [],
  errorMessage: ""
}

// plugin・接続情報・モデルが揃い、実行できる状態。
export const READY_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "idle",
  canRun: true,
  models: MODELS,
  modelsLoading: false,
  results: [],
  errorMessage: ""
}

// 完了し、原文と訳文が並ぶ。一部は未訳のまま残る場合もある。全件が 1 ページに収まる単一ページ。
export const DONE_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "done",
  canRun: true,
  models: MODELS,
  modelsLoading: false,
  results: RESULT_ROWS,
  errorMessage: "",
  paging: {
    total: RESULT_ROWS.length,
    pageNumber: 1,
    canPrev: false,
    canNext: false
  }
}

// 実行中（台詞を抽出している段階）。件数は出ず、不定バーを出す。
export const RUNNING_EXTRACT_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "running",
  canRun: false,
  models: MODELS,
  modelsLoading: false,
  results: [],
  errorMessage: "",
  progress: { stage: "extract", done: 0, total: 0 }
}

// 実行中（本文を翻訳している段階）。done/total の確定バーを出す。
export const RUNNING_TRANSLATE_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "running",
  canRun: false,
  models: MODELS,
  modelsLoading: false,
  results: [],
  errorMessage: "",
  progress: { stage: "translate", done: 34, total: 121 }
}

// 完了し、台詞の訳文と注入した口調指示文が並ぶ。話者ごとの口調差が観測できる。
// 総件数 121・ページサイズ 50 を想定した先頭ページ（次へ有効）。LINE_RESULT_ROWS は当該ページ分。
export const DONE_PERSONA_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "done",
  canRun: true,
  models: MODELS,
  modelsLoading: false,
  results: LINE_RESULT_ROWS,
  errorMessage: "",
  paging: {
    total: 121,
    pageNumber: 1,
    canPrev: false,
    canNext: true
  }
}

// 実行が失敗し、原因と対応を伝える。
export const ERROR_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "error",
  canRun: true,
  models: MODELS,
  modelsLoading: false,
  results: [],
  errorMessage:
    "翻訳 API への接続に失敗しました。エンドポイントと API キーを確認して、もう一度実行してください。"
}

// xAI（batch）を選び、まだ接続情報もモデルも入れていない。送信できない。
export const XAI_EMPTY_STATE: TranslationRunViewModel = {
  form: { ...EMPTY_FORM },
  phase: "idle",
  canRun: false,
  models: [],
  modelsLoading: false,
  results: [],
  errorMessage: "",
  provider: "xai",
  canSubmit: false,
  notice: ""
}

// xAI（batch）で接続情報とモデルが揃い、まだ状態未確認。主アクションは「送信して開始」で活性。
export const XAI_READY_STATE: TranslationRunViewModel = {
  form: XAI_FORM,
  phase: "idle",
  canRun: false,
  models: XAI_MODELS,
  modelsLoading: false,
  results: [],
  errorMessage: "",
  provider: "xai",
  canSubmit: true,
  notice: ""
}

// xAI（batch）を送信した直後。案内が出て、状態確認で進行状況を見にいく。
export const XAI_SUBMITTED_STATE: TranslationRunViewModel = {
  ...XAI_READY_STATE,
  notice: SUBMIT_NOTICE
}

// xAI（batch）で状態確認中。状態確認ボタンにスピナーが出る。
export const XAI_CHECKING_STATE: TranslationRunViewModel = {
  ...XAI_READY_STATE,
  checking: true
}

// xAI（batch）で固有名段が処理中。ステッパー現在地=固有名、処理待ちが残り主アクションはグレーアウト。
export const XAI_PROPER_PROCESSING_STATE: TranslationRunViewModel = {
  ...XAI_READY_STATE,
  canApply: false,
  batchProgress: { stage: "proper", total: 12, pending: 5, succeeded: 7, failed: 0, canApply: false }
}

// xAI（batch）で固有名段が完了。主アクションは「取り込んで本文を送信」。
export const XAI_PROPER_READY_STATE: TranslationRunViewModel = {
  ...XAI_READY_STATE,
  canApply: true,
  batchProgress: { stage: "proper", total: 12, pending: 0, succeeded: 12, failed: 0, canApply: true }
}

// xAI（batch）で本文段が処理中。ステッパー現在地=本文、処理待ちが残り主アクションはグレーアウト。
export const XAI_BODY_PROCESSING_STATE: TranslationRunViewModel = {
  ...XAI_READY_STATE,
  canApply: false,
  batchProgress: { stage: "body", total: 113, pending: 2, succeeded: 111, failed: 0, canApply: false }
}

// xAI（batch）で本文段が完了。主アクションは「取り込んで完了」。
export const XAI_BODY_READY_STATE: TranslationRunViewModel = {
  ...XAI_READY_STATE,
  canApply: true,
  batchProgress: { stage: "body", total: 113, pending: 0, succeeded: 113, failed: 0, canApply: true }
}

// xAI（batch）で取り込みが終わり、結果が入った。段=完了で、主アクションは再送信の「送信して開始」。
// batch で訳した行も同期と区別なく並ぶ。
export const XAI_DONE_STATE: TranslationRunViewModel = {
  form: XAI_FORM,
  phase: "done",
  canRun: false,
  models: XAI_MODELS,
  modelsLoading: false,
  results: RESULT_ROWS,
  errorMessage: "",
  provider: "xai",
  canSubmit: true,
  canApply: false,
  notice: "",
  batchProgress: { stage: "done", total: 113, pending: 0, succeeded: 113, failed: 0, canApply: false },
  paging: {
    total: RESULT_ROWS.length,
    pageNumber: 1,
    canPrev: false,
    canNext: false
  }
}
