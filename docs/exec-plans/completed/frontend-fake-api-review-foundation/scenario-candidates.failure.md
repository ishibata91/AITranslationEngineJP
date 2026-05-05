# Scenario Candidates: frontend-fake-api-review-foundation / failure

- `generator`: `failure`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `FFA`
- `candidate_count`: 6

## Generator Scope

- `viewpoint`: 失敗
- `included_sources`: `./task-frame.md`, `tasks/usecases/frontend-fake-api-review-foundation.yaml`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `tmp/code-map/index.json`
- `excluded_sources`: プロダクトコード, プロダクトテスト, docs 正本更新, `.codex` 変更
- `generation_notes`: fakeAPI 未選択、バインディング 混入、モックデータ欠落、本番混入、状態パターン 不整合、coverage 例外理由欠落だけを候補化する。採否、統合、最終シナリオ表は扱わない。

## Candidate Scenarios

### CAND-FFA-001 fakeAPI 起動モードで fakeAPI が選ばれない

- `根拠要件`: `task-frame.md:17`, `task-frame.md:30`, `frontend-fake-api-review-foundation.yaml:20`, `frontend-fake-api-review-foundation.yaml:33`, `tmp/code-map/index.json:932-1065`
- `viewpoint`: 設定不整合
- `candidate scenario id`: `CAND-FFA-001`
- `actor`: フロントエンドレビュー実行者
- `失敗開始条件`: レビュー起動モードで起動したが、フロントエンドの API 接続先が fakeAPI に切り替わらない。
- `拒否する操作`: 実画面レビューを Wails バインディング または バックエンド 依存なしで開始する操作を拒否する。
- `expected error`: レビュー起動モードで fakeAPI が未選択である理由を表示または検証結果に記録する。
- `観測点`: 起動時の接続先判定、root view へ注入された ゲートウェイ 種別、fakeAPI 起動モードの局所テスト結果。
- `関連詳細要求タイプ`: 起動モード、DI 境界、レビュー実行前提
- `採用判断材料`: designer は レビュー起動モードの切替失敗を、起動不能扱いにするか表示上の設定不足扱いにするか判断する。
- `競合注意`: 起動失敗を自動 フォールバック で Wails バインディング に流す候補と競合する。

### CAND-FFA-002 View または ScreenController に Wails バインディング が混入する

- `根拠要件`: `task-frame.md:27-30`, `architecture.md:17`, `architecture.md:41-47`, `architecture.md:101-117`, `coding-guidelines-frontend.md:34-38`, `tmp/code-map/index.json:3872-3963`
- `viewpoint`: 設定不整合
- `candidate scenario id`: `CAND-FFA-002`
- `actor`: フロントエンドレビュー実行者
- `失敗開始条件`: fakeAPI 起動モードでも View、ScreenController、Frontend UseCase のいずれかが 生成済み `wailsjs` または Wails adapter を直接参照する。
- `拒否する操作`: fakeAPI 起動モードで対象画面を開く操作、または状態パターンを切り替える操作を拒否する。
- `expected error`: Wails バインディング 依存の混入箇所を、利用者向けメッセージではなく内部診断または検証結果に分けて記録する。
- `観測点`: layer import の検査結果、`フロントエンド/src/main.ts` から ゲートウェイ を注入する境界、`フロントエンド/src/controller/wails/` 以外の generated バインディング 参照有無。
- `関連詳細要求タイプ`: Wails 境界、責務境界、DI 境界
- `採用判断材料`: designer は バインディング 混入をシナリオの失敗条件にするか、責務レビュー側の補助条件に残すか判断する。
- `競合注意`: architecture では `Gateway -> 生成済み `wailsjs` -> バックエンド Controller` が本番経路であるため、本番ゲートウェイ の正当な import まで拒否しないよう統合時に境界を分ける必要がある。

### CAND-FFA-003 画面固有のモックデータが欠落する

- `根拠要件`: `task-frame.md:19-21`, `frontend-fake-api-review-foundation.yaml:22-24`, `frontend-fake-api-review-foundation.yaml:35-36`, `coding-guidelines-frontend.md:52-55`
- `viewpoint`: 参照不能
- `candidate scenario id`: `CAND-FFA-003`
- `actor`: フロントエンドレビュー実行者
- `失敗開始条件`: 対象画面の状態パターンを開く時に、ユースケース task 側の モックデータが存在しない、または必要項目が欠落している。
- `拒否する操作`: 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態のいずれかを表示する操作を拒否する。
- `expected error`: 欠落した モックデータ 名、対象画面、対象 状態パターンを検証結果に記録する。
- `観測点`: モックデータ 登録一覧、状態パターンごとの テストデータ 解決結果、agent-browser で対象 状態パターンが開けたかどうか。
- `関連詳細要求タイプ`: レビューstate モックデータ、画面状態、実画面レビュー
- `採用判断材料`: designer は モックデータ欠落時の画面表示を「レビュー不能」とするか「設定不足状態」として表示確認対象に含めるか判断する。
- `競合注意`: 設定不足状態そのものを確認するシナリオと、設定不足によりレビュー不能になるシナリオが競合しうる。

