<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateTranslationOutputArtifactScreenController,
    TranslationOutputArtifactScreenControllerContract
  } from "@application/contract/translation-output-artifact"
  import type { TranslationOutputCompletedJobSummary } from "@application/gateway-contract/translation-output-artifact"

  const PLACEHOLDER_LEAD =
    "このページはまだ準備中です。上のナビゲーションまたは下の移動から別の主要ページへ進めます。"

  interface Props {
    createController: CreateTranslationOutputArtifactScreenController | null
  }

  let { createController }: Props = $props()

  function resolveController(): TranslationOutputArtifactScreenControllerContract {
    if (!createController) {
      throw new Error(
        "translation output artifact screen controller factory is not provided"
      )
    }

    return createController()
  }

  function formatCount(value: number | undefined): string {
    if (typeof value !== "number") {
      return "-"
    }

    return `${value.toLocaleString("ja-JP")} 件`
  }

  function formatDistribution(
    distribution: Record<string, number> | undefined
  ): string {
    if (!distribution || Object.keys(distribution).length === 0) {
      return "-"
    }

    return Object.entries(distribution)
      .map(([key, value]) => `${key}: ${value}`)
      .join(" / ")
  }

  function formatStatus(value: string | undefined): string {
    if (!value) {
      return "-"
    }

    return value.replaceAll("_", " ")
  }

  const controller = resolveController()
  let viewModel = $state(controller.getViewModel())

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
  })

  onMount(() => {
    void (async () => {
      await controller.mount()
    })()

    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

  async function selectJob(
    job: TranslationOutputCompletedJobSummary
  ): Promise<void> {
    await controller.setJobId(job.jobId)
  }
</script>

<section class="output-shell" id="translationOutputArtifactView">
  <section
    class="output-card hero-card"
    data-testid="output-management-output-management-summary"
  >
    <div class="hero-head">
      <div>
        <p class="eyebrow">output-management</p>
        <h2>Output Review</h2>
      </div>
      <p class="gateway-status">Gateway: {viewModel.gatewayStatus}</p>
    </div>
    <p class="lead">
      completed job、result summary、diff preview、artifact 状態を確認し、
      xTranslator XML の出力と再出力を同じ画面で行います。
    </p>
    <p class="status-copy">
      <strong>{viewModel.statusTitle}</strong>
      <span>{viewModel.statusText}</span>
    </p>
    <p class="placeholder-copy">{PLACEHOLDER_LEAD}</p>
  </section>

  <section class="output-grid">
    <section
      class="output-card job-list-card"
      aria-labelledby="outputJobListHeading"
      data-testid="output-management-output-candidate-list"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">completed job list</p>
          <h3 id="outputJobListHeading">出力候補</h3>
        </div>
        <button
          class="secondary-button"
          disabled={viewModel.isLoading || viewModel.isSubmitting}
          onclick={() => void controller.refresh()}
          type="button"
        >
          更新
        </button>
      </div>
      {#if viewModel.completedJobs.length === 0}
        <p class="empty-text">
          completed job はありません。target count 0 の job も一覧に出ません。
        </p>
      {:else}
        <div class="job-list">
          {#each viewModel.completedJobs as job (job.jobId)}
            <button
              class="job-button"
              class:is-selected={job.jobId === viewModel.selectedJobId}
              onclick={() => void selectJob(job)}
              type="button"
            >
              <div class="job-button-head">
                <strong>job #{job.jobId}</strong>
                <span class="status-pill"
                  >{formatStatus(job.artifactStatus)}</span
                >
              </div>
              <p>{formatStatus(job.jobStatus)}</p>
              <p>translated: {formatCount(job.translatedCount)}</p>
              <p>status: {formatDistribution(job.outputStatusDistribution)}</p>
            </button>
          {/each}
        </div>
      {/if}
    </section>

    <section
      class="output-card summary-card"
      aria-labelledby="outputSummaryHeading"
      data-testid="output-management-selected-job"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">selected job summary</p>
          <h3 id="outputSummaryHeading">
            {viewModel.review ? "選択中 job" : "job 選択待ち"}
          </h3>
        </div>
        <span class="status-pill" data-view-state={viewModel.viewState}>
          {formatStatus(viewModel.viewState)}
        </span>
      </div>
      {#if viewModel.review}
        <dl class="detail-grid">
          <div>
            <dt>job id</dt>
            <dd>{viewModel.review.selectedJobId}</dd>
          </div>
          <div>
            <dt>job status</dt>
            <dd>{formatStatus(viewModel.review.selectedJobStatus)}</dd>
          </div>
          <div>
            <dt>body phase</dt>
            <dd>{formatStatus(viewModel.review.bodyPhaseStatus)}</dd>
          </div>
          <div>
            <dt>readiness</dt>
            <dd>{viewModel.review.readiness ? "ready" : "not ready"}</dd>
          </div>
          <div>
            <dt>translated count</dt>
            <dd>{formatCount(viewModel.review.translatedCount)}</dd>
          </div>
          <div>
            <dt>row count</dt>
            <dd>{formatCount(viewModel.review.rowCount)}</dd>
          </div>
          <div>
            <dt>artifact status</dt>
            <dd>{formatStatus(viewModel.review.artifactStatus)}</dd>
          </div>
          <div>
            <dt>current version</dt>
            <dd>{viewModel.review.currentVersion ? "yes" : "stale"}</dd>
          </div>
          <div>
            <dt>input provenance</dt>
            <dd class="wrap-value">
              {viewModel.review.inputSnapshotDigest} / {viewModel.review
                .sourceFileDigest}
            </dd>
          </div>
        </dl>
        {#if viewModel.review.rejectionReasons.length > 0}
          <div class="notice-block warning">
            <h4>拒否理由</h4>
            <ul>
              {#each viewModel.review.rejectionReasons as reason (`${reason.errorKind}-${reason.reason}`)}
                <li>{reason.errorKind}: {reason.reason}</li>
              {/each}
            </ul>
          </div>
        {/if}
      {:else}
        <p class="empty-text">
          出力候補から completed job を選ぶと、summary と出力準備を表示します。
        </p>
      {/if}
    </section>

    <section
      class="output-card action-card"
      aria-labelledby="outputActionHeading"
      data-testid="output-management-output-actions"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">action rail</p>
          <h3 id="outputActionHeading">出力操作</h3>
        </div>
      </div>
      <label class="field-block" for="outputTargetGame">
        <span>target game</span>
        <select
          id="outputTargetGame"
          onchange={(event) => {
            const target = event.currentTarget
            if (target instanceof HTMLSelectElement) {
              controller.setTargetGame(target.value)
            }
          }}
          value={viewModel.targetGame}
        >
          <option value="skyrim_se">Skyrim SE</option>
          <option value="skyrim_le">Skyrim LE</option>
        </select>
      </label>
      <label class="field-block" for="outputPath">
        <span>output path</span>
        <input
          id="outputPath"
          oninput={(event) => {
            const target = event.currentTarget
            if (target instanceof HTMLInputElement) {
              controller.setOutputPath(target.value)
            }
          }}
          type="text"
          value={viewModel.outputPath}
        />
      </label>
      <p class="helper-text" data-path-state={viewModel.pathState}>
        {viewModel.pathReason || "出力先 path は .xml で終える必要があります。"}
      </p>
      <div class="action-row">
        <button
          class="primary-button"
          disabled={!viewModel.canGenerate || viewModel.isSubmitting}
          onclick={() => void controller.generateArtifact()}
          type="button"
        >
          XML を出力
        </button>
        <button
          class="secondary-button"
          disabled={!viewModel.canRegenerate || viewModel.isSubmitting}
          onclick={() => void controller.regenerateArtifact()}
          type="button"
        >
          再出力
        </button>
      </div>
      {#if viewModel.disabledReason}
        <p class="helper-text warning">{viewModel.disabledReason}</p>
      {/if}
      {#if viewModel.lastCommand}
        <div class="notice-block" data-testid="output-management-latest-result">
          <h4>result summary</h4>
          <dl class="detail-grid compact">
            <div>
              <dt>artifact status</dt>
              <dd>{formatStatus(viewModel.lastCommand.artifactStatus)}</dd>
            </div>
            <div>
              <dt>row count</dt>
              <dd>{formatCount(viewModel.lastCommand.rowCount)}</dd>
            </div>
            <div>
              <dt>file path</dt>
              <dd class="wrap-value">
                {viewModel.lastCommand.filePath ?? "-"}
              </dd>
            </div>
            <div>
              <dt>target game</dt>
              <dd>{viewModel.lastCommand.targetGame}</dd>
            </div>
          </dl>
          {#if viewModel.lastCommand.errorReason}
            <p class="helper-text warning">
              {viewModel.lastCommand.errorReason}
            </p>
          {/if}
        </div>
      {/if}
    </section>
  </section>

  <section
    class="output-card diff-card"
    aria-labelledby="outputDiffHeading"
    data-testid="output-management-diff-preview"
  >
    <div class="section-head">
      <div>
        <p class="eyebrow">diff preview</p>
        <h3 id="outputDiffHeading">translation unit 差分</h3>
      </div>
      <span class="mini-text">{viewModel.compatibilitySummaryText}</span>
    </div>
    {#if viewModel.diffPreview && viewModel.diffPreview.rows.length > 0}
      <div class="diff-table">
        <div class="diff-table-head">
          <span>Source</span>
          <span>Dest</span>
          <span>Status</span>
          <span>row reflection summary</span>
        </div>
        {#each viewModel.diffPreview.rows as row (row.rowDigest)}
          <button
            class="diff-row"
            onclick={() =>
              void controller.setArtifactId(
                viewModel.diffPreview?.artifactId ?? null
              )}
            type="button"
          >
            <div>
              <strong>{row.edid}</strong>
              <p>{row.sourceExcerpt || "-"}</p>
            </div>
            <div>
              <strong>{row.formId}</strong>
              <p>{row.destExcerpt || "-"}</p>
            </div>
            <div>
              <p>{row.internalOutputStatus}</p>
              <p>xTranslator: {row.xTranslatorStatus}</p>
            </div>
            <div>
              <p>{row.rowReflectionSummary}</p>
              <p>stale: {row.staleReason || "-"}</p>
              <p>re-output: {row.canRegenerate ? "可" : "不要"}</p>
            </div>
          </button>
        {/each}
      </div>
    {:else}
      <p class="empty-text">
        diff preview は未取得です。artifact がない場合、row count 0
        の場合、または gateway 未接続の場合は一覧を表示しません。
      </p>
    {/if}
  </section>
</section>

<style>
  .output-shell {
    display: grid;
    gap: 1.5rem;
  }

  .output-grid {
    display: grid;
    gap: 1rem;
    grid-template-columns: minmax(0, 1.1fr) minmax(0, 1fr) minmax(18rem, 0.9fr);
  }

  .output-card {
    border: 1px solid var(--line);
    border-radius: var(--radius-md);
    background: rgba(25, 22, 20, 0.82);
    box-shadow: var(--shadow);
    padding: 1.25rem;
  }

  .hero-card,
  .diff-card {
    display: grid;
    gap: 0.9rem;
  }

  .hero-head,
  .section-head,
  .job-button-head,
  .action-row {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    justify-content: space-between;
  }

  .eyebrow,
  .mini-text,
  .placeholder-copy {
    color: var(--muted);
    font-size: 0.85rem;
  }

  .lead,
  .status-copy,
  .helper-text,
  .empty-text,
  .job-button p,
  .diff-row p {
    line-height: 1.6;
  }

  .job-list,
  .notice-block,
  .field-block,
  .diff-table,
  .diff-row,
  .detail-grid {
    display: grid;
    gap: 0.75rem;
  }

  .job-button,
  .diff-row {
    width: 100%;
    text-align: left;
    border: 1px solid var(--line);
    border-radius: 14px;
    background: rgba(36, 31, 28, 0.92);
    color: inherit;
    padding: 0.95rem;
  }

  .job-button.is-selected {
    border-color: var(--primary);
    box-shadow: 0 0 0 1px rgba(255, 186, 56, 0.32);
  }

  .status-pill {
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    padding: 0.2rem 0.65rem;
    color: var(--primary);
    font-size: 0.82rem;
  }

  .status-pill[data-view-state="failed"],
  .helper-text.warning,
  .notice-block.warning {
    color: #ff9f7f;
  }

  .detail-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  dt {
    color: var(--muted);
    font-size: 0.82rem;
    margin-bottom: 0.2rem;
  }

  dd {
    margin: 0;
  }

  .wrap-value {
    overflow-wrap: anywhere;
  }

  .field-block span,
  .notice-block h4 {
    font-weight: 700;
  }

  .field-block input,
  .field-block select,
  .primary-button,
  .secondary-button {
    border: 1px solid var(--line);
    border-radius: 12px;
    background: rgba(15, 13, 12, 0.9);
    color: inherit;
    padding: 0.8rem 0.95rem;
  }

  .primary-button {
    background: linear-gradient(135deg, #ffb43a, #ff8743);
    color: #1a110d;
    font-weight: 700;
  }

  .secondary-button {
    cursor: pointer;
  }

  .primary-button:disabled,
  .secondary-button:disabled,
  .job-button:disabled,
  .diff-row:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .notice-block {
    border: 1px solid var(--line);
    border-radius: 12px;
    padding: 0.85rem;
    background: rgba(18, 16, 15, 0.88);
  }

  .notice-block ul {
    margin: 0;
    padding-left: 1.2rem;
  }

  .diff-table-head,
  .diff-row {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: minmax(0, 1.2fr) minmax(0, 1.2fr) minmax(
        0,
        0.7fr
      ) minmax(0, 1fr);
  }

  @media (max-width: 1080px) {
    .output-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 720px) {
    .detail-grid,
    .detail-grid.compact,
    .diff-table-head,
    .diff-row {
      grid-template-columns: 1fr;
    }

    .action-row {
      flex-direction: column;
      align-items: stretch;
    }
  }
</style>
