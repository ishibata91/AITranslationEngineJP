# Development Roadmap

関連文書: [`../docs/index.md`](../docs/index.md), [`quality-score.md`](./quality-score.md), [`tech-debt-tracker.md`](./tech-debt-tracker.md)

このファイルは、2026-03-29 時点の repository 状態をもとに、人間向けの開発順序、進捗状態、次の着手単位を整理する。

## Status Legend

- `完了`: 現物または completed plan で成立を確認できる
- `進行候補`: 次の着手対象として妥当だが、まだ開始していない
- `未完了`: 必要性は明確だが、実装または検証が不足している

## Current Snapshot

| Area | Status | Notes |
|---|---|---|
| repository 骨格 | 完了 | `src/` と `src-tauri/` の bootstrap 構成がある |
| workflow / role 契約 | 完了 | `.codex/README.md` と workflow skills が正本になっている |
| structure harness | 完了 | required path と markdown link の検査入口がある |
| design harness | 完了 | semantic checks まで含めて成立している |
| execution harness | 完了 | lint / test / build / cargo / sonar の入口がある |
| frontend 実装 | 未完了 | `AppShell` と bootstrap status の最小表示のみ |
| backend 実装 | 未完了 | bootstrap usecase と Tauri command の最小配線のみ |
| 翻訳ドメイン実装 | 未完了 | import、job、dictionary、persona、output は未着手に近い |
| 翻訳業務フロー acceptance checks | 未完了 | `tech-debt-tracker.md` の open item に残っている |

## Progress Summary

### すでに完了した土台

- `完了`: Tauri 2 + Svelte 5 + Rust + TypeScript の bootstrap
- `完了`: frontend / backend の directory contract 初期化
- `完了`: lint / test / build / sonar / cargo を束ねる execution harness の追加
- `完了`: structure harness と design harness の整備
- `完了`: `4humans/quality-score.md` と `4humans/tech-debt-tracker.md` の運用開始

### まだ終わっていない主要領域

- `未完了`: xEdit 抽出 JSON の import
- `未完了`: `PLUGIN_EXPORT`、`TRANSLATION_UNIT`、translation job の永続化
- `未完了`: ジョブ状態遷移の業務実装
- `未完了`: マスター辞書構築
- `未完了`: マスターペルソナ構築
- `未完了`: 単語翻訳フェーズ
- `未完了`: NPCペルソナ生成フェーズ
- `未完了`: 本文翻訳フェーズ
- `未完了`: LMStudio / Gemini / xAI provider 実装
- `未完了`: 標準配布形式 / xTranslator 互換形式の出力

## Roadmap Policy

- 先に `業務フローを成立させる最小縦切り` を作る
- 次に `翻訳品質に効く基盤データ` を積む
- その後に `AI provider の拡張` と `運用制御` を広げる
- 各フェーズで tests / acceptance checks / validation commands を同時に増やす
- 完了判定は文書宣言ではなく、現物、tests、completed plan で行う

## Phase 0: Foundation Stabilization

### Phase Status

- `進行候補`: 大部分は完了しているが、日常実装を始める前の残件が少しある

### Work Breakdown

- `完了`: Tauri 2 / Svelte 5 / Rust / TypeScript の bootstrap 構成を追加した
- `完了`: frontend root に `src/ui/`、`src/application/`、`src/gateway/`、`src/shared/` を置いた
- `完了`: backend root に `src-tauri/src/application/`、`domain/`、`infra/`、`gateway/` を置いた
- `完了`: structure harness を追加した
- `完了`: design harness を追加した
- `完了`: execution harness を追加した
- `完了`: `4humans` 記録の格納先を repository 内へ揃えた
- `未完了`: 初期テンプレートから実業務 feature を量産するための screen / store / usecase 雛形を増やす

### Exit Criteria

- `完了`: 新しい feature を `src/` と `src-tauri/` の正本構成に沿って追加できる

## Phase 1: Input Cache And Job Skeleton

### Phase Status

- `進行候補`: 最優先で着手すべきフェーズ

### Work Breakdown

- `未完了`: xEdit 抽出 JSON importer を追加する
- `未完了`: import 時に入力データを validation する
- `未完了`: `PLUGIN_EXPORT` 相当の保存単位を実装する
- `未完了`: translatable field を `TRANSLATION_UNIT` 相当へ正規化する
- `未完了`: translation job 作成 usecase を実装する
- `未完了`: `Draft` / `Ready` / `Running` / `Completed` の最小状態遷移を実装する
- `未完了`: UI で job 作成画面を出す
- `未完了`: UI で job 一覧と job 状態表示を出す
- `未完了`: 単一 fixture を使う import-to-job の acceptance check を追加する
- `未完了`: lossless な翻訳単位保持を test で固定する

### Exit Criteria

