# 詳細仕様: 境界結合テスト

- `detail_spec_id`: `boundary-integration-test`
- `status`: `approved`
- `source_artifacts`: `docs/exec-plans/completed/master-dictionary-boundary-test/boundary-integration-test.spec.md`, `docs/exec-plans/completed/master-dictionary-boundary-test/plan.md`
- `implementation_artifacts`: `docs/exec-plans/completed/master-dictionary-boundary-test/plan.md`（master-dictionary pilot で実装検証通過、本 task 範囲内全 pass）
- `review_artifacts`: `docs/exec-plans/completed/master-dictionary-boundary-test/plan.md`（finalization-module 正本化判断節に user 承認記録 2026-06-06）

## 親要件と仕様

### `boundary-integration-test-REQ-001` 境界結合テスト種別の定義

親要件:
frontend と backend を Wails 境界で接続する repo において、境界 API 形式の整合を「片側変更で他方が落ちる検出経路」として固定する。

仕様:
- scenario E2E は UI 起点であり、UI に露出しない境界 DTO field 値の変化を検出できないことを前提とする
- 単体テストは単一責務の検証であり、境界全体の整合を保証しないことを前提とする
- 境界結合テストは両者の間の隙間を埋めるテスト種別である
- 本仕様は複数の skill が共通の認識として参照する資料の位置付けを持つ。test-design、tests-scenario、tests-unit、design-module、implementation-module、investigation-module のいずれの skill から参照しても、同じ前提でテストを書ける形を保つ

### `boundary-integration-test-REQ-002` 上位入力からの導出原則

親要件:
境界結合テストは独立に設計しない。上位入力からの導出物として作る。

仕様:
- 形式（型、必須、null 許容、値域）は機械可読の契約成果物を上位入力とする。frontend と backend が同一の真として参照可能な形式であることを成立条件とする
- 意味（field が業務的に何を指すか）は UC（usecase 仕様）を上位入力とする
- 表示（境界に必要な field の最小集合）は画面設計を上位入力とする
- 境界結合テストの観点と case は、上記 3 系統の上位入力から機械的に展開する
- 観点表や case 表を白紙から起こさない

### `boundary-integration-test-REQ-003` 責務範囲（扱う対象）

親要件:
境界結合テストは「形式の整合」だけを検証対象として扱う。

仕様:
- 境界 contract の型、必須、null 許容、値域の固定を扱う
- 状態遷移の前後で contract 上の field がどう変化するかを扱う。UC が定める遷移を境界実装側で証明する
- 片側書き換え検出（contract または共有固定値の片側変更で他方が落ちる経路）を成立条件とする

### `boundary-integration-test-REQ-004` 責務範囲（扱わない対象）

親要件:
境界結合テストは形式の整合だけを扱い、意味、表示、内部実装、UI シナリオは扱わない。

仕様:
- 業務的意味（field が業務的に何を指すか）は扱わない。UC docs の正本確認の責務とする
- 表示文言、UI 構造、style は扱わない。画面設計の正本確認と storybook-module の責務とする
- 単一関数 / メソッドの内部実装は扱わない。unit test の責務とする
- UI 起点の業務シナリオ証明は扱わない。scenario E2E の責務とする

### `boundary-integration-test-REQ-005` 最小構成要素

親要件:
1 feature あたりの境界結合テストは、両側が同一 contract を真として参照する構成で成立する。

仕様:
- 1 feature あたり、形式 contract、backend 境界結合テスト、frontend 境界結合テストの 3 つを必須とする
- shared 固定値（golden 等）は任意とし、採用する場合は backend と frontend の両 test が同一物理 file を read-only で参照する
- backend と frontend は同一 contract を真として参照する。両側が異なる contract を参照すると片側書き換え検出経路が成立しない
- shared 固定値を使う場合、assert 期待値は test source 内のリテラルにハードコードし、固定値 file は mock 応答にのみ流用する。固定値を mock と assert 期待値の両方に使うと、固定値の片側変更で両側が同時に変わり検出経路が成立しない

### `boundary-integration-test-REQ-006` 既存テスト種別との責務境界

親要件:
境界結合テストは既存テスト種別と責務が重ならないよう、明確に位置付ける。

仕様:
- 単体テスト: 関数 / メソッドを入口として、単一責務の実装可否を扱う
- 境界結合テスト: Wails 境界を入口として、contract と両側実装の整合を扱う
- scenario E2E: UI 操作を入口として、業務シナリオの完走を扱う
- backend apitest（既存）: 公開 method を入口として、受け入れ条件の system-level 証明を扱う
- 境界結合テストは「scenario E2E と単体テストの中間層」として位置付く。UI 起点ではないため scenario E2E ではなく、単一メソッドの内部実装ではないため単体テストでもない

## 根拠

- `source_artifacts`: `docs/exec-plans/completed/master-dictionary-boundary-test/boundary-integration-test.spec.md`（5 章 / 85 行、user 承認済み 2026-06-06）
- `implementation_artifacts`: `docs/exec-plans/completed/master-dictionary-boundary-test/plan.md`（master-dictionary pilot で実装検証通過、本 task 範囲内全 pass）
- `review_artifacts`: `docs/exec-plans/completed/master-dictionary-boundary-test/plan.md`（finalization-module 正本化判断節に user 承認記録 2026-06-06）
