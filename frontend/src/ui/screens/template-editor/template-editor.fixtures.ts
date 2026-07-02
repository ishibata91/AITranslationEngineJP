// プロンプトテンプレート画面の表示用 fixture。Base 指示と、種別ごとの指示文（directive）を表示値として固定する。
// プロンプト = Base 指示 + その REC:FIELD に割り当てた指示文（変数を実行時に埋めたもの）。
import type { PromptTemplateForm } from "./template-editor-view"
import type { Directive, RecordAssignment } from "./directive-view"

// base 翻訳指示文。現状は provider のハードコード（openai_compatible.go:152 相当）。
const DEFAULT_BASE = `あなたは Skyrim Mod の翻訳者です。与えられた英語の本文を、原文の意味と語調を保った自然な日本語へ翻訳してください。訳文だけを出力し、説明や注釈は加えないでください。`

// 口調指示。現状は engine のハードコード（persona.go の buildPersonaDirective 相当）。
// {traits} に話者の性質（声質・種族の気質・所属の気風）の箇条書きを差し込む。
const DEFAULT_PERSONA = `この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。`

// 話者なし台詞の口調の既定文。汎用台詞・PC 発話それぞれの自由記述の口調指示。
const DEFAULT_GENERIC_TONE = `衛兵などの不特定多数が話す汎用的な台詞。職務的で簡潔な口調で訳す。`
const DEFAULT_PC_TONE = `プレイヤーキャラクターの選択肢。自然な口語で訳す。`

export const DEFAULT_TEMPLATE_FORM: PromptTemplateForm = {
  baseDirective: DEFAULT_BASE,
  personaTemplate: DEFAULT_PERSONA,
  // 話者なし台詞の口調。汎用・PC それぞれの自由記述と、PC の性別（既定は未指定）。
  genericToneText: DEFAULT_GENERIC_TONE,
  pcToneText: DEFAULT_PC_TONE,
  pcSex: ""
}

// 口調の自由記述を書き換え、PC 性別を男性に選んだ form。未保存＋PC 性別選択の表示確認に使う。
export const TONE_DEFAULT_EDITED_FORM: PromptTemplateForm = {
  ...DEFAULT_TEMPLATE_FORM,
  genericToneText: "衛兵の汎用台詞。威厳を込めた、やや突き放した口調で訳す。",
  pcSex: "Male"
}

// 種別ごとの指示文。口調・文体・固有名・定型句をすべて同じ形で持つ。口調だけ {traits} 変数を持つ。
export const DIRECTIVES: Directive[] = [
  {
    key: "説明体",
    instruction:
      "対象の機能・効果を簡潔に説明する文体で訳すこと。装飾を避け、ゲーム内の説明として読みやすくすること。",
    variables: []
  },
  {
    key: "書物体",
    instruction: "作中の書物として、地の文の語り口と文語調で訳すこと。",
    variables: []
  },
  {
    key: "日記体",
    instruction: "クエストの記録として、一人称の手記の語り口で訳すこと。",
    variables: []
  },
  {
    key: "世界観断片",
    instruction:
      "ロード画面に出る世界観の断章として、短く余韻のある語り口で訳すこと。",
    variables: []
  },
  {
    key: "口調",
    instruction: DEFAULT_PERSONA,
    variables: [
      {
        token: "{traits}",
        description:
          "話者の性質（声質・種族の気質・所属の気風）を箇条書きで差し込む。話者の属性（種族・声型）から定まる性質文が入る。"
      }
    ]
  },
  {
    key: "固有名",
    instruction:
      "原語を Skyrim の固有名として、カタカナ表記の名前へ訳すこと。一般名詞に訳し下げないこと。",
    variables: []
  },
  {
    key: "定型句",
    instruction:
      "ゲーム内に繰り返し出る短い定型句として、簡潔で一貫した訳にすること。",
    variables: []
  }
]

// REC:FIELD と指示文キーの対応（固定）。
function assign(recField: string, logicalName: string, directive: string): RecordAssignment {
  return { recField, logicalName, directive }
}

export const RECORD_ASSIGNMENTS: RecordAssignment[] = [
  assign("WEAP:DESC", "武器の説明", "説明体"),
  assign("ARMO:DESC", "防具の説明", "説明体"),
  assign("SPEL:DESC", "呪文の説明", "説明体"),
  assign("MGEF:DNAM", "魔法効果の説明", "説明体"),
  assign("PERK:DESC", "Perk の説明", "説明体"),
  assign("BOOK:DESC", "本文（書物の中身）", "書物体"),
  assign("QUST:CNAM", "クエストの記録（ログ）", "日記体"),
  assign("QUST:NNAM", "クエストの目標", "日記体"),
  assign("LSCR:DESC", "ロード画面の文", "世界観断片"),
  assign("INFO:NAM1", "台詞", "口調"),
  assign("INFO:RNAM", "選択肢の上書き文", "口調"),
  assign("DIAL:FULL", "選択肢の既定文", "口調"),
  assign("WEAP:FULL", "武器名", "固有名"),
  assign("ARMO:FULL", "防具名", "固有名"),
  assign("SPEL:FULL", "呪文名", "固有名"),
  assign("QUST:FULL", "クエスト名", "固有名"),
  assign("CELL:FULL", "地名（場所名）", "固有名"),
  assign("BOOK:FULL", "本のタイトル", "固有名"),
  assign("NPC_:FULL", "人物名", "固有名"),
  assign("NPC_:SHRT", "人物の短縮名", "固有名"),
  assign("RACE:FULL", "種族名", "固有名"),
  assign("FACT:FULL", "勢力名", "固有名"),
  assign("ACTI:RNAM", "起動の動作名（「開く」等）", "定型句"),
  assign("MESG:ITXT", "メッセージのボタン文言", "定型句"),
  assign("WOOP:TNAM", "龍語の語義", "定型句")
]

// 説明体の指示文を書き換えた、未保存の編集中の指示文。dirty 表示の確認に使う。
export const EDITED_DIRECTIVES: Directive[] = DIRECTIVES.map((d) =>
  d.key === "説明体"
    ? {
        ...d,
        instruction:
          "対象の機能・効果だけを 1〜2 文で簡潔に説明すること。修飾語を足さず、誇張しないこと。"
      }
    : d
)
