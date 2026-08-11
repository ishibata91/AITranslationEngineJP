import {
  emptyPrebuiltDictionaryFilters,
  prebuiltDictionaryPageRows,
  prebuiltDictionaryRows
} from "./prebuilt-dictionary-editor.fixtures"
import { defineScreenState } from "../screen-spec"

export const prebuiltDictionaryEditorScreenStates = {
  list: defineScreenState({
    storyName: "一覧",
    precondition: "用語辞書が表示されている。",
    args: {
      rows: prebuiltDictionaryPageRows,
      filters: emptyPrebuiltDictionaryFilters,
      pageNumber: 1,
      totalCount: 15917,
      canPrev: false,
      canNext: true,
      expandedRowIds: [],
      hasPendingChanges: false,
      editingRowId: ""
    },
    specs: [
      { id: "prebuilt-dictionary.list.table", statement: "テーブルUIに辞書と各カラムのフィルタを表示する。" },
      { id: "prebuilt-dictionary.list.pagination", statement: "ページ送りと1ページ50件を表示する。" },
      { id: "prebuilt-dictionary.row-edit.selected", statement: "編集を選んだ行だけに入力欄を表示する。" }
    ]
  }),
  edit: defineScreenState({
    storyName: "編集",
    precondition: "テーブルUIに表示された用語辞書を編集している。",
    args: {
      filters: emptyPrebuiltDictionaryFilters,
      pageNumber: 1,
      totalCount: 15917,
      canPrev: false,
      canNext: true,
      rows: prebuiltDictionaryRows.map((row, index) => ({ ...row, pending: index === 0 ? "edited" : undefined, deletePending: index === 1 })),
      expandedRowIds: ["2"],
      hasPendingChanges: true,
      editingRowId: "1"
    },
    specs: [
      { id: "prebuilt-dictionary.edit.form", statement: "編集する入力内容と保存操作を表示する。" },
      { id: "prebuilt-dictionary.edit.row-cancel", statement: "編集保留と削除保留の行に同じ取消操作を表示する。" },
      { id: "prebuilt-dictionary.row-edit.deleted", statement: "削除保留の行に編集操作と入力欄を表示しない。" }
    ]
  }),
  empty: defineScreenState({
    storyName: "該当なし",
    precondition: "各カラムのフィルタに一致する用語辞書がない。",
    args: {
      rows: [],
      filters: { ...emptyPrebuiltDictionaryFilters, source: "Unmatched" },
      pageNumber: 1,
      totalCount: 0,
      canPrev: false,
      canNext: false,
      expandedRowIds: [],
      hasPendingChanges: false,
      editingRowId: ""
    },
    specs: [
      { id: "prebuilt-dictionary.empty.message", statement: "一致する辞書がないことを表示する。" }
    ]
  }),
  longText: defineScreenState({
    storyName: "長い原語と訳語",
    precondition: "原語と訳語が表示幅を大幅に超えている。",
    args: {
      rows: [{
        id: "long-text",
        source: "TheUnfathomablyLongNameOfAnAncientDwemerArtifactThatWasDiscoveredDeepWithinTheRuinsOfBlackreach",
        destination: "ブラックリーチの深奥で発見された計り知れないほど長い名前を持つ古代ドゥーマーの遺物",
        partOfSpeech: "noun",
        categories: ["MISC"]
      }],
      filters: emptyPrebuiltDictionaryFilters,
      pageNumber: 1,
      totalCount: 1,
      canPrev: false,
      canNext: false,
      expandedRowIds: [],
      hasPendingChanges: false,
      editingRowId: "long-text"
    },
    specs: [
      { id: "prebuilt-dictionary.long-text.editable", statement: "長い原語と訳語を幅が制限された編集欄で表示する。" }
    ]
  })
} 
