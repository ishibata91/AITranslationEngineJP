import { render, screen } from "@testing-library/svelte"

import AppShell from "./AppShell.svelte"

describe("AppShell", () => {
  test("renders the app shell title and description", () => {
    render(AppShell, {
      title: "Architecture Skeleton"
    })

    expect(
      screen.getByRole("heading", { name: "Architecture Skeleton" })
    ).toBeInTheDocument()
    expect(
      screen.getByText(
        "`docs/architecture.md` に沿った Wails + Go + Svelte の初期骨格です。"
      )
    ).toBeInTheDocument()
  })
})
