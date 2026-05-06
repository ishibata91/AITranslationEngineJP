import { describe, expect, test } from "vitest"

import type { TranslationInputReviewItem } from "@application/gateway-contract/translation-input"

import { canOpenJobSetup } from "./translation-input.presenter"

function createItem(
  status: TranslationInputReviewItem["status"]
): TranslationInputReviewItem {
  return {
    localId: `item-${status}`,
    inputId: 41,
    fileName: "input-review.json",
    filePath: "/mods/input-review.json",
    fileHash: "hash-41",
    importTimestamp: "2026-04-26T09:30:00Z",
    status,
    accepted: status === "registered" || status === "warning",
    canRebuild: status !== "failed",
    lastAction: "import",
    errorKind: status === "failed" ? "invalid_json" : null,
    warnings:
      status === "warning"
        ? [
            {
              kind: "unknown_field_definition",
              recordType: "BOOK",
              subrecordType: "DESC",
              message: "unknown description field"
            }
          ]
        : [],
    summary: null
  }
}

describe("canOpenJobSetup", () => {
  test("registered と warning では true を返す", () => {
    expect(canOpenJobSetup(createItem("registered"))).toBe(true)
    expect(canOpenJobSetup(createItem("warning"))).toBe(true)
  })

  test("failed と rebuild-required では false を返す", () => {
    expect(canOpenJobSetup(createItem("failed"))).toBe(false)
    expect(canOpenJobSetup(createItem("rebuild-required"))).toBe(false)
  })

  test("選択 item がない時は false を返す", () => {
    expect(canOpenJobSetup(null)).toBe(false)
  })
})
