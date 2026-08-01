import { describe, expect, it } from "vitest"
import {
  BATCH_RETRY_UNTRANSLATED_LABEL,
  batchMainAction,
  batchUntranslatedNotice,
  untranslatedNotice
} from "./translation-run-presentation"

describe("untranslatedNotice", () => {
  it("未訳が複数残る場合は正確な件数と再実行対象を表示する", () => {
    expect(untranslatedNotice(3)).toBe(
      "3 件が未訳のまま残りました。もう一度実行すると、残った 3 件だけを翻訳します。"
    )
  })

  it("未訳が1件残る場合は1件の案内を表示する", () => {
    expect(untranslatedNotice(1)).toBe(
      "1 件が未訳のまま残りました。もう一度実行すると、残った 1 件だけを翻訳します。"
    )
  })

  it("未訳が0件の場合は案内を表示しない", () => {
    expect(untranslatedNotice(0)).toBe("")
  })
})

describe("batch 完了後の未訳表示", () => {
  it("未訳が複数残る場合は正確な件数と再送信操作を表示する", () => {
    expect(batchUntranslatedNotice(3)).toBe(
      "3 件が未訳のまま残りました。未訳だけを再送信できます。"
    )
    expect(
      batchMainAction(
        {
          stage: "done",
          total: 0,
          pending: 0,
          succeeded: 0,
          failed: 0,
          canApply: false,
          untranslatedCount: 3
        },
        true
      ).label
    ).toBe(BATCH_RETRY_UNTRANSLATED_LABEL)
  })

  it("未訳が1件の場合は1件の案内を表示する", () => {
    expect(batchUntranslatedNotice(1)).toBe(
      "1 件が未訳のまま残りました。未訳だけを再送信できます。"
    )
  })

  it("未訳が0件の場合は未訳案内を表示しない", () => {
    expect(batchUntranslatedNotice(0)).toBe("")
  })
})
