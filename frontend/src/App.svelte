<script lang="ts">
  // production の root。最上部の画面間ナビゲーションで翻訳実行とテンプレート編集を切り替える。
  // 現在画面の state だけを持ち、表示は AppShell（ナビ＋現在画面）と各 container に委ねる。
  import AppShell from "./ui/components/AppShell.svelte"
  import type { AppRoute, NavItem } from "./ui/components/app-nav-view"
  import TranslationRunContainer from "./ui/screens/translation-run/TranslationRunContainer.svelte"
  import TemplateEditorContainer from "./ui/screens/template-editor/TemplateEditorContainer.svelte"

  const NAV_ITEMS: NavItem[] = [
    { route: "translation", label: "翻訳実行" },
    { route: "template", label: "プロンプトテンプレート" }
  ]

  let active = $state<AppRoute>("translation")

  function onNavigate(route: AppRoute) {
    active = route
  }
</script>

<AppShell {active} items={NAV_ITEMS} {onNavigate}>
  {#if active === "translation"}
    <TranslationRunContainer />
  {:else}
    <TemplateEditorContainer />
  {/if}
</AppShell>
