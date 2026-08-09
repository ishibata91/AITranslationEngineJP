# Implementation: prebuilt-npc-name-derivation

## 変更したfile

- REF-3 の本文参考語を作る処理: 事前作成辞書の `NPC_` 氏名から、本文実行中だけで使う二つ名と姓名分割の部分形を参考語へ追加した。
- REF-12: 派生結果に派生元入力の位置を持たせ、派生元の参考語メタデータを再計算なしで引けるようにした。
- REF-3 と REF-6 の本文用人名派生: 対象 plugin の会話文だけを用法集計へ渡した。
- REF-10 と REF-11: 事前作成辞書の派生、境界、保存しないこと、batch 本文計画の参考語を検証した。
- `spec.md`: 各仕様の実テストを記録した。

## 仕様との対応

- R-1-1: `appendPrebuiltNPCDerivedReferences` が二つ名と姓名分割を `termderive.DeriveTerms` へ渡す。同期は `bodyReferences`、batch は `planBodyRequests` から同じ入口を使う。
- R-1-2: 英語と日本語の組で NPC 入力を一意化し、最初の NPC のメタデータを派生参考語へ引き継ぐ。
- R-1-3: `NPC_` 以外を入力から除き、既存の `termderive` が対象外の語形と一般語用法を除く。
- R-1-4, R-1-7: 事前作成辞書の完全形と対象 plugin の翻訳済み固有名の原語を `baseSources` に入れる。
- R-1-5: `dialogueUsageForPlugin` が対象 plugin の `INFO:NAM1` だけを集計する。
- R-1-6: 派生結果は `TranslationReference` の戻り値にだけ追加する。

## 検証結果

- `go test ./internal/engine`: 成功。
- `go test ./internal/core/termderive ./internal/engine`: 成功。
- `git diff --check`: 成功。
- `npm run verify:backend`: 成功。Go backend test、architecture lint、boundary lint が成功。

## 未確認事項と停止理由

- なし

## 人間の指摘

<人間が直接記入する。メインエージェントは記入内容を削除または言い換えない。>
