# Implementation: hyphenated-npc-name-derivation

## 変更したfile

- REF-1: ハイフンを含む姓の部分形を作る選択と、英語部分と日本語部分を対応付ける規則を追加した。
- REF-3 と REF-4: 事前作成辞書の NPC と翻訳対象 mod の NPC の派生だけで、ハイフンを含む姓の部分形を作る選択を有効にした。
- REF-7、REF-8、REF-9: 対応付け、境界、事前作成辞書、対象 mod、横断辞書の非適用を検証した。
- REF-10: 非Windowsの子プロセス設定を no-op として定義し、backend 全検証の build 条件を満たした。
- `spec.md`: 実テストとの対応を記録した。

## 仕様との対応

- R-1-1: 事前作成辞書と翻訳対象 mod の NPC で、ハイフンを含む姓の部分形を本文参考語へ加える。
- R-1-2: 複数のハイフンを含む姓へ連続する中黒部分を対応付ける。
- R-1-3: 対応する部分数が不足または過剰な場合は派生しない。
- R-1-4: 横断辞書を作る既存の人名派生では、ハイフンを含む姓の部分形を有効にしない。

## 検証結果

- `go test ./internal/core/termderive ./internal/engine`: 成功。
- `git diff --check`: 成功。
- `npm run verify:backend`: 成功。Go backend test、architecture lint、boundary lint が成功。

## 未確認事項と停止理由

- なし

## 人間の指摘

<人間が直接記入する。メインエージェントは記入内容を削除または言い換えない。>
