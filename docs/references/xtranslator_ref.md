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

## 3. XML エクスポート形式

### 3.1 基本構造

`File -> Extract Translation -> XML File` で出力される標準フォーマット。

```xml
<?xml version="1.0" encoding="utf-8"?>
<SSETranslator>
  <String>
    <EDID>RecordEditorID</EDID>
    <REC>REC_TYPE</REC>
    <FIELD>FULL</FIELD>
    <FORMID>0x00012345</FORMID>
    <Source>Original English text</Source>
    <Dest>翻訳後テキスト</Dest>
    <Status>4</Status>
  </String>
</SSETranslator>
```

### 3.2 ルート要素

| 要素名 | 説明 |
|---|---|
| `SSETranslator` | Skyrim SE 向けのルート要素 |
| `TESVTranslator` | Skyrim LE 向けのルート要素 |
| `FO4Translator` | Fallout 4 向けのルート要素 |

### 3.3 `<String>` 子要素一覧

| タグ | 型 | 必須 | 説明 |
|---|---|---|---|
| `<EDID>` | string | ○ | レコードのエディタ ID |
| `<REC>` | string (4文字) | ○ | レコードタイプ（例: `FULL`, `DESC`, `DIAL`） |
| `<FIELD>` | string | ○ | フィールド名 |
| `<FORMID>` | hex string | ○ | レコードの FormID（`0x` プレフィックス付き） |
| `<Source>` | string | ○ | 原文テキスト |
| `<Dest>` | string | ○ | 翻訳テキスト（未翻訳時は原文と同一またはブランク） |
| `<Status>` | integer | ○ | 翻訳ステータスコード（後述） |

### 3.4 Status コード

| 値 | 意味 | UI 表示色 |
|---|---|---|
| `0` | 未翻訳 (Untranslated) | 赤 |
| `1` | 翻訳済み (Translated) | 白 |
| `2` | 部分翻訳 (Partial) | オレンジ |
| `3` | 仮翻訳 (Provisional) | 紫 |
| `4` | 承認済み (Validated/Approved) | 青 |

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
