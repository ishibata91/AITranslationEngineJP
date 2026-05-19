<script lang="ts">
  interface Props {
    pageTitle: string
    pageLead: string
    onOpenInputReview: () => void
  }

  let { pageTitle, pageLead, onOpenInputReview }: Props = $props()
</script>

<section
  class="panel job-management-hero"
  data-testid="translation-job-management-screen-header"
>
  <div class="hero-copy">
    <h2>{pageTitle}</h2>
    <p class="hero-lead">{pageLead}</p>
  </div>
  <button
    class="ghost-button"
    onclick={onOpenInputReview}
    type="button"
  >
    新規翻訳を開始
  </button>
</section>

<style>
  .panel {
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
    padding: 1.1rem 1.2rem;
  }

  .job-management-hero {
    display: flex;
    gap: 1rem;
    justify-content: space-between;
    align-items: flex-start;
    flex-wrap: wrap;
  }

  h2 {
    margin: 0;
    color: #fff6ea;
  }

  .hero-lead {
    margin: 0.35rem 0 0;
    color: rgba(236, 223, 205, 0.78);
  }

  button {
    min-height: 2.8rem;
    padding: 0.65rem 1rem;
    border-radius: 14px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    font: inherit;
    cursor: pointer;
  }

  .ghost-button {
    background: rgba(255, 255, 255, 0.06);
    color: #fff6ea;
  }
</style>
