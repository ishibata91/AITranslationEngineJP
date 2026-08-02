import { defineScreenState } from "../screen-spec"
import {
  DEFAULT_TEMPLATE_FORM,
  DIRECTIVES,
  EDITED_DIRECTIVES,
  RECORD_ASSIGNMENTS,
  TONE_DEFAULT_EDITED_FORM
} from "./template-editor.fixtures"

const sharedRecordArgs = {
  activeTab: "record" as const,
  assignments: RECORD_ASSIGNMENTS,
  saving: false
}

export const templateEditorScreenStates = {
  baseTab: defineScreenState({
    storyName: "ベースタブ",
    precondition: "ベースtabを開いており、未保存の変更がない。",
    args: {
      form: DEFAULT_TEMPLATE_FORM,
      activeTab: "base" as const,
      directives: DIRECTIVES,
      assignments: RECORD_ASSIGNMENTS,
      dirty: false,
      saving: false
    },
    specs: [
      {
        id: "template-editor.base.content",
        statement: "ベースtabを選択状態にし、ベースの指示文を表示する。"
      },
      {
        id: "template-editor.base.unsaved-hidden",
        statement: "未保存の案内を表示しない。"
      },
      {
        id: "template-editor.base.actions-disabled",
        statement: "「戻す」と「保存」が無効になる。"
      }
    ]
  }),
  recordTab: defineScreenState({
    storyName: "レコード別タブ",
    precondition: "レコード別tabを開いており、未保存の変更がない。",
    args: {
      form: DEFAULT_TEMPLATE_FORM,
      ...sharedRecordArgs,
      directives: DIRECTIVES,
      dirty: false
    },
    specs: [
      {
        id: "template-editor.record.content",
        statement:
          "レコード別tabを選択状態にし、口調とレコード別の指示文を表示する。"
      },
      {
        id: "template-editor.record.unsaved-hidden",
        statement: "未保存の案内を表示しない。"
      },
      {
        id: "template-editor.record.actions-disabled",
        statement: "「戻す」と「保存」が無効になる。"
      }
    ]
  }),
  recordTabToneDefaultEdited: defineScreenState({
    storyName: "レコード別タブ（口調と PC 性別を変更）",
    precondition: "レコード別tabで口調とPC性別を変更し、保存していない。",
    args: {
      form: TONE_DEFAULT_EDITED_FORM,
      ...sharedRecordArgs,
      directives: DIRECTIVES,
      dirty: true
    },
    specs: [
      {
        id: "template-editor.tone-edited.unsaved",
        statement: "「未保存の変更」を表示する。"
      },
      {
        id: "template-editor.tone-edited.actions-enabled",
        statement: "「戻す」と「保存」が有効になる。"
      }
    ]
  }),
  recordTabDirty: defineScreenState({
    storyName: "レコード別タブ（未保存）",
    precondition: "レコード別tabで指示文を変更し、保存していない。",
    args: {
      form: DEFAULT_TEMPLATE_FORM,
      ...sharedRecordArgs,
      directives: EDITED_DIRECTIVES,
      dirty: true
    },
    specs: [
      {
        id: "template-editor.directive-edited.unsaved",
        statement: "「未保存の変更」を表示する。"
      },
      {
        id: "template-editor.directive-edited.actions-enabled",
        statement: "「戻す」と「保存」が有効になる。"
      }
    ]
  })
} as const
