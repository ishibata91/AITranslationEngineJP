import { describe, expect, test } from "vitest"

import { resolveReviewFakeApiRuntimeContext } from "./review-fake-api-runtime"

function createSearchParams(values: Record<string, string | null>) {
  return {
    get(name: string): string | null {
      return values[name] ?? null
    }
  }
}

describe("resolveReviewFakeApiRuntimeContext", () => {
  test("本番起動相当では fakeScenario を無視する", () => {
    const context = resolveReviewFakeApiRuntimeContext(
      createSearchParams({ fakeScenario: "success" }),
      { reviewModeEnabled: false }
    )

    expect(context.enabled).toBe(false)
    expect(context.scenarioId).toBe("empty")
    expect(context.overrideValue).toBeNull()
  })

  test("本番起動相当では fakeApi と fakeScenario の両方を無視する", () => {
    const context = resolveReviewFakeApiRuntimeContext(
      createSearchParams({ fakeApi: "1", fakeScenario: "success" }),
      { reviewModeEnabled: false }
    )

    expect(context.enabled).toBe(false)
    expect(context.scenarioId).toBe("empty")
    expect(context.triggerValue).toBe("1")
    expect(context.overrideValue).toBeNull()
  })

  test("fakeAPI 起動時に 6 種の状態パターン ID を解決できる", () => {
    const scenarioIds = [
      "empty",
      "loading",
      "success",
      "running",
      "error",
      "config-missing"
    ] as const

    for (const scenarioId of scenarioIds) {
      const context = resolveReviewFakeApiRuntimeContext(
        createSearchParams({ fakeApi: "1", fakeScenario: scenarioId }),
        { reviewModeEnabled: true }
      )

      expect(context.enabled).toBe(true)
      expect(context.scenarioId).toBe(scenarioId)
    }
  })

  test("未登録状態パターンは成功状態へ落ちず config-missing へ正規化する", () => {
    const context = resolveReviewFakeApiRuntimeContext(
      createSearchParams({ fakeApi: "1", fakeScenario: "unknown" }),
      { reviewModeEnabled: true }
    )

    expect(context.scenarioId).toBe("config-missing")
    expect(context.scenarioId).not.toBe("success")
  })

  test("fakeAPI 起動時にトリガー値が未登録なら config-missing へ正規化される", () => {
    const context = resolveReviewFakeApiRuntimeContext(
      createSearchParams({ fakeApi: "unknown" }),
      { reviewModeEnabled: true }
    )

    expect(context.enabled).toBe(true)
    expect(context.defaultScenarioId).toBe("config-missing")
    expect(context.scenarioId).toBe("config-missing")
  })
})
