# .codex

このディレクトリは Codex 作業流れ の正本です。
Codex は設計 作業流れ、承認済み 対象範囲 からの実装、実装後 レビュー、docs 正本化を進めます。

プロダクト仕様と設計判断の正本は `docs/` です。
作業流れ、skill、agent、引き継ぎ 契約の正本は `.codex/` に置きます。
live 作業流れ の説明本文と判断基準の正本はこの `README.md` とします。
`.codex/workflow.md` は補助図であり、live 判断を上書きしません。
Codex 本体の内蔵ブラウザ利用規約は `.codex/browser-use.md` とします。
Storybook の作成、起動、分類、確認資源、`fixture` 種類基準の共通規約は `docs/references/storybook.md` とします。

## Live Skills

### 主 skill

- 新規実装レーン (`implement-lane`): `skills/implement-lane/SKILL.md`
- 設計壁打ち: `skills/wall-discussion/SKILL.md`
- design bundle 進行: `skills/design-bundle/SKILL.md`
- 設計前調査: `skills/investigate/SKILL.md`
- 詳細仕様差分 (`detail-spec-design`): `skills/detail-spec-design/SKILL.md`
- 実装スコープ (`implementation-scope`): `skills/implementation-scope/SKILL.md`
- 探索テスト計画 (`exploration-test-planning`): `skills/exploration-test-planning/SKILL.md`
- 探索テストレーン (`exploration-test-lane`): `skills/exploration-test-lane/SKILL.md`
- 軽量変更レーン (`light-change-lane`): `skills/light-change-lane/SKILL.md`
- UX 保守レーン (`ux-maintainance-lane`): `skills/ux-maintainance-lane/SKILL.md`
- リファクタレーン (`refactor-lane`): `skills/refactor-lane/SKILL.md`
- 軽量変更計画 (`light-change-planning`): `skills/light-change-planning/SKILL.md`
- 実装時調査 (`implementation-investigate`): `skills/implementation-investigate/SKILL.md`
- 修正レーン (`fix-lane`): `skills/fix-lane/SKILL.md`
- マージレーン (`merge-lane`): `skills/merge-lane/SKILL.md`
- 修正方針判断 (`fix-decision`): `skills/fix-decision/SKILL.md`
- 実装後ブラウザ確認 (`browser-confirmation`): `skills/browser-confirmation/SKILL.md`
- Storybook レビューループ (`story-book-review-loop`): `skills/story-book-review-loop/SKILL.md`
- 観測ログ追加 (`observability-implementer`): `skills/observability-implementer/SKILL.md`
- プロダクトコード 実装 重点 skill: `skills/implement-backend/SKILL.md`、`skills/implement-frontend/SKILL.md`、`skills/implement-integration/SKILL.md`
- シナリオテスト 実装 (`tests-scenario`): `skills/tests-scenario/SKILL.md`
- 単体テスト 実装 (`tests-unit`): `skills/tests-unit/SKILL.md`
- docs 正本化: `skills/updating-docs/SKILL.md`
- 作業流れ 契約保守 (`workflow-contract-maintenance`): `skills/workflow-contract-maintenance/SKILL.md`
- 実装後 レビュー 観点: `skills/codex-review-behavior/SKILL.md`、`skills/codex-review-contract/SKILL.md`、`skills/codex-review-trust-boundary/SKILL.md`、`skills/codex-review-state-invariant/SKILL.md`、`skills/codex-review-responsibility-boundary/SKILL.md`

### 補助 skill

- 図作成補助: `skills/diagramming/SKILL.md`
- 実装時調査 重点 skill: `skills/implementation-investigate-reproduce/SKILL.md`、`skills/implementation-investigate-trace/SKILL.md`、`skills/implementation-investigate-observe/SKILL.md`、`skills/implementation-investigate-reobserve/SKILL.md`

## Agent / Skill Boundary

