import { describe, expect, it } from "vitest"
import { untranslatedNotice } from "./translation-run-presentation"

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
