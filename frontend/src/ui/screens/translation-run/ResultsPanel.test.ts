import { render, screen } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"
import ResultsPanel from "./ResultsPanel.svelte"

describe("ResultsPanel", () => {
  // R-3-2, R-3-4: 未訳 0 件でも、絞り込み前に結果があれば空表示と書き出し操作を保つ。
  it("未訳のみの選択と未訳なしを表示する", async () => {
    const user = userEvent.setup()
    const onUntranslatedOnlyChange = vi.fn()
    render(ResultsPanel, {
      phase: "done",
      results: [],
      total: 0,
      untranslatedOnly: true,
      hasUnfilteredResults: true,
      onUntranslatedOnlyChange
    })

    expect(screen.getByText("未訳はありません。")).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "xTranslator へ書き出し" })
    ).toBeInTheDocument()
    const checkbox = screen.getByRole("checkbox", { name: "未訳のみ" })
    expect(checkbox).toBeChecked()

    await user.click(checkbox)
    expect(onUntranslatedOnlyChange).toHaveBeenCalledWith(false)
  })
})
