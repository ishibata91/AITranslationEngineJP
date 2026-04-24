# Copilot report

## Placement

- `run_folder`: `work_history/runs/YYYY-MM-DD-<task-id>-run/`
- `report_file`: `./copilot.md`
- `cross_role_summary`: `./README.md`
- `do_not_write_to`: `docs/exec-plans/`, `.codex/history/`, handoff file

## Metadata

- `task_id`: `<task-id>`
- `run_date`: `<YYYY-MM-DD>`
- `lane`: `Copilot`
- `role`: `<implementation / test / review / investigate / other>`
- `status`: `<completed / partial / rerouted / failed / not-run>`

## Completion Packet Skeleton

- `copilot_work_report.report_path`: `work_history/runs/YYYY-MM-DD-<task-id>-run/copilot.md`
- `copilot_work_report.status`: `<completed / partial / rerouted / failed / not-run>`
- `copilot_work_report.改善すべきこと`: `<次回の進め方で直すこと>`
- `copilot_work_report.時間がかかったこと`: `<重かった実装、調査、検証>`
- `copilot_work_report.無駄だったこと`: `<不要だった確認、重複、遠回り>`
- `copilot_work_report.困ったこと`: `<情報不足、tool、scope、テスト失敗>`
- `copilot_work_report.次に見るべき場所`: `<path / issue / command>`

## Expected Role

- `期待された役割`: `<Copilot が担うべきこと>`
- `対象外`: `<Copilot が扱わないこと>`
- `入力`: `<implementation-scope / scenario / source artifact>`
- `完了条件`: `<何が満たされれば完了か>`

## Result

- `結果`: `<何が終わったか>`
- `未完了`: `<残ったこと or なし>`
- `触ったファイル`: `<product code / product test / なし>`
- `重要エラー`: `<重大な失敗 or なし>`

## Time Use

- `時間がかかったこと`: `<一番重かった工程>`
- `長かった理由`: `<実装 / 調査 / test / review / 不明>`
- `待ち時間`: `<tool / test / user / なし>`
- `短縮できること`: `<次回減らせる工程>`

## Problems

- `改善すべきこと`: `<次回の進め方で直すこと>`
- `時間がかかったこと`: `<重かった実装、調査、検証>`
- `無駄だったこと`: `<不要だった確認、重複、遠回り>`
- `困ったこと`: `<情報不足、tool、scope、テスト失敗>`
- `前提や指示で曖昧だったこと`: `<曖昧だった入力 or なし>`

## Waste

- `重複作業`: `<同じ確認や実装の繰り返し or なし>`
- `不要な調査`: `<使わなかった調査 or なし>`
- `不要な再実行`: `<避けられた command / test / agent run or なし>`
- `削れる待ち`: `<test 待ち、review 待ち、handoff 待ち or なし>`

## Blocked Or Confused

- `困ったこと`: `<詰まった原因>`
- `再作業・reroute の原因`: `<なぜ戻ったか or なし>`
- `implementation-scope の読み取り`: `<明確 / 曖昧 / 不足>`
- `実装分割の詰まり`: `<分割過大 / 依存不明 / なし>`
- `完了報告の不足`: `<Codex が close 判断に困る不足 or なし>`

## Validation

- `実行した確認`: `<command / test / review>`
- `検証で不足したこと`: `<足りなかった test / evidence>`
- `調査`: `<reproduce / trace / reobserve / なし>`
- `review`: `<implementation review / UI check / なし>`
- `reroute`: `<発生理由 or なし>`

## Improvements

- `次回の prompt 改善`: `<依頼文に足すこと>`
- `次回の handoff 改善`: `<implementation-scope に足すこと>`
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
