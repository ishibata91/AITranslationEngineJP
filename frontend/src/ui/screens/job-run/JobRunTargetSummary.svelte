<script lang="ts">
  import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"

  interface Props {
    target: TranslationJobManagementJobRunTarget
  }

  let { target }: Props = $props()
</script>

<section class="job-run-target-summary" data-testid="job-run-selected-job-summary">
  <div>
    <p class="eyebrow">選択中のジョブ</p>
    <h3>ジョブ #{target.jobId}</h3>
  </div>
  <dl class="target-summary-grid">
    <div>
      <dt>状態</dt>
      <dd>{target.stateLabel} / {target.stateDescription}</dd>
    </div>
    <div>
      <dt>現在の翻訳段階</dt>
      <dd>{target.currentPhaseLabel}</dd>
    </div>
    <div>
      <dt>進捗</dt>
      <dd>{target.progressLabel}</dd>
    </div>
    <div>
      <dt>入力</dt>
      <dd>{target.inputSourceLabel}</dd>
    </div>
  </dl>
  <details class="target-path">
    <summary>入力ファイル path</summary>
    <p>{target.sourcePath}</p>
  </details>
</section>

<style>
  .job-run-target-summary {
    display: grid;
    gap: 0.8rem;
    min-width: 0;
    padding: 1.15rem 1.3rem;
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
  }

  .eyebrow,
  .target-path {
    color: rgba(236, 223, 205, 0.78);
  }

  h3 {
    margin: 0.2rem 0 0;
    color: #fff6ea;
  }

  .target-summary-grid {
    display: grid;
    gap: 0.6rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin: 0;
  }

  .target-summary-grid div {
    display: grid;
    gap: 0.2rem;
    min-width: 0;
  }

  .target-summary-grid dt {
    color: rgba(236, 223, 205, 0.78);
  }

  .target-summary-grid dd {
    margin: 0;
    color: #fff6ea;
    overflow-wrap: anywhere;
  }

  @media (max-width: 900px) {
    .target-summary-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
