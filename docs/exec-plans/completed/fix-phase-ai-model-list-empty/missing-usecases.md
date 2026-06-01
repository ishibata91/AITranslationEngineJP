# 不足UC

## 概要

- 対象: fix-phase-ai-model-list-empty の UC 差分確認
- 判断: 差分なし
- 根拠: `uc-translation-management.md` の「翻訳段階を開始する」UC に、provider 選択 → モデル一覧取得 → モデル選択 → 固定 → 開始の遷移が既に定義されている。今回の修正は `FakeProviderSettingsSecretStore` を追加して secret store 抽象化を貫徹させる実装であり、画面仕様や UC の追加・変更は不要である。UC 正本への変更は不要である。

## 不足UC一覧

| ID | 関連UC | 不足箇所 | 不足内容 | 理由 | 判断 |
| --- | --- | --- | --- | --- | --- |
