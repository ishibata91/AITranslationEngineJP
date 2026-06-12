# 単語翻訳フェーズの単語対象レコード候補

関連文書: [`./index.md`](./index.md), [`../architecture.md`](../architecture.md)

この文書は、単語翻訳フェーズの単語対象レコードと周辺候補を整理する。
単語対象は、過去プロジェクト由来の一覧を根拠に固定する。
辞書項目数と現行抽出有無は、`dictionaries/dawnguard_english_japanese.xml` と `tools/extractor` を根拠にする。

## 調査対象

- `dictionaries/dawnguard_english_japanese.xml` は、`REC` に辞書側のレコード種別とフィールド名を持つ。
- `tools/extractor` は、出力 JSON の `type` に抽出側のレコード種別とフィールド名を持つ。
- 辞書側の `REC` は `NPC_:FULL` の形式で表す。
- 抽出側の `type` は `NPC_ FULL` の形式で出力されるが、この文書では比較のため `NPC_:FULL` の形式へ読み替える。

## 単語対象

単語翻訳フェーズと XML 辞書取り込みが対象とする REC は、以下の 13 種別で統一されている。
両フェーズは同一の REC 判定（`IsTermTarget`）を共有する。

| `REC` | 説明 | カテゴリ | 現行抽出 |
| --- | --- | --- | --- |
| `BOOK:FULL` | 本のアイテム名 | 書籍 | あり |
| `NPC_:FULL` | NPC 名 | NPC | あり |
| `NPC_:SHRT` | NPC の短い名前 | NPC | あり |
| `RACE:FULL` | 種族名 | NPC | あり |
| `ARMO:FULL` | 防具名 | 装備 | あり |
| `WEAP:FULL` | 武器名 | 装備 | あり |
| `LCTN:FULL` | ロケーション名 | 地名 | あり |
| `CELL:FULL` | セル名 | 地名 | あり |
| `CONT:FULL` | コンテナ名 | アイテム | あり |
| `MISC:FULL` | その他アイテム名 | アイテム | あり |
| `INGR:FULL` | 錬金素材名 | アイテム | あり |
| `ALCH:FULL` | 食料・ポーション等の錬金術アイテム名 | アイテム | あり |
| `SHOU:FULL` | シャウト名 | シャウト | あり |

`NPC_:FULL` と `NPC_:SHRT` は別 REC として扱い、同一原語でも別候補として識別できる。

### 対象外となった REC

過去に XML 辞書取り込みの許可リストに含まれていた `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` は、単語翻訳フェーズの対象外である。
XML 辞書取り込みでも同様に対象外である。
単語翻訳フェーズ対象 REC と XML 辞書取り込み対象 REC は同一の 13 種別であり、これら 3 種別は両集合から除外されている。

## 現行抽出上の翻訳対象

`tools/extractor` に出力処理があり、単語対象 13 種別に含まれない種別は翻訳対象とする。
`SLGM:FULL` は `tools/extractor` に出力処理があるが、`dictionaries/dawnguard_english_japanese.xml` には実例がない。

| `REC` | 抽出元 | 扱い |
| --- | --- | --- |
| `AMMO:FULL` | `ExtractItem` | 翻訳対象 |
| `BOOK:DESC` | `ExtractItem` | 翻訳対象 |
| `DIAL:FULL` | `ExtractDialogue` | 翻訳対象 |
| `DOOR:FULL` | `ExtractNamedRecord` | 翻訳対象。単語翻訳フェーズと XML 辞書取り込みの対象外。 |
| `ENCH:FULL` | `ExtractMagic` | 翻訳対象 |
| `FLOR:FULL` | `ExtractNamedRecord` | 翻訳対象。単語翻訳フェーズと XML 辞書取り込みの対象外。 |
| `FURN:FULL` | `ExtractNamedRecord` | 翻訳対象。単語翻訳フェーズと XML 辞書取り込みの対象外。 |
| `INFO:NAM1` | `ExtractInfo` | 翻訳対象 |
| `KEYM:FULL` | `ExtractItem` | 翻訳対象 |
| `LIGH:FULL` | `ExtractItem` | 翻訳対象 |
| `LSCR:DESC` | `ExtractLoadScreen` | 翻訳対象 |
| `MESG:DESC` | `ExtractMessage` | 翻訳対象 |
| `MGEF:FULL` | `ExtractMagic` | 翻訳対象 |
| `PERK:FULL` | `ExtractPerk` | 翻訳対象 |
| `QUST:CNAM` | `ExtractQuest` | 翻訳対象 |
| `QUST:FULL` | `ExtractQuest` | 翻訳対象 |
| `QUST:NNAM` | `ExtractQuest` | 翻訳対象 |
| `SCRL:FULL` | `ExtractMagic` | 翻訳対象 |
| `SLGM:FULL` | `ExtractItem` | 翻訳対象。辞書実例なし。 |
| `SPEL:FULL` | `ExtractMagic` | 翻訳対象 |
| `WRLD:FULL` | `ExtractLocation` | 翻訳対象 |

