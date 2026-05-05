# Scenario Candidates: frontend-fake-api-review-foundation / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `FFARF-OA`

## Generator Scope

- `viewpoint`: `operation-audit`
- `included_sources`:
  - `./plan.md`
  - `./task-frame.md`
  - `tasks/usecases/frontend-fake-api-review-foundation.yaml`
  - `docs/architecture.md`
  - `docs/coding-guidelines-frontend.md`
  - `tmp/code-map/index.json`
- `excluded_sources`:
  - プロダクトコード の変更指示
  - プロダクトテスト の変更指示
  - docs 正本の更新判断
  - `.codex` の変更判断
  - 最終シナリオ表、採否、統合、競合解消
- `generation_notes`:
  - fakeAPI は provider 選択肢ではなく、レビュー起動時の DI による差し替えとして扱う。
  - 監査対象は、後から レビュー手順、表示 状態パターン、確認証跡、局所テスト、coverage 例外、混入防止を再確認できる材料に限定する。
  - 保存候補は再現に必要な要約だけにし、secret、個人情報、過剰な本文、外部 provider 応答原文は保存しない。

## Candidate Scenarios

### CAND-FFARF-OA-001 fakeAPI レビュー起動手順を後追い確認できる

- `根拠要件`: task-frame 完了条件「起動モードで フロントエンドの API 接続先を fakeAPI に切り替えられる」「本番起動では fakeAPI が選ばれない」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-FFARF-OA-001`
- `actor`: フロントエンドレビュー実行者
- `trigger`: レビュー用起動モードで フロントエンド を起動する
- `期待結果`: 後から、起動コマンド、起動モード、fakeAPI 選択結果、本番起動との差分を確認できる
- `audit event`: レビュー起動開始、起動成功、起動失敗、本番起動での fakeAPI 非選択確認
- `stored summary`: task id、起動モード名、実行コマンド名、対象 URL、選択された ゲートウェイ 種別、開始時刻、終了結果
- `redaction rule`: 環境変数値、local path の個人名部分、secret、provider key、任意入力本文は保存しない
- `観測点`: フロントエンド bootstrap が レビュー用 DI を選び、本番ゲートウェイ が選ばれていないことを確認できる記録
- `関連詳細要求タイプ`: 運用確認、監査ログ、再現材料、混入防止
- `採用判断材料`: designer は起動手順を検証シナリオへ統合できる
- `競合注意`: 起動モード名や記録粒度は、実装範囲側で確定する可能性がある

### CAND-FFARF-OA-002 レビューstate 状態パターン一覧を再現材料として確認できる

- `根拠要件`: task-frame 完了条件「空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態を fakeAPI で再現できる」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-FFARF-OA-002`
- `actor`: フロントエンドレビュー実行者
- `trigger`: レビュー用状態パターンを切り替える
- `期待結果`: 後から、確認対象 状態パターン、画面、モックデータ 要約、期待表示状態を一覧で確認できる
- `audit event`: 状態パターン一覧生成、状態パターン選択、状態パターン 表示成功、状態パターン 表示失敗
- `stored summary`: 状態パターン id、状態種別、対象画面、モックデータ セット名、期待する主要表示、確認結果
- `redaction rule`: モックデータの本文全量、外部応答原文、個人情報に見える値、secret 形式の値は保存しない
- `観測点`: 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態が個別 状態パターン として識別できる
- `関連詳細要求タイプ`: 履歴、再現材料、運用確認
- `採用判断材料`: designer は 状態パターン一覧を acceptance 条件または人間確認観点へ統合できる
- `競合注意`: 状態パターンの画面別分割粒度は、画面固有 task 側の確認単位と競合する可能性がある

### CAND-FFARF-OA-003 agent-browser 証跡を 状態パターンごとにひも付けられる

- `根拠要件`: task-frame 完了条件「実画面を agent-browser で開き、状態パターンごとの表示を確認できる」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-FFARF-OA-003`
- `actor`: フロントエンドレビュー実行者
- `trigger`: agent-browser で fakeAPI 起動中の実画面を開く
- `期待結果`: 後から、状態パターンごとの URL、操作、表示確認結果、取得証跡を追跡できる
- `audit event`: browser open、状態パターン select、snapshot または screenshot 取得、表示確認結果記録
- `stored summary`: 対象 URL、状態パターン id、操作名、証跡ファイル名または証跡参照、確認結果、失敗時の短いエラー要約
- `redaction rule`: screenshot 内の secret、provider key、個人情報、任意本文が映る場合は保存前に伏せ字または取得対象から除外する
- `観測点`: agent-browser 証跡が 状態パターン一覧と対応し、実フロント表示を後追い確認できる
- `関連詳細要求タイプ`: 監査ログ、履歴、再現材料
- `採用判断材料`: designer は agent-browser 確認を UI 人間操作 E2E 候補へ統合できる
- `競合注意`: 証跡の保存場所と保持期間は、人間判断が必要になる可能性がある

### CAND-FFARF-OA-004 fakeAPI 起動モードの局所テスト結果を記録できる

- `根拠要件`: task-frame 完了条件「fakeAPI 起動モードが壊れていないことを局所テストで確認できる」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-FFARF-OA-004`
- `actor`: 実装レーン検証者
- `trigger`: fakeAPI 起動モードの局所テストを実行する
- `期待結果`: 後から、どの局所テストが DI による差し替え、本番非選択、状態パターン供給を確認したか分かる
- `audit event`: 局所テスト開始、局所テスト完了、局所テスト失敗、再実行
- `stored summary`: 実行コマンド、対象テスト名、確認した境界、結果、失敗時の短いエラー要約、再実行有無
- `redaction rule`: stack trace の環境固有 path、secret、長い モック payload は保存しない
- `観測点`: フロントエンド bootstrap、ゲートウェイ契約、fakeAPI 境界、本番ゲートウェイ 非選択が局所テストで確認されている
- `関連詳細要求タイプ`: 監査ログ、再現材料、混入防止
- `採用判断材料`: designer は局所テストを実装後検証候補へ統合できる
- `競合注意`: どのテスト種別へ置くかは implementation-scope 側の分割判断と競合する可能性がある