- live Codex agent は新規実装レーン 進行役 (`implement_lane`)、修正レーン 進行役 (`fix_lane`)、探索テストレーン 進行役 (`exploration_test_lane`)、軽量変更レーン 進行役 (`light_change_lane`)、UX 保守レーン 進行役 (`ux_maintainance_lane`)、リファクタレーン 進行役 (`refactor_lane`)、マージレーン 進行役 (`merge_lane`)、軽量変更計画 agent (`light_change_planner`)、探索テスト計画 agent (`exploration_test_planner`)、設計成果物 agent (`designer`)、図作成 agent (`diagrammer`)、調査 agent (`investigator`)、修正方針判断 agent (`fix_decider`)、実装時調査 agent (`implementation_investigator`)、実装後ブラウザ確認 agent (`browser_confirmation`)、backend 実装 agent (`backend_implementer`)、frontend 実装 agent (`frontend_implementer`)、統合境界実装 agent (`integration_implementer`)、観測ログ追加 agent (`observability_implementer`)、シナリオテスト 実装 agent (`implementation_scenario_tester`)、単体テスト 実装 agent (`implementation_unit_tester`)、docs 更新 agent (`docs_updater`)、観点別 レビュー agent にする
- `implement_lane` は新規実装と機能拡張の task 内成果物 DAG、HITL、引き継ぎ、branch 準備、作業 commit、マージ準備入力を管理する
- `fix_lane` は人間が確認した不具合、レビュー非通過、検証失敗の task 内成果物 DAG、起動入力、担当 agent 起動、停止、戻し、close 条件を管理する
- `exploration_test_lane` は探索テストの task 内成果物 DAG、起動入力、停止、戻し、close 条件を管理する。探索計画、探索証跡、実装、回帰確認の担当 agent を分ける
- `light_change_lane` は既存仕様の意味を大きく広げない軽い backend / frontend / integration 変更の task 内成果物 DAG、軽量変更計画、起動入力、人間確認、テスト修正、レビュー、正本化判断、close 条件を管理する
- `ux_maintainance_lane` は Storybook 上の人間指摘を起点に、frontend 修正、frontend 整理、接続整合証跡、統合メンテ、単体テストメンテ、ハーネス通過、docs 正本化判断、close 条件を管理する
- `ux_maintainance_lane` は人間の仕様変更指示がある場合だけ、画面表示以外の仕様変更と詳細仕様正本反映を扱う
- `refactor_lane` は仕様乖離整理起動入力、人間による仕様実装優先判断、構造品質調査起動入力、テスト品質調査起動入力、リファクタ範囲確認、`implementation-scope`、修正実行入力、レビュー、docs 正本化判断、close 条件を管理する
- `merge_lane` は各 active plan のマージ準備入力を読み、source branch を target branch へ local merge し、conflict 解消、merge 後検証、completed 移動、merge 結果 commit を管理する
- `light_change_planner` は人間要望、仕様製本、関連 docs、task-local 成果物、既存実装を突き合わせ、軽量変更として進めるか、設計または修正レーンへ戻すかを判断する
- `exploration_test_planner` は探索計画だけを作り、観測、ログ確認、画面確認、原因仮説の作成を扱わない
- `designer` は呼び出し元レーンから渡された設計対象と根拠参照を読み、`detail-spec-diff.md` を詳細仕様差分として作る。画面変更がある時は `screen-design-diff.<screen-id>.md` を揃える。人間レビュー後に `implementation-scope` を固定する
- `diagrammer` は `diagramming` に従い、人間設計レビュー前または軽量変更の実装着手前に、予定変更箇所だけの追加・削除差分を示すコンポーネント図を標準形として作る。シーケンス図、状態遷移図、その他の図はユーザー要求または複雑性がある場合だけ扱う。修正レーン標準形の原因箇所シーケンス図は扱わない
- `fix_decider` は `fix-decision` に従い、修正前調査から原因の原因、責務境界、採用する修正方針、禁止する修正、原因箇所シーケンス図を固定する。修正前調査、人間修正レビュー、修正実行入力は扱わない
- `browser_confirmation` は実装後ブラウザ確認の軽量実行だけを扱う。確認経路と期待値は `implement_lane`、`fix_lane`、`light_change_lane` が定義し、`browser_confirmation` は期待値の妥当性を判断しない
- `observability_implementer` は `implement_lane` の `観測ログ追加` で、最終検証前に完成済み成果物を読み、実行後に消える原因分離材料を残す恒久ログだけを追加する
- `designer`、`exploration_test_planner`、`investigator`、`docs_updater` は 文脈 を引き継がず、引き継ぎ入力 だけで動く
- `implement_lane` は承認済み 実行成果物 を実行正本にし、`designer`、`diagrammer`、`implementation_investigator`、`backend_implementer`、`frontend_implementer`、`integration_implementer`、`implementation_scenario_tester`、`implementation_unit_tester`、`observability_implementer`、`browser_confirmation` を 文脈 継承なしで直接 起動 する。設計中の `designer` と `diagrammer` は担当する設計成果物が承認済みまたは停止として記録されるまで閉じず、差し戻しまたは追加質問を会話文脈を維持した同じ agent へ返す。UI がある task では frontend 実装後に Storybook レビューループ入力を確認し、人間が別セッションで `story-book-review-loop` を実行する。`implement_lane` は Storybook レビューループの起動、Codex 内蔵ブラウザ操作、コメント収集、コメント解釈、frontend 修正、修正結果判定を扱わない。Storybook レビューループ後に UI 変更がある場合は `designer` へ戻し、plan 内の画面設計成果物を更新させる。観測ログ追加 後に最終検証を行い、実装後ブラウザ確認 後に作業 commit、マージ準備入力を揃える
- `light_change_lane` は `task 枠` と `軽量変更計画` を実行正本にし、`light_change_planner`、`diagrammer`、`backend_implementer`、`frontend_implementer`、`integration_implementer`、`implementation_scenario_tester`、`implementation_unit_tester`、`browser_confirmation`、観点別レビュー agent、`docs_updater` を 文脈 継承なしで直接 起動 する。完了済みサブエージェントは完了結果を集約した後に閉じる
- `ux_maintainance_lane` は `作業準備` と人間が立てた別セッションから返された `browser-use指摘記録` を実行正本にし、`frontend_implementer`、`integration_implementer`、`implementation_unit_tester`、`docs_updater` を 文脈 継承なしで直接 起動 する。Storybook レビューループは人間が別セッションで実行するため、`ux_maintainance_lane` は確定済みの人間指摘が返るまで停止し、Codex 内蔵ブラウザ操作、コメント収集、コメント解釈、frontend 修正入力作成、`frontend_implementer` 再起動、修正結果判定を扱わない
- `refactor_lane` は `task 枠`、`仕様実装優先判断`、`リファクタ範囲確認`、承認済み `implementation-scope` を実行正本にし、`investigator`、`designer`、`backend_implementer`、`frontend_implementer`、`integration_implementer`、`implementation_scenario_tester`、`implementation_unit_tester`、`browser_confirmation`、観点別レビュー agent、`docs_updater` を 文脈 継承なしで直接 起動 する。仕様と実装のどちらを正とするかは人間判断まで固定せず、実装または docs 正本化へ進めない
- agent は代理人であり、職責、職能、ロール、ツール権限 の 担当者 として扱う。`agents/<agent>.toml` の中で「自分は何者か」と 実行境界 を明示する
- skill は作業プロトコルであり、担当ロールが成果物を作る時の判断規約、成果物規約、完了規約、停止規約を持つ。手順、標準 型、参照タイミング一覧、知識範囲一覧は持たない
- Codex agent の人間可読な実行説明は対応する `skills/*/SKILL.md` に置き、紐づけ と `sandbox_mode` は `agents/<agent>.toml` に置き、入力、出力、完了、停止の規約は対応する `skills/*/SKILL.md` に置く
- サンドボックス外で実行してよい command prefix は `.codex/rules/default.rules` の Codex rules に置く
- remote repository を変更する command は agent が実行しない。`push`、tag push、remote branch delete が必要な場合は人間へ返す
- `.agent.md` は使わない

