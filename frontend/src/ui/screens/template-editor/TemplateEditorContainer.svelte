<script lang="ts">
  // プロンプトテンプレート編集画面の container。画面状態（$state）を保持し、gateway 経由で backend を呼ぶ。
  // 表示は TemplateEditorScreen に委ね、本 component は state と配線だけを持つ。
  // 保存失敗時は保存済みスナップショットを更新しないため dirty が残り、利用者は再保存できる（暗黙のフィードバック）。
  import { onMount } from "svelte"
  import TemplateEditorScreen from "./TemplateEditorScreen.svelte"
  import type {
    PromptTemplateForm,
    PromptTemplateField,
    TemplatePlaceholder
  } from "./template-editor-view"
  import { getPromptTemplate, savePromptTemplate } from "../../../gateway/template-gateway"
  import { createPinoFrontendDiagnosticLogger } from "../../../application/diagnostic"

  // 口調指示テンプレートに差し込めるプレースホルダ。production の固定説明。{traits} は話者の性質列。
  const PLACEHOLDERS: TemplatePlaceholder[] = [
    {
      token: "{traits}",
      description:
        "話者の性質（声質・種族の気質・所属の気風）を箇条書きで差し込む。話者の属性（種族・声型）から定まる性質文が入る。"
    }
  ]

  const log = createPinoFrontendDiagnosticLogger({ screen: "template-editor" })

  // 編集中フォームの値。
  let baseDirective = $state("")
  let personaTemplate = $state("")
  // 保存済みの値。dirty 判定と「戻す」の基準にする。
  let savedBase = $state("")
  let savedPersona = $state("")
  let saving = $state(false)

  const form: PromptTemplateForm = $derived({ baseDirective, personaTemplate })
  const dirty = $derived(baseDirective !== savedBase || personaTemplate !== savedPersona)

  function onFieldInput(field: PromptTemplateField, value: string) {
    if (field === "baseDirective") baseDirective = value
    else if (field === "personaTemplate") personaTemplate = value
  }

  // 保存済みスナップショットを現在値へ更新する（読み込み完了時・保存成功時）。
  function syncSaved(base: string, persona: string) {
    baseDirective = base
    personaTemplate = persona
    savedBase = base
    savedPersona = persona
  }

  async function load() {
    try {
      const template = await getPromptTemplate()
      syncSaved(template.baseDirective, template.personaTemplate)
    } catch (error) {
      log.error("プロンプトテンプレートの読み込みに失敗", { reason: messageOf(error) })
    }
  }

  async function onSave() {
    if (!dirty || saving) return
    saving = true
    try {
      await savePromptTemplate({ baseDirective, personaTemplate })
      syncSaved(baseDirective, personaTemplate)
    } catch (error) {
      log.error("プロンプトテンプレートの保存に失敗", { reason: messageOf(error) })
    } finally {
      saving = false
    }
  }

  // 編集を保存済みの値へ戻す。
  function onReset() {
    baseDirective = savedBase
    personaTemplate = savedPersona
  }

  function messageOf(error: unknown): string {
    if (error instanceof Error) return error.message
    if (typeof error === "string") return error
    return "予期しないエラーが発生しました。"
  }

  onMount(() => {
    void load()
  })
</script>

<TemplateEditorScreen
  {form}
  placeholders={PLACEHOLDERS}
  {dirty}
  {saving}
  {onFieldInput}
  {onSave}
  {onReset}
/>
