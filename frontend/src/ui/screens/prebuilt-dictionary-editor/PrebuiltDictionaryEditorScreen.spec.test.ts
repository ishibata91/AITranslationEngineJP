import { fireEvent, render, screen, within } from "@testing-library/svelte"
import { expect, it } from "vitest"
import { runScreenSpecHarness } from "../../../test/screen-spec-harness"
import PrebuiltDictionaryEditorScreen from "./PrebuiltDictionaryEditorScreen.svelte"
import PrebuiltDictionaryEditorPreview from "./PrebuiltDictionaryEditorPreview.svelte"
import { prebuiltDictionaryEditorScreenStates } from "./prebuilt-dictionary-editor-screen-specs"

const callbacks = {
  onFilterInput: () => {}, onSearch: () => {}, onCreate: () => {}, onDelete: () => {}, onStartEdit: () => {}, onCancelRow: () => {}, onRowInput: () => {}, onToggleCategories: () => {},
  onConfirmChanges: () => {}, onCancelChanges: () => {}, onPrev: () => {}, onNext: () => {}
}

function renderState(state: (typeof prebuiltDictionaryEditorScreenStates)[keyof typeof prebuiltDictionaryEditorScreenStates]) {
  return render(PrebuiltDictionaryEditorScreen, { ...callbacks, ...state.args })
}

export const prebuiltDictionaryEditorScreenChecks = [
  { id: "prebuilt-dictionary.list.table", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.list); expect(screen.getByRole("table")).toBeInTheDocument() } },
  { id: "prebuilt-dictionary.list.pagination", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.list); expect(screen.getByText("1ページ 50件")).toBeInTheDocument() } },
  { id: "prebuilt-dictionary.row-edit.selected", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.list); expect(screen.getAllByRole("textbox")).toHaveLength(4) } },
  { id: "prebuilt-dictionary.edit.form", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.edit); expect(screen.getByRole("button", { name: "変更を確定" })).toBeInTheDocument() } },
  { id: "prebuilt-dictionary.edit.row-cancel", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.edit); expect(screen.getAllByRole("button", { name: "取消" })).toHaveLength(2); expect(screen.getByRole("button", { name: "全て取消" })).toBeInTheDocument() } },
  { id: "prebuilt-dictionary.row-edit.deleted", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.edit); expect(screen.getAllByRole("button", { name: "編集" })).toHaveLength(1) } },
  { id: "prebuilt-dictionary.empty.message", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.empty); expect(screen.getByText("一致する辞書はありません。 ".trim())).toBeInTheDocument() } }
  ,{ id: "prebuilt-dictionary.long-text.editable", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.longText); expect(screen.getByDisplayValue(/TheUnfathomablyLongName/)).toBeInTheDocument() } }
]

const prebuiltDictionaryEditorSpecIds = Object.values(prebuiltDictionaryEditorScreenStates).flatMap((state) => state.specs.map((spec) => spec.id))

runScreenSpecHarness("用語辞書", prebuiltDictionaryEditorSpecIds, prebuiltDictionaryEditorScreenChecks)

it("削除保留の取消は先行する編集保留を残す", async () => {
  render(PrebuiltDictionaryEditorPreview)
  await fireEvent.click(screen.getAllByRole("button", { name: "編集" })[0])
  const destination = screen.getByDisplayValue("同胞団")
  await fireEvent.input(destination, { target: { value: "同胞団改" } })
  const editedRow = screen.getByDisplayValue("同胞団改").closest("tr")
  expect(editedRow).not.toBeNull()
  await fireEvent.click(within(editedRow as HTMLTableRowElement).getByRole("button", { name: "削除" }))
  const deletedRow = screen.getByText("同胞団改").closest("tr")
  expect(deletedRow).not.toBeNull()
  await fireEvent.click(within(deletedRow as HTMLTableRowElement).getByRole("button", { name: "取消" }))
  const restoredRow = screen.getByText("同胞団改").closest("tr")
  expect(restoredRow).toHaveClass("bg-warning/15")
})

it("別の行を編集すると前の行を文字表示へ戻す", async () => {
  render(PrebuiltDictionaryEditorPreview)
  const editButtons = screen.getAllByRole("button", { name: "編集" })
  await fireEvent.click(editButtons[0])
  expect(screen.getByDisplayValue("同胞団")).toBeInTheDocument()
  await fireEvent.click(screen.getAllByRole("button", { name: "編集" })[0])
  expect(screen.queryByDisplayValue("同胞団")).not.toBeInTheDocument()
  expect(screen.getByDisplayValue("仲間")).toBeInTheDocument()
})