## 形式規約

- agent は人間の代わりに task を実行する担当ロールとして定義する
- agent は自分が何者か、職責、実行境界、入力、出力、停止条件、戻し先を自分の 実行定義内に持つ
- skill は手順書ではなく作業プロトコルとして定義する
- skill は遵守すべき外部規約、判断規約、成果物規約、完了規約、停止規約を持つ
- skill には手順、網羅的な例外分岐、参照タイミング一覧、知識範囲一覧を置かない

## 責務境界

- `implement_lane` は新規実装レーンの進行役として 成果物 DAG、起動入力、人間レビュー、人間向け引き継ぎ、branch 準備、作業 commit、マージ準備入力を扱う
- `implement_lane` は run の 終了処理、停止、戻し 時に、作業 commit とマージ準備入力を判断できる根拠を作る
- `fix_lane` は人間観測、レビュー非通過、検証失敗、修正前調査、修正方針判断、原因箇所シーケンス図、人間修正レビューを読み、修正実行入力、最終検証、レビュー通過根拠を管理する。調査、修正方針判断、実装、テスト、レビューは担当 agent を起動して委任し、プロダクトコードとプロダクトテストは変更しない
- `designer` は 詳細仕様差分、画面設計差分、implementation-scope の task 内成果物を作る
- `exploration_test_lane` は探索計画と探索証跡を読み、バグ一覧、ログ、影響ファイルを集約する。プロダクトコードとプロダクトテストは変更しない
- `light_change_lane` は人間依頼、変更禁止範囲、確認したい結果を `task 枠` に固定し、軽量変更計画、実装、人間確認、テスト修正、レビュー、正本化判断を管理する。プロダクトコードとプロダクトテストは変更しない
- `ux_maintainance_lane` は人間が立てた別セッションから返された Storybook 上の人間指摘を読み、frontend 修正入力、frontend 整理、接続整合証跡、単体テストメンテ、ハーネス通過、docs 正本化判断を管理する。人間の仕様変更指示がない仕様変更、プロダクトコード、プロダクトテスト、docs 正本本文は直接変更しない
- `refactor_lane` は仕様乖離整理起動入力、仕様実装優先判断、構造品質調査起動入力、テスト品質調査起動入力、リファクタ範囲確認、`implementation-scope`、実装実行入力、検証、レビュー、docs 正本化判断を管理する。プロダクトコード、プロダクトテスト、docs 正本本文は直接変更しない
- `light_change_planner` は軽量変更計画だけを作り、プロダクトコード、プロダクトテスト、docs 正本本文を変更しない
- `exploration_test_planner` は探索テストの観測対象、探索観点、テストデータ方針、停止条件だけを固定する
- `investigator` は設計前調査、探索テスト証跡、修正前調査、リファクタレーンの仕様乖離整理、構造品質調査、テスト品質調査のために実画面や観測対象を確認し、観測事実、UI 証跡、ログ、未確認事項を返す。探索テストレーンでは探索証跡だけを担当し、探索範囲を広げる判断をしない。修正レーンでは修正前調査だけを担当し、修正実行入力を作らない。リファクタレーンでは仕様実装優先判断、リファクタ範囲確認、実装範囲を確定しない
- `fix_decider` は修正レーンの修正方針判断と原因箇所シーケンス図を担当し、原因の原因、責務境界、採用する修正方針、禁止する修正、原因箇所の呼び出し順序を分ける。実装 agent が判断し直す余地を残す実装方針、原因未確認の仮説、対症療法は修正実行入力へ進めない
- `browser_confirmation` は実装後ブラウザ確認で、呼び出し元が定義した確認 URL、操作経路、期待値、安全条件に従い、`agent-browser` CLI の `snapshot`、`errors`、必要な `screenshot` とログを残す。期待値の追加、仕様判断、原因推定、修正方針作成は扱わない
- `implement_lane` は承認済み 実行成果物 DAG に従い、実装時調査、実装、テスト、観測ログ追加、最終検証、実装後ブラウザ確認、作業 commit、マージ準備入力を進める
- `implementation_investigator` は承認済み実装範囲 内で実装時の証跡だけを扱う。承認済み実装範囲 外は、今回変更の直接影響を切り分ける観測範囲確認だけを扱う
- `backend_implementer` は 承認済み backend 実装範囲 と、今回変更の直接影響で backend 責務内プロダクトコードに閉じる影響範囲だけを変更する
- `frontend_implementer` は 承認済み frontend 実装範囲、Storybook 人間レビューに必要な story と `fixture` と関連資源、今回変更の直接影響で frontend 責務内プロダクトコードに閉じる影響範囲だけを変更する
- `integration_implementer` は 承認済み 統合境界実装範囲 と、今回変更の直接影響で統合境界責務内プロダクトコードに閉じる影響範囲だけを変更する
- `integration_implementer` は 合意済み frontend 保護 がある場合、承認済み統合境界実装範囲 または限定された影響範囲に閉じる原因箇所だけを変更する
- `observability_implementer` は 完成済み実装成果物 内で、実行時にしか確定しない値、実行後に消える中間状態、消えると原因候補を分離できない分岐理由を残す恒久ログだけを追加する
- `implementation_scenario_tester` は 承認済み詳細仕様差分、承認済み実装範囲、今回のテスト変更が直接壊した検証経路を証明する シナリオテスト と必要最小限の テスト補助 だけを変更する
- `implementation_unit_tester` は 仕様根拠、実装済み責務、承認済み実装範囲、今回のテスト変更が直接壊した検証経路を証明する 単体テスト と必要最小限の テスト補助 だけを変更する
- `docs_updater` は呼び出し元レーンの docs 正本化起動入力と human 承認済み対象範囲が分かった後だけ正本化する
- `implement_lane` は全 implementation 引き継ぎ、最終検証、実装後ブラウザ確認の完了根拠をマージ準備入力へ残す
- `light_change_lane` は軽量変更計画と実装証跡から必要なテスト追従を `implementation_scenario_tester` または `implementation_unit_tester` へ渡し、その後に観点別レビュー agent を起動する。新しい詳細仕様、状態遷移、永続仕様、公開契約、外部連携判断が必要な場合は停止して人間へ返す
- `ux_maintainance_lane` は backend API、DTO、生成物、gateway 境界、frontend 呼び出し、項目値、項目削減、リクエスト削減を確認する。frontend と backend の接続に必要な統合メンテは `integration_implementer` へ渡し、単体テストメンテは `implementation_unit_tester` へ渡す。backend プロダクトコード変更が必要な場合は停止して人間へ返す
- `merge_lane` は local merge と completed 移動だけを扱う。conflict 解消が仕様判断、設計変更、レーン外の再実装を必要とする場合は停止して人間へ返す
- 観点別 レビュー agent は挙動正しさ、契約・互換性、権限・信頼境界、状態・データ不変条件、責務境界のいずれか 1 つだけを扱い、`reviewback.<観点>.yaml` を作成、追記、解決更新、削除する
- 観点別 レビュー agent は広い ハーネス 再実行を担当せず、呼び出し元レーンから渡された検証証跡をレビュー入力として扱う
- 観点別 レビュー agent は 失敗 または 停止 の場合も `reviewback.<観点>.yaml` に結果、根拠、未解決指摘を記録する
- `reviewback.<観点>.yaml` はゲート判断用レビュー成果物とし、work_history 側に観点別の非通過 YAML は作らない
- `implement_lane`、`fix_lane`、`exploration_test_lane`、`light_change_lane`、`ux_maintainance_lane`、`refactor_lane`、`light_change_planner`、`exploration_test_planner`、`designer`、`diagrammer`、`investigator`、`browser_confirmation`、`docs_updater`、レビュー agent は プロダクトコード と プロダクトテスト を変更しない
- プロダクトコード は `backend_implementer`、`frontend_implementer`、`integration_implementer`、`observability_implementer` だけが 承認済み実装範囲 または担当 agent の責務内に閉じる限定された影響範囲で変更できる
- シナリオテスト は `implementation_scenario_tester` だけが 承認済み実装範囲 または今回のテスト変更が直接壊した担当テスト成果物の影響範囲で変更できる
- 単体テスト は `implementation_unit_tester` だけが 承認済み実装範囲 または今回のテスト変更が直接壊した担当テスト成果物の影響範囲で変更できる
- implementation レーン は docs 正本、`.codex/` 作業流れ 文書、agent 実行定義、ツール権限 を変更しない


