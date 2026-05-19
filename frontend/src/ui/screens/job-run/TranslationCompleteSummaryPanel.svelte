<script lang="ts">
  interface Props {
    jobId: number
  }

  let { jobId }: Props = $props()
</script>

<section
  class="complete-card hero-card"
  data-testid="translation-complete-completion-summary"
>
  <div class="hero-head">
    <div>
      <p class="eyebrow">translation complete</p>
      <h2>翻訳完了</h2>
    </div>
    <span class="status-pill">ジョブ #{jobId}</span>
  </div>
  <p class="lead">
    本文翻訳が完了したジョブの原文と訳文を確認します。成果物の出力処理は出力管理で扱います。
  </p>
</section>

<style>
  @import "./translation-complete-card.css";

  .hero-card {
    gap: 0.85rem;
  }

  .lead {
    color: rgba(236, 223, 205, 0.76);
  }

  h2 {
    color: #fff6ea;
    margin: 0;
  }

  .status-pill {
    border: 1px solid rgba(255, 212, 165, 0.22);
    border-radius: 999px;
    color: #ffd8ae;
    padding: 0.25rem 0.65rem;
  }
</style>
