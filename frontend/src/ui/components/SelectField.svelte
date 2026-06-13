<script lang="ts">
  // 汎用の選択欄。選択肢からの選択を扱う。trailing の補助操作（取得・更新ボタンなど）を action スロットで足せる。
  // 値の保持はしない。表示と change event の中継だけを行う。
  import type { Snippet } from "svelte"
  import Field from "./Field.svelte"

  interface Props {
    label: string
    value: string
    options: ReadonlyArray<string>
    // 選択肢があるが未選択のときの先頭プレースホルダ。
    placeholder?: string
    // 選択肢が空のときの表示。
    emptyText?: string
    hint?: string
    disabled?: boolean
    onChange: (value: string) => void
    // select の右に並べる補助操作（任意）。
    action?: Snippet
  }

  let {
    label,
    value,
    options,
    placeholder = "選択してください",
    emptyText = "選択肢がありません",
    hint = "",
    disabled = false,
    onChange,
    action
  }: Props = $props()

  function handle(event: Event) {
    onChange((event.currentTarget as HTMLSelectElement).value)
  }
</script>

<Field {label} {hint}>
  <div class="flex gap-2">
    <select
      class="select w-full bg-base-100/60 u-mono text-sm"
      {value}
      disabled={disabled || options.length === 0}
      onchange={handle}
    >
      {#if options.length === 0}
        <option value="" disabled selected>{emptyText}</option>
      {:else}
        <option value="" disabled>{placeholder}</option>
        {#each options as option (option)}
          <option value={option}>{option}</option>
        {/each}
      {/if}
    </select>
    {#if action}
      {@render action()}
    {/if}
  </div>
</Field>
