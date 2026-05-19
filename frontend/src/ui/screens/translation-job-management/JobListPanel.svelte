<script lang="ts">
  import type {
    TranslationJobManagementFilterChipViewModel,
    TranslationJobManagementJobCardViewModel,
    TranslationJobManagementOperationViewModel,
    TranslationJobManagementScreenPhase
  } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import JobCard from "./JobCard.svelte"

  interface Props {
    phase: TranslationJobManagementScreenPhase
    jobs: TranslationJobManagementJobCardViewModel[]
    searchQuery: string
    filterChips: TranslationJobManagementFilterChipViewModel[]
    listEmptyTitle: string
    listEmptyDescription: string
    listErrorTitle: string
    listErrorDescription: string
    onSearchQueryChange: (searchQuery: string) => void
    onFilterChange: (filterId: TranslationJobManagementFilterChipViewModel["id"]) => void
    onOpenJob: (job: TranslationJobManagementJobCardViewModel) => void | Promise<void>
    onOperation: (
      job: TranslationJobManagementJobCardViewModel,
      operation: TranslationJobManagementOperationViewModel
    ) => void | Promise<void>
  }

  let {
    phase,
    jobs,
    searchQuery,
    filterChips,
    listEmptyTitle,
    listEmptyDescription,
    listErrorTitle,
    listErrorDescription,
    onSearchQueryChange,
    onFilterChange,
    onOpenJob,
    onOperation
  }: Props = $props()

  const currentFilter = $derived(
    filterChips.find((chip) => chip.selected) ?? filterChips[0]
  )

  function handleFilterChange(event: Event): void {
    onFilterChange(
      (event.currentTarget as HTMLSelectElement)
        .value as TranslationJobManagementFilterChipViewModel["id"]
    )
  }
</script>

<section class="job-management-layout">
  <section
    class="panel job-list-panel"
    aria-labelledby="jobManagementListHeading"
    data-testid="translation-job-management-job-list-region"
  >
    <div class="section-head">
      <div>
        <p class="page-label">未完了ジョブ</p>
        <h3 id="jobManagementListHeading">一覧</h3>
      </div>
    </div>

    <div class="filter-toolbar">
      <label class="search-field" for="jobManagementSearch">
        <span>検索</span>
        <input
          data-testid="translation-job-management-search-field"
          id="jobManagementSearch"
          oninput={(event) =>
            onSearchQueryChange((event.currentTarget as HTMLInputElement).value)}
          placeholder="ジョブID、入力名、翻訳段階で検索"
          type="search"
          value={searchQuery}
        />
      </label>

      <label class="select-field" for="jobManagementFilter">
        <span>状態フィルタ</span>
        <select
          data-testid="translation-job-management-state-filter"
          id="jobManagementFilter"
          onchange={handleFilterChange}
          value={currentFilter?.id ?? "all"}
        >
          {#each filterChips as chip (chip.id)}
            <option value={chip.id}>{chip.label} ({chip.count})</option>
          {/each}
        </select>
      </label>
    </div>

    {#if phase === "loading"}
      <div
        class="panel-empty"
        data-testid="translation-job-management-list-status-display"
      >
        <p class="empty-title">一覧を読み込んでいます</p>
        <p class="empty-description">未完了ジョブの状態を取得しています。</p>
      </div>
    {:else if phase === "error"}
      <div
        class="panel-empty tone-danger"
        data-testid="translation-job-management-list-status-display"
      >
        <p class="empty-title">{listErrorTitle}</p>
        <p class="empty-description">{listErrorDescription}</p>
      </div>
    {:else if jobs.length === 0}
      <div
        class="panel-empty"
        data-testid="translation-job-management-list-status-display"
      >
        <p class="empty-title">{listEmptyTitle}</p>
        <p class="empty-description">{listEmptyDescription}</p>
      </div>
    {:else}
      <div
        class="job-card-list"
        data-testid="translation-job-management-list-status-display"
      >
        {#each jobs as job (job.jobId)}
          <JobCard {job} {onOpenJob} {onOperation} />
        {/each}
      </div>
    {/if}
  </section>
</section>

<style>
  .job-management-layout {
    display: grid;
    gap: 1rem;
  }

  .panel {
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
    padding: 1.1rem 1.2rem;
  }

  .section-head {
    display: flex;
    gap: 1rem;
    justify-content: space-between;
    align-items: flex-start;
    flex-wrap: wrap;
  }

  h3,
  .empty-title {
    margin: 0;
    color: #fff6ea;
  }

  .page-label,
  .empty-description {
    color: rgba(236, 223, 205, 0.78);
  }

  .filter-toolbar,
  .search-field,
  .select-field,
  .job-card-list {
    display: grid;
    gap: 0.75rem;
  }

  .job-card-list {
    gap: 1rem;
    margin-top: 1.5rem;
  }

  .filter-toolbar {
    grid-template-columns: minmax(0, 1fr) minmax(220px, 280px);
    margin-top: 1rem;
  }

  .search-field span,
  .select-field span {
    color: rgba(236, 223, 205, 0.78);
    font-size: 0.82rem;
    letter-spacing: 0.04em;
  }

  .search-field input,
  .select-field select {
    min-height: 2.8rem;
    border-radius: 14px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    font: inherit;
    padding: 0.65rem 0.9rem;
    color: #fff6ea;
    background: rgba(255, 255, 255, 0.05);
  }

  .panel-empty {
    display: grid;
    gap: 0.35rem;
    padding: 0.8rem 0;
  }

  .tone-danger {
    border-color: rgba(255, 128, 102, 0.38);
  }

  @media (max-width: 1120px) {
    .job-management-layout {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 760px) {
    .filter-toolbar {
      grid-template-columns: 1fr;
    }
  }
</style>
