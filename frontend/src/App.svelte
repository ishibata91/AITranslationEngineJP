<script lang="ts">
  // production の root。ナビは「翻訳対象プラグイン」と「プロンプトテンプレート」の 2 タブ。
  // 実行・結果はトップタブでなく、翻訳対象プラグイン画面でプラグインを選ぶ／一覧の行を開いた先の画面にする。
  // 現在画面と選択中プラグインのパスだけを state に持ち、表示は AppShell と各 container に委ねる。
  import AppShell from "./ui/components/AppShell.svelte"
  import type { AppRoute, NavItem } from "./ui/components/app-nav-view"
  import TargetPluginsContainer from "./ui/screens/target-plugins/TargetPluginsContainer.svelte"
  import TranslationRunContainer from "./ui/screens/translation-run/TranslationRunContainer.svelte"
  import TemplateEditorContainer from "./ui/screens/template-editor/TemplateEditorContainer.svelte"

  const NAV_ITEMS: NavItem[] = [
    { route: "plugins", label: "翻訳対象プラグイン" },
    { route: "template", label: "プロンプトテンプレート" }
  ]

  let active = $state<AppRoute>("plugins")
  // 実行・結果画面へ渡す、選ばれたプラグインのフルパス。
  let selectedPluginPath = $state("")

  function onNavigate(route: AppRoute) {
    active = route
  }

  // 翻訳対象プラグイン画面でプラグインを選ぶ／一覧の行を開くと、実行・結果画面へ進む。
  function openRun(pluginPath: string) {
    selectedPluginPath = pluginPath
    active = "translation"
  }
</script>

<AppShell {active} items={NAV_ITEMS} {onNavigate}>
  {#if active === "plugins"}
    <TargetPluginsContainer onOpen={openRun} />
  {:else if active === "translation"}
    <TranslationRunContainer pluginPath={selectedPluginPath} />
  {:else}
    <TemplateEditorContainer />
  {/if}
</AppShell>
