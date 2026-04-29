---
name: codex-review-contract
description: Codex 実装後 レビュー の契約・互換性グループ作業プロトコル。
---
# Codex Review Contract

## 目的

既存利用者、外部 API、内部 API、DB schema、event 入力内容 を壊していないかを見る。
コード自体が動いても契約破壊が利用者側障害になるため、diff の public 境界 を採点する。

## 対応ロール

- `review_contract` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` の レビュー 集約 とする。
- 担当成果物は `codex-review-contract` の出力規約で固定する。

## 入力規約

- API request / response
- GraphQL schema と DB migration
- public method と event 入力内容
- queue message と webhook
- エラー code、nullable / 必須、versioning
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: レビュー対象差分, implementation-scope の場所, 実装結果
- 任意入力: 最終検証結果, 変更ファイル

## 外部参照規約

- エージェント実行定義とツール権限は [review_contract.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/review_contract.toml) の 書き込み許可 / 実行許可 とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-contract/SKILL.md

## 内部参照規約

## 判断規約

互換性が 0.85 を超える場合だけ通過とする。
既存 field の意味変更や nullable / 必須 の変更は高い減点対象にする。

## 出力規約

出力は、次の情報を必ず返す。

- 観点名: 契約・互換性観点であることを返す。
- 判定: 通過、失敗、停止のいずれかを返す。
- 互換性評価: 既存利用者、外部 API、内部 API、DB schema、event 入力内容 への互換性を返す。
- 根拠十分性: 確認済み根拠の十分性を返す。
- 確認範囲: 確認した public 境界、schema、入力内容、未確認範囲を返す。
- 根拠: 判断へ使ったファイルと参照先を返す。
- 破られた不変条件: API、schema、nullable / 必須、versioning のどれが破られたかを返す。
- 原因候補: 契約破壊を生む原因候補と根拠を返す。
- 指摘: API、schema、public method、event 入力内容、nullable / 必須、versioning の問題だけを返す。
- 局所修正評価: 局所 shim で足りるか、公開接点 固定が必要かを返す。
- 追加確認範囲: 互換性確認に追加で読むべき範囲と読まない範囲を返す。
- 修正時の考慮点: 修正者が考慮すべき支配点、互換リスク、恒久修正の判断材料を返す。
- 不変条件テスト: 破れた契約を固定するテスト観点を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務、修正範囲の命令を含めない。

## 完了規約

- 対象 レビュー 観点の指摘、互換性評価、根拠、残留リスクが返却されている。
- 権限・信頼境界系の強制停止条件は、他観点の高評価で相殺せず明示されている。
- API request / response の互換性を確認した。
- DB schema、event 入力内容、queue message を確認した。
- public method と エラー code を確認した。
- nullable / 必須 と versioning を確認した。
- 破られた不変条件と原因候補を分けた。
- 局所修正評価と不変条件テスト観点を返した。
- 内部実装の綺麗さを主判定にしなかった。
- 完了判断材料として、互換性評価、破られた不変条件、原因候補、局所修正評価、根拠が返っている。
- 残留リスクとして、未確認範囲と理由が返っている。

## 停止規約

- 内部実装の綺麗さ
- テストの十分性
- 可読性
- パフォーマンス最適化
- 停止時は不足項目、衝突箇所、戻し先を返す。
