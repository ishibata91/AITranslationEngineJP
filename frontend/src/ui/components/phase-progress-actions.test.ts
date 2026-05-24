import { describe, expect, test } from "vitest"

import { selectPhaseProgressActions } from "./phase-progress-actions"
import type { PhaseActionItem } from "./phase-panel-types"

type TestActionId = "start" | "pause" | "resume" | "retry" | "next"

function action(
  id: TestActionId,
  disabled: boolean
): PhaseActionItem<TestActionId> {
  return {
    id,
    label: id,
    disabled,
    blockedReason: disabled ? "blocked" : "",
    tone: "default"
  }
}

describe("selectPhaseProgressActions", () => {
  test("開始系 action は有効な 1 件だけにし、中断は disabled でも残す", () => {
    const selected = selectPhaseProgressActions(
      [
        action("start", false),
        action("pause", true),
        action("resume", false),
        action("retry", false),
        action("next", false)
      ],
      {
        hiddenActionIds: ["next"],
        runActionIds: ["start", "resume", "retry"],
        alwaysActionIds: ["pause"]
      }
    )

    expect(selected.map((item) => item.id)).toEqual(["start", "pause"])
  })

  test("開始系 action が全部 disabled なら先頭の開始系を disabled のまま残す", () => {
    const selected = selectPhaseProgressActions(
      [
        action("start", true),
        action("pause", true),
        action("resume", true),
        action("retry", true)
      ],
      {
        runActionIds: ["start", "resume", "retry"],
        alwaysActionIds: ["pause"]
      }
    )

    expect(selected.map((item) => [item.id, item.disabled])).toEqual([
      ["start", true],
      ["pause", true]
    ])
  })
})
