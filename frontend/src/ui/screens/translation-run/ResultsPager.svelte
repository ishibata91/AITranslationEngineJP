<script lang="ts">
  // 結果一覧の keyset ページャ。順次送り（前へ・次へ）と現在ページ番号だけを出す。
  // cursor ページングは順次送りのため番号ジャンプは持たない。state は持たず、
  // 表示値（現在ページ番号・前後の可否）と操作 callback を props で受ける。
  interface Props {
    // 現在ページ番号（1 始まり）。
    pageNumber: number
    // 前ページがあるか。先頭ページでは false にして「前へ」を無効化する。
    canPrev: boolean
    // 次ページがあるか。末尾ページでは false にして「次へ」を無効化する。
    canNext: boolean
    onPrev: () => void
    onNext: () => void
  }

  let { pageNumber, canPrev, canNext, onPrev, onNext }: Props = $props()
</script>

<nav
  class="flex items-center justify-center gap-3 pt-1"
  aria-label="結果ページ送り"
>
  <button
    type="button"
    class="btn btn-sm btn-ghost"
    disabled={!canPrev}
    onclick={onPrev}
  >
    ← 前へ
  </button>
  <span class="u-mono text-xs text-base-content/60 min-w-[5rem] text-center">
    ページ {pageNumber}
  </span>
  <button
    type="button"
    class="btn btn-sm btn-ghost"
    disabled={!canNext}
    onclick={onNext}
  >
    次へ →
  </button>
</nav>
