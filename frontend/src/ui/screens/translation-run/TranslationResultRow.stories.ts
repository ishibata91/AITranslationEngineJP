import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationResultRow from "./TranslationResultRow.svelte"
import type { NarrationResultRow } from "./translation-run-view"

// 結果行を展開すると、台詞の話者の生成済み基底口調を「口調」メタデータとして強調表示する（判定結果＋性質文を大きく、根拠は小さく）。
// generic-voice-tone-fallback で、話者を解決できない汎用台詞と PC 発話へも口調を付ける表示を追加した。
// 汎用・PC は対人段階・セル・印を持たず、見出し（汎用台詞／PC発話）＋感情＋性別を出す。
// Storybook 人間レビュー承認済み（2026-06-30）。通常分類（UI Components）に置く。
const meta = {
  title: "UI Components/TranslationResultRow",
  component: TranslationResultRow,
  parameters: { layout: "padded" }
} satisfies Meta<typeof TranslationResultRow>

export default meta
type Story = StoryObj<typeof meta>

// base 指示文。provider が現在持つ翻訳指示の相当。
const BASE = `[system]
あなたは Skyrim Mod の翻訳者です。与えられた英語の本文を、原文の意味と語調を保った自然な日本語へ翻訳してください。訳文だけを出力し、説明や注釈は加えないでください。`

// 本文経路の台詞（Inigo、Khajiit、物腰やわ）。口調メタ＋種族訛り＋置換固有名＋実プロンプトが並ぶ。
const DIALOGUE_ROW: NarrationResultRow = {
  edid: "Inigo",
  recordType: { box: "台詞", recField: "INFO:NAM1" },
  source: "This one knows Riften well, my friend.",
  dest: "この者はリフテンをよく知っているよ、友よ。",
  statusLabel: "仮訳",
  speaker: { edid: "Inigo", sex: "男性", age: "成人", voice: "MaleUniqueInigo" },
  personaLabel: "口調: 物腰やわ",
  persona: {
    cell: "物腰やわ",
    trait: "柔らかく丁寧な口調。相手を立てて穏やかに述べる。",
    attitudeBand: 2,
    emotionBand: 1,
    marked: 452,
    decisionPath: "本文"
  },
  directive:
    "この台詞の話者の人物像:\n- 口調: 柔らかく丁寧な口調。相手を立てて穏やかに述べる。\n- 種族訛り: カジートの訛り。三人称で自称する（「この者は」など）。\nこの人物像に合う口調と人称で訳すこと。",
  terms: [{ source: "Riften", dest: "リフテン" }],
  prompt: `${BASE}

この台詞の話者の人物像:
- 口調: 柔らかく丁寧な口調。相手を立てて穏やかに述べる。
- 種族訛り: カジートの訛り。三人称で自称する（「この者は」など）。
この人物像に合う口調と人称で訳すこと。

[user]
This one knows リフテン well, my friend.`
}

// 声質経路の台詞（Nazeem、台詞少で声質の横柄を prior にした）。決定経路が「声質」、印が少ない（4）。
const VOICE_ROW: NarrationResultRow = {
  edid: "Nazeem",
  recordType: { box: "台詞", recField: "INFO:NAM1" },
  source: "Do you get to the Cloud District very often?",
  dest: "雲の地区にはよく行かれるのかね？",
  statusLabel: "仮訳",
  speaker: { edid: "Nazeem", sex: "男性", age: "成人", voice: "MaleCondescending" },
  personaLabel: "口調: 冷然・見下し",
  persona: {
    cell: "冷然・見下し",
    trait: "相手を見下す冷たい口調。感情を抑え、丁寧さを欠いた突き放した言い方にする。",
    attitudeBand: 0,
    emotionBand: 0,
    marked: 4,
    decisionPath: "voice"
  },
  directive:
    "この台詞の話者の人物像:\n- 口調: 相手を見下す冷たい口調。感情を抑え、丁寧さを欠いた突き放した言い方にする。\nこの人物像に合う口調と人称で訳すこと。"
}

