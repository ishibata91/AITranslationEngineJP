import { render, screen } from "@testing-library/svelte"
import { expect } from "vitest"
import { runScreenSpecHarness } from "../../../test/screen-spec-harness"
import TranslationRunScreen from "./TranslationRunScreen.svelte"
import { translationRunScreenStates } from "./translation-run-screen-specs"

const callbacks = {
  onFieldInput: () => {},
  onLoadModels: () => {},
  onRun: () => {},
  onPagePrev: () => {},
  onPageNext: () => {},
  onUntranslatedOnlyChange: () => {},
  onProviderChange: () => {},
  onSubmit: () => {}
}

function renderState(
  state: (typeof translationRunScreenStates)[keyof typeof translationRunScreenStates]
) {
  return render(TranslationRunScreen, { ...callbacks, ...state.args })
}

function expectPhaseBadge(container: HTMLElement, label: string): void {
  const badge = [...container.querySelectorAll(".badge")].find(
    (element) => element.textContent?.trim() === label
  )
  expect(badge).toBeInTheDocument()
}

function expectManualActionsHidden(): void {
  expect(
    screen.queryByRole("button", { name: /状態確認|取り込/ })
  ).not.toBeInTheDocument()
}

export const translationRunScreenChecks = [
  {
    id: "translation-run.not-started.phase-badge",
    verify: () => {
      const view = renderState(translationRunScreenStates.notStarted)
      expectPhaseBadge(view.container, "未実行")
    }
  },
  {
    id: "translation-run.not-started.main-action",
    verify: () => {
      renderState(translationRunScreenStates.notStarted)
      expect(screen.getByRole("button", { name: "バッチ実行" })).toBeEnabled()
    }
  },
  {
    id: "translation-run.not-started.manual-actions-hidden",
    verify: () => {
      renderState(translationRunScreenStates.notStarted)
      expectManualActionsHidden()
    }
  },
  {
    id: "translation-run.running.phase-badge",
    verify: () => {
      const view = renderState(translationRunScreenStates.running)
      expectPhaseBadge(view.container, "実行中")
    }
  },
  {
    id: "translation-run.running.main-action",
    verify: () => {
      const view = renderState(translationRunScreenStates.running)
      expect(screen.getByRole("button", { name: "実行中…" })).toBeDisabled()
      expect(
        view.container.querySelector(".card-actions .loading-spinner")
      ).toBeInTheDocument()
    }
  },
  {
    id: "translation-run.running.manual-actions-hidden",
    verify: () => {
      renderState(translationRunScreenStates.running)
      expectManualActionsHidden()
    }
  },
  {
    id: "translation-run.paused.phase-badge",
    verify: () => {
      const view = renderState(translationRunScreenStates.paused)
      expectPhaseBadge(view.container, "未実行")
    }
  },
  {
    id: "translation-run.paused.main-action",
    verify: () => {
      renderState(translationRunScreenStates.paused)
      expect(
        screen.getByRole("button", { name: "バッチ実行を再開" })
      ).toBeEnabled()
    }
  },
  {
    id: "translation-run.paused.manual-actions-hidden",
    verify: () => {
      renderState(translationRunScreenStates.paused)
      expectManualActionsHidden()
    }
  },
  {
    id: "translation-run.done.phase-badge",
    verify: () => {
      const view = renderState(translationRunScreenStates.done)
      expectPhaseBadge(view.container, "完了")
    }
  },
  {
    id: "translation-run.done.main-action",
    verify: () => {
      renderState(translationRunScreenStates.done)
      expect(screen.getByRole("button", { name: "完了" })).toBeDisabled()
    }
  },
  {
    id: "translation-run.done.manual-actions-hidden",
    verify: () => {
      renderState(translationRunScreenStates.done)
      expectManualActionsHidden()
    }
  },
  {
    id: "translation-run.done.notice",
    verify: () => {
      renderState(translationRunScreenStates.done)
      expect(screen.getByRole("status")).toHaveTextContent(
        "本文を取り込みました。翻訳が完了しました。"
      )
    }
  },
  {
    id: "translation-run.done-untranslated.phase-badge",
    verify: () => {
      const view = renderState(translationRunScreenStates.doneWithUntranslated)
      expectPhaseBadge(view.container, "完了")
    }
  },
  {
    id: "translation-run.done-untranslated.main-action",
    verify: () => {
      renderState(translationRunScreenStates.doneWithUntranslated)
      expect(
        screen.getByRole("button", { name: "未訳だけを再送信" })
      ).toBeEnabled()
    }
  },
  {
    id: "translation-run.done-untranslated.manual-actions-hidden",
    verify: () => {
      renderState(translationRunScreenStates.doneWithUntranslated)
      expectManualActionsHidden()
    }
  },
  {
    id: "translation-run.done-untranslated.notice",
    verify: () => {
      renderState(translationRunScreenStates.doneWithUntranslated)
      expect(screen.getByRole("status")).toHaveTextContent(
        "3 件が未訳のまま残りました。"
      )
    }
  },
  {
    id: "translation-run.failed.phase-badge",
    verify: () => {
      const view = renderState(translationRunScreenStates.failed)
      expectPhaseBadge(view.container, "失敗")
    }
  },
  {
    id: "translation-run.failed.main-action",
    verify: () => {
      renderState(translationRunScreenStates.failed)
      expect(
        screen.getByRole("button", { name: "バッチ実行を再開" })
      ).toBeEnabled()
    }
  },
  {
    id: "translation-run.failed.manual-actions-hidden",
    verify: () => {
      renderState(translationRunScreenStates.failed)
      expectManualActionsHidden()
    }
  },
  {
    id: "translation-run.failed.error",
    verify: () => {
      renderState(translationRunScreenStates.failed)
      expect(screen.getByRole("alert")).toHaveTextContent("batch_01JTEST")
      expect(screen.getByRole("alert")).toHaveTextContent(
        "token_limit_exceeded"
      )
    }
  }
]

const translationRunSpecIds = Object.values(translationRunScreenStates).flatMap(
  (state) => state.specs.map((spec) => spec.id)
)

runScreenSpecHarness(
  "翻訳実行",
  translationRunSpecIds,
  translationRunScreenChecks
)