### CAND-FFARF-OA-005 coverage 例外理由と代替確認を記録できる

- `根拠要件`: usecase 完了条件「coverage harness では fakeAPI 基盤を数値判定の例外として扱い、例外理由と局所テスト結果を記録できる」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-FFARF-OA-005`
- `actor`: 実装レーン検証者
- `trigger`: coverage harness で fakeAPI 基盤を数値判定の例外として扱う
- `期待結果`: 後から、例外対象、例外理由、代替の局所テスト結果、再確認条件を確認できる
- `audit event`: coverage 例外登録、例外理由確認、代替テスト結果確認、例外見直し
- `stored summary`: 例外対象範囲、例外理由、代替テスト名、代替テスト結果、見直し条件
- `redaction rule`: coverage report の不要な全文、環境固有 path、secret、モック payload 全量は保存しない
- `観測点`: 数値 coverage ではなく局所テストで保証する理由が明示されている
- `関連詳細要求タイプ`: 運用確認、監査ログ、人間判断候補
- `採用判断材料`: designer は coverage 例外を検証条件または残留リスク候補へ統合できる
- `競合注意`: 例外を許容する範囲と見直し期限は、人間判断が必要になる可能性がある

### CAND-FFARF-OA-006 fakeAPI と モックデータの本番混入を監査できる

- `根拠要件`: task-frame 完了条件「fakeAPI と モックデータが本番 API、永続化、本番初期状態に混入しない」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-FFARF-OA-006`
- `actor`: 実装レーン検証者
- `trigger`: 本番起動または構造確認で fakeAPI 混入防止を確認する
- `期待結果`: 後から、fakeAPI が 本番ゲートウェイ、永続化、本番初期状態 へ混入していない根拠を確認できる
- `audit event`: 本番 起動確認、import 境界確認、永続化非接続確認、初期状態 非注入確認
- `stored summary`: 確認対象ファイル群の要約、本番ゲートウェイ 選択結果、永続化接続有無、初期状態 注入有無、確認結果
- `redaction rule`: 永続化内容、ユーザー入力本文、local DB path の個人名部分、secret は保存しない
- `観測点`: `フロントエンド/src/main.ts` が composition root であり、fakeAPI が レビュー用 DI による差し替えに限定されている
- `関連詳細要求タイプ`: 混入防止、監査ログ、再現材料
- `採用判断材料`: designer は本番非混入を強い受け入れ条件へ統合できる
- `競合注意`: 構造確認を lint、局所テスト、manual check のどこへ置くかは統合時の判断対象になる

### CAND-FFARF-OA-007 画面固有モックデータの追加履歴を追跡できる

- `根拠要件`: task-frame 完了条件「画面固有のモックデータを ユースケース task 側で追加できる」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-FFARF-OA-007`
- `actor`: ユースケース task 実装者
- `trigger`: 画面固有の レビューstate モックデータを追加する
- `期待結果`: 後から、どの画面のどの 状態パターン に モックデータが追加されたか確認できる
- `audit event`: モックデータ 追加、状態パターン ひも付け、モックデータ 更新、モックデータ 削除
- `stored summary`: 対象画面、状態パターン id、モックデータ セット名、追加理由、関連 task id、本番 からの非参照確認
- `redaction rule`: モックデータの本文全量、個人情報に見える値、secret、外部 provider 応答原文は保存しない
- `観測点`: ユースケース task 側で モックデータを追加でき、本番初期状態 へ混入していない
- `関連詳細要求タイプ`: 履歴、再現材料、混入防止
- `採用判断材料`: designer は画面別 状態パターン 拡張の運用ルール候補へ統合できる
- `競合注意`: モックデータの命名規則と配置規則は、implementation-scope 側の実装分割判断と競合する可能性がある

## Open Notes

- `人間判断候補`:
  - agent-browser 証跡の保存場所、保持期間、伏せ字範囲は人間判断が必要になる可能性がある。
  - coverage 例外の許容範囲、見直し条件、代替テストの最低基準は人間判断が必要になる可能性がある。
  - 状態パターン一覧を画面共通の正本にするか、ユースケース task ごとの成果物にするかは人間判断が必要になる可能性がある。
- `統合候補`:
  - `CAND-FFARF-OA-001` と `CAND-FFARF-OA-006` は、起動モードと本番非混入の確認として統合候補になる。
  - `CAND-FFARF-OA-002` と `CAND-FFARF-OA-003` は、状態パターン一覧と実画面証跡の確認として統合候補になる。
  - `CAND-FFARF-OA-004` と `CAND-FFARF-OA-005` は、局所テストと coverage 例外の検証根拠として統合候補になる。
- `不採用候補`:
  - tool 権限、agent 実行定義、プロダクトコード変更手順に踏み込む候補は operation-audit 観点から除外する。
  - UI 表示仕様そのものを確定する候補は画面固有 task 側へ残す。
  - 最終シナリオ表、採否、統合、競合解消は designer 側へ残す。
