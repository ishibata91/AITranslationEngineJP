import { render, screen } from "@testing-library/svelte"
import { expect } from "vitest"
import { runScreenSpecHarness } from "../../../test/screen-spec-harness"
import TargetPluginsScreen from "./TargetPluginsScreen.svelte"
import { targetPluginScreenStates } from "./target-plugins-screen-specs"

const callbacks = {
  onSelectNewPlugin: () => {},
  onNewPluginPathInput: () => {},
  onProceedToRun: () => {},
  onOpenPlugin: () => {},
  onRequestDelete: () => {},
  onConfirmDelete: () => {},
  onCancelDelete: () => {}
}

function renderState(
  state: (typeof targetPluginScreenStates)[keyof typeof targetPluginScreenStates]
) {
  return render(TargetPluginsScreen, { ...callbacks, ...state.args })
}

export const targetPluginScreenChecks = [
  {
    id: "target-plugins.empty.proceed-disabled",
    verify: () => {
      renderState(targetPluginScreenStates.empty)
      expect(screen.getByRole("button", { name: "翻訳へ進む" })).toBeDisabled()
    }
  },
  {
    id: "target-plugins.empty.list-empty",
    verify: () => {
      renderState(targetPluginScreenStates.empty)
      expect(
        screen.getByText(/まだ翻訳したプラグインはありません。/)
      ).toBeInTheDocument()
    }
  },
  {
    id: "target-plugins.loading.list-loading",
    verify: () => {
      const view = renderState(targetPluginScreenStates.loading)
      expect(
        screen.getByText("プラグインを読み込んでいます。")
      ).toBeInTheDocument()
      expect(view.container.querySelector(".loading-ring")).toBeInTheDocument()
    }
  },
  {
    id: "target-plugins.list.count",
    verify: () => {
      renderState(targetPluginScreenStates.list)
      expect(screen.getByText("4 件")).toBeInTheDocument()
    }
  },
  {
    id: "target-plugins.list.progress-badges",
    verify: () => {
      const view = renderState(targetPluginScreenStates.list)
      const badges = [...view.container.querySelectorAll("li .badge")].map(
        (badge) => badge.textContent?.trim() ?? ""
      )
      expect(badges.filter((label) => label.startsWith("完了 "))).toHaveLength(
        2
      )
      expect(badges.some((label) => label.startsWith("未着手 "))).toBe(true)
      expect(badges.some((label) => label.startsWith("翻訳中 "))).toBe(true)
    }
  },
  {
    id: "target-plugins.list.row-actions",
    verify: () => {
      renderState(targetPluginScreenStates.list)
      const openButtons = screen.getAllByRole("button", { name: /結果を開く/ })
      const deleteButtons = screen.getAllByRole("button", { name: "削除" })
      expect(openButtons).toHaveLength(4)
      expect(deleteButtons).toHaveLength(4)
      for (const button of [...openButtons, ...deleteButtons]) {
        expect(button).toBeEnabled()
      }
    }
  },
  {
    id: "target-plugins.selected.proceed-enabled",
    verify: () => {
      renderState(targetPluginScreenStates.selected)
      expect(screen.getByRole("button", { name: "翻訳へ進む" })).toBeEnabled()
    }
  },
  {
    id: "target-plugins.confirm-delete.prompt",
    verify: () => {
      renderState(targetPluginScreenStates.confirmDelete)
      expect(
        screen.getAllByText("成果を消して再実行でやり直します。削除しますか？")
      ).toHaveLength(1)
    }
  },
  {
    id: "target-plugins.confirm-delete.actions-enabled",
    verify: () => {
      renderState(targetPluginScreenStates.confirmDelete)
      expect(screen.getByRole("button", { name: "削除する" })).toBeEnabled()
      expect(screen.getByRole("button", { name: "取消" })).toBeEnabled()
    }
  },
  {
    id: "target-plugins.deleting.progress",
    verify: () => {
      const view = renderState(targetPluginScreenStates.deleting)
      expect(
        screen.getByRole("button", { name: "削除中…" })
      ).toBeInTheDocument()
      expect(
        view.container.querySelector(".loading-spinner")
      ).toBeInTheDocument()
    }
  },
  {
    id: "target-plugins.deleting.actions-disabled",
    verify: () => {
      renderState(targetPluginScreenStates.deleting)
      expect(screen.getByRole("button", { name: "削除中…" })).toBeDisabled()
      expect(screen.getByRole("button", { name: "取消" })).toBeDisabled()
    }
  },
  {
    id: "target-plugins.error.message",
    verify: () => {
      renderState(targetPluginScreenStates.errored)
      expect(screen.getByRole("alert")).toHaveTextContent(
        "プラグインの翻訳成果の削除に失敗しました。もう一度お試しください。"
      )
    }
  }
]

const targetPluginSpecIds = Object.values(targetPluginScreenStates).flatMap(
  (state) => state.specs.map((spec) => spec.id)
)

runScreenSpecHarness(
  "翻訳対象プラグイン",
  targetPluginSpecIds,
  targetPluginScreenChecks
)
