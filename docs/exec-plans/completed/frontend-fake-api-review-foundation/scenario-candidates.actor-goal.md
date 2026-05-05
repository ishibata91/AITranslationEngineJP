# Scenario Candidates: frontend-fake-api-review-foundation / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `FFARF`

## Generator Scope

- `viewpoint`: actor-goal
- `included_sources`:
  - `./task-frame.md`
  - `../../../../tasks/usecases/frontend-fake-api-review-foundation.yaml`
  - `../../../../docs/architecture.md`
  - `../../../../docs/coding-guidelines-frontend.md`
  - `../../../../tmp/code-map/index.json`
- `excluded_sources`:
  - プロダクトコード
  - プロダクトテスト
  - docs 正本変更
  - `.codex` 変更
- `generation_notes`:
  - actor の目的、開始操作、成功判定から候補を作る。
  - 採否、統合、最終シナリオ表は扱わない。
  - fakeAPI は provider 選択肢ではなく、レビュー起動時の DI による差し替えとして扱う。

## Candidate Scenarios

### CAND-FFARF-001 レビュー起動で fakeAPI を選べる

- `根拠要件`: `task-frame.md` の目的、完了条件「起動モードで フロントエンドの API 接続先を fakeAPI に切り替えられる」。`frontend-fake-api-review-foundation.yaml` の `goal` と `manual_check_steps`。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-FFARF-001`
- `actor`: フロントエンドレビュー実行者
- `goal`: バックエンド と Wails バインディングを起動せず、実フロントエンドを レビュー用状態で開く。
- `trigger`: fakeAPI 起動モードを指定してアプリを起動する。
- `期待結果`: フロントエンドの API 接続先が fakeAPI に切り替わり、実画面が起動する。
- `観測点`: `agent-browser` で実画面を開き、fakeAPI 由来の状態が表示される。
- `関連詳細要求タイプ`: レビュー起動、DI による差し替え、手動確認入口
- `採用判断材料`: fakeAPI 起動基盤の主目的を確認する候補として使える。
- `競合注意`: 起動モード名、起動 command、確認対象 URL は designer の統合時に別候補と合わせる必要がある。

### CAND-FFARF-002 Wails バインディング非依存で画面を確認できる

- `根拠要件`: `task-frame.md` の設計前提「本番ゲートウェイは Wails バインディング adapter として `フロントエンド/src/controller/wails/` に閉じ込める」。`architecture.md` の Frontend Bootstrap、Gateway、Wails 境界。`coding-guidelines-frontend.md` の Wails 境界。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-FFARF-002`
- `actor`: フロントエンド実装者
- `goal`: View、ScreenController、Frontend UseCase を Wails バインディング へ直接依存させずに レビュー状態を確認する。
- `trigger`: fakeAPI 起動モードで、画面が GatewayContract 経由の fake 実装を受け取る。
- `期待結果`: 生成済み `wailsjs` がなくても レビュー用画面確認を進められる。
- `観測点`: fakeAPI 起動時の画面表示と、Wails adapter を使わない局所テスト結果で確認する。
- `関連詳細要求タイプ`: DI 境界、Wails バインディング非依存、責務境界
- `採用判断材料`: fakeAPI が 本番ゲートウェイ の横流しではなく、DI による差し替えであることを確認する候補として使える。
- `競合注意`: Wails adapter の完全未使用をどの検証で証明するかは、contract または responsibility-boundary 観点と競合する可能性がある。

### CAND-FFARF-003 画面固有のモックデータ で レビューできる