## task 種別レーン

- task run は task type ごとの レーン として扱い、各 レーン が自分の必須 成果物 DAG を持つ
- live レーン は `implement_lane`、`fix_lane`、`exploration_test_lane`、`light_change_lane`、`ux_maintainance_lane`、`refactor_lane`、`merge_lane` にする
- `implement_lane` は新規実装と機能拡張だけを扱う
- `fix_lane` は人間が確認した不具合、レビュー非通過、検証失敗の恒久修正だけを扱う
- `exploration_test_lane` は探索計画、テストデータ、探索証跡、バグ一覧、ログ、影響ファイル、実装証跡、回帰テスト証跡を扱う
- `light_change_lane` は既存仕様の意味を大きく広げない軽い backend / frontend / integration 変更だけを扱う
- `ux_maintainance_lane` は Storybook 上の人間指摘を起点に、frontend 表示修正、frontend 整理、接続整合証跡、統合メンテ、単体テストメンテ、ハーネス通過、画面設計正本反映を扱う
- `ux_maintainance_lane` は人間の仕様変更指示がある場合だけ、画面表示以外の仕様変更と詳細仕様正本反映を扱う
- `ux_maintainance_lane` は人間の仕様変更指示がない状態で画面表示以外の仕様変更または詳細仕様正本反映が必要な場合に停止する
- `refactor_lane` は仕様乖離整理起動入力、人間による仕様実装優先判断、構造品質調査起動入力、テスト品質調査起動入力、リファクタ範囲確認、`implementation-scope`、承認済み範囲の修正、テスト、レビュー、docs 正本化判断を扱う
- `refactor_lane` は仕様と実装のどちらを正にするか人間判断がない状態で、実装修正または docs 正本化へ進む場合に停止する
- `merge_lane` は active plan ごとの local merge、conflict 解消、merge 後検証、completed 移動だけを扱う
- 各 レーン は task 内成果物 DAG を持ち、順序は phase 名ではなく `依存対象` と対象 skill の完了規約で固定する
- agent は レーン そのものではなく、成果物 を作る実行主体として扱う
- 実装、修正、探索テスト、軽量変更、UX 保守、リファクタの各レーンは branch 準備、作業 commit、マージ準備入力までを必須にする
- `merge_lane` 以外のレーンは 作業計画 folder の completed 移動を扱わない


