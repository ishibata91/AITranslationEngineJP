export type ShellRouteId =
  | "dashboard"
  | "master-dictionary"
  | "master-persona"
  | "translation-management"
  | "output-management"

export interface ShellRouteContract {
  id: ShellRouteId
  label: string
  placeholderState?: string
}

const SHELL_ROUTE_CONTRACT: ReadonlyArray<ShellRouteContract> = [
  { id: "dashboard", label: "ダッシュボード" },
  {
    id: "master-dictionary",
    label: "マスター辞書",
    placeholderState: "準備中"
  },
  {
    id: "master-persona",
    label: "マスターペルソナ",
    placeholderState: "準備中"
  },
  {
    id: "translation-management",
    label: "翻訳管理",
    placeholderState: "準備中"
  },
  {
    id: "output-management",
    label: "出力管理",
    placeholderState: "準備中"
  }
]

interface ShellState {
  defaultRouteId: ShellRouteId
  routes: ShellRouteContract[]
}

export function createShellState(): ShellState {
  return {
    defaultRouteId: "dashboard",
    routes: SHELL_ROUTE_CONTRACT.map((route) => ({ ...route }))
  }
}
