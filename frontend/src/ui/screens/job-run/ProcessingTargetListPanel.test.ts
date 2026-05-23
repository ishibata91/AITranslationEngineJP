import { render, screen, within } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test } from "vitest"

import ProcessingTargetListPanel from "@ui/components/ProcessingTargetListPanel.svelte"
import { processingTargetListPanelFixtures } from "./__fixtures__/job-run-shell-fixtures"

describe("ProcessingTargetListPanel", () => {
  test("処理対象一覧の領域、見出し、table を表示し、初期表示では詳細を閉じる", () => {
    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.termTranslationFirstPage
    })

    const region = screen.getByLabelText("処理対象一覧")

    expect(region).toBeInTheDocument()
    expect(within(region).getByRole("heading", { name: "処理対象一覧" }))
      .toBeInTheDocument()
    expect(
      within(region).getByRole("table", { name: "処理対象ページ" })
    ).toBeInTheDocument()
    expect(
      within(region).queryByRole("heading", { name: "単語翻訳" })
    ).not.toBeInTheDocument()
    expect(within(region).queryByText("処理対象名")).not.toBeInTheDocument()
    expect(within(region).queryByText("処理対象詳細")).not.toBeInTheDocument()
    expect(
      within(region).queryByText(
        "AIサービスへ送り、確定訳語として翻訳ジョブ内辞書へ保存する候補。"
      )
    ).not.toBeInTheDocument()
  })

  test("既定ページサイズとして 50 件だけを描画する", () => {
    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.termTranslationFirstPage
    })

    expect(
      screen.getByRole("row", {
        name: "原語 001: Dragonborn / 訳語候補 001: ドラゴンボーン"
      })
    ).toBeInTheDocument()
    expect(
      screen.getByRole("row", {
        name: "原語 050: Dragonborn / 訳語候補 050: ドラゴンボーン"
      })
    ).toBeInTheDocument()
    expect(screen.queryByText("原語 051: Dragonborn")).not.toBeInTheDocument()
    expect(screen.getAllByText("1-50 / 125 件")).toHaveLength(2)
  })

  test("ラッパー内表示では内側見出しだけを非表示にし、件数と table を残す", () => {
    render(ProcessingTargetListPanel, {
      props: {
        ...processingTargetListPanelFixtures.personaGeneration,
        showHeadingTitle: false
      }
    })

    const region = screen.getByLabelText("処理対象一覧")

    expect(
      within(region).queryByRole("heading", { name: "処理対象一覧" })
    ).not.toBeInTheDocument()
    expect(within(region).getAllByText("1-50 / 64 件")).toHaveLength(2)
    expect(
      within(region).getByRole("table", { name: "処理対象ページ" })
    ).toBeInTheDocument()
  })

  test("単語翻訳の行タイトルに原語と訳語候補を表示する", () => {
    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.termTranslationFirstPage
    })

    expect(
      screen.getByRole("row", {
        name: "原語 001: Dragonborn / 訳語候補 001: ドラゴンボーン"
      })
    ).toHaveAttribute(
      "aria-label",
      "原語 001: Dragonborn / 訳語候補 001: ドラゴンボーン"
    )
  })

  test("本文翻訳の行タイトルに原文と訳文を表示する", () => {
    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.bodyTranslation
    })

    expect(
      screen.getByRole("row", {
        name: /原文 001: The ancient stone door remains sealed\..*訳文 001: 古代の石扉は封印されたままです。/
      })
    ).toHaveAttribute(
      "aria-label",
      "原文 001: The ancient stone door remains sealed. / 訳文 001: 古代の石扉は封印されたままです。"
    )
  })

  test("全文が見えている行タイトル部分は hover しても tooltip を表示しない", async () => {
    const user = userEvent.setup()

    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.termTranslationFirstPage
    })

    const sourceTitlePart = screen.getByRole("button", {
      name: "原語 006: Dragonborn"
    })

    Object.defineProperty(sourceTitlePart, "scrollWidth", {
      configurable: true,
      value: 120
    })
    Object.defineProperty(sourceTitlePart, "clientWidth", {
      configurable: true,
      value: 120
    })

    await user.hover(sourceTitlePart)

    expect(sourceTitlePart).not.toHaveAttribute("aria-describedby")
    expect(screen.queryByRole("tooltip")).not.toBeInTheDocument()
  })

  test("省略された行タイトル部分だけ tooltip と関連付く", async () => {
    const user = userEvent.setup()

    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.longText
    })

    const sourceTitlePart = screen.getByRole("button", {
      name: /When the sealed archive beneath the ancient stone keep opens/
    })
    const translatedTitlePart = screen.getByRole("button", {
      name: /古代の石造りの砦の地下にある封印書庫が開いたとき/
    })

    Object.defineProperty(sourceTitlePart, "scrollWidth", {
      configurable: true,
      value: 300
    })
    Object.defineProperty(sourceTitlePart, "clientWidth", {
      configurable: true,
      value: 120
    })
    Object.defineProperty(translatedTitlePart, "scrollWidth", {
      configurable: true,
      value: 320
    })
    Object.defineProperty(translatedTitlePart, "clientWidth", {
      configurable: true,
      value: 140
    })

    await user.hover(sourceTitlePart)

    expect(sourceTitlePart).toHaveAttribute(
      "aria-describedby",
      "processing-target-title-long-processing-target-0"
    )

    await user.hover(translatedTitlePart)

    expect(translatedTitlePart).toHaveAttribute(
      "aria-describedby",
      "processing-target-title-long-processing-target-1"
    )
    expect(
      screen.getByRole("tooltip")
    ).toBeInTheDocument()
  })

  test("行クリックと keyboard 操作で table 内の詳細行を 1 件だけ開く", async () => {
    const user = userEvent.setup()

    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.termTranslationFirstPage
    })

    await user.click(
      screen.getByRole("row", {
        name: "原語 002: Dragonborn / 訳語候補 002: ドラゴンボーン"
      })
    )

    const secondMetadataRow = screen.getByRole("row", {
      name: /原語 002: Dragonborn \/ 訳語候補 002: ドラゴンボーン のメタデータ/
    })
    expect(within(secondMetadataRow).getByText("term-translation-002"))
      .toBeInTheDocument()
    expect(within(secondMetadataRow).getByText("term.source.002"))
      .toBeInTheDocument()
    expect(
      within(secondMetadataRow).getByText("AIサービスへ送る訳語候補")
    ).toBeInTheDocument()

    screen
      .getByRole("row", {
        name: "原語 003: Dragonborn / 訳語候補 003: ドラゴンボーン"
      })
      .focus()
    await user.keyboard("{Enter}")

    expect(screen.queryByText("term-translation-002"))
      .not.toBeInTheDocument()
    expect(screen.getByText("term-translation-003"))
      .toBeInTheDocument()
  })

  test("初期展開 ID、詳細本文、詳細操作を表示する", async () => {
    const user = userEvent.setup()
    let actionCount = 0

    render(ProcessingTargetListPanel, {
      props: {
        items: [
          {
            id: "persona-001",
            name: "宿屋の主人",
            detail: "相手を急かさず、短い説明を好む。",
            metadata: [
              { label: "FormID", value: "FE001001" },
              { label: "EditorID", value: "SAMPLE_NPC_A" }
            ],
            actions: [
              {
                label: "編集",
                onAction: () => {
                  actionCount += 1
                }
              }
            ]
          }
        ],
        initialExpandedItemId: "persona-001"
      }
    })

    const detailRow = screen.getByRole("row", {
      name: /宿屋の主人 のメタデータ/
    })
    expect(within(detailRow).getByText("FE001001")).toBeInTheDocument()
    expect(
      within(detailRow).queryByText("相手を急かさず、短い説明を好む。")
    ).not.toBeInTheDocument()

    await user.click(within(detailRow).getByRole("button", { name: "編集" }))

    expect(actionCount).toBe(1)
  })

  test("翻訳完了確認の行タイトルに原文と訳文を表示し、詳細行から出力対象 metadata を除く", async () => {
    const user = userEvent.setup()

    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.translationComplete
    })

    const translatedRow = screen.getByRole("row", {
      name: /原文 001: The ancient stone door remains sealed\..*訳文 001: 古代の石扉は封印されたままです。/
    })

    expect(translatedRow).toHaveAttribute(
      "aria-label",
      "原文 001: The ancient stone door remains sealed. / 訳文 001: 古代の石扉は封印されたままです。"
    )

    await user.click(translatedRow)

    const metadataRow = screen.getByRole("row", {
      name: /原文 001: The ancient stone door remains sealed\. \/ 訳文 001: 古代の石扉は封印されたままです。 のメタデータ/
    })
    expect(within(metadataRow).getByText("訳文 ID")).toBeInTheDocument()
    expect(within(metadataRow).getByText("translated-text-001"))
      .toBeInTheDocument()
    expect(within(metadataRow).getByText("原文 field")).toBeInTheDocument()
    expect(within(metadataRow).getByText("source.field.001"))
      .toBeInTheDocument()
    expect(within(metadataRow).queryByText("出力対象")).not.toBeInTheDocument()
    expect(within(metadataRow).queryByText("出力管理へ渡す訳文"))
      .not.toBeInTheDocument()
    expect(within(metadataRow).getByText("確認状態")).toBeInTheDocument()
    expect(within(metadataRow).getByText("確認済み")).toBeInTheDocument()
  })

  test("ペルソナ生成の行タイトルは NPC 名だけを表示する", () => {
    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.personaGeneration
    })

    expect(screen.getByRole("row", { name: "NPC 001" })).toBeInTheDocument()
    expect(screen.queryByRole("row", { name: /NPC ペルソナ入力 001/ }))
      .not.toBeInTheDocument()
  })

  test("次へで現在ページを進め、前へで戻す", async () => {
    const user = userEvent.setup()

    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.termTranslationFirstPage
    })

    await user.click(screen.getByRole("button", { name: "次へ" }))

    expect(screen.queryByText("原語 001: Dragonborn")).not.toBeInTheDocument()
    expect(
      screen.getByRole("row", {
        name: "原語 051: Dragonborn / 訳語候補 051: ドラゴンボーン"
      })
    ).toBeInTheDocument()
    expect(
      screen.getByRole("row", {
        name: "原語 100: Dragonborn / 訳語候補 100: ドラゴンボーン"
      })
    )
      .toBeInTheDocument()
    expect(screen.queryByText("原語 101: Dragonborn")).not.toBeInTheDocument()
    expect(screen.getAllByText("51-100 / 125 件")).toHaveLength(2)

    await user.click(screen.getByRole("button", { name: "前へ" }))

    expect(
      screen.getByRole("row", {
        name: "原語 001: Dragonborn / 訳語候補 001: ドラゴンボーン"
      })
    ).toBeInTheDocument()
    expect(screen.queryByText("原語 051: Dragonborn")).not.toBeInTheDocument()
  })

  test("1 ページ目と最終ページで該当するページ操作を無効にする", () => {
    const { unmount } = render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.termTranslationFirstPage
    })

    expect(screen.getByRole("button", { name: "前へ" })).toBeDisabled()
    expect(screen.getByRole("button", { name: "次へ" })).toBeEnabled()

    unmount()

    render(ProcessingTargetListPanel, {
      props: processingTargetListPanelFixtures.finalPage
    })

    expect(screen.getByRole("row", { name: "最終ページ確認候補 051" }))
      .toBeInTheDocument()
    expect(screen.getByRole("button", { name: "前へ" })).toBeEnabled()
    expect(screen.getByRole("button", { name: "次へ" })).toBeDisabled()
  })
})
