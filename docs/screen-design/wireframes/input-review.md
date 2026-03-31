# Input Review

## Intent

- 複数 input file を 1 ジョブに束ねる前に内容確認する
- `PLUGIN_EXPORT` と、そこから生成される翻訳対象の量を観測する

## ASCII Wireframe

```text
+------------------------------------------------------------------------------------------------------+
| Input Review                                                             [Import JSON] [Remove]      |
+--------------------------------------+---------------------------------------------------------------+
| Imported Files                       | Record Preview                                                |
| -----------------------------------  | ------------------------------------------------------------ |
| +----------------------------------+ | Plugin        : ExampleMod.esp                               |
| | ExampleMod.esp   imported 09:12  | | Source JSON   : F:\imports\ExampleMod.json                 |
| | PatchAddon.esp   imported 09:14  | | Imported At   : 2026-03-31 09:14                           |
| +----------------------------------+ |                                                            |
|                                      | Category Counts                                             |
| Translation Scope Summary            | +--------------------------------------------------------+ |
| -----------------------------------  | | dialogue_groups  124   responses      611              | |
| dialogue groups       [124      ]    | | quests            19   objectives      88              | |
| translation units     [2,846    ]    | | items             73   magic           14              | |
| quests/messages       [31       ]    | | locations         22   system/messages 57              | |
| unique NPCs           [45       ]    | +--------------------------------------------------------+ |
|                                      |                                                            |
| Filter                               | Sample Translation Unit                                     |
| Record Type [dialogue       v]       | +--------------------------------------------------------+ |
| Field Name  [text           v]       | | REC    : INFO                                           | |
| Search      [EditorID_______]        | | FIELD  : response.text                                  | |
|                                      | | FORMID : 000A1234                                       | |
|                                      | | EDID   : ExampleDialogueLine01                          | |
|                                      | | Source : <Alias=Player> Welcome to Falkreath.          | |
|                                      | +--------------------------------------------------------+ |
+--------------------------------------+---------------------------------------------------------------+
| Footer: invalid import warnings | orphan quest refs | cache rebuild eligibility                     |
+------------------------------------------------------------------------------------------------------+
```

## Notes

- `TRANSLATION_UNIT` の sample preview を置き、lossless 出力に必要な項目を見せる
- record preview は確定 UI ではなく、入力確認時に必要な最小情報を示す
