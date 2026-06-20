<script lang="ts">
  // アプリの外枠。最上部に画面間ナビゲーションを置き、その下に現在画面を描く。
  // ナビの遷移ロジックは持たない。active と onNavigate は props で受け、children に現在画面を受ける。
  import AppNav from "./AppNav.svelte"
  import type { AppRoute, NavItem } from "./app-nav-view"
  import type { Snippet } from "svelte"

  interface Props {
    // 現在表示している画面。
    active: AppRoute
    // 並べるナビ項目。
    items: NavItem[]
    onNavigate: (route: AppRoute) => void
    // 現在画面の中身。
    children: Snippet
  }

  let { active, items, onNavigate, children }: Props = $props()
</script>

<div class="flex min-h-screen flex-col">
  <AppNav {active} {items} {onNavigate} />
  <div class="flex-1">
    {@render children()}
  </div>
</div>
