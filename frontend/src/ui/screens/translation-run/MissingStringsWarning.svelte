<script lang="ts">
  // 片側 Strings 欠け時の警告表示。表示専用で、判定結果を props で受けるだけにする。
  // state・API・Wails・副作用は持たない。判定と供給は container（implementation-module）が担う。
  import type { StringsPresence } from "./translation-run-view"

  interface Props {
    // 翻訳対象 plugin の Data フォルダにある Strings の言語別有無。
    // 未指定（未判定）または両方あるときは何も表示しない。
    presence?: StringsPresence
  }

  let { presence = undefined }: Props = $props()

  // 欠けパターンごとの文言。見出し＋状態と理由＋対処を欠けている側に合わせて組む。
  // 影響（既存訳を再利用できない）は全パターン共通で impact に固定する。
  const message = $derived.by(() => {
    if (!presence) return undefined
    if (!presence.english && !presence.japanese) {
      return {
        title: "英語と日本語の Strings ファイルが見つかりません",
        state:
          "翻訳対象 plugin と同じ Data フォルダの strings/ に、英語の Strings ファイル（*_english.strings / .dlstrings / .ilstrings）と日本語の Strings ファイル（*_japanese.strings / .dlstrings / .ilstrings）のどちらもありません。",
        reason:
          "英語原文と日本語既訳を突き合わせて既存訳の対を作るため、英語と日本語の両方の Strings が必要です。",
        remedy:
          "英語版と日本語版の Strings ファイルを Data フォルダの strings/ に置いてから、抽出をやり直してください。"
      }
    }
    if (!presence.japanese) {
      return {
        title: "日本語の Strings ファイルが見つかりません",
        state:
          "翻訳対象 plugin と同じ Data フォルダの strings/ に、日本語の Strings ファイル（*_japanese.strings / .dlstrings / .ilstrings）がありません。",
        reason:
          "英語原文と日本語既訳を突き合わせて既存訳の対を作るため、英語だけでなく日本語の Strings も必要です。",
        remedy:
          "日本語版の Strings ファイルを Data フォルダの strings/ に置いてから、抽出をやり直してください。"
      }
    }
    if (!presence.english) {
      return {
        title: "英語の Strings ファイルが見つかりません",
        state:
          "翻訳対象 plugin と同じ Data フォルダの strings/ に、英語の Strings ファイル（*_english.strings / .dlstrings / .ilstrings）がありません。",
        reason:
          "英語原文と日本語既訳を突き合わせて既存訳の対を作るため、日本語だけでなく英語の Strings も必要です。",
        remedy:
          "英語版の Strings ファイルを Data フォルダの strings/ に置いてから、抽出をやり直してください。"
      }
    }
    return undefined
  })

  // 影響は欠けパターンによらず同じ。
  const impact =
    "このまま翻訳すると、既存訳（参照訳・固有名の確定訳語）を再利用できず、全文を AI 翻訳します。固有名の訳が公式訳と揃わない可能性があります。"
</script>

{#if message}
  <div role="alert" class="alert alert-warning alert-soft">
    <div class="flex flex-col gap-1">
      <span class="font-semibold">{message.title}</span>
      <span class="text-sm">{message.state}{message.reason}</span>
      <span class="text-sm">{impact}</span>
      <span class="text-sm">{message.remedy}</span>
    </div>
  </div>
{/if}
