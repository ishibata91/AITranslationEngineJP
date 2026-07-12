<script lang="ts">
  // 翻訳対象プラグイン画面の container。画面状態（$state）を保持し、gateway 経由で backend を呼ぶ。
  // 表示は TargetPluginsScreen に委ね、本 component は state と配線だけを持つ。
  // プラグインを選ぶ／一覧の行を開くと、onOpen でフルパスを上（ルーティング）へ返し、実行・結果画面へ進む。
  import { onMount } from "svelte"
  import TargetPluginsScreen from "./TargetPluginsScreen.svelte"
  import type { TargetPluginRow } from "./target-plugins-view"
  import {
    selectPluginFile,
    listTargetPlugins,
    deleteTargetPlugin
  } from "../../../gateway/translation-gateway"

  // 一覧の元データ 1 行の型。gateway の戻り値から導出する（値と型を同一 import で混ぜないため型名は import しない）。
  type Summary = Awaited<ReturnType<typeof listTargetPlugins>>[number]

  // 選んだプラグインのフルパスで実行・結果画面へ進む。ルーティングは App が持つ。
  let { onOpen }: { onOpen: (pluginPath: string) => void } = $props()

  // 新規翻訳の入口で選んだプラグインのフルパス。
  let newPluginPath = $state("")
  // 一覧の元データ（フルパス付き）。行の表示は sourcePath を落として TargetPluginRow へ写す。
  let summaries = $state<Summary[]>([])
  let loading = $state(false)
  // 削除確認中のプラグイン（未設定なら通常の削除ボタン）。
  let pendingDeletePlugin = $state<string | undefined>(undefined)
  let deleting = $state(false)
  let errorMessage = $state("")

  // 表示行。sourcePath は画面に出さないため落とす。作成順（新しい順）は backend が保証する。
  const plugins: TargetPluginRow[] = $derived(
    summaries.map((s) => ({
      plugin: s.plugin,
      createdLabel: s.createdLabel,
      total: s.total,
      translated: s.translated
    }))
  )

  async function refresh() {
    loading = true
    errorMessage = ""
    try {
      summaries = await listTargetPlugins()
    } catch (error) {
      errorMessage = messageOf(error)
    } finally {
      loading = false
    }
  }

  async function onSelectNewPlugin() {
    try {
      const path = await selectPluginFile()
      // ダイアログを閉じただけ（空）のときは変更しない。
      if (path.length > 0) newPluginPath = path
    } catch (error) {
      errorMessage = messageOf(error)
    }
  }

  function onNewPluginPathInput(value: string) {
    newPluginPath = value
  }

  function onProceedToRun() {
    if (newPluginPath.length > 0) onOpen(newPluginPath)
  }

  // 一覧の行を開く。登録時のフルパス（source_path）を復元して実行・結果画面へ渡す（結果表示と再実行に使う）。
  function onOpenPlugin(plugin: string) {
    const found = summaries.find((s) => s.plugin === plugin)
    onOpen(found?.sourcePath ?? plugin)
  }

  function onRequestDelete(plugin: string) {
    pendingDeletePlugin = plugin
  }

  function onCancelDelete() {
    pendingDeletePlugin = undefined
  }

  async function onConfirmDelete(plugin: string) {
    deleting = true
    errorMessage = ""
    try {
      await deleteTargetPlugin(plugin)
      pendingDeletePlugin = undefined
      await refresh()
    } catch (error) {
      errorMessage = messageOf(error)
    } finally {
      deleting = false
    }
  }

  function messageOf(error: unknown): string {
    if (error instanceof Error) return error.message
    if (typeof error === "string") return error
    return "予期しないエラーが発生しました。"
  }

  onMount(() => {
    void refresh()
  })
</script>

<TargetPluginsScreen
  {newPluginPath}
  {onSelectNewPlugin}
  {onNewPluginPathInput}
  {onProceedToRun}
  {plugins}
  {loading}
  {pendingDeletePlugin}
  {deleting}
  {errorMessage}
  {onOpenPlugin}
  {onRequestDelete}
  {onConfirmDelete}
  {onCancelDelete}
/>
