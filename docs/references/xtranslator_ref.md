# xTranslator 出力フォーマット仕様書

## 1. 概要

**xTranslator**（別名: sseTranslator / tesvTranslator / fallout4Translator）は、Bethesda ゲーム（Skyrim、Skyrim SE、Fallout 4、Fallout 76、Starfield 等）の MOD ローカライズ作業を支援するツールです。本仕様書では、このツールが読み書きする主要なファイル・出力フォーマットを定義します。

---

## 2. 対応ファイル種別

| 拡張子 | 種別 | 説明 |
|---|---|---|
| `.esp` / `.esm` | Plugin ファイル | ゲームデータ本体。直接翻訳を書き込む |
| `.STRINGS` | 文字列ファイル | ローカライズ対応 esp が参照する汎用テキスト |
| `.DLSTRINGS` | 文字列ファイル | ダイアログ向け文字列 |
| `.ILSTRINGS` | 文字列ファイル | FUZ 音声と対応する文字列 |
| `.sst` | 辞書ファイル | xTranslator 独自バイナリ辞書 |
| `.xml` | インポート/エクスポート | 翻訳データの移植・バックアップ用 |
| `.txt` | MCM テキスト | SkyUI MCM メニュー等の UI 文字列 |
| `.pex` | Papyrus スクリプト | コンパイル済みスクリプト（翻訳可能文字列を含む） |

---

## 3. XML エクスポート/インポート形式（SSTXMLRessources）

xTranslator の `File -> Export -> xTranslator XML` が出力し、`File -> Import Translation -> XML File` が読み込む形式。本 repo が読み書きするのはこの形式で、base 辞書 `dictionaries/xTranslatorXMLs/*_english_japanese.xml` も同形式である。仕様の根拠は xTranslator ソース `TESVT_XMLFunc.pas`（GitHub MGuffin/xTranslator）と、実 base 辞書ファイルの実測。

以前この節に記していた `SSETranslator` ルート・`FIELD`/`FORMID`/`Status` 独立要素の形式は xTranslator README 由来で、実物のエクスポートと一致しなかったため、実物形式へ訂正した（`xtranslator-export` task の実画面確認で判明）。

### 3.1 基本構造

先頭に UTF-8 BOM を付け、宣言は `encoding="UTF-8" standalone="yes"`。

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<SSTXMLRessources>
  <Params>
    <Addon>PluginName</Addon>
    <Source>english</Source>
    <Dest>japanese</Dest>
    <Version>2</Version>
  </Params>
  <Content>
    <String List="0" sID="000001" Partial="1">
      <EDID>RecordEditorID</EDID>
      <REC>WEAP:FULL</REC>
      <Source>Original English text</Source>
      <Dest>翻訳後テキスト</Dest>
    </String>
  </Content>
</SSTXMLRessources>
```

### 3.2 ルート要素と Params

- ルート要素は `SSTXMLRessources`。
- `Params` は書き出しの前提情報を持つ。`Addon`（対象 plugin 名）、`Source`（原文言語、例 english）、`Dest`（訳文言語、例 japanese）、`Version`（現行は 2）。
- `Content` が `String` 群を包む。

### 3.3 `<String>` の属性と子要素

属性:

| 属性 | 説明 |
|---|---|
| `List` | 文字列の分類インデックス（0〜2。STRINGS/DLSTRINGS/ILSTRINGS の区分に対応）。インポートの照合は EDID・REC・原文で行われ、List は分類の記録。 |
| `sID` | 文字列の識別子（6 桁 hex）。string（strings ファイル直接編集）モードでのみ照合に使う。 |
| `Partial` | 訳状態フラグ（後述）。確定訳では属性ごと省略する。 |

子要素:

| タグ | 必須 | 説明 |
|---|---|---|
| `<EDID>` | ○ | レコードのエディタ ID。インポート照合の第一キー。 |
| `<REC>` | ○ | `REC:FIELD`（4文字:4文字、例 `WEAP:FULL`・`INFO:NAM1`）。インポート照合の副キー。 |
| `<Source>` | ○ | 原文テキスト。EDID・REC で照合できない場合の第三キー。 |
| `<Dest>` | ○ | 訳文テキスト。 |

`FIELD`・`FORMID`・`Status` の独立要素は持たない。フィールドは `REC` に結合し、訳状態は `Partial` 属性で表す。

### 3.4 訳状態（`Partial` 属性）

xTranslator は独立した Status 要素を持たず、`String` の `Partial` 属性で訳状態を表す。

| Partial | 意味 |
|---|---|
| （属性なし） | 確定訳（translated） |
| `1` | 未完了訳（incompleteTrans） |
| `2` | ロック訳（lockedTrans） |

---

## 4. SST 辞書形式 (`.sst`)

SST は xTranslator 専用のバイナリ辞書フォーマット。ヒューリスティック翻訳提案や一括翻訳に使用される。

- エンコーディング: UTF-8（バイナリ内部）
- 構成: 原文文字列と翻訳文字列のペアをインデックス構造で格納
- 最適化: v1.4.6 以降で最大 **10倍高速化**されたロード処理
- 自動バックアップ: `_xTranslator\UserDictionaries\[Game]\Auto\` に保存可能（v1.4.10 以降）

> **注意**: xTranslator をアップグレードすると `UserDictionaries` フォルダが初期化される場合があるため、辞書は別フォルダで管理することを推奨。

---

## 5. STRINGS / DLSTRINGS / ILSTRINGS 形式

### 5.1 ファイル配置ルール

```
[ESP ファイルと同じディレクトリ]
  └── strings/
        ├── ModName_japanese.STRINGS
        ├── ModName_japanese.DLSTRINGS
        └── ModName_japanese.ILSTRINGS
