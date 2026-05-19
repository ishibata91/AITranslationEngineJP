<script lang="ts">
  import type { TranslationInputReviewItem } from "@application/gateway-contract/translation-input"

  interface Props {
    items: TranslationInputReviewItem[]
    selectedItemId: string | null
    totalItemCountLabel: string
    emptyStateText: string
    formatStatus: (status: string) => string
    formatDate: (timestamp: string) => string
    formatErrorKind: (errorKind: string | null) => string
    onSelectItem: (localId: string) => void
  }

  let {
    items,
    selectedItemId,
    totalItemCountLabel,
    emptyStateText,
    formatStatus,
    formatDate,
    formatErrorKind,
    onSelectItem
  }: Props = $props()
</script>

<section
  class="panel list-panel"
  aria-labelledby="inputReviewListHeading"
  data-testid="translation-input-review-loaded-input-list"
>
  <div class="section-head">
    <div class="title-stack">
      <p class="eyebrow">読み込み済みデータ</p>
      <h3 id="inputReviewListHeading">読み込み済みデータ</h3>
      <p class="support-copy">
        登録済みの JSON を選び、状態と次の確認対象を見つける。
      </p>
    </div>
    <p class="count-pill">{totalItemCountLabel}</p>
  </div>

  {#if items.length === 0}
    <div class="empty-state">
      <p>{emptyStateText}</p>
    </div>
  {:else}
    <div class="review-list" role="list">
      {#each items as item (item.localId)}
        <div role="listitem">
          <button
            aria-pressed={item.localId === selectedItemId ? "true" : "false"}
            class="review-item"
            class:is-selected={item.localId === selectedItemId}
            onclick={() => onSelectItem(item.localId)}
            type="button"
          >
            <div class="review-item-head">
              <div class="name-stack">
                <strong>{item.fileName}</strong>
                <p>{item.filePath}</p>
              </div>
              <span class="status-pill">{formatStatus(item.status)}</span>
            </div>

            <dl class="review-meta">
              <div>
                <dt>登録結果</dt>
                <dd>
                  {item.accepted ? "登録済み" : "登録に失敗"}
                  <span class="sr-only">
                    {item.accepted ? "accepted" : "rejected"}
                  </span>
                </dd>
              </div>
              <div>
                <dt>再構築</dt>
                <dd>
                  {item.canRebuild ? "再構築できます" : "まだ再構築できません"}
                  <span class="sr-only">
                    {item.canRebuild ? "rebuild 可" : "不可"}
                  </span>
                </dd>
              </div>
              <div>
                <dt>読み込み日時</dt>
                <dd>{formatDate(item.importTimestamp)}</dd>
              </div>
              <div>
                <dt>内容ハッシュ</dt>
                <dd class="hash-text">{item.fileHash}</dd>
              </div>
            </dl>

            <div class="inline-note">
              <span>問題区分</span>
              <strong>{formatErrorKind(item.errorKind)}</strong>
            </div>
          </button>
        </div>
      {/each}
    </div>
  {/if}
</section>

<style>
  .panel {
    border: 1px solid var(--line);
    border-radius: 20px;
    padding: 1.25rem;
    background: rgba(28, 23, 20, 0.74);
    box-shadow: var(--shadow);
    color: var(--text);
    backdrop-filter: blur(18px);
  }

  .list-panel {
    display: grid;
    gap: 1rem;
  }

  .section-head,
  .review-item-head {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: start;
  }

  .title-stack,
  .name-stack {
    display: grid;
    gap: 0.35rem;
  }

  .eyebrow,
  .support-copy,
  .name-stack p,
  dt,
  .inline-note span {
    color: var(--muted);
  }

  .eyebrow {
    margin: 0;
    font-size: 0.76rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
  }

  h3,
  p,
  dl,
  dt,
  dd {
    margin: 0;
  }

  .count-pill,
  .status-pill {
    margin: 0;
    padding: 0.32rem 0.7rem;
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    color: var(--primary);
    font-size: 0.82rem;
  }

  .review-list {
    display: grid;
    gap: 0.8rem;
  }

  .review-item {
    width: 100%;
    display: grid;
    gap: 0.85rem;
    text-align: left;
    padding: 1rem;
    border-radius: 18px;
    border: 1px solid rgba(255, 186, 56, 0.18);
    background: rgba(22, 18, 17, 0.78);
    color: var(--text);
    cursor: pointer;
    font: inherit;
  }

  .review-item.is-selected {
    border-color: var(--primary);
    box-shadow: 0 0 0 1px rgba(255, 186, 56, 0.18);
  }

  .review-meta {
    display: grid;
    gap: 0.8rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .review-meta div {
    display: grid;
    gap: 0.22rem;
  }

  dt {
    font-size: 0.8rem;
  }

  .hash-text {
    word-break: break-all;
  }

  .inline-note,
  .empty-state {
    display: grid;
    gap: 0.25rem;
    padding: 0.95rem 1rem;
    border-radius: 18px;
    background: rgba(18, 15, 14, 0.76);
    border: 1px solid rgba(255, 186, 56, 0.12);
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
  }

  @media (max-width: 720px) {
    .section-head,
    .review-item-head {
      flex-direction: column;
    }
  }
</style>