## 集計結果

| 区分 | 種別数 | 辞書項目数 | 意味 |
| --- | ---: | ---: | --- |
| 単語対象かつ現行抽出あり | 13 | 891 | 現行入力データから単語対象にできる。 |
| 単語対象かつ抽出追加が必要 | 0 | 0 | 単語対象だが、現行入力データには同じ種別として出ない。 |
| 翻訳対象かつ辞書実例あり | 20 | 6466 | 現行入力データに出る翻訳対象である。 |
| 翻訳対象かつ辞書実例なし | 1 | 0 | 現行入力データに出るが、調査対象の辞書 XML には実例がない翻訳対象である。 |
| 単語対象外かつ現行抽出なし | 33 | 825 | 辞書には存在するが、単語対象ではない。 |

## 判断基準

- 過去プロジェクト由来の一覧から、辞書にする価値がある 13 種別を単語対象とする。
- `NPC_:SHRT` と `RACE:FULL` は、名称の再利用価値があるため単語対象とする。
- `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` は、単語翻訳フェーズと XML 辞書取り込みの対象外とする。翻訳入力に出現する場合は後続の本文翻訳フェーズの対象範囲に委ねる。
- `tools/extractor` に出力処理があり、単語対象 13 種別に含まれない種別は翻訳対象とする。
- `FULL` は、表示名、場所名、アイテム名、種族名、シャウト名を含む場合でも、単語対象一覧にない種別は単語対象外とする。
- `SHRT` は、`NPC_:SHRT` だけを単語対象とする。
- `DESC`、`NAM1`、`CNAM`、`NNAM` は、説明文、会話本文、クエスト本文を含むため本文対象候補とする。
- `RNAM`、`ITXT`、`TNAM`、`EPF2`、`EPFD` は、操作語または短い文言を含む可能性があるが、単語対象一覧にないため単語対象外とする。
- `DIAL:FULL` は会話選択肢を含むため、翻訳対象として扱う。
- `MESG:FULL` はメッセージの題名を含むが、単語対象一覧にないため単語対象外とする。

## 候補一覧

