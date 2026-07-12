import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TargetPluginsScreen from "./TargetPluginsScreen.svelte"
import {
  EMPTY_STATE,
  LOADING_STATE,
  LIST_STATE,
  SELECTED_STATE,
  CONFIRM_DELETE_STATE,
  DELETING_STATE,
  ERROR_STATE
} from "./target-plugins.fixtures"

// 翻訳対象プラグイン画面。Storybook 人間レビュー承認済み（通常分類 Screens）。
const meta = {
  title: "Screens/翻訳対象プラグイン",
  component: TargetPluginsScreen,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    onSelectNewPlugin: () => {},
    onNewPluginPathInput: () => {},
    onProceedToRun: () => {},
    onOpenPlugin: () => {},
    onRequestDelete: () => {},
    onConfirmDelete: () => {},
    onCancelDelete: () => {}
  }
} satisfies Meta<typeof TargetPluginsScreen>

export default meta
type Story = StoryObj<typeof meta>

// まだ対象プラグインが無い初期表示。
export const Empty: Story = {
  name: "空状態",
  args: { ...EMPTY_STATE }
}

// 一覧を読み込み中。
export const Loading: Story = {
  name: "読み込み中",
  args: { ...LOADING_STATE }
}

// 複数のプラグインが並ぶ。完了・翻訳中・未着手の進捗差が出る。
export const List: Story = {
  name: "一覧",
  args: { ...LIST_STATE }
}

// 新しいプラグインを選び終え、翻訳へ進める。
export const Selected: Story = {
  name: "プラグイン選択済み",
  args: { ...SELECTED_STATE }
}

// 1 件の削除確認中。対象行だけ確認 UI を出す。
export const ConfirmDelete: Story = {
  name: "削除確認中",
  args: { ...CONFIRM_DELETE_STATE }
}

// 削除実行中。確定ボタンがスピナー・無効化される。
export const Deleting: Story = {
  name: "削除実行中",
  args: { ...DELETING_STATE }
}

// 削除に失敗した。
export const Errored: Story = {
  name: "エラー",
  args: { ...ERROR_STATE }
}
