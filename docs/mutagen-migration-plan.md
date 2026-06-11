# Mutagen 抽出ツール 移行計画（作業引き継ぎ資料）

本書は、翻訳対象抽出の基盤を xEdit Pascal から Mutagen（C#/.NET）へ移す作業を、別 PC・別セッションで再開するための自己完結資料である。前提・確定済みの設計判断・やること・検証方法を記す。

## 1. 目的

- AITranslationEngineJp は Skyrim Mod 向け翻訳エンジンである。
- plugin（esm / esp）から、翻訳対象テキストと文脈（話者・条件・関連）を抽出し JSON 化する。
- その抽出基盤を、xEdit Pascal（JvInterpreter）から Mutagen（C#/.NET）へ移す。

## 2. なぜ Mutagen へ移すか

- xEdit Pascal は Skyrim.esm 抽出に 18 分かかる。JvInterpreter の制約（forward 宣言不可、文字列 index 操作不可、要素名の癖）が多く、デバッグも難しい。
- Mutagen は全 load order（CC content 含む）を 1.6 秒で読めた（素振り実測）。型安全、テスト可能、load order 自動検出、localized string 自動解決。
- 抽出は前処理でリアルタイム性が要らない。Go / Wails 本体とは別プロセスにし、JSON 出力で繋げばよい。言語が分かれても境界は JSON。

## 3. 現状（移行前の到達点）

- xEdit Pascal 版 `extractData.v2.pas` が完成済み。Dawnguard.esm で xTranslator XML と **件数 +0 一致**、`validate_extraction.py` **0 error** を達成した。
- 「何を抽出するか・何を除外するか・どの経路で取るか」は Dawnguard で確定済みである。これが移植の仕様になる。Pascal のロジックをそのまま C# に写せる。
- Mutagen 素振りは `tools/extractor/`（C# console, net8.0, Mutagen.Bethesda.Skyrim 0.53.1）に作成済み。`GameEnvironment.Typical.Skyrim(SkyrimRelease.SkyrimSE)` で自動検出し、BOOK / NPC_ / WEAP を 1.6 秒で列挙できることを確認した。

## 4. 確定した設計判断（移植で必ず守る）

1. **plugin 単位列挙**: `env.LoadOrder[targetModKey].Mod.<Records>` で、対象 plugin がそのファイルに書いた版（新規 + override）を出す。**`WinningOverrides()` は使わない**。xTranslator は plugin 単位で翻訳するため、load order 全体で勝った版ではなく、その plugin が書いた版が要る。
2. **source = その plugin**: plugin 単位列挙そのものが source の意味になる。純 master（対象 plugin が触っていない record）は `Mod.<Records>` に現れないので自然に除外される（Pascal 版の EnsureガードA と同じ効果）。
3. **参照解決**: FormLink は `env.LinkCache` で解決する。ただし出力するのは ID（FormID、版非依存）と、対象 plugin が含む本体のみ。参照先が純 master なら ID だけ残し、本体は出さない（翻訳時に元 plugin の辞書で解決する想定）。
4. **概念モデル v19 に従う**: 正本は [conceptual-model.md](conceptual-model.md)。会話は Dialogue → InfoNode（INFO）→ ResponseLine（NAM1）の 2 階層。Player は固定モデルで抽出対象外（ペルソナはユーザーが翻訳設定として与える）。話者は InfoNode 単位で解決する。
5. **JSON schema**: `scripts/validate_extraction.py` が schema を、`scripts/compare_counts.py` が件数を検証する。**まず現行 schema 互換で移植し、compare で Dawnguard +0 を再現**してから、v19 の 2 階層化（InfoNode / ResponseLine 分離）に進む。現行 Pascal 版は INFO 単位属性を NAM1 行に複製しているので、互換移植の次に 2 階層へ直す。
6. **翻訳対象 / 非対象** は §6・§7 のとおり。

## 5. 検証方法（移植の正しさはこれで担保する）

- 対象 plugin を抽出 → JSON 保存（`dictionaries/<plugin>_Export.json`）。
- 件数比較: `python scripts/compare_counts.py --json dictionaries/<plugin>_Export.json --xml dictionaries/<plugin>_english_japanese.xml`。差分 +0 を目標。
- schema: `python scripts/validate_extraction.py dictionaries/<plugin>_Export.json`。0 error を目標。
- まず Dawnguard で +0 を再現できれば、移植の最初の関門を通過。次に Skyrim.esm、mod の順。
- `compare_counts.py` の要点（移植先でも同じ前提を守る）:
    - XML 側は **Source 非空行のみ**数える（空行は翻訳対象でない）。
    - INFO:RNAM と WOOP:TNAM は **id 単位でユニーク化**（複製・複数参照の重複を排除）。2 階層化後はこの暫定処理を撤去できる。
    - JSON 側は **source == 対象 plugin** で絞る。

