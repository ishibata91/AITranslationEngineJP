# 残作業ロードマップ

本書は、業務要件を満たすために残っている大きな作業項目を列挙し、着手順を記録する正本である。
各項目の技術的な現状・未解決事実は [`known-issues.md`](./known-issues.md) が正本であり、本書は項目の列挙と着手順の記録に限定する。

## 残作業項目

1. 固有名一貫性の後続 task（残りの言及関連: 会話の流れ e7、名乗る名 e8/e14、種族名 FK e13）: 詳細は [`known-issues.md`](./known-issues.md) 1
2. 対話木の context 扱い（Dialogue tree の文脈利用）: 詳細は [`known-issues.md`](./known-issues.md) 2
3. 固有名一貫性の事後検証（注入語の保持確認）: 詳細は [`known-issues.md`](./known-issues.md) 2
4. クラウド AI プロバイダの未実装（Gemini・xAI・Claude）: 詳細は [`known-issues.md`](./known-issues.md) 3
5. xTranslator 形式への書き出し未実装: 詳細は [`known-issues.md`](./known-issues.md) 4
6. 翻訳結果表示画面の編集・絞り込み機能未実装: 詳細は [`known-issues.md`](./known-issues.md) 5
7. 機械置換辞書の誤爆対策の残り（stoplist 外の一般語 1 語、管理用勢力の判定基準）: 詳細は [`known-issues.md`](./known-issues.md) 6
8. 既存訳との完全一致置換の未実装: 詳細は [`known-issues.md`](./known-issues.md) 7
9. 本文中の実行時タグ（`<Alias=...>` 等）を保護する機構の未実装: 詳細は [`known-issues.md`](./known-issues.md) 8

## 着手順

未定。各項目の着手時に人間と相談して決める。

## 更新方針

項目が完了したら本書から除き、経緯は [`changelog.md`](./changelog.md) に残す。新たな残作業が判明したら本書へ追記する。
