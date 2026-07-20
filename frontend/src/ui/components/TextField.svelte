<script lang="ts">
  // 汎用のテキスト入力。ラベル・補助文は Field に委ね、値の保持はしない。
  // 表示と入力 event の中継だけを行い、値と変更通知は props で受け渡す。
  import Field from "./Field.svelte"

  interface Props {
    label: string
    value: string
    placeholder?: string
    hint?: string
    secret?: boolean
    // 入力の右に出す静的な単位表示（例: "k"）。無ければ単位を出さない。
    suffix?: string
    // 入力ボックスの幅クラス。既定は w-full（セル幅いっぱい）。数桁だけ入れる欄は狭い幅（例: "w-24"）を渡す。
    inputWidthClass?: string
    oninput: (value: string) => void
  }

  let {
    label,
    value,
    placeholder = "",
    hint = "",
    secret = false,
    suffix = "",
    inputWidthClass = "w-full",
    oninput
  }: Props = $props()

  function handle(event: Event) {
    oninput((event.currentTarget as HTMLInputElement).value)
  }
</script>

<Field {label} {hint}>
  <div class="flex items-center gap-2">
    <input
      class="input {inputWidthClass} min-w-0 bg-base-100/60 u-mono text-sm"
      type={secret ? "password" : "text"}
      {value}
      {placeholder}
      autocomplete="off"
      spellcheck="false"
      oninput={handle}
    />
    {#if suffix.length > 0}
      <span class="u-mono text-sm text-base-content/60 shrink-0">{suffix}</span>
    {/if}
  </div>
</Field>