## 6. 抽出する翻訳対象（REC:FIELD と対応 record）

正は `scripts/compare_counts.py` の `flatten_json` と `validate_extraction.py` の `REQUIRED_FIELDS`。Mutagen のプロパティ名は移植時に cheat sheet / IntelliSense で確認する（下表の「Mutagen record」は型の対応、プロパティは要確認）。

| 翻訳対象 | Mutagen record | xTranslator REC:FIELD | 備考 |
|---|---|---|---|
| 会話トピック | DialogTopic (DIAL) | DIAL:FULL | プレイヤー選択肢の既定文。話者は Player |
| 応答ノードの選択肢上書き | DialogResponses (INFO) RNAM | INFO:RNAM | プレイヤー選択肢の条件別上書き。話者は Player |
| NPC 応答行 | DialogResponses (INFO) Responses | INFO:NAM1 | NPC 返答。話者は InfoNode が解決 |
| Quest 名 | Quest (QUST) | QUST:FULL | |
| Quest ログ | Quest stages CNAM | QUST:CNAM | |
| Quest 目標 | Quest objectives NNAM | QUST:NNAM | |
| 種族名・説明 | Race (RACE) | RACE:FULL / RACE:DESC | |
| 派閥名 | Faction (FACT) | FACT:FULL | |
| 派閥階級称号 | Faction ranks | FACT:MNAM / FACT:FNAM | 男 / 女。**Skyrim で要追加**（Pascal 未対応） |
| アイテム名 | MiscItem/Key/Light/Container/SoulGem/Door/Furniture | <SIG>:FULL | KEYM/MISC/LIGH/CONT/SLGM/DOOR/FURN |
| 錬金器具名 | Apparatus (APPA) | APPA:FULL | **Skyrim で要追加** |
| Activator | Activator/Flora/Tree | <SIG>:FULL / <SIG>:RNAM | ACTI/FLOR/TREE 統合。RNAM=起動動作 |
| 装備 | Weapon/Armor/Ammunition | <SIG>:FULL / <SIG>:DESC | WEAP/ARMO/AMMO |
| 消耗品 | Ingestible/Scroll/Ingredient | <SIG>:FULL / <SIG>:DESC | ALCH/SCRL/INGR |
| 呪文 | Spell (SPEL) | SPEL:FULL / SPEL:DESC | |
| 付呪 | ObjectEffect (ENCH) | ENCH:FULL | |
| 魔法効果 | MagicEffect (MGEF) | MGEF:FULL / MGEF:DNAM | DNAM=効果説明 |
| 叫び | Shout (SHOU) | SHOU:FULL / SHOU:DESC | |
| 力の言葉 | WordOfPower (WOOP) | WOOP:TNAM | TNAM=意味の訳が対象。**FULL=龍語綴りは翻訳禁止** |
| 書物 | Book (BOOK) | BOOK:FULL / BOOK:DESC / BOOK:CNAM | CNAM は SSEEdit 表示名「Description」。多くは空 |
| 場所 | Location/Worldspace/Cell | <SIG>:FULL | LCTN/WRLD/CELL |
| メッセージ | Message (MESG) | MESG:FULL / MESG:DESC / MESG:ITXT | ITXT=選択肢ボタン |
| ロード画面 | LoadScreen (LSCR) | LSCR:DESC | |
| Perk | Perk (PERK) | PERK:FULL / PERK:DESC | |
| 声型 | VoiceType (VTYP) | VTYP:EDID | 翻訳テキスト無し。識別子のみ（話者推定の手掛かり） |
| 音量カテゴリ | SoundCategory (SNCT) | SNCT:FULL | **Skyrim で要追加**（設定 UI 表示） |
| 目の種類 | Eyes (EYES) | EYES:FULL | **Skyrim で要追加**（character creation） |
| 地域マップ名 | Region (REGN) | REGN:RDMP | **Skyrim で要追加**（world map） |

### 名前なしレコードの扱い（Pascal 版の重要修正）

FULL（名前）が空でも sub-field を持つ record は落とさない。Pascal 版で MESG / QUST / PERK がこれで取りこぼしていた。移植先でも「翻訳対象テキストが 1 つでもあれば出す」を徹底する。

### 話者解決（INFO 単位、r2a〜d）

INFO の Conditions（CTDA）から導出する。ANAM-Speaker はほぼ使われない（実測）。

