# 単語翻訳 retry 待ち navigation 停止修正 plan

## 状態

- task-id: `2026-05-10-term-translation-recoverable-failed-navigation`
- workflow_state: `fix-lane-active`
- 依頼要約: 単語翻訳フェーズで provider 応答不正後、画面が `retry_waiting` を表示し、開始、中断、再開が拒否され、未完了一覧から現在の翻訳段階へ進めない症状を修正する。
- 現在成果物: `レビュー通過根拠`
- 次成果物: `レビュー通過根拠`。
- 停止条件: なし。

## 決定

- 本 task は `fix-lane` で扱う。
- 対象は、単語翻訳フェーズの失敗後状態と未完了一覧 navigation の整合である。
- 仕様変更、機能追加、新しい状態値の追加は禁止する。
- 実装 agent は `backend_implementer` に固定する。
- 実装 skill は `implement-backend` に固定する。
- 修正実行入力は、人間修正レビュー承認後に作成済みである。

## 理由

- 人間観測では、単語翻訳フェーズの操作が複数の状態条件で拒否されている。
- DB では、job 本体と phase run の状態が異なる粒度で残っている。
- 未完了一覧は `progress_percent` を navigation 可否の状態判断へ使い、現在の翻訳段階へ進む導線を止めている。
- 修正対象は backend 状態投影に閉じる。

## 影響

- `TRANSLATION_JOB.state=running` と `JOB_PHASE_RUN.state=recoverable_failed` の組み合わせは、`progress_percent` の値に関係なく現在の翻訳段階へ戻れる状態として扱う必要がある。
- current phase を特定できる `pending` phase run は、`progress_percent` の値に関係なく現在の翻訳段階へ進める状態として扱う必要がある。
- `progress_percent` は表示値または進捗表示用の値として扱い、navigation 可否から分離する必要がある。
- provider 応答不正そのものは別の失敗であり、今回の主対象は失敗後の操作導線と状態投影である。
- 直リンク防止は残すが、再試行可能失敗の正当な復帰導線を止めないことを確認する必要がある。

## 成果物DAG

| 成果物ID | 状態 | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- | --- |
| `task 枠` | 完了 | `fix_lane` | `[]` | なし |
| `人間観測記録` | 完了 | `fix_lane` | `task 枠` | なし |
| `修正前調査` | 完了 | `investigator` | `人間観測記録` | `investigator` |
| `修正方針判断` | 完了 | `fix_decider` | `人間観測記録`, `修正前調査` | `fix_decider` |
| `原因箇所シーケンス図` | 完了 | `diagrammer` | `人間観測記録`, `修正前調査`, `修正方針判断` | `diagrammer` |
| `人間修正レビュー` | 完了 | human | `修正方針判断`, `原因箇所シーケンス図` | human |
| `修正実行入力` | 完了 | `fix_lane` | `人間観測記録`, `修正前調査`, `修正方針判断`, `原因箇所シーケンス図`, `人間修正レビュー` | なし |
| `実装証跡` | 完了 | `backend_implementer` | `修正実行入力` | `backend_implementer` |
| `回帰テスト証跡` | 完了 | `implementation_unit_tester` | `実装証跡` | `implementation_unit_tester` |
| `最終検証` | 完了 | `fix_lane` | `実装証跡`, `回帰テスト証跡?` | なし |
| `実装後ブラウザ確認` | 完了 | `browser_confirmation` | `最終検証` | `browser_confirmation` |
| `レビュー通過根拠` | 着手可能 | `fix_lane` | `実装後ブラウザ確認` | review agents |
| `作業レポート入力` | 完了 | `fix_lane` | 完了または停止済み成果物 | `work_reporter` |

## 作成済み成果物

- `human-observation.md`: 人間観測記録。
- `investigation.md`: 修正前調査。
- `fix-decision.md`: 修正方針判断。
- `cause-sequence.puml`: 原因箇所シーケンス図の source。
- `cause-sequence.svg`: 原因箇所シーケンス図の描画結果。
- `human-fix-review-request.md`: 人間修正レビュー依頼。
- `human-fix-review-result.md`: 人間修正レビューの承認結果。
- `fix-implementation-input.md`: 修正実行入力。
- `implementation-evidence.md`: 実装証跡。
- `regression-test-evidence.md`: 回帰テスト証跡。
- `final-verification.md`: 最終検証。
- `browser-confirmation-input.md`: 実装後ブラウザ確認起動入力。
- `browser-confirmation-result.md`: 実装後ブラウザ確認結果。
- `work-report-input.md`: 作業レポート入力。

## 禁止変更範囲

- docs 正本本文は変更しない。
- `.codex/` と `.codex/skills` は変更しない。
- 新しい状態値の追加を前提にしない。
- UI 表示だけで症状を隠す修正を採用しない。
- `recoverable_failed` だけを特定の `progress_percent` warning の例外にしない。
- `progress_percent` による navigation 状態判断を残さない。
- provider raw response、prompt、翻訳本文全文をログへ出さない。

## 次に渡す入力

- `browser-confirmation-result.md` は、人間手動確認で phase page 遷移を補足済みである。
- `work-report-input.md` は、停止扱いではなく完了扱いへ更新済みである。
- `レビュー通過根拠` は、観点別レビュー agent へ渡せる状態である。

## 未決事項

- 単語翻訳画面の `再開` と `リトライ` の責務差分を UI でどう扱うか。
- provider 応答不正の内部原因を、今回の修正対象外として扱ってよいか。
