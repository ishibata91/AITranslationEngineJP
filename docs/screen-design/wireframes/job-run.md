# Job Run

## Intent

- 実行中ジョブの phase、status、失敗回復可否を観測する
- 中断、再開、リトライ、キャンセルを状態に応じて操作する

## ASCII Wireframe

```text
+------------------------------------------------------------------------------------------------------+
| Job Run                                                     [Pause] [Resume] [Retry] [Cancel]       |
+--------------------------------------+---------------------------------------------------------------+
| Job Summary                          | Phase Timeline                                                |
| -----------------------------------  | ------------------------------------------------------------ |
| Job Name       : ExampleMod JP v1    | Draft -> Ready -> Running -> Completed                       |
| Status         : Running             |                    ^                                         |
| Current Phase  : Body Translation    |                 current                                      |
| Started At     : 2026-03-31 09:30    |                                                            |
| Provider       : Gemini Batch        | Phase Runs                                                   |
| Provider Batch : batch_01HXYZ...     | +--------------------------------------------------------+ |
|                                      | | word_translation    Completed   09:31 -> 09:34         | |
| Control State                        | | persona_generation  Completed   09:34 -> 09:36         | |
| [pause enabled]                      | | body_translation    Running     09:36 -> --            | |
| [resume disabled]                    | +--------------------------------------------------------+ |
| [retry disabled]                     |                                                            |
| [cancel enabled]                     | Translation Progress                                         |
|                                      | +--------------------------------------------------------+ |
| Recoverable Failure Panel            | | total units        2,846                                 | |
| Last Error: none                     | | completed         1,904                                  | |
| Retry Count: 0                       | | running             128                                  | |
|                                      | | queued              814                                  | |
| Job-local Persona / Dictionary       | +--------------------------------------------------------+ |
| +----------------------------------+ |                                                            |
| | generated persona: 17            | | Selected Unit Detail                                        |
| | reused dictionary: 243           | | FORMID  000A1234                                           |
| +----------------------------------+ | | Source  <Alias=Player> Welcome to Falkreath.             |
|                                      | | Dest    ファルクリースへようこそ。                       |
|                                      | | Status  translation_status_code = 2                      |
+--------------------------------------+---------------------------------------------------------------+
| Footer: AI run hashes | provider run id | last event timestamp | manual recovery guidance            |
+------------------------------------------------------------------------------------------------------+
```

## Notes

- `JOB_PHASE_RUN`、`AI_RUN`、`JOB_TRANSLATION_UNIT` の観測項目を 1 画面に集約する
- `RecoverableFailed` 時は `Last Error` と `Retry` 操作が主表示になる前提とする
