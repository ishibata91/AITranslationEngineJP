# Scenario Candidates: frontend-fake-api-review-foundation / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `FFARF-LC`
- `candidate_count`: 7

## Generator Scope

- `viewpoint`: `lifecycle`
- `included_sources`: `task-frame.md`, `tasks/usecases/frontend-fake-api-review-foundation.yaml`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `tmp/code-map/index.json`
- `excluded_sources`: プロダクトコード, プロダクトテスト, docs 正本更新, `.codex` 変更, 最終シナリオ表, 採否決定, 統合判断, implementation-scope
- `generation_notes`: 起動、接続先切替、状態パターン選択、画面確認、終了、通常起動復帰だけを候補化する。

## Candidate Scenarios

### CAND-FFARF-LC-001 レビュー起動モードで fakeAPI 接続を開始する

- `根拠要件`: `task-frame.md:5-6`, `task-frame.md:17`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:6`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:20`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-FFARF-LC-001`
- `lifecycle phase`: 起動
- `actor`: フロントエンドレビューer
- `trigger`: レビューer が fakeAPI 起動モードで フロントエンド を起動する。
- `開始条件`: 基盤データ管理が成立し、レビュー用起動モードが指定されている。
- `期待結果`: フロントエンド は Wails バインディング と バックエンド に依存せず、fakeAPI を接続先として起動する。
- `観測点`: 起動後の実画面で、バックエンド未起動でも レビュー用状態の表示確認へ進める。
- `関連詳細要求タイプ`: `success_requirement`, `state_requirement`, `testability_requirement`
- `採用判断材料`: fakeAPI レビューの開始条件として、起動モード指定と接続先確定を分けて扱える。
- `競合注意`: external-integration 観点が扱う Wails / バックエンド 非依存条件と重複する可能性がある。

### CAND-FFARF-LC-002 レビュー起動時だけ DI による差し替えで fakeAPI へ切り替える

- `根拠要件`: `task-frame.md:27-30`, `docs/architecture.md:17`, `docs/architecture.md:53-61`, `tmp/code-map/index.json:3873-3884`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-FFARF-LC-002`
- `lifecycle phase`: 接続先切替
- `actor`: フロントエンド bootstrap
- `trigger`: レビュー起動モードの起動処理が API 接続先を決める。
- `開始条件`: フロントエンドの composition root が `フロントエンド/src/main.ts` にあり、手動 DI で ゲートウェイ を注入する。
- `期待結果`: fakeAPI は provider 選択肢ではなく、レビュー起動時の DI による差し替えとして選ばれる。
- `観測点`: View、ScreenController、Frontend UseCase から 生成済み `wailsjs` への直接参照が不要なまま、画面が起動する。
- `関連詳細要求タイプ`: `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `採用判断材料`: 本番 provider 設定と レビュー起動モードを分離する候補として扱える。
- `競合注意`: responsibility-boundary 観点が扱う composition root 配置と競合または統合する可能性がある。

### CAND-FFARF-LC-003 状態パターンを選び レビュー用 モックデータを読み込む

- `根拠要件`: `task-frame.md:19-20`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:14-16`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:22-23`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-FFARF-LC-003`
- `lifecycle phase`: 状態パターン選択
- `actor`: フロントエンドレビューer
- `trigger`: レビューer が レビュー対象画面の状態パターンを選ぶ。
- `開始条件`: レビュー起動モードで fakeAPI が有効になり、状態パターン一覧が利用できる。
- `期待結果`: 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態のいずれかが モックデータ から再現される。
- `観測点`: 選択した 状態パターン に対応する画面状態が、実フロントの表示として確認できる。
- `関連詳細要求タイプ`: `alternative_success_requirement`, `data_requirement`, `testability_requirement`
- `採用判断材料`: 画面固有モックデータの追加単位と、共通 状態パターン一覧の関係を確認する候補として使える。
- `競合注意`: state-transition 観点が扱う状態名、failure 観点が扱う失敗状態と重複する可能性がある。

### CAND-FFARF-LC-004 実画面で各 状態パターンの表示を確認する

