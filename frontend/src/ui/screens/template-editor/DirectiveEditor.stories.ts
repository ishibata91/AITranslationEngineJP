import type { Meta, StoryObj } from "@storybook/svelte-vite"
import DirectiveEditor from "./DirectiveEditor.svelte"
import { DIRECTIVES, RECORD_ASSIGNMENTS } from "./template-editor.fixtures"
import { buildDirectiveSections } from "./directive-presentation"

// 種別ごとの指示文の編集表示（本文 textarea ＋ 変数 ＋ 対象）。
// Storybook 人間レビュー承認済み（2026-06-25）。通常分類（UI Components）に置く。
const meta = {
  title: "UI Components/DirectiveEditor",
  component: DirectiveEditor,
  parameters: { layout: "padded" }
} satisfies Meta<typeof DirectiveEditor>

export default meta
type Story = StoryObj<typeof meta>

const noop = () => {}

// 7 指示文（口調は {traits} 変数つき、固有名・定型句も指示文）を同じ形で並べる。
export const Default: Story = {
  name: "指示文一覧（口調・文体・固有名・定型句）",
  args: {
    sections: buildDirectiveSections(DIRECTIVES, RECORD_ASSIGNMENTS),
    onInstructionInput: noop
  }
}
