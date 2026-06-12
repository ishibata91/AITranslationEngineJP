# plan: mutagen-extractor

## 依頼要約

- 翻訳対象抽出の基盤を xEdit Pascal から Mutagen（C#/.NET）へ移す（docs/mutagen-migration-plan.md の段階 1〜3）。
- 実行環境は macOS。Skyrim 実体は無く、`dictionaries/Data/`（Skyrim の Data フォルダを模したもの）を読み込み元とする。
- 照合先の正解は `dictionaries/xTranslatorXMLs/` の xTranslator XML。
- 検証は `scripts/compare_counts.py`（件数差分 +0）と `scripts/validate_extraction.py`（0 error）を Dawnguard.esm で通すことを最初の関門とする。

## branch

- 作業 branch: claude/mutagen-extractor
- 分岐元 branch: master
- 分岐元 commit: 1ab5171f

## 軽 / 重判定

- 画面が動くか: N。抽出基盤は C# console の CLI であり、svelte 表示や layout を変えない。
- docs/architecture.md への反映が要るか: N。抽出は Wails 本体と別プロセスの前処理であり、層構成と依存方向を変えない。
- 判定: 軽 task。design-module と storybook-module を bypass し、implementation-module → finalization-module で進める。

## 環境メモ（このセッションで確認した事実）

- dotnet SDK は 10.0.301 のみ（brew 導入）。`tools/extractor`（net8.0）の build と restore は成功。実行は `DOTNET_ROLL_FORWARD=LatestMajor` で可能（net8.0 runtime は無い）。実装時に TFM を net10.0 へ上げるのが簡潔。
- macOS では `GameEnvironment.Typical.Skyrim` の自動検出が失敗する（実行して確認済み）。`dictionaries/Data/` を data folder として明示指定する構成（GameEnvironmentBuilder 経由）が要る。
- `dictionaries/Data/` の中身: Skyrim.esm、Dawnguard.esm、Dragonborn.esm、inigo.esp、Outfit Recognition Framework.esp、unofficial skyrim creation club content patch.esl、unofficial skyrim special edition patch.esp、Strings/。
- Strings は japanese のみ 321 ファイル（english 無し）。抽出時の language 指定は Japanese にする。件数比較は非空判定ベースなので english 不在でも成立する見込み。
- `dictionaries/xTranslatorXMLs/` にあるのは Skyrim / Update / Dragonborn / HearthFires / USSEP の 5 つ。**Dawnguard 用 XML（完全版）が無い**。`tests/fixtures/master-dictionary/Dawnguard_english_japanese.xml` は 642 byte のテスト fixture で照合先にならない。
- **Update.esm と HearthFires.esm が Data に無い**。Dawnguard.esm の master は Skyrim.esm + Update.esm のため、参照解決（話者解決の alias / base object）が Update.esm 内 record に当たると解決不能になる可能性がある。XML 側には Update / HearthFires があり、esm 側と非対称。

## 検証方式の変更（人間指示 2026-06-11）

- JSON を中間に挟んで python（compare_counts.py / validate_extraction.py）で検証する方式をやめる。
- compare_counts.py と同じタイプの検証を Mutagen 側（C# テスト）で持つ。xTranslator XML の解析と件数比較を C# テストとして実装し、`dotnet test` で回す。
- validate_extraction.py の schema 検証は、JSON を出さない構成では C# の型定義が代替する。非空必須などの意味検証だけテストに残す。

## 人間回答による確定事項（2026-06-11）

- Dawnguard_english_japanese.xml（2.5MB 完全版）を dictionaries/xTranslatorXMLs/ に配置済み（置き忘れだった）。最初の関門は計画どおり Dawnguard +0。
- Update.esm と HearthFires.esm を dictionaries/Data/ に配置済み。Dawnguard の master 欠落は解消。
- JSON 出力は作らない。将来は Go/Wails 本体から Mutagen 抽出を直接呼び出し、UI から esp を直指定する形を想定する（境界設計は別 task）。今回の scope は「抽出ロジックの正しさを C# テストで固める」まで。

## 実装範囲（軽 task のため本 plan に直接記録）

