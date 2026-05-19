<script lang="ts">
  import type { BodyTranslationFieldResultItem } from "@application/contract/body-translation-phase"
  import type { FieldResultListPanelProps } from "./body-phase-card-types"

  let { availabilityLabel, items }: FieldResultListPanelProps = $props()

  function rowKey(item: BodyTranslationFieldResultItem, index: number): string {
    return `${item.fieldId}-${item.fieldLabel}-${index}`
  }
</script>

<section
  class="phase-card"
  aria-labelledby="bodyFieldResultHeading"
  data-testid="body-translation-phase-field-result-list"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">field result list</p>
      <h3 id="bodyFieldResultHeading">field result</h3>
    </div>
    <span class="mini-text">{availabilityLabel}</span>
  </div>

  {#if items.length === 0}
    <p class="detail-note">
      field result 一覧はまだ返っていません。現在は summary
      だけを表示しています。
    </p>
  {:else}
    <div class="field-result-list">
      {#each items as item, index (rowKey(item, index))}
        <article class="field-result-row">
          <div class="field-result-head">
            <strong>{item.fieldLabel}</strong>
            <span>{item.fieldId}</span>
          </div>
          <dl class="field-result-grid">
            <div>
              <dt>record type</dt>
              <dd>{item.recordTypeLabel}</dd>
            </div>
            <div>
              <dt>field type</dt>
              <dd>{item.fieldTypeLabel}</dd>
            </div>
            <div>
              <dt>FormID</dt>
              <dd>{item.formIdLabel}</dd>
            </div>
            <div>
              <dt>EditorID</dt>
              <dd>{item.editorIdLabel}</dd>
            </div>
            <div>
              <dt>source excerpt</dt>
              <dd>{item.sourceExcerpt}</dd>
            </div>
            <div>
              <dt>translated text</dt>
              <dd>{item.translatedText}</dd>
            </div>
            <div>
              <dt>output status</dt>
              <dd>{item.outputStatus}</dd>
            </div>
            <div>
              <dt>protection validation</dt>
              <dd>{item.protectionValidation}</dd>
            </div>
            <div>
              <dt>retry count</dt>
              <dd>{item.retryCountLabel}</dd>
            </div>
          </dl>
        </article>
      {/each}
    </div>
  {/if}
</section>

<style>
  .phase-card {
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    display: grid;
    gap: 1rem;
    min-width: 0;
    padding: 1.4rem;
  }

  .section-head {
    align-items: flex-start;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
    min-width: 0;
  }

  .eyebrow,
  .mini-text,
  .detail-note {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
    margin: 0 0 0.25rem;
  }

  h3 {
    color: #fff6ea;
    margin: 0;
  }

  dt {
    color: rgba(236, 223, 205, 0.78);
    font-size: 0.82rem;
  }

  dd {
    color: #f1e6d6;
    margin: 0.25rem 0 0;
    overflow-wrap: anywhere;
  }

  .field-result-list,
  .field-result-grid {
    display: grid;
    gap: 0.8rem;
    min-width: 0;
  }

  .field-result-row {
    background: rgba(255, 255, 255, 0.04);
    border-radius: 16px;
    display: grid;
    gap: 0.8rem;
    min-width: 0;
    padding: 0.85rem 0.95rem;
  }

  .field-result-head {
    align-items: flex-start;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
    min-width: 0;
  }

  .field-result-grid {
    grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
  }

  .field-result-head strong,
  .field-result-head span,
  .mini-text,
  .detail-note,
  dd {
    overflow-wrap: anywhere;
  }

  @media (max-width: 900px) {
    .field-result-head {
      align-items: stretch;
      flex-direction: column;
    }
  }

  @media (max-width: 480px) {
    .phase-card {
      border-radius: 14px;
      padding: 1rem;
    }
  }
</style>
