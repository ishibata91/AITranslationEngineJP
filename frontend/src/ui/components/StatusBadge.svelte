<script lang="ts">
  // 状態を表す小さなバッジ。段階（未実行・実行中…）や訳状態（仮訳・未訳…）を共通の見た目で出す。
  // tone は意味トーンを受け取り、daisyUI class への対応は status-badge.ts が持つ。
  import { BADGE_TONE_CLASS } from "./status-badge"

  type Tone = keyof typeof BADGE_TONE_CLASS

  interface Props {
    label: string
    tone?: Tone
    loading?: boolean
  }

  let { label, tone = "neutral", loading = false }: Props = $props()
</script>

<span class="badge {BADGE_TONE_CLASS[tone]} badge-sm gap-2">
  {#if loading}
    <span class="loading loading-spinner loading-xs"></span>
  {/if}
  {label}
</span>
