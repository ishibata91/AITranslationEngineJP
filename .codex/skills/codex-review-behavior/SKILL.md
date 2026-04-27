---
name: codex-review-behavior
description: Codex 実装後 review の挙動正しさグループ知識 package。
---

# Codex Review Behavior

## 目的

変更後のコードが PR の目的どおりに振る舞うかを見る。
diff から取得した実コードを、正解の挙動ベクトルにどの程度近いかで score 化する。

## 見るもの

- 正常系の挙動
- 条件分岐と境界値
- 例外系
- 既存挙動との差分
- bug 修正の場合の原因対応

## 見ないもの

- 命名
- 関数分割
- 読みやすさ
- テスト網羅性
- コードスタイル

## 判定

`score > 0.85` を pass とする。
仕様にない入力や不明な期待値は confidence を下げ、score と混同しない。

## 出力責務

この reviewer は修正範囲を命令しない。
修正判断に必要な情報として次を返す。

- `observed_scope`: 確認した実コード、経路、未確認範囲
- `violated_invariant`: PR 目的、既存挙動、受け入れ条件のどれが破られたか
- `root_cause_hypotheses`: 症状を生む原因候補と根拠
- `local_patch_assessment`: 局所修正で閉じるか、他層へ波及するか
- `exploration_scope`: 追加で読むべき範囲と読まない範囲
- `remediation_considerations`: 修正者が考慮すべき支配点、risk、invariant tests

## Checklist

- [codex-review-behavior-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-behavior/references/checklists/codex-review-behavior-checklist.md) を参照する。
