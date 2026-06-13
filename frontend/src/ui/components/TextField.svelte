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
    oninput: (value: string) => void
  }

  let {
    label,
    value,
    placeholder = "",
    hint = "",
    secret = false,
    oninput
  }: Props = $props()

  function handle(event: Event) {
    oninput((event.currentTarget as HTMLInputElement).value)
  }
</script>

<Field {label} {hint}>
  <input
    class="input w-full bg-base-100/60 u-mono text-sm"
    type={secret ? "password" : "text"}
    {value}
    {placeholder}
    autocomplete="off"
    spellcheck="false"
    oninput={handle}
  />
</Field>
