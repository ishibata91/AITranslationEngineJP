<script lang="ts">
  interface Props {
    dictionaryLabels: string[]
    personaLabels: string[]
  }

  let { dictionaryLabels, personaLabels }: Props = $props()
</script>

<section
  class="job-setup-card"
  aria-labelledby="jobSetupFoundationHeading"
  data-testid="translation-job-setup-shared-dictionary-persona-region"
>
  <div class="section-head">
    <div>
      <h3 id="jobSetupFoundationHeading">共通辞書と共通ペルソナ</h3>
    </div>
  </div>
  <div class="foundation-grid split">
    <div class="foundation-table">
      <p class="mini-label">共通辞書</p>
      {#if dictionaryLabels.length === 0}
        <p class="empty-text">利用可能な共通辞書はありません。</p>
      {:else}
        <div class="foundation-scroll">
          <ul class="plain-list">
            {#each dictionaryLabels as label (label)}
              <li>{label}</li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
    <div class="foundation-table">
      <p class="mini-label">共通ペルソナ</p>
      {#if personaLabels.length === 0}
        <p class="empty-text">利用可能な共通ペルソナはありません。</p>
      {:else}
        <div class="foundation-scroll">
          <ul class="plain-list">
            {#each personaLabels as label (label)}
              <li>{label}</li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
  </div>
</section>

<style>
  .job-setup-card,
  .foundation-grid,
  .foundation-table {
    display: grid;
    gap: 0.75rem;
  }

  .job-setup-card {
    gap: 1rem;
    padding: 1.25rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 1.25rem;
    background: rgba(34, 26, 23, 0.82);
    box-shadow: 0 20px 40px rgba(6, 4, 3, 0.18);
    color: var(--text);
  }

  .section-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
  }

  .foundation-grid.split {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .foundation-scroll {
    max-height: 10rem;
    overflow: auto;
    padding-right: 0.35rem;
  }

  .mini-label {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.78rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .empty-text,
  .plain-list {
    color: rgba(252, 241, 232, 0.86);
  }

  .plain-list {
    margin: 0;
    padding-left: 1.1rem;
  }

  .plain-list li {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  @media (max-width: 720px) {
    .section-head {
      flex-direction: column;
    }

    .foundation-grid.split {
      grid-template-columns: minmax(0, 1fr);
    }
  }
</style>
