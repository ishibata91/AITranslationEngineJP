# frontend-backend-connection-refactor 実装引き継ぎ入力

- 作成日: 2026-05-25
- 作成者: `refactor_lane`
- 根拠: `implementation-scope.md`
- 対象 branch: `codex/frontend-backend-connection-refactor`
- 統合先 branch: `master`

## 共通禁止事項

- `frontend/wailsjs/` を手編集しない。
- docs 正本本文を変更しない。
- `.codex` を変更しない。
- 承認済み `implementation-scope` 外のプロダクトコードとプロダクトテストを変更しない。
- remote repository を変更しない。
- 他 agent の変更を revert しない。

## 実行 wave

| wave | handoff | 担当 agent | 依存対象 |
| --- | --- | --- | --- |
| `wave-1` | `FBC-FE-001` | `frontend_implementer` | なし |
| `wave-2` | `FBC-INT-001` | `integration_implementer` | `FBC-FE-001` |
| `wave-3` | `FBC-UT-FE-001` | `implementation_unit_tester` | `FBC-INT-001` |
| `wave-3` | `FBC-UT-BE-001` | `implementation_unit_tester` | `FBC-INT-001` |
| `wave-4` | `FBC-SC-001` | `implementation_scenario_tester` | `FBC-INT-001`, `FBC-UT-FE-001`, `FBC-UT-BE-001` |

## wave-1 handoff

### `FBC-FE-001`

- 担当 agent: `frontend_implementer`
- 使用 skill: `implement-frontend`
- 承認済み範囲: `SQ-FBC-003`
- 目的: screen controller factory から gateway DTO 型依存を外し、依存方向を application contract 側へ寄せる。
- 対象ファイル候補:
  - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts`
  - `frontend/src/application/gateway-contract/term-translation-phase/term-translation-phase-gateway-contract.ts`
  - `frontend/src/application/contract/term-translation-phase/`
- 変更禁止範囲:
  - `frontend/wailsjs/`
  - backend product code
  - docs 正本本文
  - `.codex/`
- 初手:
  - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts` の `createTermTranslationPhaseScreenControllerFactory` から `@controller/wails/gateway-dto/term-translation-phase` import を外す。
- 完了条件:
  - screen controller factory は Wails 呼び出し詳細と gateway DTO 型を参照しない。
  - coverage 用の型が必要な場合は application contract 側、または gateway 境界側の test helper に閉じる。
  - 単語翻訳画面の `Gateway` 状態、段階状態、操作可否の表示契約を維持する。
- 検証コマンド:
  - `python3 scripts/harness/run.py --suite frontend-local`

## 後続 handoff

後続 handoff は `FBC-FE-001` の完了結果を確認した後に起動する。
`FBC-INT-001` は shared contract を変更するため、他 handoff と同時開始しない。
`FBC-UT-FE-001` と `FBC-UT-BE-001` は `FBC-INT-001` 完了後に並列起動できる。
