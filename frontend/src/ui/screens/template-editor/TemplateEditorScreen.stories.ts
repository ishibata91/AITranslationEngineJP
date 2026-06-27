import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TemplateEditorScreen from "./TemplateEditorScreen.svelte"
import {
  DEFAULT_TEMPLATE_FORM,
  DIRECTIVES,
  EDITED_DIRECTIVES,
  RECORD_ASSIGNMENTS
} from "./template-editor.fixtures"

// プロンプトテンプレート画面。record-type-translation-expansion で「Base + 種別ごとの指示文」モデルへ作り直した。
// Storybook 人間レビュー承認済み（2026-06-25）。通常分類（Screens）に置く。
const meta = {
  title: "Screens/プロンプトテンプレート",
  component: TemplateEditorScreen,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof TemplateEditorScreen>

export default meta
type Story = StoryObj<typeof meta>

const noop = () => {}

// ベースタブ。全文に付く Base 指示だけを出す。未保存の変更なし。
export const BaseTab: Story = {
  name: "ベースタブ",
  args: {
    form: DEFAULT_TEMPLATE_FORM,
    activeTab: "base",
    directives: DIRECTIVES,
    assignments: RECORD_ASSIGNMENTS,
    dirty: false,
    saving: false,
    onFieldInput: noop,
    onInstructionInput: noop,
    onTabChange: noop,
    onSave: noop,
    onReset: noop
  }
}

// レコード別タブ。種別ごとの指示文（口調・文体・固有名・定型句）を同じ形で並べる。未保存の変更なし。
export const RecordTab: Story = {
  name: "レコード別タブ",
  args: {
    form: DEFAULT_TEMPLATE_FORM,
    activeTab: "record",
    directives: DIRECTIVES,
    assignments: RECORD_ASSIGNMENTS,
    dirty: false,
    saving: false,
    onFieldInput: noop,
    onInstructionInput: noop,
    onTabChange: noop,
    onSave: noop,
    onReset: noop
  }
}

// レコード別タブで説明体の指示文を書き換えた未保存状態。未保存バッジが出て保存・戻すが有効。
export const RecordTabDirty: Story = {
  name: "レコード別タブ（未保存）",
  args: {
    form: DEFAULT_TEMPLATE_FORM,
    activeTab: "record",
    directives: EDITED_DIRECTIVES,
    assignments: RECORD_ASSIGNMENTS,
    dirty: true,
    saving: false,
    onFieldInput: noop,
    onInstructionInput: noop,
    onTabChange: noop,
    onSave: noop,
    onReset: noop
  }
}
