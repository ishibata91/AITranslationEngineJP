<script lang="ts">
  import StickyActionFooter from "@ui/components/StickyActionFooter.svelte"

  interface Props {
    title: string
    titleId: string
    description: string
    reasons: string[]
    primaryLabel?: string
    primaryDisabled?: boolean
    showPrimary?: boolean
    showBack?: boolean
    showOutput?: boolean
    onPrimary?: () => void
    onBack?: () => void
    onOutput?: () => void
  }

  let {
    title,
    titleId,
    description,
    reasons,
    primaryLabel = "次へ進む",
    primaryDisabled = false,
    showPrimary = true,
    showBack = true,
    showOutput = false,
    onPrimary = () => undefined,
    onBack = () => undefined,
    onOutput = () => undefined
  }: Props = $props()
</script>

<div class="phase-footer-shell">
  <StickyActionFooter
    {description}
    emptyText="次の作業へ進めます。"
    onPrimary={onPrimary}
    primaryDisabled={primaryDisabled}
    primaryLabel={primaryLabel}
    reasons={reasons}
    {title}
    {titleId}
  >
    <div class="secondary-actions">
      {#if showBack}
        <button class="secondary-button" onclick={onBack} type="button">
          未完了一覧へ戻る
        </button>
      {/if}
      {#if showOutput}
        <button class="secondary-button" onclick={onOutput} type="button">
          出力管理へ進む
        </button>
      {/if}
      {#if !showPrimary}
        <span class="footer-note">この画面で進む先はありません。</span>
      {/if}
    </div>
  </StickyActionFooter>
</div>

<style>
  .phase-footer-shell {
    padding-bottom: max(1rem, env(safe-area-inset-bottom));
    min-width: 0;
    max-width: 100%;
  }

  .secondary-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.6rem;
    min-width: 0;
  }

  .secondary-button {
    border: 1px solid rgba(255, 212, 165, 0.22);
    border-radius: 0.8rem;
    background: rgba(255, 255, 255, 0.05);
    color: #fef3e8;
    cursor: pointer;
    max-width: 100%;
    padding: 0.65rem 0.85rem;
    white-space: normal;
  }

  .footer-note {
    color: rgba(252, 241, 232, 0.72);
    overflow-wrap: anywhere;
  }

  @media (max-width: 720px) {
    .secondary-actions {
      flex-direction: column;
    }

    .secondary-button {
      width: 100%;
    }
  }
</style>
