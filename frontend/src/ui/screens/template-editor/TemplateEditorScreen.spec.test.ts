import { render, screen } from "@testing-library/svelte"
import { expect } from "vitest"
import { runScreenSpecHarness } from "../../../test/screen-spec-harness"
import TemplateEditorScreen from "./TemplateEditorScreen.svelte"
import { templateEditorScreenStates } from "./template-editor-screen-specs"

const callbacks = {
  onFieldInput: () => {},
  onInstructionInput: () => {},
  onTabChange: () => {},
  onSave: () => {},
  onReset: () => {}
}

function renderState(
  state: (typeof templateEditorScreenStates)[keyof typeof templateEditorScreenStates]
) {
  return render(TemplateEditorScreen, { ...callbacks, ...state.args })
}

function expectActionsDisabled(disabled: boolean): void {
  const reset = screen.getByRole("button", { name: "戻す" })
  const save = screen.getByRole("button", { name: "保存" })
  if (disabled) {
    expect(reset).toBeDisabled()
    expect(save).toBeDisabled()
    return
  }
  expect(reset).toBeEnabled()
  expect(save).toBeEnabled()
}

export const templateEditorScreenChecks = [
  {
    id: "template-editor.base.content",
    verify: () => {
      renderState(templateEditorScreenStates.baseTab)
      expect(screen.getByRole("tab", { name: "ベース" })).toHaveClass(
        "tab-active"
      )
      expect(
        screen.getByRole("heading", { name: "base 翻訳指示" })
      ).toBeInTheDocument()
      expect(screen.getByRole("textbox")).toHaveValue(
        templateEditorScreenStates.baseTab.args.form.baseDirective
      )
    }
  },
  {
    id: "template-editor.base.unsaved-hidden",
    verify: () => {
      renderState(templateEditorScreenStates.baseTab)
      expect(screen.queryByText("未保存の変更")).not.toBeInTheDocument()
    }
  },
  {
    id: "template-editor.base.actions-disabled",
    verify: () => {
      renderState(templateEditorScreenStates.baseTab)
      expectActionsDisabled(true)
    }
  },
  {
    id: "template-editor.record.content",
    verify: () => {
      renderState(templateEditorScreenStates.recordTab)
      expect(screen.getByRole("tab", { name: "レコード別" })).toHaveClass(
        "tab-active"
      )
      expect(
        screen.getByRole("heading", { name: "話者なし台詞の口調" })
      ).toBeInTheDocument()
      expect(
        screen.getByRole("heading", { name: "種別ごとの指示文" })
      ).toBeInTheDocument()
    }
  },
  {
    id: "template-editor.record.unsaved-hidden",
    verify: () => {
      renderState(templateEditorScreenStates.recordTab)
      expect(screen.queryByText("未保存の変更")).not.toBeInTheDocument()
    }
  },
  {
    id: "template-editor.record.actions-disabled",
    verify: () => {
      renderState(templateEditorScreenStates.recordTab)
      expectActionsDisabled(true)
    }
  },
  {
    id: "template-editor.tone-edited.unsaved",
    verify: () => {
      renderState(templateEditorScreenStates.recordTabToneDefaultEdited)
      expect(screen.getByText("未保存の変更")).toBeInTheDocument()
    }
  },
  {
    id: "template-editor.tone-edited.actions-enabled",
    verify: () => {
      renderState(templateEditorScreenStates.recordTabToneDefaultEdited)
      expectActionsDisabled(false)
    }
  },
  {
    id: "template-editor.directive-edited.unsaved",
    verify: () => {
      renderState(templateEditorScreenStates.recordTabDirty)
      expect(screen.getByText("未保存の変更")).toBeInTheDocument()
    }
  },
  {
    id: "template-editor.directive-edited.actions-enabled",
    verify: () => {
      renderState(templateEditorScreenStates.recordTabDirty)
      expectActionsDisabled(false)
    }
  }
]

const templateEditorSpecIds = Object.values(templateEditorScreenStates).flatMap(
  (state) => state.specs.map((spec) => spec.id)
)

runScreenSpecHarness(
  "プロンプトテンプレート",
  templateEditorSpecIds,
  templateEditorScreenChecks
)
