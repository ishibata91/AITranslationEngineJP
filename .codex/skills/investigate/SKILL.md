---
name: investigate
description: Codex 側の設計前調査、探索テスト証跡、修正レーンの UC 差分候補、リファクタ調査の作業プロトコル。観測事実、UI 証跡、ログ、未確認事項を根拠 first で扱う判断基準を提供する。
---
# Investigate

## 目的

`investigate` は作業プロトコルである。
`investigator` agent が設計前に必要な証拠、探索テストレーンの探索証跡、修正レーンの UC 差分候補、リファクタレーンの仕様乖離整理、構造品質調査、テスト品質調査を集めるための、観測事実、UI 証跡、ログ、仮説、残り 不足 の分け方を提供する。

設計前調査では UI check 専用 skill / agent は置かない。
設計前の画面設計根拠は `investigator` が `investigate` の一部として扱う。

## 対応ロール

- `investigator` が使う。
- 返却先は 呼び出し元 または次 agent とする。
- 担当成果物は `investigate` の出力規約で固定する。
- 探索テストレーンでは担当成果物を `探索証跡` に限定する。
- 修正レーンでは担当成果物を `UC 差分候補` に限定する。
- リファクタレーンでは担当成果物を `仕様乖離整理`、`構造品質調査`、`テスト品質調査` に限定する。

## 呼び出し元から渡される情報

- 必須入力: 呼び出し元、investigation_goal、known_context を受け取る。
- 非必須入力: investigation_mode、reproduction_steps、candidate_paths、探索計画、テストデータを受け取る。
- 非必須調査種別: `investigation_mode` は `再現`、`画面設計根拠`、`trace`、`リスク報告`、`探索テスト証跡`、`UC 差分候補`、`仕様乖離整理`、`構造品質調査`、`テスト品質調査` のいずれかを受け取る。
- 必須成果物: active task 文脈 または 呼び出し元提供 investigation 文脈を受け取る。

## 作業前に読む正本

- エージェント実行定義と実行境界は [investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/investigator.toml) に従う。
- エージェント実行定義: [investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/investigator.toml)
- 実行境界: エージェント実行定義に従う
- 画面設計は [screen-design](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/README.md) に従う。
- active plan の画面設計差分は `docs/exec-plans/active/<task-id>/screen-design-diff.<screen-id>.md` とする。
- 画面設計差分と docs の画面設計が両方ある場合は、active plan の画面設計差分を優先して読む。
- `agent-browser` CLI の利用規約は [agent-browser.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/references/agent-browser.md) に従う。
- 観測ログ仕様は [observability-logging.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/observability-logging.md) に従う。
- 探索テストレーンの探索計画は [exploration-test-planning](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/exploration-test-planning/SKILL.md) に従う。
- 探索テスト証跡の雛形は [exploration-test-evidence.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/investigate/assets/exploration-test-evidence.md) とする。
- 探索テスト証跡の task 内 artifact は `docs/exec-plans/active/<task-id>/exploration-test-evidence.md` とする。
- リファクタレーンは [refactor-lane](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/refactor-lane/SKILL.md) に従う。
- リファクタ分類表雛形は [refactor-classification.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/refactor-lane/assets/refactor-classification.md) とする。
- 構造設計正本は [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) とする。
- コーディング規約入口は [coding-guidelines.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines.md) とする。
- テスト規約は [coding-guidelines-tests.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-tests.md) とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/investigate/SKILL.md

## skill 内の拘束条件

### 拘束観点

- `再現`、`画面設計根拠`、`trace`、`リスク報告` の観点
- 観測済み事実、画面設計根拠、仮説 の分離
- 探索計画、テストデータ、探索証跡 の分離
- 仕様参照、実装参照、仕様乖離、人間判断待ち の分離
- 構造品質観点、未使用コード、根拠参照、対象範囲、変更不要範囲 の分離
- テスト品質観点、テスト参照、仕様参照、変更不要テスト範囲 の分離
- 根拠 path と再現条件の残し方
- 画面設計を操作経路と期待値の補助参照にする条件
- 画面要素を操作できない理由の分類
- 設計を止める 残留リスク の表現

## 担当ロールが判断してよい範囲

- 根拠 のない結論を書かない
- 観測事実と仮説を混ぜない
- 設計前の画面設計根拠は `agent-browser` CLI で確認する
- 画面設計根拠を集める時は、関連する `docs/screen-design/screens/*.md` または active plan の `screen-design-diff.<screen-id>.md` を、操作経路と期待値の補助参照として確認する
- 画面設計根拠は画面状態、console、screenshot、操作条件を分けて残す
- 画面要素を操作できない時は、操作経路に対応する画面設計がない、画面ID または セレクタ（`aria-label`）が不足、起動状態が不足、操作経路または期待値と実画面の対応を確認できない、のいずれかで理由を返す
- frontend log は browser console の観測事実として残す
- backend log は `tmp/logs/wails-dev.log` の観測事実として残す
- frontend log と backend log は同じ根拠 path に混ぜない
- 実装 レーン の調査は Codex implementation レーンへ戻す
- 探索テスト証跡は探索計画とテストデータを超えない
- 探索テスト証跡では探索範囲を広げる判断をしない
- 探索テスト証跡では原因仮説を固定しない
- `UC 差分候補` は既存期待を説明するユースケース記述の不足候補だけを扱い、新規仕様判断を確定しない
- 仕様乖離整理では仕様と実装の差だけを記録し、どちらを正にするかを判断しない
- 構造品質調査では責務過多、責務分離不足、コーディング規約逸脱、構造設計不整合、未使用コードを根拠別に分ける
- 構造品質調査では変更不要範囲を実装範囲候補に含めない
- テスト品質調査では `coding-guidelines-tests.md` の良いテストの品質観点に従う
- テスト品質調査ではテスト規約違反と仕様整合性の不足を根拠別に分ける
- テスト品質調査では変更不要テスト範囲を実装範囲候補に含めない

