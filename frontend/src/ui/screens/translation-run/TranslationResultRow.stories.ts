import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationResultRow from "./TranslationResultRow.svelte"
import type { NarrationResultRow } from "./translation-run-view"

// Storybook 人間レビュー承認済み。通常分類（UI Components）に置く。
// 結果行を展開すると、台詞の話者の生成済み基底口調を「口調」メタデータとして強調表示する（判定結果＋性質文を大きく、根拠は小さく）。
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

// 叙述文。口調は無く（台詞ではない）、置換固有名と実プロンプトだけが並ぶ。
const NARRATION_ROW: NarrationResultRow = {
  edid: "DLC1BookSerana",
  source:
    "I have walked these halls for centuries, and still the cold of Castle Volkihar finds me.",
  dest: "私は何世紀もこの広間を歩いてきたが、それでもヴォルキハル城の冷気は私を捉えて離さない。",
  statusLabel: "仮訳",
  terms: [{ source: "Castle Volkihar", dest: "ヴォルキハル城" }],
  prompt: `${BASE}

[user]
I have walked these halls for centuries, and still the cold of ヴォルキハル城 finds me.`
}

// 話者を解決できなかった行（口調も置換も無い）。畳んだ行は「口調なし」を控えめに示す。
const PLAIN_ROW: NarrationResultRow = {
  edid: "HonorhallDoorActivate",
  source: "The door is locked.",
  dest: "扉には鍵がかかっている。",
  statusLabel: "仮訳"
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

// 叙述文を展開。口調メタ節は出ず、置換固有名と実プロンプトが並ぶ。
export const ExpandedNarration: Story = {
  name: "展開（叙述文・口調なし）",
  args: { row: NARRATION_ROW, defaultOpen: true }
}

// 口調も置換も無い行。畳んだ行で「口調なし」を控えめに示す。
export const CollapsedWithoutPersona: Story = {
  name: "畳む（口調なし）",
  args: { row: PLAIN_ROW, defaultOpen: false }
}
