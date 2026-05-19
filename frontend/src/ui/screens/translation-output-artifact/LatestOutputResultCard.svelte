<script lang="ts">
  import type { TranslationOutputArtifactCommandSnapshot } from "@application/contract/translation-output-artifact/translation-output-artifact-screen-types"

  import { formatCount, formatStatus } from "./output-artifact-formatters"

  interface Props {
    lastCommand: TranslationOutputArtifactCommandSnapshot | null
  }

  let { lastCommand }: Props = $props()
</script>

{#if lastCommand}
  <div class="notice-block" data-testid="output-management-latest-result">
    <h4>result summary</h4>
    <dl class="detail-grid compact">
      <div>
        <dt>artifact status</dt>
        <dd>{formatStatus(lastCommand.artifactStatus)}</dd>
      </div>
      <div>
        <dt>row count</dt>
        <dd>{formatCount(lastCommand.rowCount)}</dd>
      </div>
      <div>
        <dt>file path</dt>
        <dd class="wrap-value">{lastCommand.filePath ?? "-"}</dd>
      </div>
      <div>
        <dt>target game</dt>
        <dd>{lastCommand.targetGame}</dd>
      </div>
    </dl>
    {#if lastCommand.errorReason}
      <p class="helper-text warning">{lastCommand.errorReason}</p>
    {/if}
  </div>
{/if}

<style>
  .notice-block,
  .detail-grid {
    display: grid;
    gap: 0.75rem;
  }

  .notice-block {
    border: 1px solid var(--line);
    border-radius: 12px;
    padding: 0.85rem;
    background: rgba(18, 16, 15, 0.88);
    color: var(--text);
  }

  .notice-block h4 {
    font-weight: 700;
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

  .helper-text {
    line-height: 1.6;
  }

  .helper-text.warning {
    color: #ff9f7f;
  }

  @media (max-width: 720px) {
    .detail-grid.compact {
      grid-template-columns: 1fr;
    }
  }
</style>
