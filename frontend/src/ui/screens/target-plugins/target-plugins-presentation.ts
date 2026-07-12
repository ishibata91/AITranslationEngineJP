// 翻訳対象 plugin 画面の表示用の値と純関数。型は target-plugins-view.ts に置く。
import type { TargetPluginRow } from "./target-plugins-view"
import type { BadgeTone } from "@ui/components/status-badge"

// plugin の翻訳進捗を、一覧に出すバッジのラベルとトーンへ変換した表示値。
export interface PluginProgressBadge {
  label: string
  tone: BadgeTone
}

// 訳済み行数と総行数から進捗バッジを組む。
// 総行数 0 は対象なし、全行訳済は完了、1 件も訳していなければ未着手、途中は翻訳中を出す。
export function pluginProgressBadge(row: TargetPluginRow): PluginProgressBadge {
  if (row.total <= 0) return { label: "対象なし", tone: "ghost" }
  const count = `${row.translated} / ${row.total}`
  if (row.translated >= row.total) return { label: `完了 ${count}`, tone: "success" }
  if (row.translated <= 0) return { label: `未着手 ${count}`, tone: "neutral" }
  return { label: `翻訳中 ${count}`, tone: "info" }
}