## 実装レーン成果物DAG

新規実装レーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `task 枠` | `implement_lane` | `[]` | なし |
| `詳細仕様差分` | `designer` | `task 枠` | `designer` |
| `画面設計差分` | `designer` | `詳細仕様差分` | `designer` |
| `設計差分図` | `diagrammer` | `詳細仕様差分`, `画面設計差分?` | `diagrammer` |
| `人間設計レビュー` | human | `詳細仕様差分`, `画面設計差分?`, `設計差分図` | human |
| `実装範囲` | `designer` | `人間設計レビュー` | `designer` |
| `実装引き継ぎ入力` | `implement_lane` | `実装範囲` | なし |
| `frontend 実装` | `frontend_implementer` / `implement-frontend` | `実装引き継ぎ入力` | `frontend_implementer` |
| `Storybookレビューループ入力確認` | `implement_lane` | `frontend 実装` | なし |
| `Storybookレビューループ完了証跡` | 人間が立てた別セッション / `story-book-review-loop` | `Storybookレビューループ入力確認` | なし |
| `frontend 実装後人間レビュー` | `implement_lane` | `Storybookレビューループ完了証跡` | なし |
| `Storybook後画面設計差分整合` | `designer` | `frontend 実装後人間レビュー` | `designer` |
| `合意済みfrontend保護` | `implement_lane` | `frontend 実装後人間レビュー`, `Storybook後画面設計差分整合?` | なし |
| `backend 実装` | `backend_implementer` / `implement-backend` | `実装引き継ぎ入力`, `合意済みfrontend保護?` | `backend_implementer` |
| `統合境界実装` | `integration_implementer` / `implement-integration` | `backend 実装`, `合意済みfrontend保護?` | `integration_implementer` |
| `シナリオテスト` | `implementation_scenario_tester` | `backend 実装?`, `合意済みfrontend保護?`, `統合境界実装?` | `implementation_scenario_tester` |
| `単体テスト` | `implementation_unit_tester` | `backend 実装?`, `合意済みfrontend保護?`, `統合境界実装?` | `implementation_unit_tester` |
| `観測ログ追加` | `observability_implementer` / `observability-implementer` | `backend 実装?`, `frontend 実装?`, `合意済みfrontend保護?`, `統合境界実装?`, `シナリオテスト?`, `単体テスト?` | `observability_implementer` |
| `最終検証` | `implement_lane` | `観測ログ追加` | なし |
| `実装後ブラウザ確認` | `browser_confirmation` | `最終検証` | `browser_confirmation` |
| `正本化判断` | `implement_lane` | `最終検証`, `実装後ブラウザ確認` | `docs_updater?` |
| `詳細仕様正本反映` | `docs_updater` | `正本化判断` | `docs_updater?` |
| `branch 準備` | `implement_lane` | `task 枠` | なし |
| `作業 commit` | `implement_lane` | 全完了または停止済み 成果物 | なし |
| `マージ準備入力` | `implement_lane` | `作業 commit` | `merge_lane` |

