- workflow: impl
- status: completed
- lane_owner: codex
- scope: `docs/architecture.md` の frontend concrete layout に沿って、将来の feature / slice directory contract を受ける実ディレクトリを追加する

## Request Summary

- `docs/architecture.md` に従い、frontend 側の実ディレクトリを先に追加する

## Decision Basis

- `docs/architecture.md` では frontend root として `src/ui/`、`src/application/`、`src/gateway/`、`src/shared/` を正本としている
- `src/` 配下には最小 bootstrap 実装はあるが、architecture に記載された UI 内部責務を受ける directory contract はまだ薄い
- 今回は実装コードの大移動ではなく、以後の実装と lint 拡張の土台になる directory skeleton を追加する

## UI

- ランタイムの画面表示や操作フローは変えない
- directory contract としては、`src/ui/` 配下の責務名を `app-shell/`、`screens/`、`views/`、`stores/` で明示する
- 既存の `src/ui/app-shell/` と `src/ui/screens/` は維持し、今回の追加は `Presenter / View` と `Screen Store` の受け皿を先に固定することに留める

## Scenario

- N/A

## Logic

- additive-only で固定する frontend skeleton は以下とする
  - `src/ui/views/`
  - `src/ui/stores/`
  - `src/application/usecases/`
  - `src/application/ports/input/`
  - `src/application/ports/gateway/`
  - `src/gateway/tauri/invoke/`
  - `src/gateway/tauri/events/`
- 既存の `src/ui/app-shell/`、`src/ui/screens/`、`src/application/bootstrap/`、`src/gateway/tauri/`、`src/shared/contracts/` は移動しない
- `src/shared/` は `contracts/` がすでに architecture の DTO / contract root を満たしているため、この task では追加 root を増やさない
- frontend bootstrap phase では TS 側に `src/domain/` / `src/infra/` を作らない
- `UI Service` はこの task では standalone root を増やさず、後続の feature / slice contract で配置を確定する
- 空ディレクトリの保持は `.gitkeep` など最小の方法で行う

## Implementation Plan

- `src/ui/views/` と `src/ui/stores/` を additive に追加する
- `src/application/usecases/`、`src/application/ports/input/`、`src/application/ports/gateway/` を additive に追加する
- `src/gateway/tauri/invoke/` と `src/gateway/tauri/events/` を additive に追加する
- 既存 bootstrap 実装は移動せず、空ディレクトリ保持は `.gitkeep` で行う
- subdirectory 名を canonical にするため、必要最小限の directory map を `docs/architecture.md` に同期する

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`

## Required Evidence

- 追加したディレクトリが `docs/architecture.md` の責務名と対応していること
- 既存 frontend / harness が壊れていないこと

## Docs Sync

- `docs/architecture.md` に frontend concrete layout の下位 directory map を追記した

## Outcome

- `src/ui/views/`
- `src/ui/stores/`
- `src/application/usecases/`
- `src/application/ports/input/`
- `src/application/ports/gateway/`
- `src/gateway/tauri/invoke/`
- `src/gateway/tauri/events/`
- 以上を `.gitkeep` 付きで追加し、既存 bootstrap code は移動していない
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
