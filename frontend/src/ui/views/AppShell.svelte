<script lang="ts">
  import { onMount } from "svelte"

  import type { CreateBodyTranslationPhaseScreenController } from "@application/contract/body-translation-phase"
  import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
  import type { CreateMasterPersonaScreenController } from "@application/contract/master-persona"
  import type { CreatePersonaGenerationPhaseScreenController } from "@application/contract/persona-generation-phase"
  import type { CreateProviderSettingsScreenController } from "@application/contract/provider-settings"
  import type { CreateTermTranslationPhaseScreenController } from "@application/contract/term-translation-phase"
  import type { CreateTranslationJobManagementScreenController } from "@application/contract/translation-job-management"
  import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import type { CreateTranslationOutputArtifactScreenController } from "@application/contract/translation-output-artifact"
  import type { CreateTranslationInputScreenController } from "@application/contract/translation-input"
  import AppHeader from "@ui/views/dashboard/AppHeader.svelte"
  import CurrentPageHero from "@ui/views/dashboard/CurrentPageHero.svelte"
  import DashboardEntryGrid from "@ui/views/dashboard/DashboardEntryGrid.svelte"
  import ProviderSettingsPage from "@ui/screens/provider-settings/ProviderSettingsPage.svelte"
  import MasterDictionaryPage from "@ui/screens/master-dictionary/MasterDictionaryPage.svelte"
  import MasterPersonaPage from "@ui/screens/master-persona/MasterPersonaPage.svelte"
  import JobRunPage from "@ui/screens/job-run/JobRunPage.svelte"
  import TranslationManagementStepper from "@ui/screens/translation-job-management/TranslationManagementStepper.svelte"
  import TranslationJobManagementPage from "@ui/screens/translation-job-management/TranslationJobManagementPage.svelte"
  import TranslationOutputArtifactPage from "@ui/screens/translation-output-artifact/TranslationOutputArtifactPage.svelte"
  import InputReviewPage from "@ui/screens/translation-input/InputReviewPage.svelte"
  import type {
    ShellRouteContract,
    ShellRouteId,
    TranslationManagementViewContract,
    TranslationManagementViewId
  } from "@ui/stores/shell-state"

  interface Props {
    defaultRouteId: ShellRouteId
    defaultTranslationManagementViewId: TranslationManagementViewId
    routes: ShellRouteContract[]
    translationManagementViews: TranslationManagementViewContract[]
    createBodyTranslationPhaseScreenController: CreateBodyTranslationPhaseScreenController | null
    createMasterDictionaryScreenController: CreateMasterDictionaryScreenController | null
    createMasterPersonaScreenController: CreateMasterPersonaScreenController | null
    createPersonaGenerationPhaseScreenController: CreatePersonaGenerationPhaseScreenController | null
    createProviderSettingsScreenController: CreateProviderSettingsScreenController | null
    createTermTranslationPhaseScreenController: CreateTermTranslationPhaseScreenController | null
    createTranslationJobManagementScreenController?: CreateTranslationJobManagementScreenController | null
    createTranslationOutputArtifactScreenController: CreateTranslationOutputArtifactScreenController | null
    createTranslationInputScreenController: CreateTranslationInputScreenController | null
  }

  let {
    defaultRouteId,
    defaultTranslationManagementViewId,
    routes,
    translationManagementViews,
    createBodyTranslationPhaseScreenController,
    createMasterDictionaryScreenController,
    createMasterPersonaScreenController,
    createPersonaGenerationPhaseScreenController,
    createProviderSettingsScreenController,
    createTermTranslationPhaseScreenController,
    createTranslationJobManagementScreenController = null,
    createTranslationOutputArtifactScreenController,
    createTranslationInputScreenController
  }: Props = $props()

  const PLACEHOLDER_LEAD =
    "このページはまだ準備中です。上のナビゲーションまたは下の移動から別の主要ページへ進めます。"

  const routeById = $derived(
    new Map(routes.map((route) => [route.id, route] as const))
  )

  let currentRouteId = $state<ShellRouteId>("dashboard")
  let selectedTranslationManagementViewId =
    $state<TranslationManagementViewId | null>(null)
  let selectedJobRunTarget =
    $state<TranslationJobManagementJobRunTarget | null>(null)
  let currentJobRunPhaseViewId =
    $state<TranslationManagementViewId>("term-translation")
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
  const currentTranslationManagementViewId = $derived(
    selectedTranslationManagementViewId === "job-run"
      ? currentJobRunPhaseViewId
      : (selectedTranslationManagementViewId ??
          defaultTranslationManagementViewId)
  )
  const renderedTranslationManagementViewId = $derived(
    selectedTranslationManagementViewId ?? defaultTranslationManagementViewId
  )
  const dashboardEntryRoutes = $derived(
    routes.filter((route) => route.id !== "dashboard")
  )

  function normalizeRouteId(hashValue: string): ShellRouteId {
    const routeId = hashValue.replace(/^#/, "").split("/")[0] as ShellRouteId
    return routeById.has(routeId) ? routeId : defaultRouteId
  }

  function isBaseTranslationManagementHash(hashValue: string): boolean {
    return hashValue === "#translation-management"
  }

  function syncRouteFromHash(): void {
    if (typeof window === "undefined") {
      currentRouteId = defaultRouteId
      return
    }

    const nextRouteId = normalizeRouteId(window.location.hash)
    if (window.location.hash !== `#${nextRouteId}`) {
      const normalizedHash = `#${nextRouteId}`
      const keepsKnownSubRoute = window.location.hash.startsWith(
        `${normalizedHash}/`
      )
      if (!keepsKnownSubRoute) {
        window.history.replaceState(null, "", normalizedHash)
      }
    }
    currentRouteId = nextRouteId
    if (isBaseTranslationManagementHash(window.location.hash)) {
      selectedJobRunTarget = null
      selectedTranslationManagementViewId = "job-management"
      currentJobRunPhaseViewId = "term-translation"
    } else if (
      window.location.hash.startsWith("#translation-management/job-run")
    ) {
      selectedJobRunTarget = null
      selectedTranslationManagementViewId = "job-management"
      currentJobRunPhaseViewId = "term-translation"
      window.history.replaceState(null, "", "#translation-management")
    }
    isMobileNavOpen = false
  }

  function selectRoute(routeId: ShellRouteId): void {
    currentRouteId = routeId
    if (routeId === "translation-management") {
      selectedJobRunTarget = null
      selectedTranslationManagementViewId = "job-management"
      currentJobRunPhaseViewId = "term-translation"
    }
    isMobileNavOpen = false
  }

  function toggleMobileNav(): void {
    isMobileNavOpen = !isMobileNavOpen
  }

  function selectTranslationManagementView(
    viewId: TranslationManagementViewId
  ): void {
    if (viewId !== "job-management") {
      return
    }

    openTranslationJobManagement()
  }

  function openTranslationInputReview(): void {
    selectedTranslationManagementViewId = "input-review"
  }

  function openTranslationJobRun(
    target: TranslationJobManagementJobRunTarget | null = null
  ): void {
    if (target) {
      selectedJobRunTarget = target
    }

    if (selectedJobRunTarget === null) {
      selectedTranslationManagementViewId = "job-management"
      return
    }

    selectedTranslationManagementViewId = "job-run"
    currentJobRunPhaseViewId = resolvePhaseViewId(
      selectedJobRunTarget.currentPhase
    )
    if (typeof window !== "undefined") {
      window.history.pushState(null, "", "#translation-management/job-run")
    }
  }

  function resolvePhaseViewId(
    currentPhase: string
  ): TranslationManagementViewId {
    if (currentPhase === "persona_generation") {
      return "persona-generation"
    }
    if (currentPhase === "body_translation") {
      return "body-translation"
    }
    return "term-translation"
  }

  function openTranslationJobManagement(): void {
    selectedJobRunTarget = null
    selectedTranslationManagementViewId = "job-management"
    currentJobRunPhaseViewId = "term-translation"
    if (typeof window !== "undefined") {
      window.history.pushState(null, "", "#translation-management")
    }
  }

  function openOutputManagementFromJobRun(): void {
    selectedJobRunTarget = null
    selectedTranslationManagementViewId = "job-management"
    currentJobRunPhaseViewId = "term-translation"
    currentRouteId = "output-management"
    if (typeof window !== "undefined") {
      window.history.pushState(null, "", "#output-management")
    }
  }

  function syncJobRunTarget(
    target: TranslationJobManagementJobRunTarget | null
  ): void {
    selectedJobRunTarget = target
  }

  function syncJobRunPhaseView(viewId: TranslationManagementViewId): void {
    currentJobRunPhaseViewId = viewId
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
  <AppHeader
    {currentRoute}
    {isMobileNavOpen}
    {routes}
    {selectRoute}
    {toggleMobileNav}
  />

  <section
    class="page"
    data-testid={isDashboard ? "dashboard-dashboard-content" : undefined}
  >
    <CurrentPageHero
      {currentRoute}
      dataTestId={isDashboard
        ? "dashboard-current-page-description"
        : currentRoute.id === "translation-management"
          ? "translation-management-translation-management-header"
          : undefined}
    />

    {#if isDashboard}
      <DashboardEntryGrid routes={dashboardEntryRoutes} {selectRoute} />
    {/if}

    {#if !isDashboard && currentRoute.id === "provider-settings"}
      <ProviderSettingsPage
        createController={createProviderSettingsScreenController}
      />
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

    {#if !isDashboard && currentRoute.id === "translation-management"}
      <section
        class="translation-management-shell"
        data-testid="translation-management-translation-management-shell"
      >
        <section
          class="panel section-switcher"
          data-testid={renderedTranslationManagementViewId === "job-run"
            ? "job-run-translation-management-stepper"
            : "translation-management-translation-management-stepper"}
        >
          <TranslationManagementStepper
            currentViewId={currentTranslationManagementViewId}
            onSelect={selectTranslationManagementView}
            views={translationManagementViews}
          />
        </section>

        {#if renderedTranslationManagementViewId === "input-review"}
          <InputReviewPage
            createController={createTranslationInputScreenController}
            onOpenJobRun={openTranslationJobRun}
          />
        {/if}

        {#if renderedTranslationManagementViewId === "job-management"}
          <TranslationJobManagementPage
            createController={createTranslationJobManagementScreenController}
            onJobRunTargetChange={syncJobRunTarget}
            onOpenInputReview={openTranslationInputReview}
            onOpenJobRun={openTranslationJobRun}
          />
        {/if}

        {#if renderedTranslationManagementViewId === "job-run"}
          <JobRunPage
            createBodyController={createBodyTranslationPhaseScreenController}
            createController={createTermTranslationPhaseScreenController}
            createPersonaController={createPersonaGenerationPhaseScreenController}
            onOpenJobManagement={openTranslationJobManagement}
            onOpenOutputManagement={openOutputManagementFromJobRun}
            onPhaseViewChange={syncJobRunPhaseView}
            selectedJobTarget={selectedJobRunTarget}
          />
        {/if}
      </section>
    {/if}

    {#if !isDashboard && currentRoute.id === "output-management"}
      <TranslationOutputArtifactPage
        createController={createTranslationOutputArtifactScreenController}
      />
    {/if}

    {#if !isDashboard && currentRoute.id !== "provider-settings" && currentRoute.id !== "master-dictionary" && currentRoute.id !== "master-persona" && currentRoute.id !== "translation-management" && currentRoute.id !== "output-management"}
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

  .page-label {
    font-size: 12px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--muted);
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

  .placeholder-content p {
    color: var(--muted);
    line-height: 1.7;
    font-size: 14px;
  }

  .translation-management-shell {
    display: grid;
    gap: 1.5rem;
  }

  .placeholder-content {
    padding: 24px;
  }

  .action-grid {
    display: grid;
    gap: 14px;
  }

  .section-switcher {
    padding: 0;
    overflow: hidden;
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

  @media (max-width: 860px) {
    .action-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
