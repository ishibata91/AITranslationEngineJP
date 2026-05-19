<script lang="ts">
  import type { TranslationJobSetupValidationResponse } from "@application/gateway-contract/translation-job-setup"

  interface Props {
    canValidate: boolean
    dirty: boolean
    validationResult: TranslationJobSetupValidationResponse | null
    formatDate: (timestamp: string) => string
    resolveValidationLabel: (status: string) => string
    onRunValidation: () => void
  }

  let {
    canValidate,
    dirty,
    validationResult,
    formatDate,
    resolveValidationLabel,
    onRunValidation
  }: Props = $props()
</script>

<section
  class="job-setup-card"
  aria-labelledby="jobSetupValidationHeading"
  data-testid="translation-job-setup-compatibility-precheck-region"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">pre-check</p>
      <h3 id="jobSetupValidationHeading">作成前確認</h3>
    </div>
    <button
      class="button-secondary"
      disabled={!canValidate}
      onclick={onRunValidation}
      type="button"
    >
      確認を実行
    </button>
  </div>
  <dl class="detail-grid compact">
    <div>
      <dt>状態</dt>
      <dd>
        {#if validationResult}
          {resolveValidationLabel(validationResult.status)}
        {:else}
          未実行
        {/if}
      </dd>
    </div>
    <div>
      <dt>確認日時</dt>
      <dd>{formatDate(validationResult?.validatedAt ?? "")}</dd>
    </div>
    <div>
      <dt>作成できない理由</dt>
      <dd class="wrap-value">
        {validationResult?.blockingFailureCategory ?? "-"}
      </dd>
    </div>
    <div>
      <dt>再確認</dt>
      <dd>{dirty ? "再確認が必要" : "確認済み"}</dd>
    </div>
  </dl>
  <div class="tag-list">
    {#each validationResult?.targetSlices ?? [] as slice (slice)}
      <span class="tag warning">{slice}</span>
    {/each}
    {#each validationResult?.passSlices ?? [] as slice (slice)}
      <span class="tag success">{slice}</span>
    {/each}
  </div>
</section>

<style>
  .job-setup-card,
  .detail-grid div {
    display: grid;
    gap: 0.75rem;
  }

  .job-setup-card {
    gap: 1rem;
    padding: 1.25rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 1.25rem;
    background: rgba(34, 26, 23, 0.82);
    box-shadow: 0 20px 40px rgba(6, 4, 3, 0.18);
    color: var(--text);
  }

  .section-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
  }

  .eyebrow,
  dt {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.9rem;
  }

  .eyebrow {
    font-size: 0.78rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .detail-grid {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .tag-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .tag {
    padding: 0.25rem 0.6rem;
    border-radius: 999px;
    background: rgba(255, 190, 126, 0.14);
    color: #ffd8ae;
    font-size: 0.82rem;
  }

  .tag.success {
    background: rgba(145, 208, 134, 0.16);
    color: #b8f0ad;
  }

  .tag.warning {
    background: rgba(255, 213, 128, 0.16);
    color: #ffe0a3;
  }

  .button-secondary {
    padding: 0.8rem 1rem;
    border-radius: 0.9rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    background: rgba(255, 241, 227, 0.08);
    color: #ffe2bf;
    cursor: pointer;
  }

  button:disabled {
    opacity: 0.56;
    cursor: not-allowed;
  }

  .wrap-value {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  @media (max-width: 720px) {
    .section-head {
      flex-direction: column;
    }
  }
</style>
