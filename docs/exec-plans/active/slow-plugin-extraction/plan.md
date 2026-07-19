# Task Plan: slow-plugin-extraction

`plan.md` は branch 情報と、この task でやること・やらないことの要点を持つ。
設計判断、判断履歴、検証結果、実装結果は持たない。設計は `design.md`、恒久的に残す判断は `docs/changelog.md` に書く。

## やること

- 大きめ plugin（例 Outfit Recognition Framework）の抽出フェーズが遅いという申告を調査し、翻訳前区間のパフォーマンスを改善する。
- 翻訳前区間（C# 抽出子 + Go 後段）のどこが律速かを実測で確定する。
- 起動毎のビルド評価と初回のビルド/restore 崖を消すため、C# 抽出子を publish 済み DLL 直実行へ切り替える。
- 翻訳前区間の無音区間（進捗イベントが出ない区間）に進捗表示を出し、遅いときに固まって見えないようにする。
- master 多数 mod への備えとして、C# 抽出子の `OwnsRecord` が呼ぶ `Normalize` の再計算を memoize する（抽出結果は変えない予防的効率化）。

## branch 情報

- `execution_branch`: `claude/slow-plugin-extraction`
- `target_branch`: `master`
- `source_commit`: `9c993d09`

## やらないこと

- 翻訳フェーズ（LLM 呼び出し）の性能。遅いのは翻訳前区間の申告のため。
- 抽出する行の増減（抽出結果の正しさ変更）。性能だけを対象にする。
- available データ範囲外（master 依存が数十件の mod、巨大 localized strings）の網羅ベンチ。手元に該当データが無いため、非線形リスクは静的解析で指摘し memoize で予防する範囲にとどめる。
- 配布 app での C# 抽出子の同梱（self-contained publish、`wails build` 成果への同梱、配布フロー新設）。規模が大きいため `docs/known-issues.md` 課題6 に残し、`feature-workflow` の別 task で扱う。
