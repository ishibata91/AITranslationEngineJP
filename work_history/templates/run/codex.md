# Codex report

## Placement

- `run_folder`: `work_history/runs/YYYY-MM-DD-<task-id>-run/`
- `report_file`: `./codex.md`
- `cross_role_summary`: `./README.md`
- `do_not_write_to`: `docs/exec-plans/`, `.codex/history/`, handoff file

## Metadata

- `task_id`: `<task-id>`
- `run_date`: `<YYYY-MM-DD>`
- `lane`: `Codex`
- `role`: `<plan / design / handoff / docs canonicalization / other>`
- `status`: `<completed / partial / rerouted / failed / not-run>`

## Expected Role

- `期待された役割`: `<Codex が担うべきこと>`
- `対象外`: `<Codex が扱わないこと>`
- `入力`: `<user request / plan / handoff / completion report>`
- `完了条件`: `<何が満たされれば完了か>`

## Result

- `結果`: `<何が終わったか>`
- `未完了`: `<残ったこと or なし>`
- `変更ファイル`: `<docs / workflow / template / なし>`
- `重要エラー`: `<重大な失敗 or なし>`

## Time Use

- `時間がかかったこと`: `<一番重かった工程>`
- `長かった理由`: `<調査 / 判断 / tool / review / 不明>`
- `待ち時間`: `<user / tool / test / なし>`
- `短縮できること`: `<次回減らせる工程>`

## Problems

- `改善すべきこと`: `<次回の進め方で直すこと>`
- `時間がかかったこと`: `<重かった作業や判断>`
- `無駄だったこと`: `<不要だった確認、重複、遠回り>`
- `困ったこと`: `<情報不足、tool、役割境界、判断不能>`
- `前提や指示で曖昧だったこと`: `<曖昧だった入力 or なし>`

## Waste

- `重複作業`: `<同じ確認や説明の繰り返し or なし>`
- `不要な調査`: `<使わなかった調査 or なし>`
- `不要な再実行`: `<避けられた command / test / spawn or なし>`
- `削れる待ち`: `<人間確認、tool 待ち、handoff 待ち or なし>`

## Blocked Or Confused

- `困ったこと`: `<詰まった原因>`
- `再作業・reroute の原因`: `<なぜ戻ったか or なし>`
- `設計判断の詰まり`: `<判断できなかった論点 or なし>`
- `HITL の詰まり`: `<human review / approval の不足 or なし>`
- `docs 正本化判断`: `<必要 / 不要 / 未確認>`

## Validation

- `実行した確認`: `<command / review / inspection>`
- `検証で不足したこと`: `<足りなかった evidence>`
- `handoff packet の不足`: `<implementation-scope / prompt の不足 or なし>`
- `spawn や調査の必要判定`: `<適切 / 過剰 / 不足 / 未確認>`

## Improvements

- `次回の prompt 改善`: `<依頼文に足すこと>`
- `次回の handoff 改善`: `<Copilot に渡す情報の改善>`
- `次回の template 改善`: `<この template に足すこと or なし>`
- `人間が次に見るべき場所`: `<path / issue / command>`

## Follow-up

- `必要な follow-up`: `<plan / issue / docs / なし>`
- `owner`: `<human / Codex / Copilot / unknown>`
- `期限`: `<date / next run / none>`
- `再実行コマンド`: `<command or なし>`

## SUMMARY

- `変更ファイル`: `<このランで変更した主要 file>`
- `重要エラー`: `<重大な失敗 or なし>`
- `次に見るべき場所`: `<path / issue / command>`
- `再実行コマンド`: `<command or なし>`
