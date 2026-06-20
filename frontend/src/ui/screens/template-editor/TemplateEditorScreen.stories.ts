import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TemplateEditorScreen from "./TemplateEditorScreen.svelte"
import {
  DEFAULT_TEMPLATE_FORM,
  EDITED_TEMPLATE_FORM,
  PLACEHOLDERS
} from "./template-editor.fixtures"

// Slice 4（テンプレート編集）の Storybook 人間レビュー中。作業中分類（Review/Changed Screens）に置く。
// 承認後に Screens/プロンプトテンプレート へ戻す。
const meta = {
  title: "Review/Changed Screens/プロンプトテンプレート",
  component: TemplateEditorScreen,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof TemplateEditorScreen>

export default meta
type Story = StoryObj<typeof meta>

const noop = () => {}
const noopInput = () => {}

// 既定テンプレートを表示した初期状態。未保存の変更が無く、保存・戻すは無効。
export const Default: Story = {
  name: "初期（既定テンプレート）",
  args: {
    form: DEFAULT_TEMPLATE_FORM,
    placeholders: PLACEHOLDERS,
    dirty: false,
    saving: false,
    onFieldInput: noopInput,
    onSave: noop,
    onReset: noop
  }
}

// 口調指示テンプレートを書き換えた未保存状態。未保存バッジが出て、保存・戻すが有効。
export const Dirty: Story = {
  name: "編集中（未保存）",
  args: {
    form: EDITED_TEMPLATE_FORM,
    placeholders: PLACEHOLDERS,
    dirty: true,
    saving: false,
    onFieldInput: noopInput,
    onSave: noop,
    onReset: noop
  }
}
