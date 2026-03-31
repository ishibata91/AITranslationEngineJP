<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { JobCreateScreenState } from "@application/usecases/job-create";

  export let state = undefined as unknown as JobCreateScreenState;

  const dispatch = createEventDispatcher<{
    resetResult: void;
    submit: void;
    updateSourceGroupField: {
      field: "sourceJsonPath" | "targetPlugin";
      groupIndex: number;
      value: string;
    };
    updateTranslationUnitField: {
      field:
        | "sourceEntityType"
        | "formId"
        | "editorId"
        | "recordSignature"
        | "fieldName"
        | "extractionKey"
        | "sourceText"
        | "sortKey";
      groupIndex: number;
      unitIndex: number;
      value: string;
    };
  }>();

  function primaryGroup() {
    return state.request.sourceGroups[0];
  }

  function primaryUnit() {
    return primaryGroup().translationUnits[0];
  }
</script>

<section class="panel">
  <header class="header">
    <div>
      <p class="eyebrow">Phase 1 UI Path</p>
      <h1>Job Create</h1>
    </div>
    <button
      type="button"
      class="submit"
      disabled={state.isSubmitting}
      on:click={() => dispatch("submit")}
    >
      {state.isSubmitting ? "Creating job..." : "Create job"}
    </button>
  </header>

  <p class="lede">
    Edit the minimal grouped request payload for the Phase 1 create path and submit one backend-aligned
    job create request.
  </p>

  <div class="grid">
    <section class="card">
      <h2>Source Group</h2>

      <label>
        <span>Source JSON Path</span>
        <input
          type="text"
          value={primaryGroup().sourceJsonPath}
          disabled={state.isSubmitting}
          on:input={(event) =>
            dispatch("updateSourceGroupField", {
              field: "sourceJsonPath",
              groupIndex: 0,
              value: event.currentTarget.value
            })}
        />
      </label>

      <label>
        <span>Target Plugin</span>
        <input
          type="text"
          value={primaryGroup().targetPlugin}
          disabled={state.isSubmitting}
          on:input={(event) =>
            dispatch("updateSourceGroupField", {
              field: "targetPlugin",
              groupIndex: 0,
              value: event.currentTarget.value
            })}
        />
      </label>
    </section>

    <section class="card">
      <h2>Translation Unit</h2>

      <div class="field-grid">
        <label>
          <span>Entity Type</span>
          <input
            type="text"
            value={primaryUnit().sourceEntityType}
            disabled={state.isSubmitting}
            on:input={(event) =>
              dispatch("updateTranslationUnitField", {
                field: "sourceEntityType",
                groupIndex: 0,
                unitIndex: 0,
                value: event.currentTarget.value
              })}
          />
        </label>

        <label>
          <span>Form ID</span>
          <input
            type="text"
            value={primaryUnit().formId}
            disabled={state.isSubmitting}
            on:input={(event) =>
              dispatch("updateTranslationUnitField", {
                field: "formId",
                groupIndex: 0,
                unitIndex: 0,
                value: event.currentTarget.value
              })}
          />
        </label>

        <label>
          <span>Editor ID</span>
          <input
            type="text"
            value={primaryUnit().editorId}
            disabled={state.isSubmitting}
            on:input={(event) =>
              dispatch("updateTranslationUnitField", {
                field: "editorId",
                groupIndex: 0,
                unitIndex: 0,
                value: event.currentTarget.value
              })}
          />
        </label>

        <label>
          <span>Record Signature</span>
          <input
            type="text"
            value={primaryUnit().recordSignature}
            disabled={state.isSubmitting}
            on:input={(event) =>
              dispatch("updateTranslationUnitField", {
                field: "recordSignature",
                groupIndex: 0,
                unitIndex: 0,
                value: event.currentTarget.value
              })}
          />
        </label>

        <label>
          <span>Field Name</span>
          <input
            type="text"
            value={primaryUnit().fieldName}
            disabled={state.isSubmitting}
            on:input={(event) =>
              dispatch("updateTranslationUnitField", {
                field: "fieldName",
                groupIndex: 0,
                unitIndex: 0,
                value: event.currentTarget.value
              })}
          />
        </label>

        <label>
          <span>Extraction Key</span>
          <input
            type="text"
            value={primaryUnit().extractionKey}
            disabled={state.isSubmitting}
            on:input={(event) =>
              dispatch("updateTranslationUnitField", {
                field: "extractionKey",
                groupIndex: 0,
                unitIndex: 0,
                value: event.currentTarget.value
              })}
          />
        </label>

        <label>
          <span>Sort Key</span>
          <input
            type="text"
            value={primaryUnit().sortKey}
            disabled={state.isSubmitting}
            on:input={(event) =>
              dispatch("updateTranslationUnitField", {
                field: "sortKey",
                groupIndex: 0,
                unitIndex: 0,
                value: event.currentTarget.value
              })}
          />
        </label>
      </div>

      <label>
        <span>Source Text</span>
        <textarea
          rows="4"
          disabled={state.isSubmitting}
          on:input={(event) =>
            dispatch("updateTranslationUnitField", {
              field: "sourceText",
              groupIndex: 0,
              unitIndex: 0,
              value: event.currentTarget.value
            })}
        >{primaryUnit().sourceText}</textarea>
      </label>
    </section>
  </div>

  {#if state.error}
    <div class="status error">
      <strong>Create failed</strong>
      <p>{state.error}</p>
    </div>
  {/if}

  {#if state.result}
    <aside class="status success">
      <div>
        <p class="status-label">Created Job</p>
        <h2>{state.result.jobId}</h2>
        <p>Observable state: {state.result.state}</p>
      </div>
      <button type="button" on:click={() => dispatch("resetResult")}>Clear result</button>
    </aside>
  {/if}
</section>

<style>
  .panel {
    width: min(100%, 72rem);
    display: grid;
    gap: 1.25rem;
    padding: 2rem;
    border: 1px solid #d7e1eb;
    border-radius: 1.25rem;
    background:
      linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(244, 248, 252, 0.96)),
      #fff;
    box-shadow: 0 24px 60px rgba(20, 32, 51, 0.12);
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: start;
    gap: 1rem;
  }

  .eyebrow {
    margin: 0 0 0.5rem;
    font-size: 0.8rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: #5b6a7d;
  }

  h1,
  h2,
  p {
    margin: 0;
  }

  h1 {
    font-size: clamp(2rem, 4vw, 3rem);
  }

  .lede {
    line-height: 1.6;
    color: #334155;
  }

  .submit {
    padding: 0.85rem 1.2rem;
    border: 0;
    border-radius: 999px;
    background: #142033;
    color: #fff;
    cursor: pointer;
  }

  .submit:disabled {
    opacity: 0.7;
    cursor: progress;
  }

  .grid {
    display: grid;
    gap: 1rem;
  }

  .card {
    display: grid;
    gap: 1rem;
    padding: 1.25rem;
    border-radius: 1rem;
    background: #f8fafc;
    border: 1px solid #d7e1eb;
  }

  .field-grid {
    display: grid;
    gap: 0.85rem;
  }

  label {
    display: grid;
    gap: 0.35rem;
  }

  span {
    font-size: 0.85rem;
    color: #526277;
  }

  input,
  textarea {
    width: 100%;
    box-sizing: border-box;
    border: 1px solid #c7d2df;
    border-radius: 0.8rem;
    padding: 0.7rem 0.85rem;
    background: #fff;
    color: inherit;
  }

  textarea {
    resize: vertical;
  }

  .status {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    padding: 1rem 1.2rem;
    border-radius: 1rem;
  }

  .error {
    background: #fff1f2;
    color: #9f1239;
  }

  .success {
    background: #edf7f1;
    color: #166534;
    align-items: center;
  }

  .status-label {
    margin-bottom: 0.35rem;
    font-size: 0.8rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .success button {
    padding: 0.7rem 1rem;
    border: 1px solid currentColor;
    border-radius: 999px;
    background: transparent;
    color: inherit;
    cursor: pointer;
  }

  @media (min-width: 900px) {
    .grid {
      grid-template-columns: minmax(0, 1fr) minmax(0, 1.4fr);
      align-items: start;
    }

    .field-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }

  @media (max-width: 700px) {
    .header,
    .status {
      flex-direction: column;
    }
  }
</style>
