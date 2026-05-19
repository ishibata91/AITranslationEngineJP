<script lang="ts">
  interface Props {
    errorMessage: string
  }

  let { errorMessage }: Props = $props()
</script>

<section
  class="job-setup-card hero-card"
  data-testid="translation-job-setup-screen-purpose-header"
>
  <h2>翻訳段階ごとの AI 設定</h2>
  <p class="lead">
    入力済みデータを確認し、3 つの翻訳段階で使う
    AIサービスとモデルを選びます。
  </p>
  <p class="error-text" hidden={!errorMessage}>
    {errorMessage}
  </p>
</section>

<style>
  .job-setup-card {
    display: grid;
    gap: 1rem;
    padding: 1.25rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 1.25rem;
    background: rgba(34, 26, 23, 0.82);
    box-shadow: 0 20px 40px rgba(6, 4, 3, 0.18);
    color: var(--text);
  }

  .hero-card {
    gap: 0.6rem;
  }

  .lead {
    color: rgba(252, 241, 232, 0.86);
  }

  .error-text {
    color: #ffb4ab;
    overflow-wrap: anywhere;
  }
</style>
