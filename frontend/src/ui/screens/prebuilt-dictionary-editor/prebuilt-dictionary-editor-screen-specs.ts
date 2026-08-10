import {
  emptyPrebuiltDictionaryFilters,
  prebuiltDictionaryRows
} from "./prebuilt-dictionary-editor.fixtures"
import { defineScreenState } from "../screen-spec"

export const prebuiltDictionaryEditorScreenStates = {
  list: defineScreenState({
    storyName: "一覧",
    precondition: "用語辞書が表示されている。",
    args: {
      rows: prebuiltDictionaryRows,
      filters: emptyPrebuiltDictionaryFilters,
      pageNumber: 1,
      totalCount: 15917,
      canPrev: false,
      canNext: true,
      expandedRowIds: [],
      hasPendingChanges: false
    },
    specs: [
      { id: "prebuilt-dictionary.list.table", statement: "テーブルUIに辞書と各カラムのフィルタを表示する。" },
      { id: "prebuilt-dictionary.list.pagination", statement: "ページ送りと1ページ50件を表示する。" }
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
      rows: prebuiltDictionaryRows.map((row, index) => ({ ...row, pending: index === 0 ? "edited" : index === 1 ? "deleted" : undefined })),
      expandedRowIds: ["2"],
      hasPendingChanges: true
    },
    specs: [
      { id: "prebuilt-dictionary.edit.form", statement: "編集する入力内容と保存操作を表示する。" }
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
      hasPendingChanges: false
    },
    specs: [
      { id: "prebuilt-dictionary.empty.message", statement: "一致する辞書がないことを表示する。" }
    ]
  })
} 