```

### 5.2 エンコーディング設定

言語ごとのエンコーディングは `xTranslator\Data\[Game]\codepage.txt` で定義する。

```ini
english=1252
japanese=utf8
korean=utf8
chinese=utf8
russian=1251
```

- Skyrim SE / Fallout 4 以降はデフォルトで UTF-8 対応
- Skyrim LE では `english=utf8,1252` への変更が必要な場合がある

---

## 6. MCM / Translate テキスト形式 (`.txt`)

SkyUI MCM や UI テキストに使用されるフォーマット。

- カスタムテキスト定義は `xTranslator\misc\customTxtDefinition.txt` で設定
- 1 行 1 エントリの構造（シングルライン定義が必要）
- Hybrid Mode での翻訳が推奨（Strings Only Mode は非推奨）

---

## 7. Papyrus スクリプト (`.pex`)

xTranslator 内蔵のデコンパイラを使用して翻訳可能な文字列を抽出・編集する。

- 内部変数はロック（編集不可）
- カスタムコードページ対応: Advanced Options -> Script タブで設定
- 64 ビット版 `.pex`（Skyrim SE）は v1.2.1 以降で対応

---

## 8. XML インポート仕様

`File -> Import Translation -> XML File (xTranslator)` でのインポート時オプション：

| オプション | 説明 |
|---|---|
| **Overwrite: Entire Line** | 対象レコードを全行上書き |
| **Overwrite: Dest Only** | 翻訳テキスト（`<Dest>`）のみ上書き |
| **Mode: Use FormID Reference** | FormID でレコードを照合（推奨） |
| **Mode: Use String Reference** | 原文文字列でレコードを照合 |

---

## 9. エンコーディングフォールバック機構

v1.1.6 以降、フォールバックエンコーディングが実装されている。

```
プライマリ: UTF-8 でデコード試行
  -> 失敗した場合: codepage.txt に定義されたフォールバックコードページを適用
```

---

## 10. 注意事項・制約

| 項目 | 内容 |
|---|---|
| 最大文字列サイズ | 約 1,000,000 バイト（エンコード後）。VMAD / PEX 内の文字列は 65,565 バイト上限 |
| 翻訳非推奨フィールド | `WOOP`（シャウト発音データ）は翻訳すると文字化けの恐れあり |
| RACE 先頭スペース | 一部 RACE レコードは翻訳テキストの先頭にスペースが必要（省略すると CTD） |
| 末尾スペース | ALCH / ARMO / ENCH 等の名前に末尾スペースを含むと「■」が追加される |
| 一括変換の長文バグ | Batch Replace は長い文章の後半が削除される既知のバグあり |

---

## 11. 設定ファイルパス

| ファイル | パス | 用途 |
|---|---|---|
| `prefs.ini` | `xTranslator\UserPrefs\[GameName]\prefs.ini` | オプション設定の保存先 |
| `codepage.txt` | `xTranslator\Data\[GameName]\codepage.txt` | 言語別エンコーディング定義 |
| `customTxtDefinition.txt` | `xTranslator\misc\customTxtDefinition.txt` | MCM カスタムテキスト解析ルール |
| `ApiTranslator.txt` | `xTranslator\Misc\ApiTranslator.txt` | DeepL / MS Translator エンドポイント設定 |
| `res.ini` | `xTranslator\Res\[Language]\res.ini` | UI ローカライズ文字列 |

---

*本仕様書は xTranslator v1.5.8（2024年8月）の GitHub README および namu.wiki ドキュメントを元に作成。*
