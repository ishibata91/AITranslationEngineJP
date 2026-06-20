// Wails generated bindings のラッパ（プロンプトテンプレート編集）。
// generated wailsjs の import はこの gateway 境界にだけ閉じ込める。
// View・Container からは本 gateway 経由で backend を呼ぶ。
import { GetPromptTemplate, SavePromptTemplate } from "../../wailsjs/go/api/App"
import { api } from "../../wailsjs/go/models"

// 編集対象のプロンプトテンプレート。base 翻訳指示文と口調指示テンプレート。
export interface PromptTemplate {
  baseDirective: string
  personaTemplate: string
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
