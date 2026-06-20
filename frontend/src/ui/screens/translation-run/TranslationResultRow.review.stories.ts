import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationResultRow from "./TranslationResultRow.svelte"

// Slice 2（実プロンプト参照）の Storybook 人間レビュー中。作業中分類（Review/Changed Components）に置く。
// 承認後に「実プロンプト」状態を UI Components/TranslationResultRow へ統合し、このファイルを削除する。
const meta = {
  title: "Review/Changed Components/TranslationResultRow",
  component: TranslationResultRow,
  parameters: { layout: "padded" }
} satisfies Meta<typeof TranslationResultRow>

export default meta
type Story = StoryObj<typeof meta>

// base 指示文。provider が現在ハードコードで持つ翻訳指示（openai_compatible.go:152 相当）。
const BASE = `[system]
あなたは Skyrim Mod の翻訳者です。与えられた英語の本文を、原文の意味と語調を保った自然な日本語へ翻訳してください。訳文だけを出力し、説明や注釈は加えないでください。`

// 台詞。口調指示が system に入り、置換済み原文（Dark Brotherhood → 闇の一党）が user に入る実プロンプト全文。
const PERSONA_TERM_PROMPT_ROW = {
  edid: "AventusAretinoRitual",
  source: "I knew the Dark Brotherhood would answer my prayer.",
  dest: "闇の一党が僕の祈りに応えてくれるって、分かってたんだ。",
  statusLabel: "仮訳",
  personaLabel: "声質: 幼い少年の声",
  directive:
    "この台詞の話者の人物像:\n- 声質: 幼い少年の声\n- 種族の気質: ノルド（実直で粘り強い北方の気質）\nこの人物像に合う口調と人称で訳すこと。",
  terms: [{ source: "Dark Brotherhood", dest: "闇の一党" }],
  prompt: `${BASE}

この台詞の話者の人物像:
- 声質: 幼い少年の声
- 種族の気質: ノルド（実直で粘り強い北方の気質）
この人物像に合う口調と人称で訳すこと。

[user]
I knew the 闇の一党 would answer my prayer.`
}

// 叙述文。口調指示は無く（台詞ではない）、system は base 指示のみ。置換済み原文（Castle Volkihar → ヴォルキハル城）が user に入る。
const NARRATION_PROMPT_ROW = {
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

// 台詞を展開。置換した固有名と実プロンプトの 2 節が並ぶ。口調指示の独立節は廃し、口調は実プロンプトの system 部分で確かめる。
export const ExpandedPersonaTermsPrompt: Story = {
  name: "展開（口調あり台詞・置換・実プロンプト）",
  args: { row: PERSONA_TERM_PROMPT_ROW, defaultOpen: true }
}

// 叙述文を展開。口調指示は無く、実プロンプトの system は base 指示だけになる。置換済み原文が user に入る。
export const ExpandedNarrationPrompt: Story = {
  name: "展開（叙述文・口調なし・実プロンプト）",
  args: { row: NARRATION_PROMPT_ROW, defaultOpen: true }
}
