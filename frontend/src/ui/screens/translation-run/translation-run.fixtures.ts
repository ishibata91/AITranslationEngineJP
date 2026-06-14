// 翻訳実行画面の表示用 fixture。
// 画面単位の代表状態（空・入力済み・実行中・完了・エラー）を固定する。
import type {
  TranslationRunForm,
  NarrationResultRow,
  RunPhase,
  RunProgress
} from "./translation-run-view"

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

// getModels で取得した一覧の表示用サンプル。
const MODELS = ["gpt-4o-mini", "gpt-4o", "gpt-4.1-mini", "o4-mini"]

const RESULT_ROWS: NarrationResultRow[] = [
  {
    edid: "DLC1BookSerana",
    source:
      "I have walked these halls for centuries, and still the cold of Castle Volkihar finds me.",
    dest: "私は何世紀もこの広間を歩いてきたが、それでもヴォルキハル城の冷気は私を捉えて離さない。",
    statusLabel: "仮訳"
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
      "この台詞の話者の人物像:\n- 声質: 幼い少年の声\n- 種族の気質: ノルド（実直で粘り強い北方の気質）\nこの人物像に合う口調と人称で訳すこと。"
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
    dest: "",
    statusLabel: "未訳"
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

// 完了し、原文と訳文が並ぶ。一部は未訳のまま残る場合もある。
export const DONE_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "done",
  canRun: true,
  models: MODELS,
  modelsLoading: false,
  results: RESULT_ROWS,
  errorMessage: ""
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
export const DONE_PERSONA_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "done",
  canRun: true,
  models: MODELS,
  modelsLoading: false,
  results: LINE_RESULT_ROWS,
  errorMessage: ""
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
