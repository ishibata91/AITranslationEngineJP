# Job Setup

## Intent

- 翻訳ジョブを `Draft` から `Ready` に進める
- 基盤データ、AI 基盤、実行方式、対象 input file を 1 画面で確認する

## ASCII Wireframe

```text
+------------------------------------------------------------------------------------------------------+
| Job Setup                                                                    [Save Draft] [Create Job] |
+--------------------------------------+----------------------------------------------------------------+
| Job Basics                           | AI Runtime                                                     |
| -----------------------------------  | ------------------------------------------------------------- |
| Job Name       [ExampleMod JP v1__]  | Provider        ( ) LMStudio  (x) Gemini  ( ) xAI            |
| Input Files    [2 selected       ]   | Execution Mode  ( ) Single    (x) Batch API                  |
| Target Units   [2,846            ]   | Failure Recovery [enabled]                                    |
|                                      | Pause / Resume   [enabled]                                    |
| Foundation References                |                                                              |
| -----------------------------------  | Phase Configuration                                           |
| Master Persona   [Base NPC v1   v]   | +----------------------------------------------------------+ |
| Master Dictionary [TES Dict v3   v]  | | [x] Word Translation Phase                               | |
|                                      | | [x] NPC Persona Generation Phase                         | |
| Job Local Metadata                   | | [x] Body Translation Phase                               | |
| -----------------------------------  | +----------------------------------------------------------+ |
| Reuse previous job persona [off]     |                                                              |
| Reuse previous job dictionary [off]  | Prompt / Instruction Summary                                 |
|                                      | +----------------------------------------------------------+ |
| Validation                           | | record-type aware instruction template preview           | |
| +----------------------------------+ | | protected elements handling note                         | |
| | all inputs imported      [pass]  | | dictionary reuse note                                    | |
| | foundation data selected [pass]  | +----------------------------------------------------------+ |
| | output formats selected   [pass] |                                                              |
| +----------------------------------+ | Output Formats: [x] Standard package  [x] xTranslator      |
+--------------------------------------+----------------------------------------------------------------+
```

## Notes

- `AI_PROVIDER.supports_batch` の有無を前提に execution mode を切り替える
- `translation instruction` の具体文面ではなく、構成要素の確認に留める
