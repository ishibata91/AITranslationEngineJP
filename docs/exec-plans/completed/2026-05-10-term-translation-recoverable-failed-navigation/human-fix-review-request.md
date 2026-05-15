# 人間修正レビュー依頼

## レビュー対象

- `fix-decision.md`: 修正方針判断。
- `cause-sequence.puml`: 原因箇所シーケンス図の source。
- `cause-sequence.svg`: 原因箇所シーケンス図の描画結果。

## 承認対象

- 原因の原因: 未完了一覧の backend 状態投影が、`progress_percent` を navigation 可否の状態判断に使っている。
- 責務境界: 主修正対象は `TranslationJobManagementService` の progress warning 作成と navigation 可否判定である。
- 採用方針: `progress_percent` は表示値、または進捗表示用の値として扱う。
- 採用方針: navigation 可否は、phase run state と current phase を特定できるかで判断する。
- 採用方針: current phase を特定できる `recoverable_failed` phase run は、`progress_percent` の値に関係なく現在の翻訳段階へ戻れる状態として扱う。
- 採用方針: current phase を特定できる `pending` phase run は、`progress_percent` の値に関係なく現在の翻訳段階へ進める状態として扱う。
- 採用方針: `phase_progress_aggregation_failed` は navigation block 理由として使わない。
- 実装 agent: `backend_implementer`。
- 実装 skill: `implement-backend`。

## 差し戻し対象

- `recoverable_failed` を新しい job-level 状態へ変更する方針。
- `pending` を新しい job-level 状態へ変更する方針。
- `recoverable_failed` だけを特定の `progress_percent` warning の例外にする方針。
- `progress_percent` による navigation 状態判断を残す方針。
- frontend だけで「現在の翻訳段階へ進む」を enabled にする方針。
- `phase_progress_aggregation_failed` warning を全体で無効化する方針。
- provider 応答不正そのものを今回の修正対象へ含める方針。
- 単語翻訳画面の action reason 表示だけで navigation 停止を隠す方針。

## 確認してほしい点

- 上記の原因の原因、責務境界、採用方針で `修正実行入力` へ進めてよいか。
- 単語翻訳画面の開始、中断、再開の拒否理由表示を、今回の修正範囲へ含めない判断でよいか。
- provider 応答不正の内部原因を、今回の修正対象外として扱ってよいか。
