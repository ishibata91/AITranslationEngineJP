<script lang="ts">
  // ラベル・補助文・control を 1 組にする汎用の入力枠。control の中身は children スロットで受ける。
  // これにより TextField・SelectField・FileSelectField がラベルと補助文の見た目を共有する。
  import type { Snippet } from "svelte"

  interface Props {
    label: string
    hint?: string
    children: Snippet
  }

  let { label, hint = "", children }: Props = $props()
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend text-base-content/80">{label}</legend>
  {@render children()}
  {#if hint.length > 0}
    <p class="label text-base-content/50">{hint}</p>
  {/if}
</fieldset>
