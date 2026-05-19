<script lang="ts">
  import ActionButton from "@ui/components/ActionButton.svelte"
  import ButtonGroup from "@ui/components/ButtonGroup.svelte"
  import InlineFeedback from "@ui/components/InlineFeedback.svelte"
  import SelectField from "@ui/components/SelectField.svelte"
  import TextInputField from "@ui/components/TextInputField.svelte"

  interface Props {
    targetGame: string
    outputPath: string
    pathState: string
    pathReason: string
    canGenerate: boolean
    canRegenerate: boolean
    isSubmitting: boolean
    disabledReason?: string
    onTargetGameChange: (value: string) => void
    onOutputPathInput: (value: string) => void
    onGenerate: () => void
    onRegenerate: () => void
  }

  let {
    targetGame,
    outputPath,
    pathState,
    pathReason,
    canGenerate,
    canRegenerate,
    isSubmitting,
    disabledReason = "",
    onTargetGameChange,
    onOutputPathInput,
    onGenerate,
    onRegenerate
  }: Props = $props()

  const targetGameOptions = [
    { value: "skyrim_se", label: "Skyrim SE" },
    { value: "skyrim_le", label: "Skyrim LE" }
  ]

  const pathHelp = $derived(
    pathReason || "出力先 path は .xml で終える必要があります。"
  )
</script>

<section
  class="output-card action-card"
  aria-labelledby="outputActionHeading"
  data-testid="output-management-output-actions"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">action rail</p>
      <h3 id="outputActionHeading">出力操作</h3>
    </div>
  </div>
  <SelectField
    id="outputTargetGame"
    label="target game"
    options={targetGameOptions}
    value={targetGame}
    onChange={onTargetGameChange}
  />
  <TextInputField
    id="outputPath"
    label="output path"
    value={outputPath}
    help={pathHelp}
    error={pathState === "invalid" ? pathHelp : ""}
    onInput={onOutputPathInput}
  />
  <ButtonGroup ariaLabel="出力操作" align="stretch">
    <ActionButton
      label="XML を出力"
      variant="primary"
      disabled={!canGenerate || isSubmitting}
      busy={isSubmitting}
      onClick={onGenerate}
    />
    <ActionButton
      label="再出力"
      variant="secondary"
      disabled={!canRegenerate || isSubmitting}
      onClick={onRegenerate}
    />
  </ButtonGroup>
  {#if disabledReason}
    <InlineFeedback tone="warning" message={disabledReason} />
  {/if}
</section>

<style>
  .output-card {
    border: 1px solid var(--line);
    border-radius: var(--radius-md);
    background: rgba(25, 22, 20, 0.82);
    box-shadow: var(--shadow);
    color: var(--text);
    padding: 1.25rem;
  }

  .action-card {
    display: grid;
    gap: 0.75rem;
  }

  .section-head {
    align-items: center;
    display: flex;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .eyebrow {
    color: var(--muted);
    font-size: 0.85rem;
  }
</style>
