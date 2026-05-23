<script lang="ts">
  interface Props {
    onOpenJobManagement?: () => void
  }

  let { onOpenJobManagement = () => undefined }: Props = $props()
</script>

<section
  class="job-run-target-summary"
  data-testid="job-run-job-unselected-guidance"
>
  <div>
    <p class="eyebrow">ジョブ未選択</p>
    <h3>未完了ジョブ一覧でジョブを選んでください</h3>
  </div>
  <p class="selector-copy">
    翻訳段階の画面は選択済みジョブだけを対象にします。ジョブID
    の手入力はこのページでは扱いません。
  </p>
  <button class="secondary-button" onclick={onOpenJobManagement} type="button">
    一覧へ戻る
  </button>
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
  .selector-copy {
    color: rgba(236, 223, 205, 0.78);
  }

  h3 {
    margin: 0.2rem 0 0;
    color: #fff6ea;
  }

  .secondary-button {
    justify-self: start;
    min-height: 2.8rem;
    padding: 0.65rem 0.9rem;
    border-radius: 14px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    background: rgba(255, 255, 255, 0.05);
    color: #fff6ea;
    cursor: pointer;
  }
</style>