- `未完了`: 単一の入力データを import して job 作成から完了まで追跡できる
- `未完了`: `FormID`、`EditorID`、レコード種別、フィールド種別、原文、訳文、出力ステータスを lossless に保持できる
- `未完了`: 最初の fixture-based acceptance check が execution harness から実行できる

## Phase 2: Dictionary And Persona Foundation

### Phase Status

- `未完了`: Phase 1 後に着手

### Work Breakdown

- `未完了`: xTranslator 形式 importer を追加する
- `未完了`: マスター辞書保存構造を追加する
- `未完了`: 辞書 entry の検索 / 再利用 port を定義する
- `未完了`: ベースゲーム NPC 入力からマスターペルソナ構築処理を追加する
- `未完了`: ジョブ単位ペルソナ保持を `MASTER_PERSONA` と分離して保存する
- `未完了`: UI からマスター辞書を観測できるようにする
- `未完了`: UI からマスターペルソナを観測できるようにする
- `未完了`: 基盤データ再構築の validation command を用意する

### Exit Criteria

- `未完了`: マスター辞書とマスターペルソナを個別に再構築できる
- `未完了`: translation job が基盤データを参照できる
- `未完了`: UI から基盤データの観測結果を確認できる

## Phase 3: Translation Flow MVP

### Phase Status

- `未完了`: Phase 2 後に着手

### Work Breakdown

- `未完了`: 翻訳レコード種別ごとの翻訳指示構成を実装する
- `未完了`: 単語翻訳フェーズを実装する
- `未完了`: NPCペルソナ生成フェーズを実装する
- `未完了`: 本文翻訳フェーズを実装する
- `未完了`: 単語翻訳結果を本文翻訳フェーズで再利用する
- `未完了`: `<10gold>` などの埋め込み要素保持を実装する
- `未完了`: 翻訳結果 preview を UI で観測できるようにする
- `未完了`: 代表シナリオの fixture-based regression check を追加する

### Exit Criteria

- `未完了`: 単語翻訳フェーズから本文翻訳フェーズへの再利用が成立する
- `未完了`: 埋め込み要素を壊さない回帰 check がある
- `未完了`: 代表的な翻訳レコード種別で scenario regression が回る

## Phase 4: Provider And Execution Expansion

### Phase Status

- `未完了`: Translation Flow MVP 後に着手

### Work Breakdown

- `未完了`: LMStudio provider adapter を追加する
- `未完了`: Gemini provider adapter を追加する
- `未完了`: xAI provider adapter を追加する
- `未完了`: provider 選択 port と設定保持を追加する
- `未完了`: 単発実行と Batch API 実行の切替を追加する
- `未完了`: `Paused` / `RecoverableFailed` / `Failed` / `Canceled` の遷移を実装する
- `未完了`: 再開、リトライ、キャンセルの UI 操作を追加する
- `未完了`: 進捗観測と失敗理由表示を追加する
- `未完了`: provider failure / retry の acceptance checks を追加する

### Exit Criteria

- `未完了`: provider を差し替えても application / domain の方針が崩れない
- `未完了`: job の中断、再開、失敗回復が acceptance check で確認できる
- `未完了`: provider ごとの接続失敗や再試行条件が test で固定されている

## Phase 5: Output And Release Readiness

### Phase Status

- `未完了`: provider と job 制御の基本成立後に着手

### Work Breakdown

- `未完了`: 標準配布形式 writer を実装する
- `未完了`: xTranslator 互換形式 writer を実装する
- `未完了`: `EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` の再構成を固定する
- `未完了`: job 完了後の出力成果物記録を追加する
- `未完了`: 未完了 job 参照がない入力キャッシュ削除を実装する
- `未完了`: 再取り込み可能性を壊さない cleanup check を追加する
- `未完了`: contract-level tests と scenario regression を execution harness に統合する

### Exit Criteria

- `未完了`: 出力成果物を再利用可能な形で生成できる
- `未完了`: 未完了 job 参照がない入力キャッシュだけを削除できる
- `未完了`: execution harness が業務フローの主要シナリオを継続監視できる

## Immediate Next Slice

### Priority 1

- `進行候補`: `入力取込 -> 翻訳単位正規化 -> job 作成 -> job 一覧表示` の縦切り

### Priority 2

- `進行候補`: 上記縦切りに対応する最初の fixture-based acceptance check

## Current Risks

- `未完了`: 翻訳ドメイン固有の tests / fixtures / acceptance checks が薄く、仕様固定より実装が先に走る危険がある
- `未完了`: provider や補助メタデータへ先に広げると、job 骨格が固まる前に複雑性が上がる
- `未完了`: phase 単位の進捗管理は始めたが、各 phase の owner task はまだ未分解である

## Done Definition Per Phase

- `未完了`: 実装コードだけでなく、対応する tests / acceptance checks / validation commands が同じ変更に入っている
- `未完了`: `4humans/` の品質記録と負債記録に必要な差分が同期されている
- `未完了`: 完了した非自明タスクは completed plan に結果が残っている
