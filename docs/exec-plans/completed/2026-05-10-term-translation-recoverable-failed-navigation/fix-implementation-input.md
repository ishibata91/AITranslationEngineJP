# 修正実行入力

## 判断結果

- 判定: 完了。
- 対象成果物: `修正実行入力`。
- 実装 agent: `backend_implementer`。
- 実装 skill: `implement-backend`。
- 戻し先: `fix_lane`。

## 満たされた依存対象

- 人間観測記録: `human-observation.md`。
- 修正前調査: `investigation.md`。
- 修正方針判断: `fix-decision.md`。
- 原因箇所シーケンス図: `cause-sequence.puml` と `cause-sequence.svg`。
- 人間修正レビュー: `approved`。

## 実装対象

- 主対象: `internal/service/translation_job_management_service.go`。
- 確認候補: `internal/usecase/translation_job_management_usecase.go`。
- 変更してよい範囲: 未完了一覧 read model の progress warning 作成と phase navigation 可否判定に閉じる backend プロダクトコード。
- 変更してよい範囲: 今回変更の直接影響として backend 責務内プロダクトコードに閉じる補助関数。

## 修正方針

- 未完了一覧の backend 投影では、`progress_percent` を navigation 可否の状態判断に使わない。
- `progress_percent` は表示値、または進捗表示用の値として扱う。
- 現在の翻訳段階へ進めるかどうかは、phase run state と current phase を特定できるかで判断する。
- current phase を特定できる `recoverable_failed` phase run は、`progress_percent` の値に関係なく現在の翻訳段階へ戻れる状態として扱う。
- current phase を特定できる `pending` phase run は、`progress_percent` の値に関係なく現在の翻訳段階へ進める状態として扱う。
- phase run が存在しない `ready` job は、既定の開始 phase を特定できる場合に現在の翻訳段階へ進める状態として扱う。
- `phase_progress_aggregation_failed` は navigation block 理由として使わない。
- `phase_progress_aggregation_failed` は進捗表示の補足情報が不足する場合の warning に限定する。

## 禁止変更範囲

- 新しい状態値を追加しない。
- `recoverable_failed` を別名状態へ置き換えない。
- `pending` を別名状態へ置き換えない。
- `recoverable_failed` だけを特定の `progress_percent` warning の例外にしない。
- `progress_percent` による navigation 状態判断を残さない。
- job 本体の state を変更するだけで navigation block を回避しない。
- `phase_progress_aggregation_failed` warning を全体で無効化しない。
- `phase_progress_aggregation_failed` を navigation block 理由として使い続けない。
- frontend を変更しない。
- プロダクトテスト、検証データ、スナップショット、test helper を変更しない。
- provider 応答不正そのもの、provider raw response、prompt、翻訳本文全文を扱わない。
- docs 正本本文、`.codex/`、`.codex/skills` を変更しない。

## 回帰確認観点

- current phase を特定できる `recoverable_failed` phase run は、`progress_percent` の値に関係なく `canOpenPhase=true` へ投影される。
- current phase を特定できる `pending` phase run は、`progress_percent` の値に関係なく `canOpenPhase=true` へ投影される。
- current phase を特定できない状態は、navigation block として残る。
- `phase_progress_aggregation_failed` warning は navigation block 理由ではなく、進捗表示の補足 warning として扱われる。
- frontend は backend の `canOpenPhase` を尊重するため、backend read model の修正だけで未完了一覧の導線が復帰する。

## 検証コマンド

```bash
python3 scripts/harness/run.py --suite backend-local
```

## 実装後に返す内容

- 判断結果。
- 変更した backend プロダクトコード。
- 根拠参照。
- backend-local 検証結果。
- 未実行項目または残留リスク。
- 次に `fix_lane` が判断する材料。
