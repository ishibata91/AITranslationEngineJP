// テンプレート編集画面の表示用 fixture。現状コードが持つ既定テンプレートを表示値として固定する。
import type { PromptTemplateForm, TemplatePlaceholder } from "./template-editor-view"

// base 翻訳指示文。現状は provider のハードコード（openai_compatible.go:152 相当）。
const DEFAULT_BASE = `あなたは Skyrim Mod の翻訳者です。与えられた英語の本文を、原文の意味と語調を保った自然な日本語へ翻訳してください。訳文だけを出力し、説明や注釈は加えないでください。`

// 口調指示テンプレート。現状は engine のハードコード（persona.go の buildPersonaDirective 相当）。
// {traits} に話者の性質（声質・種族の気質・所属の気風）の箇条書きを差し込む。
const DEFAULT_PERSONA = `この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。`

export const DEFAULT_TEMPLATE_FORM: PromptTemplateForm = {
  baseDirective: DEFAULT_BASE,
  personaTemplate: DEFAULT_PERSONA
}

// 口調指示を強めへ書き換えた、未保存の編集中フォーム。dirty 表示の確認に使う。
export const EDITED_TEMPLATE_FORM: PromptTemplateForm = {
  baseDirective: DEFAULT_BASE,
  personaTemplate: `あなたはこの話者になりきって訳します。
話者の人物像:
{traits}
人物像の口調・人称・語尾を訳文へ強く反映し、説明的に崩さないこと。`
}

export const PLACEHOLDERS: TemplatePlaceholder[] = [
  {
    token: "{traits}",
    description:
      "話者の性質（声質・種族の気質・所属の気風）を箇条書きで差し込む。話者の属性（種族・声型）から定まる性質文が入る。"
  }
]