// 保留経路の台詞（印が不足で固有声質に気質も無く、薄い本文値を低信頼で保持）。決定経路が「保留」。
const HELD_ROW: NarrationResultRow = {
  edid: "Galmar",
  recordType: { box: "台詞", recField: "INFO:NAM1" },
  source: "Stand and fight, you milk-drinker!",
  dest: "立って戦え、この臆病者が！",
  statusLabel: "仮訳",
  speaker: { edid: "Galmar", sex: "男性", age: "成人", voice: "MaleUniqueGalmar" },
  personaLabel: "口調: ぞんざい",
  persona: {
    cell: "ぞんざい",
    trait: "ぶっきらぼうで乱暴な口調。命令的に言い、相手を立てない。",
    attitudeBand: 0,
    emotionBand: 1,
    marked: 9,
    decisionPath: "保留"
  },
  directive:
    "この台詞の話者の人物像:\n- 口調: ぶっきらぼうで乱暴な口調。命令的に言い、相手を立てない。\nこの人物像に合う口調と人称で訳すこと。"
}

// 汎用台詞（話者を特定できない衛兵の汎用セリフ）。自由記述の汎用口調指示文＋本文 1 行の感情段階（中）
// ＋条件由来の性別（男性）。対人段階・セル・印は持たない。話者欄は出ず、口調メタは「感情 中 ・ 性別 男性」。
// 単発の感嘆符で激情へ振れないよう、感情段階は「中」に留まる（渋めしきい値）。
const GENERIC_ROW: NarrationResultRow = {
  edid: "RiftenGuardGenericHalt",
  recordType: { box: "台詞", recField: "INFO:NAM1" },
  source: "Halt! State your business.",
  dest: "止まれ！ 用件を述べろ。",
  statusLabel: "仮訳",
  personaLabel: "口調: 汎用台詞",
  persona: {
    emotionBand: 1,
    decisionPath: "汎用",
    sex: "男性"
  },
  directive:
    "この台詞の話者の人物像:\n- 口調: 衛兵などの不特定多数が話す汎用的な台詞。職務的で簡潔な口調で訳す。\n- 人称と言い回し: 一人称は「俺」。\nこの人物像に合う口調と人称で訳すこと。"
}

// PC 発話（プレイヤーの選択肢文、DIAL:FULL）。自由記述の PC 口調指示文＋本文 1 行の感情段階（抑制）
// ＋利用者選択の PC 性別（男性）。対人段階・セル・印は持たない。口調メタは「感情 抑制 ・ 性別 男性」。
const PC_ROW: NarrationResultRow = {
  edid: "DialogueRiftenThievesGuildJoin",
  recordType: { box: "台詞", recField: "DIAL:FULL" },
  source: "I'm looking to join the Thieves Guild.",
  dest: "盗賊ギルドに入りたいんだが。",
  statusLabel: "仮訳",
  personaLabel: "口調: PC発話",
  persona: {
    emotionBand: 0,
    decisionPath: "PC",
    sex: "男性"
  },
  directive:
    "この台詞の話者の人物像:\n- 口調: プレイヤーキャラクターの選択肢。自然な口語で訳す。\n- 人称と言い回し: 一人称は「俺」。\nこの人物像に合う口調と人称で訳すこと。"
}

// 叙述文。口調は無く（台詞ではない）、置換固有名と実プロンプトだけが並ぶ。
const NARRATION_ROW: NarrationResultRow = {
  edid: "DLC1BookSerana",
  recordType: { box: "叙述文", recField: "BOOK:DESC" },
  source:
    "I have walked these halls for centuries, and still the cold of Castle Volkihar finds me.",
  dest: "私は何世紀もこの広間を歩いてきたが、それでもヴォルキハル城の冷気は私を捉えて離さない。",
  statusLabel: "仮訳",
  terms: [{ source: "Castle Volkihar", dest: "ヴォルキハル城" }],
  prompt: `${BASE}

[user]
I have walked these halls for centuries, and still the cold of ヴォルキハル城 finds me.`
}

// 定型句（起動動作）。話者は無く、種別バッジで定型句と分かる。
const PLAIN_ROW: NarrationResultRow = {
  edid: "HonorhallDoorActivate",
  recordType: { box: "定型句", recField: "ACTI:RNAM" },
  source: "The door is locked.",
  dest: "扉には鍵がかかっている。",
  statusLabel: "仮訳"
}