- observed、画面設計根拠、inferred を分ける
- 証跡 path と再現条件を優先する
- 設計継続可否に効く 不足 を残す
- active 規約 は agent に対して 1 ファイルだけ置く。調査種別は selector で扱う。

## skill が扱わない対象

- implementation-scope 承認後の再現、再観測、実装時調査は扱わない。
- 恒久修正、プロダクトテスト追加、implementation レビューは扱わない。
- 承認済み実装範囲や対象 file は確定しない。
- 修正レーンの修正方針判断、修正実行入力、E2E テスト観点差分、レビュー通過根拠は扱わない。
- リファクタレーンの仕様実装優先判断、リファクタ範囲確認、implementation-scope は扱わない。
- テスト修正方針、テスト追加方針、テスト削除方針は確定しない。
- 設計前調査で UI check 専用 agent を前提にしない。
- 探索計画の作成、バグ一覧の集約、影響ファイルの確定は扱わない。

## 返す成果物

- 判断結果: 設計前調査の完了、未完了、停止の判定を返す。
- 根拠参照: 調査判断に使った資料、画面、観測結果を返す。
- 不足情報: 設計判断に不足している項目を返す。
- 次判断材料: 次 agent が判断できる材料を返す。
- 引き継ぎ先: `designer` を返す。
- 渡す対象範囲: 観測済み事実、仮説、残り 不足、残留 risks を返す。
- 調査 mode: 実施した調査の種類を返す。
- 観測事実: 観測済み事実だけを返す。
- UI 証跡: UI を確認した場合は証跡と参照先を返す。
- ログ証跡: ログを確認した場合は証跡と参照先を返す。
- 仮説: 事実と分けて原因候補を返す。
- 観測点: 確認した入口、経路、対象を返す。
- 操作不能理由: UI 要素を操作できない場合の分類と根拠を返す。
- 探索証跡: 探索テスト証跡の場合は探索計画とテストデータに対応する観測事実を返す。
- UC 差分候補: 修正レーンの場合は既存期待を説明するユースケース記述の差分有無、記述不足候補、新規判断必要箇所を返す。
- 仕様乖離整理: リファクタレーンの場合は仕様参照、実装参照、差分内容、影響範囲、人間判断待ちを返す。
- 構造品質調査: リファクタレーンの場合は責務過多、責務分離不足、コーディング規約逸脱、構造設計不整合、未使用コード、変更不要範囲を返す。
- テスト品質調査: リファクタレーンの場合はテスト規約観点別結果、仕様整合性、変更不要テスト範囲を返す。
- 残り不足: 未確認事項と理由を返す。
- 残留リスク: 設計判断に残る リスク を返す。
- 推奨 next step: 設計継続、追加調査、停止のどれが妥当かを返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 作業を完了できる条件

- 出力規約を満たし、次の 実行者 が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- 観測事実、画面設計根拠、仮説、未観測 不足 を分けた。
- 根拠 path、再現条件、UI check 対象範囲 を残した。
- 画面設計根拠を扱った場合は、参照した画面設計または画面設計差分を返した。
- 操作できない画面要素がある場合は、操作不能理由が分類されている。
- 探索テスト証跡の場合は、探索計画、テストデータ、観測事実、UI 証跡、ログ証跡、未確認事項を分けた。
- 探索テスト証跡の場合は、`exploration-test-evidence.md` に証跡が記録されている。
- `UC 差分候補` の場合は、差分なし、記述不足、新規判断必要を分けた。
- 仕様乖離整理の場合は、仕様参照、実装参照、差分内容、影響範囲、人間判断待ちを分けた。
- 構造品質調査の場合は、責務過多、責務分離不足、コーディング規約逸脱、構造設計不整合、未使用コード、変更不要範囲を分けた。
- テスト品質調査の場合は、テスト規約観点別結果、仕様整合性、変更不要テスト範囲を分けた。
- design continuation に必要な リスク を返した。
- 必須 根拠: 観測済み事実 根拠, 画面設計根拠 when mode is 画面設計根拠, reproduction condition, 根拠 path when used
- 完了判断材料: designer が設計継続か停止かを判断できる。
- 残留リスク: 設計判断に残る リスク が返っている。

## 作業を止める条件

- implementation-scope 承認後の再現や再観測を扱う時
- 恒久修正や プロダクトテスト 追加が必要な時
- implementation レビュー が主目的の時
- 観測条件が不足する場合は停止する。
- 必須入力または必須成果物が不足する場合は停止する。
- 恒久修正が必要なら `designer` へ戻す。
- 実装時調査なら、Codex implementation レーン [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate/SKILL.md) を使う前提で `designer` へ戻す。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 根拠 なしの結論を書く必要がある場合は停止する。
- 設計前調査で UI check 専用 agent を前提にする場合は停止する。
- implementation-time investigation を扱う場合は停止する。
- 探索計画またはテストデータが不足した探索テスト証跡を扱う場合は停止する。
- 探索テスト証跡で探索範囲を広げる必要がある場合は停止する。
- `UC 差分候補` で新規仕様判断を確定する必要がある場合は停止する。
- 仕様乖離整理で仕様と実装のどちらを正にするか判断する必要がある場合は停止する。
- 構造品質調査でリファクタ範囲確認または実装範囲を確定する必要がある場合は停止する。
- テスト品質調査でリファクタ範囲確認、実装範囲、テスト修正方針を確定する必要がある場合は停止する。
- 拒否条件: implementation-time investigation
- 拒否条件: permanent fix request
- 拒否条件: 根拠成果物 不足