| `REC` | 辞書項目数 | 現行抽出 | 候補判定 | 辞書実例 |
| --- | ---: | --- | --- | --- |
| `****:****` | 89 | なし | 本文対象候補。単語対象からは分離する候補である。 | `00000000|00000000|0|0` / `Alchemy laboratory. (<Global=HDMarkarthAlchemy> gold)` |
| `ACTI:FULL` | 63 | なし | 単語対象外。 | `DLC1TrapCrossbow` / `Crossbow Mount` |
| `ACTI:RNAM` | 39 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1_WESC08AshPile` / `Search` |
| `ALCH:FULL` | 4 | あり | 単語対象。現行抽出がある。 | `DLC1RedwaterDenSkooma` / `Redwater Skooma` |
| `AMMO:DESC` | 8 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1ElvenArrowBlessed` / `Causes sunburst attacks to nearby targets if shot at the sun with Auriel's Bow.` |
| `AMMO:FULL` | 24 | あり | 翻訳対象。 | `DLC1BoltSteel` / `Steel Bolt` |
| `ARMO:DESC` | 28 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1VampireNightPowerNecklace3` / `While wearing this amulet, your Night Power while in Vampire Lord Form will be Bats.` |
| `ARMO:FULL` | 121 | あり | 単語対象。現行抽出がある。 | `DLC1ArmorAurielsShield` / `Auriel's Shield` |
| `AVIF:DESC` | 2 | なし | 本文対象候補。単語対象からは分離する候補である。 | `AVHealRatePowerMod` / `New perks become available after feeding on enough corpses.` |
| `AVIF:FULL` | 2 | なし | 単語対象外。 | `AVMagickaRateMod` / `Vampire Lord` |
| `BOOK:CNAM` | 10 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1VQSaintPage01` / 空文字 |
| `BOOK:DESC` | 65 | あり | 翻訳対象。 | `DLC1_WESC08Orders` / `<font face='$HandwrittenFont' size='16'>` |
| `BOOK:FULL` | 65 | あり | 単語対象。現行抽出がある。 | `DLC1_WESC06Orders` / `Dawnguard Orders - Hakar` |
| `BPTD:BPTN` | 11 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DragonBodyPartData` / `Saddle` |
| `CELL:FULL` | 98 | あり | 単語対象。現行抽出がある。 | `VampireDungeon01` / `Quick Test Cell` |
| `CLAS:FULL` | 5 | なし | 単語対象外。 | `DLC1_EncClassChaurusFlyer` / `ChaurusFlyer` |
| `CLFM:FULL` | 1 | なし | 単語対象外。 | `HairColor16AlbinoWhite` / `Albino` |
| `CONT:FULL` | 34 | あり | 単語対象。現行抽出がある。 | `DLC1_WESC_DawnguardCache` / `Dawnguard Cache` |
| `DIAL:FULL` | 940 | あり | 翻訳対象。 | `DLC01VQ01TolanOffersQuest3` / `Any idea why the vampires attacked the Vigilants?` |
| `DOOR:FULL` | 30 | あり | 単語対象外。XML 辞書取り込みも対象外。 | `CasExFreeLgDoor01` / `Door` |
| `ENCH:FULL` | 28 | あり | 翻訳対象。 | `DLC1EnchArmorReflectingShield` / `Auriel's Shield Knockback` |
| `EXPL:FULL` | 6 | なし | 単語対象外。 | `DLC1AurielsBowExp01` / `DLC1AurielsBow Exp01` |
| `FACT:FULL` | 17 | なし | 単語対象外。 | `DLC1VampireIntroEnemyFaction` / `DLC1 Vampire Intro Quest Enemy Faction` |
| `FLOR:FULL` | 3 | あり | 単語対象外。XML 辞書取り込みも対象外。 | `DLC01Gleamblossom01` / `Gleamblossom` |
| `FLOR:RNAM` | 3 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC01Gleamblossom01` / `Harvest` |
| `FURN:FULL` | 23 | あり | 単語対象外。XML 辞書取り込みも対象外。 | `DLC1_SnowElfThrone` / `Throne` |
| `GMST:DATA` | 5 | なし | 本文対象候補。単語対象からは分離する候補である。 | `sRSMFinishedWarning` / `Finished?` |
| `HDPT:FULL` | 41 | なし | 単語対象外。 | `BrowsMaleSnowElf` / `BrowsMaleSnowElf` |
| `INFO:NAM1` | 4034 | あり | 翻訳対象。 | `DLC01_WESC07_StandAside` / `Stand aside!` |
| `INFO:RNAM` | 109 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1LD_KatriaLongStory` / `How did you end up here?` |
| `INGR:FULL` | 5 | あり | 単語対象。現行抽出がある。 | `DLC01GlowPlant01Ingredient` / `Gleamblossom` |
| `KEYM:FULL` | 5 | あり | 翻訳対象。 | `DLC1_WESC_DawnguardCacheKey` / `Dawnguard Cache Key` |
| `LCTN:FULL` | 56 | あり | 単語対象。現行抽出がある。 | `DLC1VampireCastleGuildhallLocation` / `Volkihar Keep` |
| `LIGH:FULL` | 1 | あり | 翻訳対象。 | `DLC1Torch` / `Torch Bright` |
| `LSCR:DESC` | 19 | あり | 翻訳対象。 | `DLC1WerewolfPerkTree` / `As a werewolf, the inventory button will bring up the werewolf perk tree.` |
| `MESG:DESC` | 73 | あり | 翻訳対象。 | `DLC1NightPowerShrineMessage` / `Choose your Night Power (power button)` |
| `MESG:FULL` | 45 | なし | 単語対象外。 | `DLC1_WESC02_VigilantName` / `Vigilant of Stendarr` |
| `MESG:ITXT` | 16 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1BloodMagicShrineMessage` / `Conjure Gargoyle` |
| `MGEF:DNAM` | 131 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1VampireChangeFXEffect` / `Transform into the vampire lord.` |
| `MGEF:FULL` | 260 | あり | 翻訳対象。 | `DLC1TrapSpotLightConcAimed` / `Trap Sunlight` |
| `MISC:FULL` | 58 | あり | 単語対象。現行抽出がある。 | `DLC1RH05DwarvenTechEnhancedCrossbow` / `Enhanced Crossbow Schematic` |
| `NPC_:FULL` | 255 | あり | 単語対象。現行抽出がある。 | `DLC1EncDeerGlowing` / `Vale Deer` |
| `NPC_:SHRT` | 25 | あり | 単語対象。現行抽出がある。 | `DLC1SorineJurard` / `Sorine` |
| `PERK:DESC` | 42 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1VampiricGrip` / `Blood Magic: Can pull a creature to you from a distance` |
| `PERK:EPF2` | 4 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1VampireTurnPerk` / `Turn into Vampire` |
| `PERK:EPFD` | 2 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1PlayerWerewolfSavageFeeding` / `Feed` |
| `PERK:FULL` | 52 | あり | 翻訳対象。 | `DLC1AnimalVigorPerk` / `Animal Vigor` |
| `PROJ:FULL` | 15 | なし | 単語対象外。 | `DLC1BoltSteelProjectile` / `Steel Bolt` |
| `QUST:CNAM` | 190 | あり | 翻訳対象。 | `DLC1VQ03Vampire` / `I'm one step behind the Moth Priest.` |
| `QUST:FULL` | 134 | あり | 翻訳対象。 | `DLC1_WESC06` / `Dawnguard Novice (06 / Whiterun Caches)` |
| `QUST:NNAM` | 262 | あり | 翻訳対象。 | `DLC1RV07` / `Return to <Alias=QuestGiver>` |
| `RACE:DESC` | 12 | なし | 本文対象候補。単語対象からは分離する候補である。 | `SnowElfRace` / `Also known as "Altmer" in their homeland of Summerset Isle` |
| `RACE:FULL` | 46 | あり | 単語対象。現行抽出がある。 | `DLC1VampireBeastRace` / `Vampire Lord` |
| `REFR:FULL` | 17 | なし | 単語対象外。 | `DLC1VolkiharFerryRef` / `Icewater Jetty` |
| `SCRL:FULL` | 1 | あり | 翻訳対象。 | `DLC1dunRedwaterDenTelekinesisScroll` / `Scroll of Telekinesis` |
| `SHOU:DESC` | 11 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1SummonDragonShout` / `Durnehviir! Hear my Voice and come forth from the Soul Cairn.` |
| `SHOU:FULL` | 24 | あり | 単語対象。現行抽出がある。 | `DLC1SummonDragonShout` / `Summon Durnehviir` |
| `SPEL:DESC` | 29 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1VampireDetectLife` / `Sends out a quick pulse that detects any creature.` |
| `SPEL:FULL` | 306 | あり | 翻訳対象。 | `DLC1AbVampireLord` / `Vampire Lord Abilities` |
| `TACT:FULL` | 3 | なし | 単語対象外。 | `DLC1LD_KatriaTalkingActivator` / `Ghostly Voice` |
| `TREE:FULL` | 25 | なし | 単語対象外。 | `TreeDeadVineLongAsh` / `Ash Vine` |
| `WEAP:DESC` | 10 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1LD_AetherialStaff` / 空文字 |
| `WEAP:FULL` | 100 | あり | 単語対象。現行抽出がある。 | `DLC1AurielsBow` / `Auriel's Bow` |
| `WOOP:FULL` | 12 | なし | 単語対象外。 | `DLC1DragonSummon1Dur` / `D6` |
| `WOOP:TNAM` | 12 | なし | 本文対象候補。単語対象からは分離する候補である。 | `DLC1DrainVitality2Lah` / `Magicka` |
| `WRLD:FULL` | 16 | あり | 翻訳対象。 | `DLC1HunterHQWorld` / `Dayspring Canyon` |

## 現行抽出から見える不足候補

`tools/extractor` は、単語対象 13 種別を現行出力に含める。
追加 5 種別のうち、`NPC_:SHRT` と `RACE:FULL` は単語対象として現行出力に含める。
追加 5 種別のうち、`DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` は `tools/extractor` が出力に含めるが、単語翻訳フェーズと XML 辞書取り込みの対象外である。

## 実装判断で確認すること

- 単語翻訳フェーズは、単語対象 13 種別だけを対象にする。
- `tools/extractor` に出力処理があり、単語対象 13 種別に含まれない種別は翻訳対象へ寄せる。
- `NPC_:SHRT` と `RACE:FULL` は単語対象として `tools/extractor` に追加済みである。
- `DOOR:FULL`、`FLOR:FULL`、`FURN:FULL` は翻訳対象として `tools/extractor` に追加済みである。
