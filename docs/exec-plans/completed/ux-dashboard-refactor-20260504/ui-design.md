# UI改善契約: ux-dashboard-refactor-20260504

- `skill`: `ui-design`
- `status`: `approved`
- `target_screen`: ダッシュボード
- `source_plan`: [plan.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/ux-dashboard-refactor-20260504/plan.md)
- `existing_screen_evidence`: [existing-screen-evidence.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/ux-dashboard-refactor-20260504/existing-screen-evidence.md)
- `human_ui_review`: 承認

## 既存画面根拠

対象: `AppShell` の `#dashboard` 表示に限定する。
理由: `plan.md` は対象画面をダッシュボードに限定し、既存表示項目、導線、状態表示、禁止表示の維持を求めている。

- 既存画面確認は `http://127.0.0.1:34115/#dashboard` で pass と記録されている。
- 共通ヘッダ、グローバルナビゲーション、現在地表示、本文入口カードは `AppShell.svelte` の構造に存在する。
- 入口カードの表示内容は `shell-state.ts` の `SHELL_ROUTE_CONTRACT` から供給される。
- シナリオ正本はダッシュボード既定表示、主要ページ遷移、ジョブ一覧と進捗サマリの非表示を要求する。
- `SCN-DAS-003` の手順文はダッシュボード入口カードを含むように読めるが、既存画面根拠と実装はダッシュボード自身の入口カードを出さない。

根拠行:

- [existing-screen-evidence.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/ux-dashboard-refactor-20260504/existing-screen-evidence.md:17)
- [AppShell.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/views/AppShell.svelte:142)
- [shell-state.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts:28)
- [dashboard-and-app-shell.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/dashboard-and-app-shell.md:10)

## 維持対象チェックリスト

表示項目:

- [ ] 共通ヘッダに `AITranslationEngineJp` と `翻訳エンジン` を表示する。
- [ ] グローバルナビゲーションに `ダッシュボード`、`AIサービス設定`、`マスター辞書`、`マスターペルソナ`、`翻訳管理`、`出力管理` を表示する。
- [ ] ページ見出しに `現在のページ` と `ダッシュボード` を表示する。
- [ ] リード文に `最初に移動したい作業を選び、共通ナビゲーションからいつでも別の主要ページへ切り替えられます。` を表示する。
- [ ] ダッシュボード本文に `主要ページ` と `作業を選ぶ` を表示する。
- [ ] 入口カードに `AIサービス設定`、`マスター辞書`、`マスターペルソナ`、`翻訳管理`、`出力管理` を表示する。

入口カード:

- [ ] `AIサービス設定`: 状態 `設定状態を確認`、説明 `エンドポイントと APIキー状態を AIサービスごとに確認します。`、操作 `開く` を表示する。
- [ ] `マスター辞書`: 状態 `確認可能`、説明 `用語と訳語の基盤データを確認します。`、操作 `開く` を表示する。
- [ ] `マスターペルソナ`: 状態 `確認可能`、説明 `翻訳に使うペルソナ設定を確認します。`、操作 `開く` を表示する。
- [ ] `翻訳管理`: 状態 `Body Phase UI 追加`、説明 `入力確認、validation、ready job 作成、term phase、persona phase、body phase の実行状況をまとめて確認します。`、操作 `開く` を表示する。
- [ ] `出力管理`: 状態 `確認可能`、説明 `生成物と書き出し結果を確認します。`、操作 `開く` を表示する。

導線:

- [ ] グローバルナビゲーションの各リンクは対応する `#<route id>` へ遷移する。
- [ ] 入口カードの各リンクは同じラベルのグローバルナビゲーションと同じ `#<route id>` へ遷移する。
- [ ] `#dashboard` 表示中は入口カードだけを表示し、遷移先ページの業務 UI を混ぜない。
- [ ] 860px 以下では `主要ページ` ボタンでグローバルナビゲーションの開閉状態を確認できる。

## 禁止表示チェックリスト

- [ ] ダッシュボード本文に `ジョブ一覧` を表示しない。
- [ ] ダッシュボード本文に `進捗サマリ` を表示しない。
- [ ] ダッシュボード入口カードに `ダッシュボード` 自身を表示しない。
- [ ] 既存 route 以外の新しい主要ページ入口を追加しない。
- [ ] backend 連携、永続化、公開契約変更を前提にした表示を追加しない。

## UI改善契約

目的: ダッシュボードを主要ページ入口として整理し、既存の全表示項目、導線、状態表示、禁止表示を欠落させない。
範囲: `AppShell` のダッシュボード本文、共通ヘッダ、グローバルナビゲーション、入口カードの表示維持に限定する。

表示契約:

- `dashboardEntryRoutes` は `dashboard` を除外した route 一覧として扱う。
- 各入口カードは `label`、`state`、`description`、操作 `開く` を 1 組で表示する。
- `card-tag` に route id を出す場合、既存 id だけを表示し、新しい業務概念を作らない。
- 状態値と説明文は `shell-state.ts` の既存文字列から意味を広げない。
- 人間UIレビューの追加条件により、既存状態値 `準備中` は `確認可能` へ変更する。

