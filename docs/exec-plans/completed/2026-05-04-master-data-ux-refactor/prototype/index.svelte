<script lang="ts">
  import GenerationSetupPanel from "./GenerationSetupPanel.svelte"
  import PersonaActionModal from "./PersonaActionModal.svelte"
  import PersonaReviewPanel from "./PersonaReviewPanel.svelte"
  import RunStatusPanel from "./RunStatusPanel.svelte"
  import prototypeData from "../mock-data/master-persona-prototype.json"

  type PreviewState = "empty" | "ready" | "running" | "complete" | "error"
  type PersonaModalMode = "edit" | "delete" | null
  type StateTabId = PreviewState | "edit-modal" | "delete-modal"
  type StateTab = { id: StateTabId; label: string }

  const stateTabs: StateTab[] = [
    { id: "empty", label: "入力前" },
    { id: "ready", label: "選択済み" },
    { id: "running", label: "生成中" },
    { id: "complete", label: "完了" },
    { id: "error", label: "エラー" },
    { id: "edit-modal", label: "編集モーダル" },
    { id: "delete-modal", label: "削除モーダル" }
  ]

  let previewState = $state<PreviewState>("ready")
  let selectedIdentityKey = $state(prototypeData.personas[0]?.identityKey ?? "")
  let provider = $state("gemini")
  let model = $state("gemini-2.5-flash")
  let executionMethod = $state("single_request")
  let refreshNotice = $state("")
  let modalMode = $state<PersonaModalMode>(null)
  let personaPage = $state(prototypeData.personaPage.currentPage)

  const selectedPersona = $derived(
    prototypeData.personas.find((persona) => persona.identityKey === selectedIdentityKey) ??
      prototypeData.personas[0] ??
      null
  )
  const personaPageInfo = $derived({
    ...prototypeData.personaPage,
    currentPage: personaPage
  })

  function chooseFile(): void {
    previewState = "ready"
    refreshNotice = "JSON を選択しました。ペルソナを作成できます。"
  }

  function startGeneration(): void {
    previewState = "running"
    refreshNotice = "生成を開始しました。"
  }

  function finishGeneration(): void {
    previewState = "complete"
    refreshNotice = "生成が完了しました。"
  }

  function pauseGeneration(): void {
    refreshNotice = "一時停止しました。"
  }

  function cancelGeneration(): void {
    previewState = "ready"
    refreshNotice = "生成を中止しました。"
  }

  function refreshModels(): void {
    refreshNotice = "モデル一覧を更新しました。"
  }

  function selectStateTab(id: StateTabId): void {
    if (id === "edit-modal") {
      modalMode = "edit"
      return
    }

    if (id === "delete-modal") {
      modalMode = "delete"
      return
    }

    previewState = id
    modalMode = null
  }

  function isStateTabActive(id: StateTabId): boolean {
    if (id === "edit-modal") return modalMode === "edit"
    if (id === "delete-modal") return modalMode === "delete"
    return previewState === id && modalMode === null
  }

  function selectPersona(identityKey: string): void {
    selectedIdentityKey = identityKey
    modalMode = null
  }

  function changePersonaPage(page: number): void {
    personaPage = page
    modalMode = null
    refreshNotice = `${page} ページ目を表示しています。`
  }

  function closeModal(): void {
    modalMode = null
  }

  function savePersona(): void {
    modalMode = null
    refreshNotice = "編集内容を保存しました。"
  }

  function deletePersona(): void {
    modalMode = null
    refreshNotice = "ペルソナを削除しました。"
  }
</script>

<svelte:head>
  <title>マスターペルソナ生成 UI プロトタイプ</title>
</svelte:head>

