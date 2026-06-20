<script lang="ts">
  // 画面間ナビゲーションバー。翻訳実行とテンプレート編集を切り替える導線を出す。
  // 遷移ロジックは持たない。現在画面は active props、タブ押下は onNavigate callback prop で上へ返す。
  import type { AppRoute, NavItem } from "./app-nav-view"

  interface Props {
    // 現在表示している画面。タブの強調に使う。
    active: AppRoute
    // 並べるナビ項目。
    items: NavItem[]
    onNavigate: (route: AppRoute) => void
  }

  let { active, items, onNavigate }: Props = $props()
</script>

<nav class="w-full border-b border-base-300/60 bg-base-200/40 backdrop-blur">
  <div class="mx-auto flex h-14 max-w-5xl items-center px-6">
    <div class="flex items-center gap-1">
      {#each items as item (item.route)}
        <button
          type="button"
          class="rounded-box px-3 py-1.5 text-sm transition-colors {item.route === active
            ? 'bg-primary/15 text-primary'
            : 'text-base-content/60 hover:bg-base-100/50 hover:text-base-content'}"
          aria-current={item.route === active ? "page" : undefined}
          onclick={() => onNavigate(item.route)}
        >
          {item.label}
        </button>
      {/each}
    </div>
  </div>
</nav>