操作契約:

- グローバルナビゲーションと入口カードは hash route のみを変更する。
- 入口カードの主要操作はカード全体のリンクと `開く` 表示だけにする。
- ダッシュボード上で作成、実行、削除、保存、再試行の操作を追加しない。
- モバイルナビゲーションは `主要ページ` ボタンで開閉し、遷移後は閉じる。

配置契約:

- 共通ヘッダ、現在地表示、リード文、入口カード群の順序を維持する。
- 入口カードはデスクトップで複数列、狭い幅で 1 列に落とす。
- 入口カード内ではラベル、状態、説明、操作を同じカード内に近接させる。
- 長い説明文はカード内で折り返し、横スクロールや重なりを発生させない。

部品化判断:

- 部品化対象: 入口カードは画面専用部品として切り出せる。
- 配置先判断: 共有部品化は人間UIレビュー後の実装修正入力で判断する。
- 分けない対象: `AppShell` の route 初期化、hash 同期、グローバルナビゲーション開閉はダッシュボード入口カードへ移さない。
- 判断理由: 入口カードは表示単位として独立するが、route 同期は共通シェル責務である。

UIプロトタイプ:

- `prototype_required`: `no`
- `prototype_path`: なし
- `reason`: 既存画面変更の土台は実画面と既存証跡で足りる。新しい見た目体系を提案しないため、task-local UIプロトタイプは作らない。
- `confirmation_url`: なし
- `server_command`: なし
- `human_review_designer_agent_required`: `no`
- `human_feedback_route`: `ux_refactor_lane`

## UX標準確認

画面構造:

- 画面目的: pass。ダッシュボードは主要ページ入口に限定する。
- 主要CTA: pass。各入口カードの操作 `開く` だけを扱う。
- 画面責務: pass。ジョブ一覧、進捗サマリ、実行操作を混ぜない。
- 状態表示: pass。各入口カードの状態値を既存 route 契約から表示する。
- 禁止条件: pass。禁止表示チェックリストで非表示対象を固定する。

配置とレスポンシブ:

- 情報階層: pass。現在地、リード文、主要ページ入口の順序を維持する。
- 関連情報の近接: pass。カード内でラベル、状態、説明、操作を同居させる。
- ブレークポイント: pass。860px 以下はナビゲーションを開閉式にし、入口カードを 1 列にする。
- 長文耐性: 要実装後確認。説明文と状態値がカード外へはみ出さないことを確認する。
- 証跡: 要実装後確認。デスクトップとモバイルの画面証跡を保存する。

表示文言:

- 状態ラベルは短い状態値として維持する。
- `準備中` は利用者の次操作を妨げるため、入口カードでは `確認可能` に変更する。
- 操作文言は `開く` に限定する。
- 説明文は既存の日本語文を維持する。
- `Body Phase UI 追加`、`validation`、`ready job`、`term phase`、`persona phase`、`body phase` は既存表示としてのみ維持する。

## 実装後確認観点

画面確認:

- `agent-browser open http://127.0.0.1:34115/#dashboard`
- `agent-browser snapshot`
- `agent-browser screenshot tmp/agent-browser/ux-dashboard-refactor-after-desktop.png`
- `agent-browser errors`

維持確認:

- [ ] 維持対象チェックリストの全項目を目視または snapshot で確認する。
- [ ] 禁止表示チェックリストの全項目が表示されないことを確認する。
- [ ] 5 つの入口カードの遷移先がグローバルナビゲーションと一致することを確認する。
- [ ] 860px 以下で `主要ページ` ボタン、ヘッダ、見出し、入口カード、操作 `開く` が確認できる。
- [ ] console error が存在する場合は、UIリファクタ由来か既存由来かを分けて記録する。

候補コマンド:

- `npm --prefix frontend run test -- AppShell`

## agent-browser 確認結果

- `designer_agent_browser`: 未実行
- `reason`: designer は実画面 observation を新規に実施せず、`ux_refactor_lane` が固定した既存画面証跡を根拠に UI改善契約を作成した。
- `existing_result`: `agent-browser open`、`snapshot`、`screenshot` は pass と記録済み。
- `existing_screenshot`: `tmp/agent-browser/ux-dashboard-refactor-before.png`
- `console_errors`: 具体的な console error text は既存証跡から取得できていない。

## 人間UIレビュー状態

- `human_ui_review`: 承認
- `required_before_implementation`: yes
- `review_target`: この `ui-design.md`
- `review_result`: `準備中` の文言変更を追加すれば承認。
- `review_decision`: 入口カードの状態値 `準備中` を `確認可能` へ変更する。
- `after_review`: 承認済み UI改善契約として `ux_refactor_lane` へ戻す。