- `根拠要件`: `task-frame.md` の完了条件「画面固有のモックデータを ユースケース task 側で追加できる」。`frontend-fake-api-review-foundation.yaml` の `outputs`。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-FFARF-003`
- `actor`: ユースケース task 作成者
- `goal`: 画面固有の レビューデータ を追加し、実画面で表示差分を確認できるようにする。
- `trigger`: ユースケース task 側で画面固有のモックデータを fakeAPI 基盤へ渡す。
- `期待結果`: 追加した モックデータが該当画面だけの レビュー状態として表示される。
- `観測点`: 実画面の表示内容が モックデータと一致し、別画面や 本番初期状態 へ広がらない。
- `関連詳細要求タイプ`: モックデータ 拡張、画面固有 レビュー、混入防止
- `採用判断材料`: 後続 ユースケース task が fakeAPI 基盤を再利用できるかを確認する候補として使える。
- `競合注意`: モックデータの置き場所と命名規則は implementation-scope で確定する必要がある。

### CAND-FFARF-004 空状態を レビューできる

- `根拠要件`: `task-frame.md` と `frontend-fake-api-review-foundation.yaml` の完了条件「空状態を fakeAPI で再現できる」。`coding-guidelines-frontend.md` の UX 一般規約。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-FFARF-004`
- `actor`: フロントエンドレビュー実行者
- `goal`: データがない画面で、ユーザーが次に何をすべきか判断できる表示を確認する。
- `trigger`: fakeAPI の状態パターンを空状態へ切り替える。
- `期待結果`: 実画面が空状態専用の表示を出し、成功状態や失敗状態と混同しない。
- `観測点`: `agent-browser` で空状態の表示、主要操作、メッセージを確認できる。
- `関連詳細要求タイプ`: 状態パターン、空状態、UI 確認
- `採用判断材料`: 状態パターン群の基本候補として使える。
- `競合注意`: 個別画面の空状態文言は画面固有 task 側の判断候補になる。

### CAND-FFARF-005 読み込み中と進行中を レビューできる

- `根拠要件`: `task-frame.md` の完了条件「読み込み中、進行中状態を fakeAPI で再現できる」。`coding-guidelines-frontend.md` の UX 一般規約。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-FFARF-005`
- `actor`: フロントエンドレビュー実行者
- `goal`: 待機が必要な状態で、ユーザーが作業中か取得中かを区別できる表示を確認する。
- `trigger`: fakeAPI の状態パターンを読み込み中または進行中へ切り替える。
- `期待結果`: 実画面が読み込み中と進行中を区別して表示する。
- `観測点`: `agent-browser` で待機表示、進捗表示、操作可否を確認できる。
- `関連詳細要求タイプ`: 状態パターン、読み込み中、進行中
- `採用判断材料`: レビュー実行者の成功判定が画面表示に直結する候補として使える。
- `競合注意`: 読み込み中と進行中を 1 候補に束ねるか分けるかは designer の統合判断に残す。

### CAND-FFARF-006 成功状態と失敗状態を レビューできる

- `根拠要件`: `task-frame.md` の完了条件「成功状態、失敗状態を fakeAPI で再現できる」。`frontend-fake-api-review-foundation.yaml` の `manual_check_steps`。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-FFARF-006`
- `actor`: フロントエンドレビュー実行者
- `goal`: 正常完了時と失敗時の画面体験を同じ実フロントエンド 上で比較する。
- `trigger`: fakeAPI の状態パターンを成功状態または失敗状態へ切り替える。
- `期待結果`: 成功状態では完了結果が表示され、失敗状態では 利用者向けメッセージ が表示される。
- `観測点`: `agent-browser` で成功表示、失敗表示、再試行または次操作の導線を確認できる。
- `関連詳細要求タイプ`: 状態パターン、成功状態、失敗状態
- `採用判断材料`: レビューに必要な代表的正常系と代表的異常系を actor-goal から押さえる候補として使える。
- `競合注意`: 失敗原因の網羅は failure 観点に残す。

### CAND-FFARF-007 設定不足状態を レビューできる

