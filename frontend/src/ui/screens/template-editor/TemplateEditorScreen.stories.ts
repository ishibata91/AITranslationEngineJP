import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TemplateEditorScreen from "./TemplateEditorScreen.svelte"
import {
  DEFAULT_TEMPLATE_FORM,
  DIRECTIVES,
  EDITED_DIRECTIVES,
  RECORD_ASSIGNMENTS,
  TONE_DEFAULT_EDITED_FORM
} from "./template-editor.fixtures"

// プロンプトテンプレート画面。Base + 種別ごとの指示文モデルに、話者なし台詞の口調設定を足した。
// generic-voice-tone-fallback で、汎用台詞・PC 発話の口調を自由記述で書き、PC の性別を選べるようにした。
// translation-prompt-revision で Base 指示文を 4 段落へ、指示文を 9 種へ変えた。
// Storybook 人間レビュー承認済み（2026-07-26）。通常分類（Screens）に置く。
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

// レコード別タブ。話者なし台詞の口調（汎用・PC の自由記述と PC 性別）と、種別ごとの指示文を並べる。
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

// レコード別タブで汎用口調を書き換え、PC 性別を男性に選んだ状態。自由記述の編集と PC 性別選択が反映される。
export const RecordTabToneDefaultEdited: Story = {
  name: "レコード別タブ（口調と PC 性別を変更）",
  args: {
    form: TONE_DEFAULT_EDITED_FORM,
    activeTab: "record",
    directives: DIRECTIVES,
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

// レコード別タブで物品説明の指示文を書き換えた未保存状態。未保存バッジが出て保存・戻すが有効。
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