// 追加した叙述文の種別（武器の説明）。種別バッジで叙述文 ・ WEAP:DESC と分かる。
const WEAPON_DESC_ROW: NarrationResultRow = {
  edid: "DragonbaneDesc",
  recordType: { box: "叙述文", recField: "WEAP:DESC" },
  source: "A blade forged to slay dragons, humming with stored lightning.",
  dest: "竜を討つために鍛えられた刃。蓄えた雷の力で唸りを上げる。",
  statusLabel: "仮訳",
  terms: [{ source: "Dragonbane", dest: "竜の災い" }]
}

// 固有名（武器名）の AI 訳。本文より先に確定し proper_noun へ入った訳。話者・口調は無い。
const PROPER_NOUN_ROW: NarrationResultRow = {
  edid: "Dragonbane",
  recordType: { box: "固有名", recField: "WEAP:FULL" },
  source: "Dragonbane",
  dest: "竜の災い",
  statusLabel: "訳済"
}

// 畳んだ台詞。口調チップ（基底口調セル）と「固有名 N」チップで、一覧のまま口調差と置換を観測する。
export const CollapsedDialogue: Story = {
  name: "畳む（口調あり台詞）",
  args: { row: DIALOGUE_ROW, defaultOpen: false }
}

// 本文経路の台詞を展開。口調メタ（物腰やわ＋性質文を強調、根拠は小さく）＋置換固有名＋実プロンプトが並ぶ。
export const ExpandedDialoguePersona: Story = {
  name: "展開（本文・口調メタ）",
  args: { row: DIALOGUE_ROW, defaultOpen: true }
}

// 声質経路の台詞を展開。決定経路が「声質」、印が少ない。
export const ExpandedVoicePersona: Story = {
  name: "展開（声質・口調メタ）",
  args: { row: VOICE_ROW, defaultOpen: true }
}

// 保留経路の台詞を展開。決定経路が「保留」（低信頼）。
export const ExpandedHeldPersona: Story = {
  name: "展開（保留・口調メタ）",
  args: { row: HELD_ROW, defaultOpen: true }
}

// 汎用台詞を展開。話者欄は出ず、決定経路が「汎用既定」、根拠に「性別 男性」を出し、印は出さない。
export const ExpandedGenericPersona: Story = {
  name: "展開（汎用既定・口調メタ）",
  args: { row: GENERIC_ROW, defaultOpen: true }
}

// PC 発話を展開。話者欄は出ず、決定経路が「PC既定」、性別・印は出ない。
export const ExpandedPcPersona: Story = {
  name: "展開（PC既定・口調メタ）",
  args: { row: PC_ROW, defaultOpen: true }
}

// 汎用台詞の畳んだ行。口調チップ（基底口調セル）が付き、話者を特定できなくても口調が決まることを示す。
export const CollapsedGeneric: Story = {
  name: "畳む（汎用台詞）",
  args: { row: GENERIC_ROW, defaultOpen: false }
}

// 叙述文を展開。口調メタ節は出ず、置換固有名と実プロンプトが並ぶ。
export const ExpandedNarration: Story = {
  name: "展開（叙述文・口調なし）",
  args: { row: NARRATION_ROW, defaultOpen: true }
}

// 定型句の畳んだ行。種別バッジで「定型句 ・ ACTI:RNAM」と分かる。口調なし。
export const CollapsedSetPhrase: Story = {
  name: "畳む（定型句）",
  args: { row: PLAIN_ROW, defaultOpen: false }
}

// 追加した叙述文の種別（武器の説明）の畳んだ行。種別バッジで「叙述文 ・ WEAP:DESC」と分かる。
export const CollapsedWeaponDesc: Story = {
  name: "畳む（叙述文・WEAP:DESC）",
  args: { row: WEAPON_DESC_ROW, defaultOpen: false }
}

// 固有名の AI 訳の畳んだ行。種別バッジで「固有名 ・ WEAP:FULL」と分かる。本文より先に確定した訳。
export const CollapsedProperNoun: Story = {
  name: "畳む（固有名）",
  args: { row: PROPER_NOUN_ROW, defaultOpen: false }
}
