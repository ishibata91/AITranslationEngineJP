import type { Meta, StoryObj } from "@storybook/svelte-vite"
import DirectiveEditor from "./DirectiveEditor.svelte"
import { DIRECTIVES, RECORD_ASSIGNMENTS } from "./template-editor.fixtures"
import { buildDirectiveSections } from "./directive-presentation"

// 種別ごとの指示文の編集表示（本文 textarea ＋ 変数 ＋ 対象）。
// translation-prompt-revision で指示文を 7 種から 9 種へ増やした。
// Storybook 人間レビュー承認済み（2026-07-26）。通常分類（UI Components）に置く。
const meta = {
  title: "UI Components/DirectiveEditor",
  component: DirectiveEditor,
  parameters: { layout: "padded" }
} satisfies Meta<typeof DirectiveEditor>

export default meta
type Story = StoryObj<typeof meta>

const noop = () => {}

// 9 指示文を同じ形で並べる。叙述文の 5 種（物品説明・効果説明・世界観断片・書物体・日記体）、
// 固有名、短文の 2 種（操作名・語義）、口調（{traits} 変数つき）。行の増減以外に表示構造は変えない。
export const Default: Story = {
  name: "指示文一覧（9 種）",
  args: {
    sections: buildDirectiveSections(DIRECTIVES, RECORD_ASSIGNMENTS),
    onInstructionInput: noop
  }
}
