# Foundation Data

## Intent

- 実行前と実行後の両方で、マスターペルソナとマスター辞書を観測する
- 基盤データの build source と build time を確認する
- エントリ単位の追加、修正、無効化を行える感触を持たせる

## ASCII Wireframe

```text
+------------------------------------------------------------------------------------------------------+
| Foundation Data                                      [New Persona] [New Dictionary Term] [Rebuild] |
+--------------------------------------+----------------------------------------------------------------+
| Master Persona List                  | Entry Editor                                                   |
| -----------------------------------  | ------------------------------------------------------------- |
| Search [Whiterun Guard___________]   | Dataset        [Base Game NPC Persona                    v]    |
| Filter [source: base game      v]    | Source Type    [base_game______________________________]       |
| [x] show editable only               | Built At       [2026-03-31 09:00______________________]       |
|                                      |                                                               |
| +----------------------------------+ | Selected Persona Entry                                        |
| | * 00013BA1 Lydia    edited       | | NPC Name      [Lydia__________________________________]      |
| |   0001A696 Balgruuf clean        | | NPC FormID    [00013BA1______________________________]      |
| | ! 00013BBD Farengar draft        | | Race          [Nord_______________________________]         |
| +----------------------------------+ | | Sex           [Female_____________________________]         |
|                                      | | Voice         [FemaleCommander_____________________]         |
| Master Dictionary List               | | Persona Text                                               |
| -----------------------------------  | | +-------------------------------------------------------+ |
| Search [dragon______________]        | | | Formal, reliable, battle-hardened tone...            | |
| [Import CSV] [Bulk Replace]          | | +-------------------------------------------------------+ |
| +----------------------------------+ | | [Save Entry] [Duplicate] [Disable] [Discard Changes]    |
| | * dragon priest  ドラゴン...    | |                                                               |
| |   Thu'um         シャウト        | | Draft Changes                                               |
| | ! hold position  その場を守れ     | | +-------------------------------------------------------+ |
| +----------------------------------+ | | | source_text      old: hold position                  | |
|                                      | | | dest_text        new: その場を守れ                  | |
|                                      | | +-------------------------------------------------------+ |
+--------------------------------------+----------------------------------------------------------------+
| Status: persona 12,842 | dictionary 8,304 | dirty entries 3 | job-local persona is not shown here     |
+------------------------------------------------------------------------------------------------------+
```

## Notes

- `MASTER_PERSONA` と `MASTER_DICTIONARY` を別 pane で同時観測する
- job 単位で生成された `JOB_PERSONA_ENTRY` は別画面で扱う
- 右 pane は read-only detail ではなく entry editor として扱う
