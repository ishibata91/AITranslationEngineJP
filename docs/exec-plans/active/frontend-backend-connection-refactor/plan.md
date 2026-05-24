# フロントエンドとバックエンド接続境界リファクタ

## task 枠

- task-id: `frontend-backend-connection-refactor`
- 作業計画フォルダ: `docs/exec-plans/active/frontend-backend-connection-refactor/`
- 作業場所: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- 作業ブランチ: `codex/frontend-backend-connection-refactor`
- 統合先ブランチ: `master`
- 作成日: 2026-05-25

## リファクタ目的

フロントエンドとバックエンドの接続部分を中心に、既存仕様と既存実装の差分を整理する。
その後、人間判断で正とする側を固定し、責務境界とテスト品質の調査へ進める。

この task は新規機能を作らない。
この task は承認済み `implementation-scope` が作られるまでプロダクトコードとプロダクトテストを変更しない。
人間に求める判断は、要件、詳細仕様、画面仕様に関係する差分だけに限定する。
architecture、coding guideline、test guideline 起点の差分は、構造品質調査またはテスト品質調査の材料として扱い、人間へ仕様実装優先判断を求めない。

## 対象仕様参照

状態: 固定。

固定理由:

人間は 2026-05-25 に「広くてもやって」と指示した。
この指示により、前回の候補を広い調査範囲として固定する。

固定範囲:

- `docs/spec.md`: 接続境界で守るべき恒久要件
- `docs/detail-specs/`: 接続境界に関係する詳細仕様正本
- `docs/screen-design/`: 接続境界に関係する画面仕様正本

判断対象外の参照:

- `docs/architecture.md`: 構造品質調査の材料として扱う。
- `docs/coding-guidelines-frontend.md`: 構造品質調査の材料として扱う。
- `docs/coding-guidelines-backend.md`: 構造品質調査の材料として扱う。
- `docs/coding-guidelines-tests.md`: テスト品質調査の材料として扱う。

## 対象実装範囲

状態: 固定。

固定理由:

人間は 2026-05-25 に「広くてもやって」と指示した。
この指示により、前回の候補を広い調査範囲として固定する。

固定範囲:

- `frontend/src/application/`: frontend 側の共有 gateway contract
- `frontend/src/controller/wails/`: generated binding wrapper と Wails adapter
- `frontend/src/main.ts`: frontend 側の手動 DI 入口
- `frontend/wailsjs/`: generated bindings。読む候補であり、手編集しない候補
- `internal/controller/`: backend 側の Wails bind 入口と DTO 写像
- `internal/bootstrap/`: backend 側の wiring

## 対象テスト範囲

状態: 固定。

固定理由:

人間は 2026-05-25 に「広くてもやって」と指示した。
この指示により、対象実装範囲に対応する既存テストを広く調査範囲として固定する。

固定範囲:

- frontend 側 gateway、controller、画面操作の既存テスト
- backend 側 controller、usecase 境界、DTO 写像の既存テスト
- 接続境界を保護するシナリオテストまたは統合テスト

## 変更禁止範囲

人間指定: 未指定。

既定の禁止範囲:

- 承認済み `implementation-scope` が作られるまで、プロダクトコードを変更しない。
- 承認済み `implementation-scope` が作られるまで、プロダクトテストを変更しない。
- docs 正本本文を `refactor-lane` が直接変更しない。
- `frontend/wailsjs/` は generated bindings として手編集しない候補にする。
- remote repository を変更しない。

## 検証要件

状態: 固定。

固定範囲:

- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite structure`
- UI 表示または frontend と backend の接続境界を変更した場合は、実装後ブラウザ確認を行う。

## 成果物DAG 状態

| 成果物ID | 状態 | 根拠 |
| --- | --- | --- |
| `task 枠` | 完了 | この `plan.md` に task-id、作業計画フォルダ、目的、対象参照候補、未固定入力を記録した。 |
| `branch 準備` | 完了 | `codex/frontend-backend-connection-refactor` を作成し、`master` の存在を確認した。 |
| `仕様乖離整理` | 完了 | `investigator` が `spec-implementation-drift.md` を更新し、要件、詳細仕様、画面仕様に基づく人間判断待ちは `なし` と記録した。 |
| `仕様実装優先判断` | 完了 | 人間判断待ちがないため、空集合として完了扱いにする。 |
| `構造品質調査` | 完了 | `investigator` が `structure-quality-investigation.md` を作成し、`SQ-FBC-001` から `SQ-FBC-003` までを記録した。 |
| `テスト品質調査` | 完了 | `investigator` が `test-quality-investigation.md` を作成し、`TQI-FBC-001` から `TQI-FBC-003` までを記録した。 |
| `リファクタ範囲確認` | 完了 | 人間が「全部承認」と指示したため、`SQ-FBC-001` から `SQ-FBC-003`、`TQI-FBC-001` から `TQI-FBC-003` までを承認済みとした。 |
| `実装範囲` | 完了 | `designer` が `implementation-scope.md` を作成し、`FBC-FE-001`、`FBC-INT-001`、`FBC-UT-FE-001`、`FBC-UT-BE-001`、`FBC-SC-001` へ分割した。 |
| `実装引き継ぎ入力` | 完了 | `implementation-handoff-input.md` に wave と `FBC-FE-001` の起動入力を記録した。 |
| `frontend リファクタ` | 完了 | `frontend_implementer` が `FBC-FE-001` を実装し、`frontend-local` が通過した。 |
| `統合境界リファクタ` | 完了 | `integration_implementer` が `FBC-INT-001` を実装した。`backend-local` は通過し、`frontend-local` は後続 `FBC-UT-FE-001` の test 整理待ちで失敗した。 |
| `単体テスト` | 完了 | `implementation_unit_tester` が `FBC-UT-FE-001` と `FBC-UT-BE-001` を実装した。`frontend-local`、`backend-local`、`coverage` が通過した。 |
| `シナリオテスト` | 完了 | `implementation_scenario_tester` が `FBC-SC-001` を追加した。sandbox 内 system-test は Wails dev 起動制約で失敗し、sandbox 外 system-test は通過した。 |
| `最終検証` | 完了 | `frontend-local`、`backend-local`、`structure`、sandbox 外 `system-test`、`git diff --check` が通過した。 |
| `実装後ブラウザ確認` | 完了 | `browser_confirmation` が provider settings 実画面で `Gateway: 接続準備済み`、provider 3 件、`Health(): {"status":"ok"}`、秘匿値非表示を確認した。 |
| `レビュー通過根拠` | 着手可能 | 最終検証と実装後ブラウザ確認が完了したため、観点別レビュー agent を起動できる。 |

## 未固定入力

- 変更禁止範囲: 今回触ってはいけない path、生成物、画面、API

補足:

変更禁止範囲は既定の禁止範囲を採用する。
追加の禁止範囲がある場合は、人間が後続で指定する。

## 次アクション

`refactor_lane` は観点別レビュー agent を起動する。
レビューで修正必須指摘がなければ、docs正本化判断へ進む。
