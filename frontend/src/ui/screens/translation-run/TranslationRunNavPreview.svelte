<script lang="ts">
  // Storybook プレビュー専用。翻訳実行画面の最上部に画面間ナビゲーションが乗った状態を見せる。
  // 本番のルーティング（タブ押下での画面切り替え）は implementation-module が組む。
  // ここでは「翻訳実行画面からテンプレート編集へ飛ぶ導線がある」という表示の意図だけ固定する。
  import AppShell from "@ui/components/AppShell.svelte"
  import type { NavItem } from "@ui/components/app-nav-view"
  import TranslationRunScreen from "./TranslationRunScreen.svelte"
  import { DONE_PERSONA_STATE } from "./translation-run.fixtures"

  const items: NavItem[] = [
    { route: "translation", label: "翻訳実行" },
    { route: "template", label: "プロンプトテンプレート" }
  ]
  const noop = () => {}
</script>

<AppShell active="translation" {items} onNavigate={noop}>
  <TranslationRunScreen
    {...DONE_PERSONA_STATE}
    onFieldInput={noop}
    onLoadModels={noop}
    onRun={noop}
    onPagePrev={noop}
    onPageNext={noop}
  />
</AppShell>
