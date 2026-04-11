<script lang="ts">
  import { onMount } from "svelte"

  import type {
    ShellRouteContract,
    ShellRouteId
  } from "@ui/stores/shell-state"

  interface Props {
    defaultRouteId: ShellRouteId
    routes: ShellRouteContract[]
  }

  let { defaultRouteId, routes }: Props = $props()

  const PLACEHOLDER_LEAD =
    "このページはまだ準備中です。上のナビゲーションまたは下の移動から別の主要ページへ進めます。"

  const routeById = $derived(
    new Map(routes.map((route) => [route.id, route] as const))
  )
  let currentRouteId = $state<ShellRouteId>("dashboard")
  const currentRoute = $derived(
    routeById.get(currentRouteId) ?? routes[0] ?? { id: "dashboard", label: "ダッシュボード" }
  )
  const isDashboard = $derived(currentRoute.id === "dashboard")

  function normalizeRouteId(hashValue: string): ShellRouteId {
    const routeId = hashValue.replace(/^#/, "") as ShellRouteId
    return routeById.has(routeId) ? routeId : defaultRouteId
  }

  function syncRouteFromHash(): void {
    if (typeof window === "undefined") {
      currentRouteId = defaultRouteId
      return
    }

    const nextRouteId = normalizeRouteId(window.location.hash)
    if (window.location.hash !== `#${nextRouteId}`) {
      window.history.replaceState(null, "", `#${nextRouteId}`)
    }
    currentRouteId = nextRouteId
  }

  function selectRoute(routeId: ShellRouteId): void {
    currentRouteId = routeId
  }

  onMount(() => {
    syncRouteFromHash()

    const onHashChange = (): void => {
      syncRouteFromHash()
    }

    window.addEventListener("hashchange", onHashChange)
    return () => {
      window.removeEventListener("hashchange", onHashChange)
    }
  })
</script>

<main class="shell">
  <header class="shell-bar">
    <p class="eyebrow">AITranslationEngineJp</p>
    <nav aria-label="グローバルナビゲーション" class="global-nav">
      {#each routes as route (route.id)}
        <a
          aria-current={route.id === currentRoute.id ? "page" : undefined}
          class="nav-link"
          href={`#${route.id}`}
          onclick={() => selectRoute(route.id)}
        >
          {route.label}
        </a>
      {/each}
    </nav>
    <p class="location-chip"><span>現在地</span><strong>{currentRoute.label}</strong></p>
  </header>

  <section class="hero panel">
    <h1>{currentRoute.label}</h1>
    <p class="description">
      ダッシュボードと共通ナビゲーションを起点に、主要ページへの移動導線を提供します。
    </p>
  </section>

  {#if isDashboard}
    <section class="dashboard panel">
      <h2>作業を選ぶ</h2>
      <div class="actions">
        {#each routes as route (route.id)}
          <a class="action-link" href={`#${route.id}`} onclick={() => selectRoute(route.id)}>
            {route.label}
          </a>
        {/each}
      </div>
    </section>
  {/if}

  {#if !isDashboard}
    <section class="placeholder panel">
      <p class="placeholder-state">{currentRoute.placeholderState ?? "準備中"}</p>
      <p class="placeholder-title">{currentRoute.label}</p>
      <p>{PLACEHOLDER_LEAD}</p>
      <div class="actions">
        {#each routes as route (route.id)}
          <a class="action-link" href={`#${route.id}`} onclick={() => selectRoute(route.id)}>
            {route.label}
          </a>
        {/each}
      </div>
    </section>
  {/if}
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: "Noto Sans JP", "Hiragino Sans", sans-serif;
    color: #f2e5de;
    background:
      radial-gradient(circle at top left, rgba(255, 186, 56, 0.14), transparent 30%),
      linear-gradient(180deg, #1c1715 0%, #161311 100%);
  }

  .shell {
    min-height: 100vh;
    max-width: 1120px;
    margin: 0 auto;
    padding: 24px;
    display: grid;
    gap: 16px;
  }

  .panel {
    padding: 20px;
    border-radius: 16px;
    border: 1px solid rgba(255, 186, 56, 0.2);
    background: rgba(32, 26, 24, 0.8);
    backdrop-filter: blur(16px);
  }

  .eyebrow {
    margin: 0;
    font-size: 0.8rem;
    letter-spacing: 0.2em;
    text-transform: uppercase;
    color: #e2c7a9;
  }

  .shell-bar {
    display: grid;
    gap: 12px;
  }

  .global-nav {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .nav-link {
    padding: 8px 12px;
    border-radius: 999px;
    color: #f2e5de;
    border: 1px solid rgba(255, 186, 56, 0.22);
    text-decoration: none;
  }

  .location-chip {
    margin: 0;
    display: inline-flex;
    gap: 8px;
    align-items: center;
    color: #e2c7a9;
  }

  .location-chip strong {
    color: #f2e5de;
  }

  h1,
  h2,
  p {
    margin: 0;
  }

  .hero {
    display: grid;
    gap: 10px;
  }

  h1 {
    font-size: clamp(2rem, 6vw, 3.2rem);
    line-height: 1.1;
  }

  .description {
    color: #e2c7a9;
    line-height: 1.6;
  }

  .dashboard,
  .placeholder {
    display: grid;
    gap: 14px;
  }

  .placeholder-state {
    font-size: 0.85rem;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #ffba38;
  }

  .placeholder-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: #f2e5de;
  }

  .actions {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 10px;
  }

  .action-link {
    padding: 10px 12px;
    border-radius: 10px;
    border: 1px solid rgba(255, 186, 56, 0.24);
    color: #f2e5de;
    text-decoration: none;
  }
</style>
