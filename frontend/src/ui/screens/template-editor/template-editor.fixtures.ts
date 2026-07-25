// プロンプトテンプレート画面の表示用 fixture。Base 指示と、種別ごとの指示文（directive）を表示値として固定する。
// プロンプト = Base 指示 + その REC:FIELD に割り当てた指示文（変数を実行時に埋めたもの）。
import type { PromptTemplateForm } from "./template-editor-view"
import type { Directive, RecordAssignment } from "./directive-view"

// base 翻訳指示文。全リクエストの先頭に付く。役割・機械置換済み固有名の保持・出力の崩れ方の禁止・
// 口調と原文尊重の優先順位を、1 段落 1 論点で並べる（translation-prompt-revision の層1）。
// 出力形（訳文だけを返す）は構造化出力の schema が強制するため、指示文では述べない。
const DEFAULT_BASE = `あなたは The Elder Scrolls V: Skyrim の Mod を日本語へ訳す翻訳者である。与えられた英語の本文を、意味を変えずに自然な日本語へ訳す。

本文には日本語へ置き換え済みの固有名が混ざる場合がある。日本語で書かれた部分はそのまま残し、訳し直したり表記を変えたりしない。

原文の改行の数と位置を保つ。原文に無い鍵括弧・句点・感嘆符を足さない。英単語を訳さずに残さない。

続く指示で口調を指定する場合、口調は語の選び方と語尾に反映する。原文の意味を変える理由にはしない。`

// 口調指示。指示文「口調」の既定文（db/migrations/0006_record_type_translation.sql の seed 相当）。
// {traits} に話者の性質（基底口調・役割語・口調の例文・種族訛り）の箇条書きを差し込む。
const DEFAULT_PERSONA = `この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。`

// 話者なし台詞の口調の既定文。汎用台詞・PC 発話それぞれの自由記述の口調指示。
const DEFAULT_GENERIC_TONE = `衛兵などの不特定多数が話す汎用的な台詞。職務的で簡潔な口調で訳す。`
const DEFAULT_PC_TONE = `プレイヤーキャラクターの選択肢。自然な口語で訳す。`

export const DEFAULT_TEMPLATE_FORM: PromptTemplateForm = {
  baseDirective: DEFAULT_BASE,
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

// 種別ごとの指示文（9 種）。文体・固有名・操作名・語義・口調をすべて同じ形で持つ。
// 口調だけ {traits} 変数を持つ。並びは design.md 層2 の表と同じで、叙述文の 5 種 → 固有名 →
// 短文の 2 種 → 口調 の順。旧「説明体」は物品説明と効果説明へ、旧「定型句」は操作名と語義へ割った。
export const DIRECTIVES: Directive[] = [
  {
    key: "物品説明",
    instruction:
      "これは武器・防具・薬・巻物などの品物の説明文です。用途と特徴を正確に保ち、簡潔で読みやすい説明口調の日本語へ訳すこと。",
    variables: []
  },
  {
    key: "効果説明",
    instruction:
      "これは呪文・付呪・特典・シャウト・魔法効果の説明文です。数値と <> で囲まれた実行時タグは原文のまま残し、増減や順序を変えないこと。体言止めで短くまとめること。",
    variables: []
  },
  {
    key: "世界観断片",
    instruction:
      "これはロード画面や種族の解説です。作品世界を語る地の文として、簡潔で含みのある説明口調の日本語へ訳すこと。",
    variables: []
  },
  {
    key: "書物体",
    instruction:
      "これは書物の本文です。文章の格調と語り口を保ち、読み物として自然な日本語へ訳すこと。",
    variables: []
  },
  {
    key: "日記体",
    instruction:
      "これはクエストの進行ログまたは目標です。主人公の視点で、簡潔な記録口調の日本語へ訳すこと。",
    variables: []
  },
  {
    key: "固有名",
    instruction:
      "これは固有名詞です。カタカナ音写を基本とし、意味訳へ置き換えないこと。ゲーム UI に収まるよう、原語より長くしないこと。既存の確定訳語があればそれに合わせること。",
    variables: []
  },
  {
    key: "操作名",
    instruction:
      "これは調べる・採取するなどの操作を表す短い語句です。動詞の終止形で短く訳し、UI の表示幅に収まる長さにすること。",
    variables: []
  },
  {
    key: "語義",
    instruction:
      "これは龍語など、語そのものの意味を述べる語釈です。語の意味を短く言い切る名詞句または短文で訳すこと。説明を足さないこと。",
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
  }
]

// REC:FIELD と指示文キーの対応（固定）。
function assign(recField: string, logicalName: string, directive: string): RecordAssignment {
  return { recField, logicalName, directive }
}

// 割り当ての正本は db/migrations/0006_record_type_translation.sql の record_type_master seed。
// ここは表示確認用の抜粋で、9 種すべてが対象を 1 件以上持つように選ぶ。logical_name は seed に合わせる。
export const RECORD_ASSIGNMENTS: RecordAssignment[] = [
  assign("WEAP:DESC", "武器の説明", "物品説明"),
  assign("ARMO:DESC", "防具の説明", "物品説明"),
  assign("ALCH:DESC", "薬・食料の説明", "物品説明"),
  assign("MESG:DESC", "メッセージ本文", "物品説明"),
  assign("SPEL:DESC", "呪文の説明", "効果説明"),
  assign("ENCH:DESC", "付呪効果の説明", "効果説明"),
  assign("PERK:DESC", "特典の説明", "効果説明"),
  assign("SHOU:DESC", "シャウトの説明", "効果説明"),
  assign("MGEF:DNAM", "魔法効果の説明", "効果説明"),
  assign("LSCR:DESC", "ロード画面の解説", "世界観断片"),
  assign("RACE:DESC", "種族の説明", "世界観断片"),
  assign("BOOK:DESC", "本の本文", "書物体"),
  assign("QUST:CNAM", "クエストログ", "日記体"),
  assign("QUST:NNAM", "クエスト目標", "日記体"),
  assign("WEAP:FULL", "武器の名前", "固有名"),
  assign("ARMO:FULL", "防具の名前", "固有名"),
  assign("SPEL:FULL", "呪文の名前", "固有名"),
  assign("BOOK:FULL", "本の題名", "固有名"),
  assign("CELL:FULL", "セルの名前", "固有名"),
  assign("QUST:FULL", "クエスト名", "固有名"),
  assign("RACE:FULL", "種族名", "固有名"),
  assign("FACT:FULL", "勢力名", "固有名"),
  assign("NPC_:FULL", "NPC の氏名", "固有名"),
  assign("NPC_:SHRT", "NPC の短縮名", "固有名"),
  assign("ACTI:RNAM", "オブジェクトの操作名", "操作名"),
  assign("FLOR:RNAM", "採取植物の操作名", "操作名"),
  assign("TREE:RNAM", "樹木の操作名", "操作名"),
  assign("MESG:ITXT", "メッセージのボタン", "操作名"),
  assign("WOOP:TNAM", "龍語の語義", "語義"),
  assign("INFO:NAM1", "NPC の返答", "口調"),
  assign("INFO:RNAM", "選択肢の条件別上書き", "口調"),
  assign("DIAL:FULL", "選択肢の既定文", "口調")
]

// 物品説明の指示文を書き換えた、未保存の編集中の指示文。dirty 表示の確認に使う。
export const EDITED_DIRECTIVES: Directive[] = DIRECTIVES.map((d) =>
  d.key === "物品説明"
    ? {
        ...d,
        instruction:
          "対象の機能・効果だけを 1〜2 文で簡潔に説明すること。修飾語を足さず、誇張しないこと。"
      }
    : d
)
