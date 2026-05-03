import { describe, expect, test } from "vitest"

import { TranslationJobSetupStore } from "./translation-job-setup.store"

describe("TranslationJobSetupStore", () => {
  test("snapshot は null の targetSlices と passSlices を空配列へ正規化する", () => {
    const store = new TranslationJobSetupStore()

    store.update((draft) => {
      draft.validationResult = {
        status: "pass",
        targetSlices: null as unknown as string[],
        validatedAt: "2026-05-03T06:58:30Z",
        canCreate: true,
        passSlices: null as unknown as string[]
      }
      draft.validationState = "fresh"
      draft.phase = "ready"
    })

    expect(store.snapshot().validationResult).toEqual({
      status: "pass",
      targetSlices: [],
      validatedAt: "2026-05-03T06:58:30Z",
      canCreate: true,
      passSlices: []
    })
  })
})
