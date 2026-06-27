// Wails generated bindings のラッパ（プロンプトテンプレート編集）。
// generated wailsjs の import はこの gateway 境界にだけ閉じ込める。
// View・Container からは本 gateway 経由で backend を呼ぶ。
import {
  GetPromptTemplate,
  SavePromptTemplate,
  GetDirectiveEditing,
  SaveDirective
} from "../../wailsjs/go/api/App"
import { api } from "../../wailsjs/go/models"

// 編集対象のプロンプトテンプレート。base 翻訳指示文と口調指示テンプレート。
export interface PromptTemplate {
  baseDirective: string
  personaTemplate: string
}

// 指示文が実行時に埋める変数 1 件（例: {traits}）。
export interface TemplateVariableRow {
  token: string
  description: string
}

// 指示文 1 件（key・指示文・変数宣言）。レコード別タブで指示文を編集する単位。
export interface DirectiveRow {
  key: string
  instruction: string
  variables: TemplateVariableRow[]
}

// REC:FIELD → 指示文キーの割り当て 1 件（読み取り専用の対象一覧）。
export interface RecordAssignmentRow {
  recField: string
  logicalName: string
  directive: string
}

// レコード別タブの初期表示データ。指示文一覧と REC:FIELD 割り当て。
export interface DirectiveEditing {
  directives: DirectiveRow[]
  assignments: RecordAssignmentRow[]
}

// 現在保存されているプロンプトテンプレートを取得する（編集画面の初期表示用）。
export async function getPromptTemplate(): Promise<PromptTemplate> {
  const view = await GetPromptTemplate()
  return {
    baseDirective: view.baseDirective,
    personaTemplate: view.personaTemplate
  }
}

// 編集したプロンプトテンプレートを保存する。保存後の翻訳と実プロンプト参照へ反映される。
export async function savePromptTemplate(template: PromptTemplate): Promise<void> {
  await SavePromptTemplate(api.PromptTemplateView.createFrom(template))
}

// レコード別タブの初期表示用に、指示文一覧と REC:FIELD 割り当てを取得する。
export async function getDirectiveEditing(): Promise<DirectiveEditing> {
  const view = await GetDirectiveEditing()
  return {
    directives: (view.directives ?? []).map((d) => ({
      key: d.key,
      instruction: d.instruction,
      variables: (d.variables ?? []).map((v) => ({
        token: v.token,
        description: v.description
      }))
    })),
    assignments: (view.assignments ?? []).map((a) => ({
      recField: a.recField,
      logicalName: a.logicalName,
      directive: a.directive
    }))
  }
}

// 編集した指示文を保存する。割り当て（record_type_master）は固定のため変えず、指示文だけ書き換える。
// 保存後の翻訳と実プロンプト参照へ反映される。
export async function saveDirective(key: string, instruction: string): Promise<void> {
  await SaveDirective(key, instruction)
}
