<script lang="ts">
  // プロンプトテンプレート編集画面の container。画面状態（$state）を保持し、gateway 経由で backend を呼ぶ。
  // 表示は TemplateEditorScreen に委ね、本 component は state と配線だけを持つ。
  // ベースタブは Base 指示、レコード別タブは種別ごとの指示文（口調・文体・固有名・定型句）を編集する。
  // 保存失敗時は保存済みスナップショットを更新しないため dirty が残り、利用者は再保存できる（暗黙のフィードバック）。
  import { onMount } from "svelte"
  import TemplateEditorScreen from "./TemplateEditorScreen.svelte"
  import type {
    PromptTemplateForm,
    PromptTemplateField,
    TemplateTab
  } from "./template-editor-view"
  import type { Directive, RecordAssignment } from "./directive-view"
  import {
    getPromptTemplate,
    savePromptTemplate,
    getDirectiveEditing,
    saveDirective
  } from "../../../gateway/template-gateway"
  import { createPinoFrontendDiagnosticLogger } from "../../../application/diagnostic"

  const log = createPinoFrontendDiagnosticLogger({ screen: "template-editor" })

  // 表示中のサブタブ。
  let activeTab = $state<TemplateTab>("base")

  // 編集中の base 値。personaTemplate は口調 directive へ畳んだため編集せず、保存時の素通し用に保持する。
  let baseDirective = $state("")
  let personaTemplate = $state("")
  let savedBase = $state("")

  // 編集中の指示文（レコード別タブ）。saved は dirty 判定と「戻す」の基準。
  let directives = $state<Directive[]>([])
  let savedDirectives = $state<Directive[]>([])
  // REC:FIELD → 指示文キーの割り当て（読み取り専用）。
  let assignments = $state<RecordAssignment[]>([])

  let saving = $state(false)

  const form: PromptTemplateForm = $derived({ baseDirective, personaTemplate })
  const baseDirty = $derived(baseDirective !== savedBase)
  // 指示文の dirty は、いずれかの key の instruction が保存済みと異なるか。
  const directivesDirty = $derived(
    directives.some((d) => d.instruction !== savedInstructionOf(d.key))
  )
  const dirty = $derived(baseDirty || directivesDirty)

  function savedInstructionOf(key: string): string {
    return savedDirectives.find((d) => d.key === key)?.instruction ?? ""
  }

  function onFieldInput(field: PromptTemplateField, value: string) {
    if (field === "baseDirective") baseDirective = value
    else if (field === "personaTemplate") personaTemplate = value
  }

  function onInstructionInput(key: string, instruction: string) {
    directives = directives.map((d) =>
      d.key === key ? { ...d, instruction } : d
    )
  }

  function onTabChange(tab: TemplateTab) {
    activeTab = tab
  }

  // 指示文配列の深いコピー（編集と保存済みスナップショットを別実体にする）。
  function cloneDirectives(rows: Directive[]): Directive[] {
    return rows.map((d) => ({
      key: d.key,
      instruction: d.instruction,
      variables: d.variables.map((v) => ({ ...v }))
    }))
  }

  async function load() {
    try {
      const template = await getPromptTemplate()
      baseDirective = template.baseDirective
      personaTemplate = template.personaTemplate
      savedBase = template.baseDirective
    } catch (error) {
      log.error("プロンプトテンプレートの読み込みに失敗", { reason: messageOf(error) })
    }
    try {
      const editing = await getDirectiveEditing()
      directives = cloneDirectives(editing.directives)
      savedDirectives = cloneDirectives(editing.directives)
      assignments = editing.assignments
    } catch (error) {
      log.error("指示文の読み込みに失敗", { reason: messageOf(error) })
    }
  }

  async function onSave() {
    if (!dirty || saving) return
    saving = true
    try {
      if (baseDirty) {
        await savePromptTemplate({ baseDirective, personaTemplate })
        savedBase = baseDirective
      }
      // 変更された指示文だけを保存する（割り当ては固定のため触らない）。
      for (const d of directives) {
        if (d.instruction !== savedInstructionOf(d.key)) {
          await saveDirective(d.key, d.instruction)
        }
      }
      savedDirectives = cloneDirectives(directives)
    } catch (error) {
      log.error("プロンプトテンプレートの保存に失敗", { reason: messageOf(error) })
    } finally {
      saving = false
    }
  }

  // 編集を保存済みの値へ戻す（base と指示文の両方）。
  function onReset() {
    baseDirective = savedBase
    directives = cloneDirectives(savedDirectives)
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
  {activeTab}
  {directives}
  {assignments}
  {dirty}
  {saving}
  {onFieldInput}
  {onInstructionInput}
  {onTabChange}
  {onSave}
  {onReset}
/>