<main class="prototype-shell" data-ui-prototype-sample-data-root>
  <header class="shell-bar">
    <div class="brand">
      <span>AITranslationEngineJp</span>
      <strong>マスターペルソナ</strong>
    </div>
    <nav aria-label="状態確認" class="state-nav">
      {#each stateTabs as tab}
        <button
          aria-current={isStateTabActive(tab.id) ? "page" : undefined}
          class:is-active={isStateTabActive(tab.id)}
          onclick={() => selectStateTab(tab.id)}
          type="button"
        >
          {tab.label}
        </button>
      {/each}
    </nav>
  </header>

  <section class="page">
    <section class="hero-panel" aria-labelledby="pageHeading">
      <h1 id="pageHeading">マスターペルソナ作成</h1>
      <p>
        ベースゲームや大型 Mod の NPC を対象に、翻訳前の準備としてペルソナをまとめて作成します。
        作成後は同じ画面で一覧と詳細を確認できます。
      </p>
    </section>

    {#if refreshNotice}
      <p class="notice-banner" role="status">{refreshNotice}</p>
    {/if}

    <GenerationSetupPanel
      aiProviders={prototypeData.aiProviders}
      executionMethods={prototypeData.executionMethods}
      modelOptions={prototypeData.models}
      {provider}
      {model}
      {executionMethod}
      {previewState}
      changeExecutionAction={(event) => {
        const target = event.currentTarget
        if (target instanceof HTMLSelectElement) executionMethod = target.value
      }}
      changeModelAction={(event) => {
        const target = event.currentTarget
        if (target instanceof HTMLSelectElement) model = target.value
      }}
      changeProviderAction={(event) => {
        const target = event.currentTarget
        if (target instanceof HTMLSelectElement) provider = target.value
      }}
      refreshModelsAction={refreshModels}
      chooseFileAction={chooseFile}
      startGenerationAction={startGeneration}
    />

    <RunStatusPanel
      {previewState}
      pauseGenerationAction={pauseGeneration}
      cancelGenerationAction={cancelGeneration}
    />

    {#if previewState === "running"}
      <div class="complete-row">
        <button class="button-primary" onclick={finishGeneration} type="button">
          完了状態を確認
        </button>
      </div>
    {/if}

    <PersonaReviewPanel
      personas={prototypeData.personas}
      {selectedIdentityKey}
      pageInfo={personaPageInfo}
      openEditAction={() => (modalMode = "edit")}
      openDeleteAction={() => (modalMode = "delete")}
      changePageAction={changePersonaPage}
      selectPersonaAction={selectPersona}
    />
  </section>

  <PersonaActionModal
    mode={modalMode}
    persona={selectedPersona}
    closeAction={closeModal}
    saveAction={savePersona}
    deleteAction={deletePersona}
  />
</main>

<style>
  :global(body) {
    margin: 0;
  }

  :global(button),
  :global(input),
  :global(select) {
    font: inherit;
  }

  :global(h1),
  :global(h2),
  :global(h3),
  :global(p) {
    margin: 0;
  }

  .prototype-shell {
    --surface: rgba(32, 31, 29, 0.86);
    --line: rgba(222, 196, 130, 0.28);
    --line-strong: rgba(222, 196, 130, 0.48);
    --text: #f3eee8;
    --muted: #d7cec3;
    --primary: #e8c76f;
    --accent: #8ed8c6;
    --shadow: 0 22px 58px rgba(0, 0, 0, 0.25);
    color: var(--text);
    min-height: 100vh;
  }

  .shell-bar {
    align-items: center;
    background: rgba(20, 20, 19, 0.9);
    border-bottom: 1px solid var(--line);
    display: flex;
    gap: 18px;
    justify-content: space-between;
    padding: 14px clamp(16px, 4vw, 44px);
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .brand {
    display: grid;
    gap: 3px;
    min-width: 180px;
  }

  .brand span {
    color: var(--primary);
    font-size: 0.78rem;
  }

  .state-nav {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    justify-content: flex-end;
  }

  .state-nav button {
    background: transparent;
    border: 1px solid transparent;
    border-radius: 999px;
    color: var(--muted);
    cursor: pointer;
    min-height: 38px;
    padding: 8px 12px;
  }

  .state-nav button.is-active {
    background: rgba(142, 216, 198, 0.12);
    border-color: rgba(142, 216, 198, 0.5);
    color: var(--text);
  }

  .page {
    display: grid;
    gap: 14px;
    margin: 0 auto;
    max-width: 1280px;
    padding: 22px clamp(14px, 3vw, 34px) 56px;
  }

  .hero-panel {
    background: var(--surface);
    border: 1px solid var(--line);
    border-radius: 8px;
    box-shadow: var(--shadow);
    display: grid;
    gap: 10px;
    padding: clamp(18px, 3vw, 28px);
  }

  h1 {
    font-size: clamp(2rem, 4vw, 3.1rem);
    line-height: 1.1;
  }

  .notice-banner {
    background: rgba(142, 216, 198, 0.1);
    border: 1px solid rgba(142, 216, 198, 0.32);
    border-radius: 8px;
    color: var(--accent);
    line-height: 1.7;
    margin: 0;
    padding: 10px 14px;
  }

  .hero-panel p {
    color: var(--muted);
    line-height: 1.7;
    max-width: 760px;
  }

  .complete-row {
    display: flex;
    justify-content: flex-end;
  }

  .button-primary {
    background: var(--primary);
    border: 1px solid var(--primary);
    border-radius: 7px;
    color: #201309;
    cursor: pointer;
    min-height: 40px;
    padding: 9px 13px;
  }

  @media (max-width: 760px) {
    .shell-bar {
      align-items: stretch;
      flex-direction: column;
    }

    .state-nav {
      justify-content: flex-start;
    }
  }
</style>
