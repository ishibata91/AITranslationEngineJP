<script lang="ts">
  // 汎用のファイル選択欄。選択ボタンと、選択結果（ファイル名・フルパス）または未選択表示を出す。
  // ファイルダイアログ自体は持たず、押下を onSelect で中継する。実体は呼び出し側が配線する。
  import Field from "./Field.svelte"

  interface Props {
    label: string
    buttonLabel: string
    path: string
    emptyText?: string
    hint?: string
    onSelect: () => void
  }

  let {
    label,
    buttonLabel,
    path,
    emptyText = "ファイルが未選択です。",
    hint = "",
    onSelect
  }: Props = $props()

  // フルパスからファイル名だけを取り出す。Windows と Unix の区切りに対応する。
  function fileName(value: string): string {
    const segments = value.split(/[\\/]/)
    return segments[segments.length - 1] ?? value
  }
</script>

<Field {label} {hint}>
  <div class="flex flex-wrap items-center gap-4">
    <button
      type="button"
      class="btn btn-outline btn-primary gap-2"
      onclick={onSelect}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="size-4"
        aria-hidden="true"
      >
        <path d="M14 3v4a1 1 0 0 0 1 1h4" />
        <path d="M5 3h9l5 5v11a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z" />
      </svg>
      {buttonLabel}
    </button>
    {#if path.length > 0}
      <div class="flex min-w-0 flex-col">
        <span class="u-mono text-sm text-base-content">{fileName(path)}</span>
        <span class="u-mono text-xs text-base-content/45 truncate">{path}</span>
      </div>
    {:else}
      <span class="text-base-content/45">{emptyText}</span>
    {/if}
  </div>
</Field>
