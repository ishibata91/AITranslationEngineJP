# Harness Engineering Foundation

- Date: 2026-03-28
- Status: Completed

## Goal

開発基盤向けのハーネスエンジニアリング方針を、実際のリポジトリ構造と検証入口へ落とし込む。

## Scope

- `AGENTS.md` の新設
- `docs/` の索引、原則、品質、負債、計画置き場の整備
- `docs/external-design/` の参照切れ解消
- `scripts/harness/` の最小検証入口追加

## Acceptance Criteria

- `AGENTS.md` から読む順序が一意に辿れる
- `docs/external-design/*` の参照切れがなくなる
- 構造、設計、実行の 3 系統の検証入口が存在する
- 品質スコアと負債一覧に、現時点の制約が反映される

## Outcome

- 入口文書と索引文書を追加した
- 記録系ドキュメントの保存先を固定した
- 外部設計の初期スタブを追加した
- PowerShell ベースの最小ハーネスを追加した

## Follow-up

- `external-design/` を実装計画に応じて肉付けする
- `docs/api-refrences/` の段階移行方針を決める
- 翻訳品質評価ハーネスを次フェーズで設計する
