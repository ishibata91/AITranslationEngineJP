import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationResultRow from "./TranslationResultRow.svelte"

// Storybook 人間レビュー承認済み。通常分類（UI Components）に置く。
const meta = {
  title: "UI Components/TranslationResultRow",
  component: TranslationResultRow,
  parameters: { layout: "padded" }
} satisfies Meta<typeof TranslationResultRow>

export default meta
type Story = StoryObj<typeof meta>

const CHILD_ROW = {
  edid: "AventusAretinoInnocenceLost",
  source: "Mother? Are you there? It's so cold without you.",
  dest: "母さん？ そこにいるの？ あなたがいないと、こんなに寒いんだ。",
  statusLabel: "仮訳",
  personaLabel: "声質: 幼い少年の声",
  directive:
    "この台詞の話者の人物像:\n- 声質: 幼い少年の声\n- 種族の気質: ノルド（実直で粘り強い北方の気質）\nこの人物像に合う口調と人称で訳すこと。"
}

// 既定の畳んだ 1 行。口調チップ（声質）で口調が分かり、全文は出さない。
export const CollapsedWithPersona: Story = {
  name: "畳む（口調あり・子供）",
  args: { row: CHILD_ROW, defaultOpen: false }
}

// 行を展開した状態。原文・訳文・口調指示の全文が出る。
export const ExpandedWithPersona: Story = {
  name: "展開（口調あり・子供）",
  args: { row: CHILD_ROW, defaultOpen: true }
}

// 別の話者（老女）。畳んだ一覧でも口調チップが子供と違い、口調差が分かる。
export const CollapsedOldWoman: Story = {
  name: "畳む（口調あり・老女）",
  args: {
    row: {
      edid: "GrelodTheKindScold",
      source:
        "You are all here to serve and honor your matron. Now off to bed, the lot of you!",
      dest: "お前たちは院母に仕え、敬うためにここにいるんだよ。さあ、とっとと寝な、お前たち全員！",
      statusLabel: "仮訳",
      personaLabel: "声質: しわがれた老女の声",
      directive:
        "この台詞の話者の人物像:\n- 声質: しわがれた老女の声\n- 種族の気質: ノルド（実直で粘り強い北方の気質）\nこの人物像に合う口調と人称で訳すこと。"
    },
    defaultOpen: false
  }
}

// 話者を解決できなかった行。口調チップは出さず「口調なし」を控えめに示す。
export const CollapsedWithoutPersona: Story = {
  name: "畳む（口調なし）",
  args: {
    row: {
      edid: "HonorhallDoorActivate",
      source: "The door is locked.",
      dest: "扉には鍵がかかっている。",
      statusLabel: "仮訳"
    },
    defaultOpen: false
  }
}