- `根拠要件`: `task-frame.md:21`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:24`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:35-36`, `docs/coding-guidelines-frontend.md:48-55`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-FFARF-LC-004`
- `lifecycle phase`: 画面確認
- `actor`: フロントエンドレビューer
- `trigger`: レビューer が `agent-browser` で実フロントを開き、選択済み 状態パターンの表示を確認する。
- `開始条件`: fakeAPI 起動モードで フロントエンド が起動し、対象 状態パターンが選択されている。
- `期待結果`: 読み込み中、空状態、エラー、完了、設定不足などの画面状態が区別できる。
- `観測点`: 実画面の文言、状態表示、主要操作の見え方が 状態パターンごとに確認できる。
- `関連詳細要求タイプ`: `success_requirement`, `observability_requirement`, `testability_requirement`
- `採用判断材料`: UI 人間操作 E2E の入口候補として、状態パターンごとの表示確認へ接続できる。
- `競合注意`: 表示文言や配置の採否は画面固有 task 側に残す。

### CAND-FFARF-LC-005 状態パターンを切り替えて同じ起動セッション内で再確認する

- `根拠要件`: `task-frame.md:19-21`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:32-36`, `docs/architecture.md:101-108`, `docs/architecture.md:119-123`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-FFARF-LC-005`
- `lifecycle phase`: 更新
- `actor`: フロントエンドレビューer
- `trigger`: レビューer が一つの 状態パターン確認後、別の 状態パターンを選び直す。
- `開始条件`: 同じ レビュー起動セッションで、ScreenController と Store が画面状態を保持している。
- `期待結果`: 選び直した 状態パターンの状態へ表示が更新され、前の 状態パターンの画面状態が混入しない。
- `観測点`: Store 由来の 表示モデル が更新され、画面上の状態表示が選び直した 状態パターン と一致する。
- `関連詳細要求タイプ`: `consistency_requirement`, `state_requirement`, `testability_requirement`
- `採用判断材料`: 起動し直し不要の 状態パターン確認を要求に含めるか判断する材料になる。
- `競合注意`: 状態パターン切替方法が未決の場合、CAND-FFARF-LC-003 と統合される可能性がある。

### CAND-FFARF-LC-006 レビュー起動を終了して モックデータを残さない

- `根拠要件`: `task-frame.md:22`, `task-frame.md:32-37`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:25`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-FFARF-LC-006`
- `lifecycle phase`: 終了
- `actor`: フロントエンドレビューer
- `trigger`: レビューer が fakeAPI レビューセッションを終了する。
- `開始条件`: fakeAPI 起動モードで 状態パターン確認が完了している。
- `期待結果`: fakeAPI と モックデータ は本番 API、永続化、本番初期状態に混入しない。
- `観測点`: 終了後に永続化された本番データや 初期状態 へ レビュー用 モックデータが残らない。
- `関連詳細要求タイプ`: `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `採用判断材料`: レビュー用データの寿命を、画面確認セッション内に閉じる候補として扱える。
- `競合注意`: operation-audit 観点が扱う レビュー証跡保存と、モックデータ 非永続化が衝突する可能性がある。

### CAND-FFARF-LC-007 通常起動へ戻ると fakeAPI が選ばれない

- `根拠要件`: `task-frame.md:18`, `task-frame.md:22-23`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:21`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:26-27`, `tasks/usecases/frontend-fake-api-review-foundation.yaml:37`, `docs/architecture.md:60-61`, `docs/coding-guidelines-frontend.md:40-46`
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-FFARF-LC-007`
- `lifecycle phase`: 通常起動復帰
- `actor`: フロントエンド bootstrap
- `trigger`: fakeAPI レビュー終了後、通常の本番起動で フロントエンド を起動する。
- `開始条件`: レビュー起動モードが指定されていない。
- `期待結果`: 本番ゲートウェイ が選ばれ、fakeAPI と レビューモックデータ は接続先にも初期状態にも選ばれない。
- `観測点`: 本番起動で fakeAPI が選ばれないことを確認でき、局所テストでも fakeAPI 起動モードの破損と本番混入防止条件を確認できる。
- `関連詳細要求タイプ`: `compatibility_requirement`, `security_requirement`, `testability_requirement`
- `採用判断材料`: 本番起動の保護条件として、レビュー起動モードの終了後確認に接続できる。
- `競合注意`: trust-boundary 観点が扱う本番混入防止条件と統合される可能性がある。

## Open Notes

- `人間判断候補`: レビュー起動モードの指定方法は未決である。候補は URL クエリ、環境変数、専用 npm script、レビュー専用 entry のどれでも成立するが、設計担当が採否前に確認する必要がある。
- `人間判断候補`: 状態パターン選択の操作面は未決である。候補は画面内の選択 UI、URL クエリ、開発用パネル、起動時固定値のどれでも成立するが、UI 要件としての採否は扱わない。
- `統合候補`: CAND-FFARF-LC-003 と CAND-FFARF-LC-005 は、状態パターン切替方法を一つに固定する場合に統合候補になる。
- `統合候補`: CAND-FFARF-LC-006 と CAND-FFARF-LC-007 は、終了後の本番復帰確認を一つの受け入れシナリオにまとめられる可能性がある。
- `不採用候補`: 同じ起動セッション内の 状態パターン 再選択が不要と判断される場合、CAND-FFARF-LC-005 は不採用候補になる。
