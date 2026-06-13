// StatusBadge の意味トーンと、daisyUI badge class への対応。
// 型を .svelte の外（.ts）に置き、tsc から named import できるようにする。

// 状態を表す意味トーン。daisyUI の class 名を呼び出し側へ漏らさず、意味で指定させる。
export type BadgeTone =
  | "neutral"
  | "primary"
  | "accent"
  | "success"
  | "warning"
  | "danger"
  | "info"
  | "ghost"
  | "outline"

// 意味トーンを daisyUI の badge トーン class へ対応づける。
export const BADGE_TONE_CLASS: Record<BadgeTone, string> = {
  neutral: "badge-neutral",
  primary: "badge-primary",
  accent: "badge-accent",
  success: "badge-success",
  warning: "badge-warning",
  danger: "badge-error",
  info: "badge-info",
  ghost: "badge-ghost",
  outline: "badge-outline"
}
