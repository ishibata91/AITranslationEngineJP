# <YYYY-MM-DD> <task-id> run

## Placement

- `run_folder`: `work_history/runs/YYYY-MM-DD-<task-id>-run/`
- `codex_report`: `./codex.md`
- `copilot_report`: `./copilot.md`
- `cross_role_summary`: `./README.md`
- `do_not_write_to`: `docs/exec-plans/`, `.codex/history/`, handoff file

## Run Metadata

- `task_id`: `<task-id>`
- `run_date`: `<YYYY-MM-DD>`
- `related_plan`: `<path or N/A>`
- `related_handoff`: `<path or N/A>`
- `final_status`: `<completed / partial / rerouted / failed>`

## Outcome

- `結果`: `<何が終わったか>`
- `未完了`: `<残ったこと or なし>`
- `重要エラー`: `<重大な失敗 or なし>`
- `次に見るべき場所`: `<path / command / issue>`

## Timeline

- `開始`: `<時刻 or 不明>`
- `終了`: `<時刻 or 不明>`
- `時間がかかったこと`: `<一番重かった工程>`
- `待ち時間`: `<tool / review / test / user decision / なし>`
- `再作業`: `<reroute / re-run / rollback / なし>`

## Role Reports

- `Codex`: `./codex.md`
- `Copilot`: `./copilot.md`
- `Codex status`: `<completed / partial / not-run>`
- `Copilot status`: `<completed / partial / not-run>`

## Cross-Role Findings

- `改善すべきこと`: `<両役割を見て改善すべき運用>`
- `時間がかかったこと`: `<設計、handoff、実装、検証の遅延要因>`
- `無駄だったこと`: `<重複作業、不要な調査、不要な再実行>`
- `困ったこと`: `<役割境界、前提、tool、情報不足>`
- `検証で不足したこと`: `<足りなかった test / check / evidence>`

## Next Improvements

- `prompt 改善`: `<次回の依頼や指示で変えること>`
- `handoff 改善`: `<implementation-scope や完了報告で増やすこと>`
- `template 改善`: `<この template に足すべき項目 or なし>`
- `人間が次に見るべき場所`: `<path / issue / command>`

## SUMMARY

- `変更ファイル`: `<このランで変更した主要 file>`
- `重要エラー`: `<重大な失敗 or なし>`
- `次に見るべき場所`: `<path / issue / command>`
- `再実行コマンド`: `<command or なし>`