### CAND-FFA-004 fakeAPI または モックデータが本番起動へ混入する

- `根拠要件`: `task-frame.md:18`, `task-frame.md:22`, `task-frame.md:34-36`, `frontend-fake-api-review-foundation.yaml:21`, `frontend-fake-api-review-foundation.yaml:25`, `frontend-fake-api-review-foundation.yaml:37`, `tmp/code-map/index.json:932-1065`
- `viewpoint`: 設定不整合
- `candidate scenario id`: `CAND-FFA-004`
- `actor`: フロントエンドレビュー実行者
- `失敗開始条件`: 本番起動で fakeAPI、レビューstate モックデータ、または 本番初期状態 への モックデータ 注入が有効になる。
- `拒否する操作`: 本番起動で対象画面を開く操作、または本番 API 経路で モックデータを読む操作を拒否する。
- `expected error`: 本番起動で fakeAPI が選ばれた事実、混入した モックデータの出所、本番初期状態 への影響有無を検証結果に記録する。
- `観測点`: 本番起動時の ゲートウェイ 選択結果、永続化経路への write 有無、本番初期状態 の初期値、局所テスト結果。
- `関連詳細要求タイプ`: 本番 API 混入防止、永続化境界、本番初期状態
- `採用判断材料`: designer は本番混入を hard failure として扱うか、レビュー起動モードの失敗とは別シナリオに分離するか判断する。
- `競合注意`: fakeAPI 起動モードの利便性のために モックデータを共通 初期状態 へ置く候補と競合する。

### CAND-FFA-005 状態パターンの一覧と実画面表示が一致しない

- `根拠要件`: `task-frame.md:19-21`, `frontend-fake-api-review-foundation.yaml:16`, `frontend-fake-api-review-foundation.yaml:22-24`, `frontend-fake-api-review-foundation.yaml:35-36`, `architecture.md:119-123`
- `viewpoint`: 失敗入力
- `candidate scenario id`: `CAND-FFA-005`
- `actor`: フロントエンドレビュー実行者
- `失敗開始条件`: 状態パターン一覧には存在する状態が実画面で表示できない、または実画面に表示された状態が 状態パターン一覧に存在しない。
- `拒否する操作`: 状態パターン切替、対象状態の agent-browser 確認、または表示証跡の完了記録を拒否する。
- `expected error`: 不一致になった 状態パターン 名、期待表示、実際表示、対象 画面状態 を検証結果に記録する。
- `観測点`: 状態パターン一覧、Store の 画面状態、Presenter が作る 表示モデル、agent-browser の表示確認結果。
- `関連詳細要求タイプ`: 実画面レビュー用状態パターン、画面状態、Presenter / Store
- `採用判断材料`: designer は 状態パターン一覧の不足を入力不備として扱うか、表示実装の失敗として扱うか判断する。
- `競合注意`: 画面ごとの 状態パターン 名を統一する候補と、画面固有の状態名を許容する候補が競合しうる。

### CAND-FFA-006 coverage harness の例外理由または局所テスト結果が欠落する

- `根拠要件`: `frontend-fake-api-review-foundation.yaml:26-27`, `task-frame.md:23`, `frontend-fake-api-review-foundation.yaml:34`
- `viewpoint`: 参照不能
- `candidate scenario id`: `CAND-FFA-006`
- `actor`: フロントエンドレビュー実行者
- `失敗開始条件`: fakeAPI 基盤を coverage harness の数値判定から例外扱いにしたが、例外理由または局所テスト結果が記録されていない。
- `拒否する操作`: fakeAPI 起動モードの完了判定、または coverage 例外を含む検証証跡の受理を拒否する。
- `expected error`: coverage 例外理由の欠落、局所テスト結果の欠落、対象範囲の欠落を検証結果に記録する。
- `観測点`: coverage harness の例外記録、fakeAPI 起動モードの局所テスト結果、完了条件との対応。
- `関連詳細要求タイプ`: coverage 例外理由、局所テスト、完了根拠
- `採用判断材料`: designer は coverage 例外理由をシナリオ必須条件へ含めるか、実装レーンの検証証跡条件へ残すか判断する。
- `競合注意`: coverage の数値判定を優先する候補と、fakeAPI 基盤を局所テストで証明する候補が競合しうる。

## Open Notes

- `人間判断候補`: fakeAPI 未選択時の扱いを、起動不能、設定不足表示、検証失敗のどれにするかは未確定である。
- `人間判断候補`: モックデータ欠落時に、設定不足状態として表示対象に含めるか、レビュー不能として止めるかは未確定である。
- `統合候補`: `CAND-FFA-001` と `CAND-FFA-004` は ゲートウェイ 選択失敗として統合できる可能性がある。
- `統合候補`: `CAND-FFA-003` と `CAND-FFA-005` は 状態パターン 表示不能として統合できる可能性がある。
- `不採用候補`: バインディング 混入は scenario ではなく責務境界レビュー条件へ移す判断もありうる。
- `不採用候補`: coverage 例外理由欠落はシナリオではなく完了根拠条件へ移す判断もありうる。