1. tools/extractor を再構成する。TFM を net10.0 へ上げ、`dictionaries/Data/` を data folder として明示指定する経路（自動検出に依存しない）にする。
2. 抽出ロジックを library として実装する。extractData.v2.pas を移植元、docs/mutagen-migration-plan.md §4〜§7 を仕様とする。plugin 単位列挙（WinningOverrides 不使用）。
3. xUnit テストプロジェクトを追加する。xTranslator XML を解析し、REC:FIELD ごとの件数を抽出結果と比較する（compare_counts.py 相当）。
4. validate_extraction.py 相当の意味検証（必須フィールド非空など）をテストに含める。

## 実装で確定した抽出仕様（移行計画からの差分）

- 抽出言語は英語を既定にする。翻訳エンジンの入力は英語 mod で、xTranslator 辞書の Source も英語。
  日本語など他言語は `--language` で指定できる。
- 会話は最初から v19 の 2 階層（DialogueTopic → InfoNode → ResponseLine）で実装した。
  現行 JSON schema 互換の段階は、JSON 廃止により不要になった。
- WOOP は shout への埋め込みでなく top-level で列挙する（compare の id ユニーク化が不要になる）。
- Skyrim 追加 record（APPA / FACT:MNAM / FNAM / SNCT / EYES / REGN:RDMP）に加え、
  HAZD:FULL（Dragonborn の辞書に存在）と TES4:CNAM / SNAM（非 localized plugin の header）を対象に追加した。
- 翻訳所有（override の field をどの plugin の翻訳対象とするか）は次の規則で判定する:
  - 新規 record は常に所有。
  - override は「正規化 record data」が master のどの版とも一致しない場合に所有。
    正規化 = string ID を strings table の本文（英語）へ解決し、VMAD は byte 整列、
    CK 保存ノイズ（XCLW / XLCN / XSCL / XLOC）を除外し、subrecord を整列したバイト列。
  - DIAL は plugin 内に所有 INFO を持てば FULL も所有（container 規則）。
  - CELL は所有された配置 ref（REFR / ACHR）を持てば FULL も所有（NAVM / LAND は対象外）。
- 非 localized plugin の内蔵文字列は UTF-8 優先（cp1252 fallback）で decode する
  （翻訳適用済み esp が日本語 UTF-8 を内蔵するため）。

## 検証結果（2026-06-12、dotnet test 11/11 通過）

| plugin | XML 比較対象 | 差分 |
|---|---:|---|
| Skyrim.esm | 64,571 | +0（完全一致） |
| Update.esm | 1,391 | +26（既知差分 6 種） |
| Dawnguard.esm | 7,917 | +1（既知差分 2 種） |
| HearthFires.esm | 2,249 | −12（既知差分 3 種） |
| Dragonborn.esm | 10,834 | +21（既知差分 4 種） |
| USSEP | 18,837 | +1（既知差分 1 種） |

- 既知差分は全 98,000 strings 中 65（0.066%）。tools/extractor.Tests/CountParityTests.cs の
  KnownDeltas 表に plugin × REC:FIELD 単位で固定し、表に無い乖離・表より改善した乖離の両方を fail にする。
- 既知差分の根拠: xTranslator 辞書側に plugin 内容から導出できない帰属の揺れがある。
  同一内容・同一構造の record 対で帰属が逆になる実測例を確認した
  （CELL: HallOfTheVigilant01 採用 vs KarthwastenSanuarachMine 非対象、
    QUST: DA07 採用 vs TG07 / MQ105 非対象。どちらの対も差分は同種）。
  公式 5 esm の XML（2026-02 生成）と USSEP の XML（2026-06 生成）でセッション・設定も異なる。
- 旧 python 検証（scripts/compare_counts.py / validate_extraction.py）は C# テストが置き換えた。
  JSON 中間出力は作らない（人間指示）。

## finalization

- 正本化判断: docs/architecture.md と docs/screen-design/ への反映対象なし。
  抽出基盤は Wails 本体と別プロセスの CLI / library で、層構成・依存方向・画面を変えない。
- 作業 commit: d3bea2a6（13 files、tools/extractor + tools/extractor.Tests + 本 plan）。
- 検証: dotnet test tools/extractor.Tests → 11/11 通過（件数比較 6 + 意味検証 5）。
- 残留リスク: 既知差分 65 strings は xTranslator 辞書側の帰属揺れ。抽出を変えた場合は
  KnownDeltas 表の検証が fail して検知される。