- 名指し（r2a → speaker_ids）: `GetIsID`（Base Object が NPC_/TACT）、`GetIsAliasRef`（quest alias を ALFR/ALUA で解決）。
- 集合（r2b → faction_ids）: `GetInFaction`。
- 集合（r2c → race_ids）: `GetIsRace`。
- 集合（r2d → voice_type_ids）: `GetIsVoiceType`。
- 前後関係（r7）: PNAM（previous、単一）、TCLT（next、複数）。FormID で持つ。

Speaker（NPC_/TACT）のペルソナ名解決順: (1) 自身の FULL、(2) TPLT 連鎖の本体 FULL、(3) EditorID 由来の役割名（プール / 声型名）。

## 7. 非抽出対象（翻訳不要・除外）

`compare_counts.py` の `NON_TARGET_REC_BASES` / `NON_TARGET_REC_FIELDS` が正。

- record 種別ごと除外: HDPT, PROJ, BPTD, EXPL, CLFM, CLAS, AVIF, GMST, REFR, COLL（collision、内部）, WATR（water、内部）, `****`（script fragment）。
- field 単位除外: PERK:EPF2, PERK:EPFD, WOOP:FULL（龍語綴り、翻訳禁止）。
- 要表示確認（保留）: HAZD:FULL（hazard。通常クロスヘアに出ない推測）、REFR:FULL（配置参照の name override）。

## 8. やること（段階）

1. **plugin 単位抽出の枠組み**: ModKey を引数に、`env.LoadOrder[ModKey].Mod` の各 record を列挙し、参照は `env.LinkCache` で解決、JSON を出力する CLI を作る。
2. **現行 schema 互換で移植**: 単純カテゴリ（Book/Race/Faction/Item/Activator/Equipment/Consumable/Magic/MagicEffect/Enchantment/Shout+Word/Message/LoadingScreen/Perk/Location/VoiceType）→ Speaker（NPC_/TACT）→ Dialogue/INFO + 話者解決 + 条件 + PNAM/TCLT。
3. **Dawnguard で compare +0・validate 0 error を再現**（移植の正しさを既存テストで担保）。
4. **Skyrim 追加 record 対応**: APPA / FACT:MNAM / FACT:FNAM / SNCT / EYES / REGN を追加。COLL / WATR は非対象に。
5. **Skyrim.esm と mod で compare +0 を確認**。
6. **v19 の 2 階層化**: InfoNode / ResponseLine を分離（INFO 単位属性を NAM1 行に複製しない）。compare の id ユニーク化（暫定）を撤去。schema を 2 階層へ更新。

## 9. 環境・前提

- dotnet: 9.0.205 / 8.0.422（SDK 両方あり）。プロジェクトは net8.0。
- NuGet: `Mutagen.Bethesda.Skyrim` 0.53.1。
- プロジェクト位置: `tools/extractor/`（C# console）。
- Skyrim SE は `GameEnvironment.Typical.Skyrim(SkyrimRelease.SkyrimSE)` で自動検出（registry / plugins.txt、CC content 含む）。素振り実測 1.6 秒。
- 最小 API: `env.LoadOrder.PriorityOrder.<Record>().WinningOverrides()`（素振り用。**本番は使わない**）／ 本番は `env.LoadOrder[ModKey].Mod.<Records>`。
- localized string は `.Name?.String` 等で自動解決。

## 10. 参照ファイル

- [conceptual-model.md](conceptual-model.md): v19 概念モデル（抽出仕様の正本）。
- `extractData.v2.pas`: xEdit Pascal 版。移植元ロジック（話者解決・条件・emit 条件・FULL 空対応の実装が全部入っている）。
- `scripts/compare_counts.py`: JSON vs xTranslator XML の件数比較（翻訳対象・非対象・source フィルタ・ユニーク化の正）。
- `scripts/validate_extraction.py`: JSON schema 検証（REQUIRED_FIELDS の正）。
- `dictionaries/`: 対象 esm の `*_Export.json`（抽出結果）と `*_english_japanese.xml`（xTranslator 辞書、照合先）。
- `tools/extractor/`: Mutagen 素振り（Program.cs）。

## 11. 既知の検証実績（Pascal 版、移植先で再現すべき数値）

- Dawnguard.esm: 51 REC:FIELD 全一致、差分 +0、validate 0 error（warning 8 は dangling = 純 master 参照先を出さない、意図通り）。
- Skyrim.esm: Pascal 版で −204（§6 の Skyrim 追加 record 未対応分）。移植 + 追加対応で +0 を目指す。
