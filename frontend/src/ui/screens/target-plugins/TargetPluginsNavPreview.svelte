<script lang="ts">
  // Storybook プレビュー専用。翻訳対象プラグイン画面の最上部に画面間ナビゲーションが乗った状態を見せる。
  // 本番のルーティング（タブ押下での画面切り替え、行を選んで実行・結果画面へ進む）は implementation-module が組む。
  // ここでは「翻訳対象プラグインがナビの入口にあり、他画面へ飛べる」という表示の意図だけ固定する。
  import AppShell from "@ui/components/AppShell.svelte"
  import type { NavItem } from "@ui/components/app-nav-view"
  import TargetPluginsScreen from "./TargetPluginsScreen.svelte"
  import { LIST_STATE } from "./target-plugins.fixtures"

  const items: NavItem[] = [
    { route: "plugins", label: "翻訳対象プラグイン" },
    { route: "template", label: "プロンプトテンプレート" }
  ]
  const noop = () => {}
</script>

<AppShell active="plugins" {items} onNavigate={noop}>
  <TargetPluginsScreen
    {...LIST_STATE}
    onSelectNewPlugin={noop}
    onNewPluginPathInput={noop}
    onProceedToRun={noop}
    onOpenPlugin={noop}
    onRequestDelete={noop}
    onConfirmDelete={noop}
    onCancelDelete={noop}
  />
</AppShell>
