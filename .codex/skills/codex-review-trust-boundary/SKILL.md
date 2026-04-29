---
name: codex-review-trust-boundary
description: Codex 実装後 レビュー の権限・信頼境界グループ作業プロトコル。必須判定条件 として扱う。
---
# Codex Review Trust Boundary

## 目的

ユーザー、tenant、role、外部入力、secret の境界を越えていないかを見る。
他観点の高評価で相殺してはいけないため、強制停止条件として扱う。

## 対応ロール

- `review_trust_boundary` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` の レビュー 集約 とする。
- 担当成果物は `codex-review-trust-boundary` の出力規約で固定する。

## 入力規約

- 認証と認可
- tenant isolation と admin 権限
- user-controlled 入力
- secret 漏洩と PII
- SQL injection、XSS、SSRF、file upload、外部 URL
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: レビュー対象差分, implementation-scope の場所, 実装結果
- 任意入力: 最終検証結果, 変更ファイル

## 外部参照規約

- エージェント実行定義とツール権限は [review_trust_boundary.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/review_trust_boundary.toml) の 書き込み許可 / 実行許可 とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-trust-boundary/SKILL.md

## 内部参照規約

## 判断規約

安全性が 0.85 を超える場合だけ通過とする。
この観点の失敗は、他観点の平均評価で相殺しない。

## 出力規約

出力は、次の情報を必ず返す。

- 観点名: 権限・信頼境界観点であることを返す。
- 判定: 通過、失敗、停止のいずれかを返す。
- 安全性評価: 権限・信頼境界の安全性を返す。
- 根拠十分性: 確認済み根拠の十分性を返す。
- 強制停止扱い: この観点の失敗を他観点で相殺しないことを返す。
- 確認範囲: 確認した信頼境界、外部入力、secret、認可経路、未確認範囲を返す。
- 根拠: 判断へ使ったファイルと参照先を返す。
- 破られた不変条件: 認証、認可、secret、tenant、外部入力のどれが破られたかを返す。
- 原因候補: 信頼境界破壊を生む原因候補と根拠を返す。
- 指摘: 認証、認可、tenant isolation、外部入力、secret、SQL injection、XSS、SSRF、外部 URL、admin 権限、PII の問題だけを返す。
- 局所修正評価: 局所ガードで足りるか、境界の再固定が必要かを返す。
- 追加確認範囲: 認可、入力、secret、外部接続で追加確認すべき範囲と読まない範囲を返す。
- 修正時の考慮点: 修正者が考慮すべき支配点、強制停止リスク、恒久修正の判断材料を返す。
- 不変条件テスト: 破れた信頼境界を固定するテスト観点を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務、修正範囲の命令を含めない。

## 完了規約

- 対象 レビュー 観点の指摘、安全性評価、根拠、残留リスクが返却されている。
- 権限・信頼境界系の強制停止条件は、他観点の高評価で相殺せず明示されている。
- 認証、認可、tenant isolation を確認した。
- user-controlled 入力 と外部 URL を確認した。
- secret、admin 権限、PII を確認した。
- SQL injection、XSS、SSRF、file upload を確認した。
- 破られた不変条件と原因候補を分けた。
- 局所修正評価と不変条件テスト観点を返した。
- 強制停止条件の失敗を他観点で相殺しなかった。
- 完了判断材料として、安全性評価、破られた不変条件、原因候補、局所修正評価、根拠が返っている。
- 残留リスクとして、未確認範囲と理由が返っている。

## 停止規約

- 実装の短さ
- 読みやすさ
- 性能
- テスト妥当性
- 停止時は不足項目、衝突箇所、戻し先を返す。
