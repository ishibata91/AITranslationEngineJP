import { render, screen } from "@testing-library/svelte"
import { expect } from "vitest"
import { runScreenSpecHarness } from "../../../test/screen-spec-harness"
import PrebuiltDictionaryEditorScreen from "./PrebuiltDictionaryEditorScreen.svelte"
import { prebuiltDictionaryEditorScreenStates } from "./prebuilt-dictionary-editor-screen-specs"

const callbacks = {
  onFilterInput: () => {}, onCreate: () => {}, onDelete: () => {}, onToggleCategories: () => {},
  onConfirmChanges: () => {}, onCancelChanges: () => {}, onPrev: () => {}, onNext: () => {}
}

function renderState(state: (typeof prebuiltDictionaryEditorScreenStates)[keyof typeof prebuiltDictionaryEditorScreenStates]) {
  return render(PrebuiltDictionaryEditorScreen, { ...callbacks, ...state.args })
}

export const prebuiltDictionaryEditorScreenChecks = [
  { id: "prebuilt-dictionary.list.table", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.list); expect(screen.getByRole("table")).toBeInTheDocument() } },
  { id: "prebuilt-dictionary.list.pagination", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.list); expect(screen.getByText("1ページ 50件")).toBeInTheDocument() } },
  { id: "prebuilt-dictionary.edit.form", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.edit); expect(screen.getByRole("button", { name: "変更を確定" })).toBeInTheDocument() } },
  { id: "prebuilt-dictionary.empty.message", verify: () => { renderState(prebuiltDictionaryEditorScreenStates.empty); expect(screen.getByText("一致する辞書はありません。 ".trim())).toBeInTheDocument() } }
]

const prebuiltDictionaryEditorSpecIds = Object.values(prebuiltDictionaryEditorScreenStates).flatMap((state) => state.specs.map((spec) => spec.id))

runScreenSpecHarness("用語辞書", prebuiltDictionaryEditorSpecIds, prebuiltDictionaryEditorScreenChecks)
