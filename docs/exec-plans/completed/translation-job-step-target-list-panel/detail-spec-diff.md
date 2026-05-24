# 詳細仕様差分: 翻訳ジョブステップ処理対象一覧表示パネル

- `skill`: detail-spec-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `detail_spec_target`: `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `screen_design_diff`: `./screen-design-diff.job-run.md`, `./screen-design-diff.term-translation-phase.md`, `./screen-design-diff.persona-generation-phase.md`, `./screen-design-diff.body-translation-phase.md`, `./screen-design-diff.translation-complete.md`
- `component_diagram`: `./design-diff.translation-job-step-target-list-panel.md`

## 詳細仕様差分

### `translation-job-management-REQ-006` 利用者向け情報を主要目的へ絞る

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/translation-job-management.md`

親要件:
利用者は未完了候補、操作可否、実行対象を主要情報として判断できる。

仕様:
- 処理対象は、現在段階で処理、生成、確認する実体を表す。
- 処理対象名は、処理対象を短く表す実体名である。
- 処理対象詳細は、その実体が現在段階でどう処理、生成、確認されるかを表す利用者向け情報である。
- 利用者は、選択した翻訳ジョブの現在段階で処理するものと、その詳細を判断できる。
- 処理対象一覧は、50 件程度を既定ページサイズとして扱う。
- 利用者は、処理対象一覧のページを切り替え、現在ページの表示範囲を操作できる。
- 数万件レベルの処理対象でも、画面要素は現在ページの表示範囲に限定し、利用者操作を継続できる。

未決:
- なし

回答:
- なし

### `term-translation-phase-REQ-007` 操作と状態を利用者が判別できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
利用者は開始、一時停止、再開、再試行、取り消しの可否と理由を判断できる。

仕様:
- 利用者は、単語翻訳フェーズの処理対象名が共通辞書対象外の用語と固有名詞であり、AIサービスへ送り、確定訳語として翻訳ジョブ内辞書へ保存するものであることを判断できる。

未決:
- なし

回答:
- なし

### `persona-generation-phase-REQ-007` 操作と状態を利用者が判別できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/persona-generation-phase.md`

親要件:
利用者はフェーズの進行、生成対象、生成結果、ペルソナ参照状態、本文翻訳フェーズの開始可否を判断できる。

仕様:
- 利用者は、NPC ペルソナ生成フェーズの処理対象名が NPC ごとのペルソナ生成入力であり、NPC の原文発話、NPC 属性、会話文脈、共通ペルソナ参照からペルソナ参照情報を作るものであることを判断できる。

未決:
- なし

回答:
- なし

### `body-translation-phase-REQ-006` 進行と結果を判断できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/body-translation-phase.md`

親要件:
利用者は本文翻訳フェーズの進行、結果、失敗理由、成果物出力条件を判断できる。

仕様:
- 利用者は、本文翻訳フェーズの処理対象名が辞書置換対象外の翻訳項目であり、辞書とペルソナを参照して訳文を作るものであることを判断できる。

未決:
- なし

回答:
- なし

### `body-translation-phase-REQ-005` 翻訳ジョブ全体の完了と成果物出力条件を作る

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/body-translation-phase.md`

親要件:
本文翻訳フェーズの完了時点で翻訳ジョブ全体は `Completed` になり、完了済み翻訳ジョブから成果物を出力できる。

仕様:
- 利用者は、翻訳完了確認の処理対象名が翻訳項目単位の訳文であり、本文翻訳で保持された訳文として出力管理へ進む前に確認するものであることを判断できる。

未決:
- なし

回答:
- なし

## 根拠

- `source`: `docs/spec.md` の翻訳ジョブ実行進捗、翻訳フロー、各フェーズ、翻訳ジョブ状態。
- `source`: `docs/detail-specs/translation-job-management.md` の `translation-job-management-REQ-006`。
- `source`: `docs/detail-specs/term-translation-phase.md` の `term-translation-phase-REQ-007`。
- `source`: `docs/detail-specs/persona-generation-phase.md` の `persona-generation-phase-REQ-007`。
- `source`: `docs/detail-specs/body-translation-phase.md` の `body-translation-phase-REQ-005` と `body-translation-phase-REQ-006`。
- `source`: `task-frame.md` の依頼、対象段階、対象情報。
- `review`: 人間差し戻しを反映し、人間設計レビュー待ち。
- `validation`: 文書整形確認済み。