## 修正レーン成果物DAG

修正レーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `人間観測記録` | `fix_lane` | `task 枠` | なし |
| `修正前調査` | `investigator` | `人間観測記録` | `investigator` |
| `修正方針判断` | `fix_decider` | `人間観測記録`, `修正前調査` | `fix_decider` |
| `原因箇所シーケンス図` | `fix_decider` | `人間観測記録`, `修正前調査`, `修正方針判断` | なし |
| `人間修正レビュー` | human | `修正方針判断`, `原因箇所シーケンス図` | human |
| `修正実行入力` | `fix_lane` | `人間観測記録`, `修正前調査`, `修正方針判断`, `原因箇所シーケンス図`, `人間修正レビュー` | なし |
| `実装証跡` | 実装種別別 agent / `implement-backend` または `implement-frontend` または `implement-integration` | `修正実行入力` | `backend_implementer` または `frontend_implementer` または `integration_implementer` |
| `回帰テスト証跡` | `implementation_scenario_tester` または `implementation_unit_tester` | `実装証跡` | `implementation_scenario_tester` または `implementation_unit_tester` |
| `最終検証` | `fix_lane` | `実装証跡`, `回帰テスト証跡?` | なし |
| `実装後ブラウザ確認` | `browser_confirmation` | `最終検証` | `browser_confirmation` |
| `レビュー通過根拠` | `fix_lane` | `人間観測記録`, `修正前調査`, `修正方針判断`, `原因箇所シーケンス図`, `人間修正レビュー`, `修正実行入力`, `実装証跡`, `回帰テスト証跡?`, `最終検証`, `実装後ブラウザ確認` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `branch 準備` | `fix_lane` | `人間観測記録` | なし |
| `作業 commit` | `fix_lane` | `レビュー通過根拠` | なし |
| `マージ準備入力` | `fix_lane` | `作業 commit` | `merge_lane` |

## 探索テストレーン成果物DAG

探索テストレーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `探索計画` | `exploration_test_planner` | `task 枠` | `exploration_test_planner` |
| `テストデータ` | `exploration_test_lane` | `探索計画` | なし |
| `探索証跡` | `investigator` | `探索計画`, `テストデータ` | `investigator` |
| `バグ一覧とログ、影響ファイル` | `exploration_test_lane` | `探索証跡` | なし |
| `実装証跡` | 実装種別別 agent | `バグ一覧とログ、影響ファイル` | `backend_implementer` または `frontend_implementer` または `integration_implementer` |
| `回帰テスト証跡` | `implementation_scenario_tester` または `implementation_unit_tester` | `実装証跡` | `implementation_scenario_tester` または `implementation_unit_tester` |
| `レビュー通過根拠` | `exploration_test_lane` | `探索計画`, `探索証跡`, `バグ一覧とログ、影響ファイル`, `実装証跡?`, `回帰テスト証跡?` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `branch 準備` | `exploration_test_lane` | `探索計画` | なし |
| `作業 commit` | `exploration_test_lane` | `レビュー通過根拠` | なし |
| `マージ準備入力` | `exploration_test_lane` | `作業 commit` | `merge_lane` |

