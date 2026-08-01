import { render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { beforeEach, describe, expect, it, vi } from "vitest"

const gateway = vi.hoisted(() => ({
  fetchModels: vi.fn(),
  fetchStringsPresence: vi.fn(),
  fetchXaiModels: vi.fn(),
  runExtractAndTranslate: vi.fn(),
  submitBatchTranslation: vi.fn(),
  refreshBatchTranslations: vi.fn(),
  getBatchProgress: vi.fn(),
  listResultsPage: vi.fn(),
  onRunProgress: vi.fn(),
  exportXTranslatorXml: vi.fn()
}))

vi.mock("../../../gateway/translation-gateway", () => gateway)

import TranslationRunContainer from "./TranslationRunContainer.svelte"

const oldPage = {
  total: 1,
  unfilteredTotal: 1,
  results: [
    {
      edid: "OldRow",
      source: "old source",
      dest: "既存の訳",
      statusLabel: "訳済"
    }
  ],
  nextCursor: "",
  hasMore: false
}

describe("TranslationRunContainer", () => {
  beforeEach(() => {
    vi.clearAllMocks()
    gateway.fetchStringsPresence.mockResolvedValue({
      english: true,
      japanese: true
    })
    gateway.fetchXaiModels.mockResolvedValue(["grok"])
    gateway.onRunProgress.mockReturnValue(() => {})
    gateway.listResultsPage.mockResolvedValue(oldPage)
  })

  // R-3-5: 条件変更用の取得に失敗しても、選択状態と現在の一覧を変更しない。
  it("未訳のみの取得失敗時に選択前の一覧とページを保つ", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })

    expect((await screen.findAllByText("old source")).length).toBeGreaterThan(0)
    const checkbox = screen.getByRole("checkbox", { name: "未訳のみ" })
    gateway.listResultsPage.mockRejectedValueOnce(
      new Error("未訳ページを取得できません")
    )

    await user.click(checkbox)

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "未訳ページを取得できません"
    )
    expect(checkbox).not.toBeChecked()
    expect(screen.getAllByText("old source").length).toBeGreaterThan(0)
    expect(screen.getByText("ページ 1")).toBeInTheDocument()
  })

  // R-2-4: 外部 batch なしで送信が完了しても、一覧と状態の両方を取得できるまで以前の表示を保つ。
  it("送信成功後の画面更新失敗で以前の表示を保つ", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })

    expect((await screen.findAllByText("old source")).length).toBeGreaterThan(0)
    await user.click(screen.getByRole("button", { name: "xAI（batch）" }))
    await user.click(screen.getByRole("button", { name: "取得" }))
    await waitFor(() => {
      expect(screen.getByRole("button", { name: "送信して開始" })).toBeEnabled()
    })

    gateway.submitBatchTranslation.mockResolvedValue({
      reusedPreparation: true,
      completedWithoutExternalBatch: true
    })
    gateway.listResultsPage.mockRejectedValueOnce(
      new Error("更新後の一覧を取得できません")
    )
    gateway.getBatchProgress.mockResolvedValue({
      stage: "done",
      total: 0,
      pending: 0,
      succeeded: 0,
      failed: 0,
      canApply: false,
      untranslatedCount: 0
    })

    await user.click(screen.getByRole("button", { name: "送信して開始" }))

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "batch の処理は完了しましたが、画面の更新に失敗しました"
    )
    expect(screen.getAllByText("old source").length).toBeGreaterThan(0)
    expect(
      screen.queryByText("保存済みの準備を使って未訳だけを処理しました。")
    ).not.toBeInTheDocument()
  })

  // R-2-3: 準備未完了の初回送信では、保存済みの準備を使った案内を表示しない。
  it("初回のbatch送信を未訳だけの再送信として表示しない", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })

    expect((await screen.findAllByText("old source")).length).toBeGreaterThan(0)
    await user.click(screen.getByRole("button", { name: "xAI（batch）" }))
    await user.click(screen.getByRole("button", { name: "取得" }))
    await waitFor(() => {
      expect(screen.getByRole("button", { name: "送信して開始" })).toBeEnabled()
    })
    gateway.submitBatchTranslation.mockResolvedValue({
      reusedPreparation: false,
      completedWithoutExternalBatch: false
    })

    await user.click(screen.getByRole("button", { name: "送信して開始" }))

    expect(await screen.findByRole("status")).toHaveTextContent(
      "batch を送信しました"
    )
    expect(
      screen.queryByText("保存済みの準備を使って未訳だけを処理しました。")
    ).not.toBeInTheDocument()
  })

  // R-1-4: 本文の取り込みが成功した後に未訳件数を取得できなくても、完了済みを取り消さず以前の表示を保つ。
  it("取り込み成功後の未訳件数取得失敗で以前の表示を保つ", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })

    expect((await screen.findAllByText("old source")).length).toBeGreaterThan(0)
    await user.click(screen.getByRole("button", { name: "xAI（batch）" }))
    gateway.getBatchProgress.mockResolvedValueOnce({
      stage: "body",
      total: 1,
      pending: 0,
      succeeded: 1,
      failed: 0,
      canApply: true,
      untranslatedCount: 0
    })
    await user.click(screen.getByRole("button", { name: "状態確認" }))
    await waitFor(() => {
      expect(
        screen.getByRole("button", { name: "取り込んで完了" })
      ).toBeEnabled()
    })

    gateway.refreshBatchTranslations.mockResolvedValue(undefined)
    gateway.getBatchProgress.mockRejectedValueOnce(
      new Error("未訳件数を取得できません")
    )
    await user.click(screen.getByRole("button", { name: "取り込んで完了" }))

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "batch の取り込みは完了しましたが、画面の更新に失敗しました"
    )
    expect(screen.getAllByText("old source").length).toBeGreaterThan(0)
  })
})
