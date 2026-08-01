import { fireEvent, render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { tick } from "svelte"
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"

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

const updatedPage = {
  ...oldPage,
  results: [{ ...oldPage.results[0], source: "updated source" }]
}

const properPending = {
  stage: "proper" as const,
  total: 1000,
  pending: 1000,
  succeeded: 0,
  failed: 0,
  canApply: false,
  untranslatedCount: 0
}

const properReady = {
  ...properPending,
  pending: 0,
  succeeded: 1000,
  canApply: true
}

const bodyPending = {
  ...properPending,
  stage: "body" as const,
  total: 20,
  pending: 20
}

const bodyReady = {
  ...bodyPending,
  pending: 0,
  succeeded: 20,
  canApply: true
}

const done = {
  ...bodyReady,
  stage: "done" as const,
  canApply: false
}

function deferred<T>() {
  let resolve!: (value: T) => void
  let reject!: (reason: unknown) => void
  const promise = new Promise<T>((resolvePromise, rejectPromise) => {
    resolve = resolvePromise
    reject = rejectPromise
  })
  return { promise, resolve, reject }
}

async function selectXai(user: ReturnType<typeof userEvent.setup>) {
  await user.click(screen.getByRole("button", { name: "xAI（batch）" }))
  await user.click(screen.getByRole("button", { name: "取得" }))
  await waitFor(() => {
    expect(screen.getByRole("button", { name: "バッチ実行" })).toBeEnabled()
  })
}

async function flushAsyncState() {
  await Promise.resolve()
  await Promise.resolve()
  await tick()
}

describe("TranslationRunContainer", () => {
  beforeEach(() => {
    vi.useRealTimers()
    vi.clearAllMocks()
    gateway.fetchStringsPresence.mockResolvedValue({
      english: true,
      japanese: true
    })
    gateway.fetchXaiModels.mockResolvedValue(["grok"])
    gateway.fetchModels.mockResolvedValue(["sync-model"])
    gateway.onRunProgress.mockReturnValue(() => {})
    gateway.listResultsPage.mockResolvedValue(oldPage)
    gateway.submitBatchTranslation.mockResolvedValue({
      reusedPreparation: false,
      completedWithoutExternalBatch: false
    })
    gateway.refreshBatchTranslations.mockResolvedValue(undefined)
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  // 既存の一覧取得に失敗しても、選択状態と現在の一覧を変更しない。
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
  })

  // R-2-3: 進行が無い時だけ新しい外部 batch を一度送り、送信直後の進行を表示する。
  it("進行なしでは一度だけ送信して直後に進行を取得する", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    gateway.getBatchProgress
      .mockResolvedValueOnce(undefined)
      .mockResolvedValueOnce(properPending)

    await user.click(screen.getByRole("button", { name: "バッチ実行" }))

    await waitFor(() =>
      expect(gateway.submitBatchTranslation).toHaveBeenCalledTimes(1)
    )
    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(2)
    expect(screen.getByRole("button", { name: "実行中…" })).toBeDisabled()
    expect(screen.getByText(/処理待ち\s+1000/)).toBeInTheDocument()
  })

  // R-1-1/R-1-2: 状態確認と取り込みを直列に繰り返し、同じ段の次チャンクから本文完了まで進む。
  it("10秒ごとに複数チャンクと固有名から本文を完了まで進める", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    gateway.getBatchProgress
      .mockResolvedValueOnce(undefined)
      .mockResolvedValueOnce(properPending)
      .mockResolvedValueOnce(properReady)
      .mockResolvedValueOnce(properPending)
      .mockResolvedValueOnce(properReady)
      .mockResolvedValueOnce(bodyPending)
      .mockResolvedValueOnce(bodyReady)
      .mockResolvedValueOnce(done)
    gateway.listResultsPage.mockResolvedValue(updatedPage)
    vi.useFakeTimers()

    await fireEvent.click(screen.getByRole("button", { name: "バッチ実行" }))
    await flushAsyncState()
    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(2)

    await vi.advanceTimersByTimeAsync(10_000)
    await flushAsyncState()
    expect(gateway.refreshBatchTranslations).toHaveBeenCalledTimes(1)
    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(4)

    await vi.advanceTimersByTimeAsync(10_000)
    await flushAsyncState()
    expect(gateway.refreshBatchTranslations).toHaveBeenCalledTimes(2)
    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(6)

    await vi.advanceTimersByTimeAsync(10_000)
    await flushAsyncState()
    expect(gateway.refreshBatchTranslations).toHaveBeenCalledTimes(3)
    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(8)
    expect(screen.getByRole("button", { name: "完了" })).toBeDisabled()
    expect(screen.getAllByText("updated source").length).toBeGreaterThan(0)
  })

  // R-1-1: 応答待ちの間は次の状態確認を重ねず、応答完了後から10秒を数える。
  it("状態確認の応答待ちには次の状態確認を開始しない", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    const polling = deferred<typeof properPending>()
    gateway.getBatchProgress
      .mockResolvedValueOnce(undefined)
      .mockResolvedValueOnce(properPending)
      .mockReturnValueOnce(polling.promise)
      .mockResolvedValueOnce(properPending)
    vi.useFakeTimers()

    await fireEvent.click(screen.getByRole("button", { name: "バッチ実行" }))
    await flushAsyncState()
    await vi.advanceTimersByTimeAsync(20_000)
    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(3)

    polling.resolve(properPending)
    await flushAsyncState()
    await vi.advanceTimersByTimeAsync(9_999)
    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(3)
    await vi.advanceTimersByTimeAsync(1)
    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(4)
  })

  // R-1-1: 一つの自動状態確認では開始時の接続先を固定し、入力欄の後続変更を混ぜない。
  it("自動状態確認中は開始時の接続情報を使い続ける", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    gateway.getBatchProgress
      .mockResolvedValueOnce(undefined)
      .mockResolvedValueOnce(properPending)
      .mockResolvedValueOnce(properPending)
    vi.useFakeTimers()

    await fireEvent.click(screen.getByRole("button", { name: "バッチ実行" }))
    await flushAsyncState()
    await fireEvent.input(screen.getByRole("textbox"), {
      target: { value: "https://changed.example" }
    })
    await vi.advanceTimersByTimeAsync(10_000)

    expect(gateway.getBatchProgress).toHaveBeenLastCalledWith(
      "A.esp",
      "xai",
      expect.objectContaining({ endpoint: "https://api.x.ai" })
    )
  })

  // R-2-2: 保存済みの途中進行は新規送信せず、その場から自動状態確認を再開する。
  it("保存済みの進行を新規送信せず再開する", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    gateway.getBatchProgress.mockResolvedValue(properPending)

    await user.click(screen.getByRole("button", { name: "バッチ実行" }))

    await waitFor(() =>
      expect(gateway.getBatchProgress).toHaveBeenCalledTimes(1)
    )
    expect(gateway.submitBatchTranslation).not.toHaveBeenCalled()
    expect(screen.getByRole("button", { name: "実行中…" })).toBeDisabled()
  })

  // R-3-4: 完了済みで未訳が無ければ送信せず終了し、未訳があれば未訳再送信へ進む。
  it("完了済みの未訳件数に応じて終了または未訳再送信を選ぶ", async () => {
    const user = userEvent.setup()
    const first = render(TranslationRunContainer, {
      pluginPath: "/data/A.esp"
    })
    await selectXai(user)
    gateway.getBatchProgress.mockResolvedValueOnce(done)

    await user.click(screen.getByRole("button", { name: "バッチ実行" }))

    await waitFor(() => {
      expect(screen.getByRole("button", { name: "完了" })).toBeDisabled()
    })
    expect(gateway.submitBatchTranslation).not.toHaveBeenCalled()
    first.unmount()

    vi.clearAllMocks()
    gateway.fetchStringsPresence.mockResolvedValue({
      english: true,
      japanese: true
    })
    gateway.fetchXaiModels.mockResolvedValue(["grok"])
    gateway.onRunProgress.mockReturnValue(() => {})
    gateway.listResultsPage.mockResolvedValue(oldPage)
    gateway.submitBatchTranslation.mockResolvedValue({
      reusedPreparation: true,
      completedWithoutExternalBatch: false
    })
    const doneWithUntranslated = { ...done, untranslatedCount: 2 }
    gateway.getBatchProgress
      .mockResolvedValueOnce(doneWithUntranslated)
      .mockResolvedValueOnce(properPending)
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)

    await user.click(screen.getByRole("button", { name: "バッチ実行" }))

    await waitFor(() => {
      expect(gateway.submitBatchTranslation).toHaveBeenCalledTimes(1)
    })
    expect(screen.getByRole("button", { name: "実行中…" })).toBeDisabled()
  })

  // R-2-4: provider 不一致の確認エラーでは送信せず、理由を表示して開始操作へ戻す。
  it("provider不一致では新規送信を行わない", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    gateway.getBatchProgress.mockRejectedValue(
      new Error("保存済み provider は openai です")
    )

    await user.click(screen.getByRole("button", { name: "バッチ実行" }))

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "保存済み provider は openai です"
    )
    expect(gateway.submitBatchTranslation).not.toHaveBeenCalled()
    expect(screen.getByRole("button", { name: "バッチ実行" })).toBeEnabled()
  })

  // R-3-2: 開始確認の応答待ちでは主操作を無効にし、重ねて押しても呼び出しを増やさない。
  it("開始処理中の多重押下を防ぐ", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    const initial = deferred<undefined>()
    gateway.getBatchProgress.mockReturnValue(initial.promise)
    const button = screen.getByRole("button", { name: "バッチ実行" })

    await fireEvent.click(button)
    await fireEvent.click(button)

    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(1)
    expect(screen.getByRole("button", { name: "実行中…" })).toBeDisabled()
    initial.resolve(undefined)
  })

  // R-3-3: 状態確認エラーで予約を止め、途中進行を残したまま再開操作を有効にする。
  it("状態確認エラーで停止して再開操作を有効にする", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    gateway.getBatchProgress
      .mockResolvedValueOnce(properPending)
      .mockRejectedValueOnce(
        new Error("batch batch_failed: token_limit_exceeded")
      )
    vi.useFakeTimers()

    await fireEvent.click(screen.getByRole("button", { name: "バッチ実行" }))
    await flushAsyncState()
    await vi.advanceTimersByTimeAsync(10_000)
    await flushAsyncState()

    expect(screen.getByRole("alert")).toHaveTextContent(
      "batch batch_failed: token_limit_exceeded"
    )
    expect(
      screen.getByRole("button", { name: "バッチ実行を再開" })
    ).toBeEnabled()
    await vi.advanceTimersByTimeAsync(30_000)
    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(2)
  })

  // R-2-1: 画面終了時に予約を破棄し、後から状態確認を始めない。
  it("画面終了後は予約済みの状態確認を行わない", async () => {
    const user = userEvent.setup()
    const view = render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    gateway.getBatchProgress.mockResolvedValue(properPending)
    vi.useFakeTimers()

    await fireEvent.click(screen.getByRole("button", { name: "バッチ実行" }))
    await flushAsyncState()
    view.unmount()
    await vi.advanceTimersByTimeAsync(20_000)

    expect(gateway.getBatchProgress).toHaveBeenCalledTimes(1)
  })

  // R-2-5: provider 変更前の遅延応答は変更後の画面へ反映しない。
  it("provider変更後は変更前の遅延応答を無視する", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)
    const initial = deferred<typeof properPending>()
    gateway.getBatchProgress.mockReturnValue(initial.promise)

    await fireEvent.click(screen.getByRole("button", { name: "バッチ実行" }))
    await user.click(screen.getByRole("button", { name: "OpenAI（batch）" }))
    initial.resolve(properPending)
    await flushAsyncState()

    expect(gateway.submitBatchTranslation).not.toHaveBeenCalled()
    expect(screen.getByText("未開始")).toBeInTheDocument()
  })

  // R-2-5: 変更前に始めた結果一覧取得も、provider 変更後の一覧へ反映しない。
  it("provider変更後は変更前の遅延した結果一覧を無視する", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    expect((await screen.findAllByText("old source")).length).toBeGreaterThan(0)
    await selectXai(user)
    const delayedPage = deferred<typeof updatedPage>()
    gateway.listResultsPage.mockReturnValueOnce(delayedPage.promise)
    gateway.getBatchProgress
      .mockResolvedValueOnce(properReady)
      .mockResolvedValueOnce(bodyPending)

    await fireEvent.click(screen.getByRole("button", { name: "バッチ実行" }))
    await waitFor(() => {
      expect(gateway.refreshBatchTranslations).toHaveBeenCalledTimes(1)
    })
    await user.click(screen.getByRole("button", { name: "OpenAI（batch）" }))
    delayedPage.resolve(updatedPage)
    await flushAsyncState()

    expect(screen.queryByText("updated source")).not.toBeInTheDocument()
    expect(screen.getAllByText("old source").length).toBeGreaterThan(0)
  })

  // R-2-5: plugin 変更前の遅延応答は変更後の進行状況へ反映しない。
  it("plugin変更後は変更前の遅延応答を無視する", async () => {
    const user = userEvent.setup()
    const view = render(TranslationRunContainer, {
      pluginPath: "/data/A.esp"
    })
    await selectXai(user)
    const initial = deferred<typeof properPending>()
    gateway.getBatchProgress.mockReturnValue(initial.promise)

    await fireEvent.click(screen.getByRole("button", { name: "バッチ実行" }))
    await view.rerender({ pluginPath: "/data/B.esp" })
    initial.resolve(properPending)
    await flushAsyncState()

    expect(gateway.submitBatchTranslation).not.toHaveBeenCalled()
    expect(screen.getByText("未開始")).toBeInTheDocument()
  })

  // R-1-3/R-3-5: 同期実行は従来の一回実行だけを行い、batch の状態確認を予約しない。
  it("同期実行ではbatchの状態確認を開始しない", async () => {
    const user = userEvent.setup()
    gateway.runExtractAndTranslate.mockResolvedValue({
      translatedCount: 1,
      untranslatedCount: 0
    })
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await user.click(screen.getByRole("button", { name: "取得" }))
    await waitFor(() => {
      expect(screen.getByRole("button", { name: "実行" })).toBeEnabled()
    })
    vi.useFakeTimers()

    await fireEvent.click(screen.getByRole("button", { name: "実行" }))
    await flushAsyncState()
    await vi.advanceTimersByTimeAsync(20_000)

    expect(gateway.runExtractAndTranslate).toHaveBeenCalledTimes(1)
    expect(gateway.getBatchProgress).not.toHaveBeenCalled()
  })

  // R-3-1: batch 操作は主操作一つだけを表示し、旧状態確認と手動取り込みは表示しない。
  it("batch画面に状態確認と手動取り込みを表示しない", async () => {
    const user = userEvent.setup()
    render(TranslationRunContainer, { pluginPath: "/data/A.esp" })
    await selectXai(user)

    expect(screen.getByRole("button", { name: "バッチ実行" })).toBeVisible()
    expect(
      screen.queryByRole("button", { name: "状態確認" })
    ).not.toBeInTheDocument()
    expect(screen.queryByText(/取り込んで/)).not.toBeInTheDocument()
  })
})
