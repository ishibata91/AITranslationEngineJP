# Output Review

## Intent

- 完了ジョブの翻訳結果と出力成果物を観測する
- 現時点で対応している XML 出力を確認する
- 完了後の手修正、再出力、差分確認ができる感触を持たせる

## ASCII Wireframe

```text
+------------------------------------------------------------------------------------------------------+
| Output Review                                   [Open Folder] [Save Edits] [Regenerate Selected]    |
+--------------------------------------+----------------------------------------------------------------+
| Completed Job List                   | Result Editor                                                  |
| -----------------------------------  | ------------------------------------------------------------- |
| Search [ExampleMod___________]       | Job Name      : ExampleMod JP v1                               |
| +----------------------------------+ | Job Status    : Completed                                      |
| | ExampleMod JP v1  09:58          | | Finished At   : 2026-03-31 09:58                             |
| | PatchAddon JP v2  10:12          | |                                                              |
| +----------------------------------+ | XML Artifact                                                   |
|                                      | +----------------------------------------------------------+ |
| Result Summary                       | | format_code : xtranslator_xml                            | |
| total units      [2,846        ]     | | status      : generated                                  | |
| completed units  [2,846        ]     | | file_path   : dist/ExampleMod_ja.xml                     | |
| failed units     [0            ]     | +----------------------------------------------------------+ |
| preserved tags   [118          ]     | [Download XML] [Replace File] [Regenerate XML]               |
| edited after run [14           ]     | Selected Translation Unit                                      |
| xml ready        [yes          ]     |                                                              |
|                                      | EDID          [ExampleDialogueLine01____________________]     |
| Filters                              | REC / FIELD   [INFO] [response.text____________________]     |
| Record Type [dialogue         v]     | FORMID        [000A1234________________________________]     |
| Status Code [all              v]     | Source                                                       |
| Search Text [Falkreath_______]       | +----------------------------------------------------------+ |
| [show edited only]                   | | <Alias=Player> Welcome to Falkreath.                    | |
| +----------------------------------+ | +----------------------------------------------------------+ |
| | * edited rows list               | | Dest                                                         |
| |   translated rows list           | | +----------------------------------------------------------+ |
| +----------------------------------+ | | | ファルクリースへようこそ。                             | |
|                                      | | +----------------------------------------------------------+ |
|                                      | | [Apply Terminology] [Mark Needs Review] [Reset To AI]      |
|                                      | |                                                              |
|                                      | | Diff Preview                                                 |
|                                      | | - Welcome to Falkreath.                                     |
|                                      | | + ファルクリースへようこそ。                               |
+--------------------------------------+----------------------------------------------------------------+
| Footer: export warnings | missing outputs | completed job cache cleanup candidate                       |
+------------------------------------------------------------------------------------------------------+
```

## Notes

- 現時点の wireframe では `JOB_OUTPUT_ARTIFACT` のうち XML だけを主表示にする
- 翻訳 preview は xTranslator 再構成に必要な列を意識して配置する
- 完了後の修正は translation unit 単位で保持し、artifact は再生成できる前提の wireframe にする
- 標準配布形式など他形式は planned scope として、この画面では未表示にする
