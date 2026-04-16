<script lang="ts">
  import { onMount } from "svelte"

  import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
  import type { CreateMasterPersonaScreenController } from "@application/contract/master-persona"
  import MasterDictionaryPage from "@ui/screens/master-dictionary/MasterDictionaryPage.svelte"
  import MasterPersonaPage from "@ui/screens/master-persona/MasterPersonaPage.svelte"
  import type { ShellRouteContract, ShellRouteId } from "@ui/stores/shell-state"

  interface Props {
    defaultRouteId: ShellRouteId
    routes: ShellRouteContract[]
    createMasterDictionaryScreenController: CreateMasterDictionaryScreenController | null
    createMasterPersonaScreenController: CreateMasterPersonaScreenController | null
  }

  let {
    defaultRouteId,
    routes,
    createMasterDictionaryScreenController,
    createMasterPersonaScreenController
  }: Props = $props()

  const PLACEHOLDER_LEAD =
    "このページはまだ準備中です。上のナビゲーションまたは下の移動から別の主要ページへ進めます。"

  const routeById = $derived(
    new Map(routes.map((route) => [route.id, route] as const))
  )

  let currentRouteId = $state<ShellRouteId>("dashboard")
  let isMobileNavOpen = $state(false)

  const fallbackRoute: ShellRouteContract = {
    id: "dashboard",
    label: "ダッシュボード",
    state: "既定表示",
    lead: "最初に移動したい作業を選び、共通ナビゲーションからいつでも別の主要ページへ切り替えられます。",
    description: "主要ページへの入口をまとめて確認します。"
  }

  const currentRoute = $derived(
    routeById.get(currentRouteId) ?? routes[0] ?? fallbackRoute
  )
  const isDashboard = $derived(currentRoute.id === "dashboard")
  const dashboardEntryRoutes = $derived(
  
    routes.filter((route) => route.id !== "dashboard")
  )

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
    isMobileNavOpen = false
  }

  function selectRoute(routeId: ShellRouteId): void {
    currentRouteId = routeId
    isMobileNavOpen = false
  }

  function toggleMobileNav(): void {
    isMobileNavOpen = !isMobileNavOpen
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
  <header class="shell-bar" class:is-open={isMobileNavOpen}>
    <div class="shell-bar-inner">
      <div class="brand">
        <p class="brand-eyebrow">AITranslationEngineJp</p>
        <strong>翻訳エンジン</strong>
      </div>
      <nav
        aria-label="グローバルナビゲーション"
        class="global-nav"
        id="globalNav"
      >
        {#each routes as route (route.id)}
          <a
            aria-current={route.id === currentRoute.id ? "page" : undefined}
            class="nav-link"
            class:is-active={route.id === currentRoute.id}
            href={`#${route.id}`}
            onclick={() => selectRoute(route.id)}
          >
            {route.label}
          </a>
        {/each}
      </nav>
      <div class="bar-status">
        <button
          aria-controls="globalNav"
          aria-expanded={isMobileNavOpen ? "true" : "false"}
          class="nav-toggle"
          onclick={toggleMobileNav}
          type="button"
        >
          主要ページ
        </button>
      </div>
    </div>
  </header>

  <section class="page">
    <section class="panel hero-panel">
      <div class="hero-top">
        <div>
          <p class="page-label">現在のページ</p>
          <h1>{currentRoute.label}</h1>
        </div>
      </div>
      <p class="hero-lead">{currentRoute.lead}</p>
    </section>

    {#if isDashboard}
      <section class="hero-grid" id="dashboardView">
        <section class="panel entry-panel">
          <div class="section-head">
            <div>
              <p class="page-label">主要ページ</p>
              <h2>作業を選ぶ</h2>
            </div>
          </div>
          <div class="entry-grid">
            {#each dashboardEntryRoutes as route (route.id)}
              <a
                class="entry-card"
                href={`#${route.id}`}
                onclick={() => selectRoute(route.id)}
              >
                <div class="entry-head">
                  <div>
                    <p class="card-tag">{route.id}</p>
                    <h3>{route.label}</h3>
                  </div>
                  <span class="entry-state">{route.state}</span>
                </div>
                <p>{route.description}</p>
                <span class="entry-action">開く</span>
              </a>
            {/each}
          </div>
        </section>
      </section>
    {/if}

    {#if !isDashboard && currentRoute.id === "master-dictionary"}
      <MasterDictionaryPage
        createController={createMasterDictionaryScreenController}
      />
    {/if}

    {#if !isDashboard && currentRoute.id === "master-persona"}
      <MasterPersonaPage
        createController={createMasterPersonaScreenController}
      />
    {/if}

    {#if !isDashboard && currentRoute.id !== "master-dictionary" && currentRoute.id !== "master-persona"}
      <section class="placeholder-content" id="placeholderView">
        <section class="panel placeholder-card">
          <p class="page-label">現在のページ</p>
          <h2>{currentRoute.label}</h2>
          <p>{PLACEHOLDER_LEAD}</p>
          <div class="action-grid">
            {#each routes as route (route.id)}
              <a
                class="action-link"
                href={`#${route.id}`}
                onclick={() => selectRoute(route.id)}
              >
                {route.label}
              </a>
            {/each}
          </div>
        </section>
      </section>
    {/if}
  </section>
</main>

<style>
  :global(body) {
    margin: 0;
    min-height: 100vh;
    color: var(--text);
    font-family: "Noto Serif JP", serif;
    background:
      radial-gradient(
        circle at top left,
        rgba(255, 186, 56, 0.16),
        transparent 28%
      ),
      radial-gradient(
        circle at 85% 18%,
        rgba(255, 104, 63, 0.14),
        transparent 24%
      ),
      linear-gradient(180deg, #1c1715 0%, var(--bg) 100%);
  }

  :global(*) {
    box-sizing: border-box;
  }

  :global(a) {
    color: inherit;
    text-decoration: none;
  }

  :global(button) {
    font: inherit;
  }

  :global(h1),
  :global(h2),
  :global(h3),
  :global(p) {
    margin: 0;
  }

  :global(:root) {
    --bg: #161311;
    --bg-strong: #110d0c;
    --surface: rgba(35, 31, 29, 0.78);
    --text: #eae1dd;
    --muted: #d8c3ae;
    --primary: #ffba38;
    --line: rgba(255, 186, 56, 0.18);
    --line-strong: rgba(255, 186, 56, 0.32);
    --shadow: 0 24px 64px rgba(0, 0, 0, 0.42);
    --radius-lg: 22px;
    --radius-md: 16px;
    --transition: 180ms ease;
  }

  .shell {
    min-height: 100vh;
  }

  .shell-bar {
    position: sticky;
    top: 0;
    z-index: 10;
    backdrop-filter: blur(38px);
    background: rgba(22, 19, 17, 0.78);
    border-bottom: 0.5px solid var(--line);
    box-shadow: 0 20px 50px rgba(0, 0, 0, 0.38);
  }

  .shell-bar-inner {
    max-width: 1440px;
    margin: 0 auto;
    padding: 18px 24px;
    display: grid;
    grid-template-columns: auto 1fr auto;
    gap: 18px;
    align-items: center;
  }

  .brand {
    display: grid;
    gap: 2px;
  }

  .brand-eyebrow,
  .page-label,
  .card-tag {
    font-size: 12px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--muted);
  }

  .brand strong {
    font-size: 20px;
    color: var(--primary);
    font-style: italic;
    font-weight: 900;
  }

  .nav-toggle {
    display: none;
    min-height: 42px;
    padding: 0 16px;
    border: 0;
    border-radius: 999px;
    color: #432c00;
    background: linear-gradient(135deg, var(--primary) 0%, #f3a114 100%);
    cursor: pointer;
  }

  .global-nav {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    gap: 10px;
  }

  .nav-link {
    min-height: 42px;
    padding: 10px 14px;
    border-radius: 999px;
    color: var(--muted);
    border: 0.5px solid transparent;
    transition:
      background var(--transition),
      color var(--transition),
      border-color var(--transition),
      transform var(--transition);
  }

  .nav-link:hover,
  .nav-link:focus-visible {
    color: var(--text);
    border-color: rgba(255, 186, 56, 0.14);
    background: rgba(255, 255, 255, 0.04);
    outline: none;
  }

  .nav-link.is-active {
    color: var(--bg-strong);
    background: linear-gradient(135deg, var(--primary) 0%, #e8a31d 100%);
    box-shadow: 0 0 20px rgba(255, 186, 56, 0.18);
  }

  .bar-status {
    display: flex;
    justify-content: flex-end;
    flex-wrap: wrap;
    gap: 10px;
  }

  .page {
    max-width: 1440px;
    margin: 0 auto;
    padding: 32px 24px 56px;
    display: grid;
    gap: 20px;
  }

  .panel {
    background: var(--surface);
    border: 0.5px solid var(--line);
    border-radius: var(--radius-lg);
    backdrop-filter: blur(38px);
    box-shadow: var(--shadow);
  }

  .hero-panel {
    position: relative;
    overflow: hidden;
    padding: 28px;
    display: grid;
    gap: 16px;
  }

  .hero-panel::before {
    content: "";
    position: absolute;
    inset: 0 auto auto 0;
    width: 220px;
    height: 220px;
    background: radial-gradient(
      circle,
      rgba(255, 186, 56, 0.14) 0%,
      transparent 72%
    );
    pointer-events: none;
  }

  .hero-top,
  .section-head {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 16px;
    flex-wrap: wrap;
  }

  .hero-panel h1 {
    font-size: clamp(40px, 6vw, 72px);
    line-height: 0.96;
    letter-spacing: -0.05em;
    font-style: italic;
  }

  .hero-lead,
  .entry-card p,
  .placeholder-content p {
    color: var(--muted);
    line-height: 1.7;
    font-size: 14px;
  }

  .hero-grid {
    display: grid;
    gap: 20px;
    grid-template-columns: 1fr;
  }

  .entry-panel,
  .placeholder-content {
    padding: 24px;
  }

  .section-head {
    margin-bottom: 18px;
  }

  .section-head h2 {
    font-size: 24px;
  }

  .entry-grid,
  .action-grid {
    display: grid;
    gap: 14px;
  }

  .entry-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .entry-card {
    position: relative;
    overflow: hidden;
    padding: 18px;
    border-radius: var(--radius-md);
    border: 0.5px solid rgba(255, 186, 56, 0.12);
    background: rgba(17, 13, 12, 0.34);
    transition:
      transform var(--transition),
      border-color var(--transition),
      background var(--transition),
      box-shadow var(--transition);
  }

  .entry-card::before {
    content: "✦";
    position: absolute;
    top: 14px;
    right: 16px;
    color: rgba(255, 186, 56, 0.26);
    font-size: 12px;
  }

  .entry-card:hover,
  .entry-card:focus-visible {
    transform: translateY(-2px);
    border-color: rgba(255, 186, 56, 0.24);
    background: rgba(255, 186, 56, 0.08);
    box-shadow: 0 16px 40px rgba(0, 0, 0, 0.26);
    outline: none;
  }

  .entry-head {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    align-items: flex-start;
    margin-bottom: 10px;
  }

  .entry-card h3 {
    font-size: 20px;
  }

  .entry-state {
    display: inline-flex;
    align-items: center;
    min-height: 28px;
    padding: 0 10px;
    border-radius: 999px;
    border: 0.5px solid rgba(255, 186, 56, 0.14);
    background: rgba(255, 255, 255, 0.03);
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .entry-action {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    margin-top: 14px;
    color: var(--primary);
    font-size: 13px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .placeholder-content {
    display: grid;
    gap: 20px;
  }

  .placeholder-card {
    padding: 28px;
    display: grid;
    gap: 16px;
  }

  .placeholder-card h2 {
    font-size: clamp(30px, 4vw, 48px);
    line-height: 1.05;
    font-style: italic;
  }

  .action-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .action-link {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: 48px;
    padding: 0 16px;
    border-radius: 999px;
    border: 0.5px solid rgba(255, 186, 56, 0.12);
    background: rgba(255, 255, 255, 0.02);
    color: var(--text);
    text-align: center;
    transition:
      background var(--transition),
      border-color var(--transition),
      transform var(--transition);
  }

  .action-link:hover,
  .action-link:focus-visible {
    background: rgba(255, 186, 56, 0.08);
    border-color: rgba(255, 186, 56, 0.24);
    transform: translateY(-1px);
    outline: none;
  }

  @media (max-width: 1120px) {
    .shell-bar-inner,
    .hero-grid {
      grid-template-columns: 1fr;
    }

    .bar-status {
      justify-content: flex-start;
    }
  }

  @media (max-width: 860px) {
    .nav-toggle {
      display: inline-flex;
      align-items: center;
      justify-content: center;
    }

    .global-nav {
      display: none;
      justify-content: flex-start;
    }

    .shell-bar.is-open .global-nav {
      display: flex;
    }

    .entry-grid,
    .action-grid {
      grid-template-columns: 1fr;
    }
  }
</style>