## 軽量変更レーン成果物DAG

軽量変更レーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `task 枠` | `light_change_lane` | `[]` | なし |
| `軽量変更計画` | `light_change_planner` | `task 枠` | `light_change_planner` |
| `設計差分図` | `diagrammer` | `軽量変更計画` | `diagrammer` |
| `実装証跡` | 実装種別別 agent / `implement-backend` または `implement-frontend` または `implement-integration` | `軽量変更計画`, `設計差分図` | `backend_implementer` または `frontend_implementer` または `integration_implementer` |
| `人間確認` | human | `実装証跡` | human |
| `テスト修正証跡` | `implementation_scenario_tester` または `implementation_unit_tester` | `実装証跡`, `人間確認?` | `implementation_scenario_tester` または `implementation_unit_tester` |
| `実装後ブラウザ確認` | `browser_confirmation` | `実装証跡`, `人間確認?`, `テスト修正証跡?` | `browser_confirmation` |
| `レビュー通過根拠` | `light_change_lane` | `軽量変更計画`, `実装証跡`, `人間確認?`, `テスト修正証跡?`, `実装後ブラウザ確認` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `正本化判断` | `light_change_lane` | `レビュー通過根拠` | `docs_updater?` |
| `詳細仕様正本反映` | `docs_updater` | `正本化判断` | `docs_updater?` |
| `branch 準備` | `light_change_lane` | `task 枠` | なし |
| `作業 commit` | `light_change_lane` | `詳細仕様正本反映?`, `レビュー通過根拠` | なし |
| `マージ準備入力` | `light_change_lane` | `作業 commit` | `merge_lane` |

## UX 保守レーン成果物DAG

UX 保守レーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `作業準備` | `ux_maintainance_lane` | `[]` | なし |
| `Storybookレビューループ完了証跡` | 人間が立てた別セッション / `story-book-review-loop` | `作業準備` | なし |
| `browser-use指摘記録` | human | `Storybookレビューループ完了証跡` | human |
| `frontend修正入力` | `ux_maintainance_lane` | `browser-use指摘記録`, `Storybookレビューループ完了証跡` | なし |
| `frontend修正証跡` | `frontend_implementer` | `frontend修正入力` | `frontend_implementer` |
| `frontend整理証跡` | `frontend_implementer` | `frontend修正証跡` | `frontend_implementer` |
| `接続整合証跡` | `ux_maintainance_lane` または `integration_implementer` | `frontend整理証跡` | `integration_implementer?` |
| `単体テストメンテ証跡` | `implementation_unit_tester` | `frontend整理証跡`, `接続整合証跡` | `implementation_unit_tester?` |
| `docs正本化判断` | `ux_maintainance_lane` | `接続整合証跡`, `単体テストメンテ証跡` | `docs_updater?` |
| `画面設計正本反映` | `docs_updater` | `docs正本化判断` | `docs_updater?` |
| `詳細仕様正本反映` | `docs_updater` | `docs正本化判断` | `docs_updater?` |
| `ハーネス通過` | `ux_maintainance_lane` | `frontend整理証跡`, `接続整合証跡`, `単体テストメンテ証跡`, `画面設計正本反映?`, `詳細仕様正本反映?` | なし |
| `作業 commit` | `ux_maintainance_lane` | `ハーネス通過` | なし |
| `マージ準備入力` | `ux_maintainance_lane` | `作業 commit` | `merge_lane` |

## リファクタレーン成果物DAG

リファクタレーンの成果物DAGは次を標準形にする。
順序は `依存対象` と対象 skill の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `task 枠` | `refactor_lane` | `[]` | なし |
| `branch 準備` | `refactor_lane` | `task 枠` | なし |
| `仕様乖離整理` | `investigator` / `investigate` | `task 枠` | `investigator` |
| `仕様実装優先判断` | human | `仕様乖離整理` | human |
| `構造品質調査` | `investigator` / `investigate` | `仕様実装優先判断` | `investigator` |
| `テスト品質調査` | `investigator` / `investigate` | `仕様実装優先判断` | `investigator` |
| `リファクタ範囲確認` | human | `構造品質調査`, `テスト品質調査` | human |
| `実装範囲` | `designer` / `implementation-scope` | `仕様実装優先判断`, `リファクタ範囲確認` | `designer` |
| `実装引き継ぎ入力` | `refactor_lane` | `実装範囲` | なし |
| `backend リファクタ` | `backend_implementer` / `implement-backend` | `実装引き継ぎ入力` | `backend_implementer?` |
| `frontend リファクタ` | `frontend_implementer` / `implement-frontend` | `実装引き継ぎ入力` | `frontend_implementer?` |
| `統合境界リファクタ` | `integration_implementer` / `implement-integration` | `backend リファクタ?`, `frontend リファクタ?`, `実装引き継ぎ入力` | `integration_implementer?` |
| `シナリオテスト` | `implementation_scenario_tester` | `backend リファクタ?`, `frontend リファクタ?`, `統合境界リファクタ?` | `implementation_scenario_tester?` |
| `単体テスト` | `implementation_unit_tester` | `backend リファクタ?`, `frontend リファクタ?`, `統合境界リファクタ?` | `implementation_unit_tester?` |
| `最終検証` | `refactor_lane` | `backend リファクタ?`, `frontend リファクタ?`, `統合境界リファクタ?`, `シナリオテスト?`, `単体テスト?` | なし |
| `実装後ブラウザ確認` | `browser_confirmation` | `最終検証` | `browser_confirmation?` |
| `レビュー通過根拠` | `refactor_lane` | `最終検証`, `実装後ブラウザ確認?` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `docs正本化判断` | `refactor_lane` | `仕様実装優先判断`, `レビュー通過根拠` | `docs_updater?` |
| `詳細仕様正本反映` | `docs_updater` | `docs正本化判断` | `docs_updater?` |
| `作業 commit` | `refactor_lane` | `レビュー通過根拠`, `詳細仕様正本反映?` | なし |
| `マージ準備入力` | `refactor_lane` | `作業 commit` | `merge_lane` |