- `根拠要件`: `task-frame.md` と `frontend-fake-api-review-foundation.yaml` の完了条件「設定不足状態を fakeAPI で再現できる」。`coding-guidelines-frontend.md` の 利用者向けメッセージ 規約。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-FFARF-007`
- `actor`: フロントエンドレビュー実行者
- `goal`: 設定不足の画面で、ユーザーが不足項目と次操作を判断できる表示を確認する。
- `trigger`: fakeAPI の状態パターンを設定不足状態へ切り替える。
- `期待結果`: 実画面が不足状態を失敗状態と分けて表示し、内部診断を露出しない。
- `観測点`: `agent-browser` で不足項目の表示、操作導線、内部診断の非表示を確認できる。
- `関連詳細要求タイプ`: 状態パターン、設定不足、表示責務
- `採用判断材料`: 設定不足が product error と混ざらないか確認する候補として使える。
- `競合注意`: 設定不足の対象項目は画面別要件と統合する必要がある。

### CAND-FFARF-008 本番起動で fakeAPI が選ばれない

- `根拠要件`: `task-frame.md` の完了条件「本番起動では fakeAPI が選ばれない」「fakeAPI と モックデータが本番 API、永続化、本番初期状態に混入しない」。`frontend-fake-api-review-foundation.yaml` の `outputs` と `manual_check_steps`。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-FFARF-008`
- `actor`: release 前確認者
- `goal`: レビュー用 fakeAPI と モックデータが本番起動へ混入していないことを確認する。
- `trigger`: 本番起動相当の起動モードでアプリを起動する。
- `期待結果`: 本番ゲートウェイ が選ばれ、fakeAPI と モックデータが API 接続、永続化、初期状態 に現れない。
- `観測点`: 実画面と局所テストで、fakeAPI 選択なし、モックデータ 露出なし、本番初期状態 汚染なしを確認できる。
- `関連詳細要求タイプ`: 本番混入防止、本番 wiring、永続化非混入
- `採用判断材料`: レビュー基盤の安全側の成功体験を actor-goal から確認する候補として使える。
- `競合注意`: 本番起動相当の判定方法は実装範囲で固定する必要がある。

### CAND-FFARF-009 fakeAPI 基盤の局所テスト結果を レビュー根拠にできる

- `根拠要件`: `task-frame.md` の完了条件「fakeAPI 起動モードが壊れていないことを局所テストで確認できる」。`frontend-fake-api-review-foundation.yaml` の 完了条件「coverage harness では fakeAPI 基盤を数値判定の例外として扱い、例外理由と局所テスト結果を記録できる」。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-FFARF-009`
- `actor`: implement_lane 進行役
- `goal`: fakeAPI 基盤の完了判断に、広い coverage 数値ではなく局所テスト結果を使えるようにする。
- `trigger`: fakeAPI 起動モードの局所テストを実行する。
- `期待結果`: fakeAPI 起動モード、DI による差し替え、本番非選択の確認結果が記録される。
- `観測点`: 局所テストの実行結果と、coverage harness 例外理由の記録で確認できる。
- `関連詳細要求タイプ`: 局所検証、coverage 例外、完了根拠
- `採用判断材料`: actor-goal 観点では、進行役が完了根拠を判断できる候補として使える。
- `競合注意`: coverage 例外の扱いは最終検証や report 側の規約と統合する必要がある。

## Open Notes

- `人間判断候補`:
  - fakeAPI 起動モード名、起動 command、確認 URL。
  - 状態パターンの切り替え UI または切り替え手段。
  - 画面固有モックデータの配置単位と命名規則。
  - 本番起動相当の判定方法。
- `統合候補`:
  - CAND-FFARF-004、CAND-FFARF-005、CAND-FFARF-006、CAND-FFARF-007 は状態パターン確認として統合できる可能性がある。
  - CAND-FFARF-001 と CAND-FFARF-002 は fakeAPI 起動と DI による差し替え確認として統合できる可能性がある。
  - CAND-FFARF-008 と CAND-FFARF-009 は本番混入防止の検証根拠として統合できる可能性がある。
- `不採用候補`:
  - UI 文言や個別画面導線だけを決める候補は、画面固有 task 側へ回す。
  - 失敗原因の網羅だけを目的にする候補は、failure 観点へ回す。
  - 状態遷移の禁止条件だけを目的にする候補は、state-transition 観点へ回す。

## Completion Evidence

- `task_id`: `frontend-fake-api-review-foundation`
- `artifact_path`: `docs/exec-plans/completed/frontend-fake-api-review-foundation/scenario-candidates.actor-goal.md`
- `viewpoint`: `actor-goal`
- `candidate_count`: 9
- `remaining_risk`: 起動 command、状態パターン 切り替え手段、モックデータ 配置、本番起動相当の判定方法は AI だけで確定していない。
