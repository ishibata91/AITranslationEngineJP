import { defineScreenState } from "../screen-spec"
import {
  EMPTY_STATE,
  LOADING_STATE,
  LIST_STATE,
  SELECTED_STATE,
  CONFIRM_DELETE_STATE,
  DELETING_STATE,
  ERROR_STATE
} from "./target-plugins.fixtures"

export const targetPluginScreenStates = {
  empty: defineScreenState({
    storyName: "空状態",
    precondition: "pluginを選んでおらず、翻訳したpluginもない。",
    args: EMPTY_STATE,
    specs: [
      {
        id: "target-plugins.empty.proceed-disabled",
        statement: "pluginを選んでいない時は「翻訳へ進む」が無効になる。"
      },
      {
        id: "target-plugins.empty.list-empty",
        statement: "pluginの一覧が空であることを案内する。"
      }
    ]
  }),
  loading: defineScreenState({
    storyName: "読み込み中",
    precondition: "pluginの一覧を読み込んでいる。",
    args: LOADING_STATE,
    specs: [
      {
        id: "target-plugins.loading.list-loading",
        statement:
          "読み込み中の表示と「プラグインを読み込んでいます。」を表示する。"
      }
    ]
  }),
  list: defineScreenState({
    storyName: "一覧",
    precondition: "進捗が異なる複数のpluginが一覧にある。",
    args: LIST_STATE,
    specs: [
      {
        id: "target-plugins.list.count",
        statement: "一覧にあるpluginの件数を表示する。"
      },
      {
        id: "target-plugins.list.progress-badges",
        statement:
          "各行の件数に対応して「完了」「未着手」「翻訳中」の状態バッジを表示する。"
      },
      {
        id: "target-plugins.list.row-actions",
        statement: "各行に結果を開くbuttonと削除buttonを有効な状態で表示する。"
      }
    ]
  }),
  selected: defineScreenState({
    storyName: "プラグイン選択済み",
    precondition: "新しいpluginを選び終えている。",
    args: SELECTED_STATE,
    specs: [
      {
        id: "target-plugins.selected.proceed-enabled",
        statement: "pluginを選んだ時は「翻訳へ進む」が有効になる。"
      }
    ]
  }),
  confirmDelete: defineScreenState({
    storyName: "削除確認中",
    precondition: "一覧にある一つのpluginについて削除確認中である。",
    args: CONFIRM_DELETE_STATE,
    specs: [
      {
        id: "target-plugins.confirm-delete.prompt",
        statement: "削除対象の行だけに削除確認を表示する。"
      },
      {
        id: "target-plugins.confirm-delete.actions-enabled",
        statement: "「削除する」と「取消」を有効にする。"
      }
    ]
  }),
  deleting: defineScreenState({
    storyName: "削除実行中",
    precondition: "削除確認したpluginを削除している。",
    args: DELETING_STATE,
    specs: [
      {
        id: "target-plugins.deleting.progress",
        statement: "「削除中…」と処理中の表示を出す。"
      },
      {
        id: "target-plugins.deleting.actions-disabled",
        statement: "「削除中…」と「取消」が無効になる。"
      }
    ]
  }),
  errored: defineScreenState({
    storyName: "エラー",
    precondition: "pluginの翻訳成果の削除に失敗している。",
    args: ERROR_STATE,
    specs: [
      {
        id: "target-plugins.error.message",
        statement: "エラーの内容を画面に表示する。"
      }
    ]
  })
} as const