## マージレーン成果物DAG

マージレーンの成果物DAGは次を標準形にする。
順序は `依存対象` と `merge-lane` の完了規約で固定し、phase 名では固定しない。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `マージ準備確認` | `merge_lane` | `マージ準備入力` | なし |
| `local merge` | `merge_lane` | `マージ準備確認` | なし |
| `conflict 解消` | `merge_lane` | `local merge` | なし |
| `merge 後検証` | `merge_lane` | `local merge`, `conflict 解消?` | なし |
| `completed 移動` | `merge_lane` | `merge 後検証` | なし |
| `merge 結果 commit` | `merge_lane` | `completed 移動` | なし |

## 実行計画 folder

- 新規 task は `docs/exec-plans/active/<task-id>/` に folder として作る
- `plan.md` は索引、状態、HITL、検証、終了処理 だけを書く
- 各 skill の資料は同じ folder の skill 名つき file に分ける
- AI は最初に `plan.md` だけ読み、必要な資料だけ追加で読む
- 各レーンの作業完了後も folder は `docs/exec-plans/active/<task-id>/` に残す
- `merge_lane` が local merge と merge 後検証を完了した時だけ、folder ごと `docs/exec-plans/completed/<task-id>/` へ移す

## Docs 正本化

- docs 正本化は呼び出し元レーンの docs 正本化起動入力と human 承認済み対象範囲が分かった後に扱う
- docs 正本化は Codex 側だけで扱う
- human 承認済みの 成果物 だけ `docs_updater` が `updating-docs` を参照して正本へ反映する
- `detail-spec-diff.md`、`screen-design-diff.<screen-id>.md`、実装結果のいずれかに仕様変更または仕様追加が少しでも含まれる場合は、`implement_lane` が `正本化判断` を必ず記録する
- `refactor_lane` は `実装が正` と人間判断された仕様乖離がある場合、docs 正本化判断と docs 正本化起動入力を必ず記録する
- 仕様変更または仕様追加が human 承認済みの恒久仕様である場合は、`docs_updater` が `詳細仕様正本反映` を必ず完了または停止理由付きで返す
- task 内 詳細仕様差分、画面設計差分、ブラウザ確認結果は task folder に置く
- UI の確認は、承認済み画面設計差分と実画面確認結果で扱う
- UI の細かな visual polish は実装後の実物確認で差分を扱う
- `implementation-scope` は 引き継ぎ 履歴であり docs 正本へ昇格しない
- `detail-specs` は詳細仕様正本とし、`detail-spec-diff.md`、`screen-design-diff.<screen-id>.md`、実装結果から human 承認済みの恒久仕様だけを製本する

## 非 live 扱い

- 旧 `design` は詳細仕様差分、画面設計差分、`implementation-scope` に再整理した
- 旧 flat file 形式の exec-plan は legacy とし、新規 task では使わない
- 設計前調査では UI check 専用 agent を置かず、設計前の画面設計根拠を `investigator` が扱う
- 作業前の影響範囲、実行計画、検証方法の確認は `AGENTS.md` の入口規約に集約する
- Codex 側の人間可読な 実行定義 説明は skill へ集約し、`.codex/agents/*.agent.md` は持たない
- `.codex/workflow.md` は補助図として残し、live 作業流れ の正本にはしない
- 旧 skill / agent の退避物は live 作業流れ に残さない

## 作業計画

- 非自明な変更は `docs/exec-plans/active/<task-id>/` に置く
- マージ完了後は `merge_lane` が `docs/exec-plans/completed/<task-id>/` へ移す
- completed plan は履歴として残す
