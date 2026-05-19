<script lang="ts">
  interface Props {
    headingId: string
    testId: string
    errorKindLabel: string
    errorReasonLabel: string
    retryableLabel: string
  }

  let {
    headingId,
    testId,
    errorKindLabel,
    errorReasonLabel,
    retryableLabel
  }: Props = $props()
</script>

<section class="phase-card" aria-labelledby={headingId} data-testid={testId}>
  <div class="section-head">
    <div>
      <p class="eyebrow">error summary</p>
      <h3 id={headingId}>失敗情報</h3>
    </div>
  </div>
  <dl class="detail-grid compact">
    <div>
      <dt>error kind</dt>
      <dd class="wrap-value">{errorKindLabel}</dd>
    </div>
    <div>
      <dt>短い理由</dt>
      <dd class="wrap-value">{errorReasonLabel}</dd>
    </div>
    <div>
      <dt>retryable</dt>
      <dd>{retryableLabel}</dd>
    </div>
  </dl>
</section>

<style>
  .phase-card {
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    display: grid;
    gap: 1rem;
    padding: 1.4rem;
  }

  .section-head {
    align-items: flex-start;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
  }

  .eyebrow {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
    margin: 0 0 0.25rem;
  }

  h3 {
    color: #fff6ea;
    margin: 0;
  }

  dt {
    color: rgba(236, 223, 205, 0.8);
    margin: 0 0 0.18rem;
  }

  dd {
    color: rgba(250, 242, 232, 0.9);
    margin: 0;
  }

  .detail-grid {
    display: grid;
    gap: 0.9rem 1rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin: 0;
  }

  .detail-grid div {
    display: grid;
    gap: 0.25rem;
    min-width: 0;
  }

  .wrap-value {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  @media (max-width: 900px) {
    .detail-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
