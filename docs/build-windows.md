# Windows ビルド手順

AITranslationEngineJp を Windows 実機でビルドし、配布フォルダを組み立てる手順を定める。macOS からのクロスビルドは Wails が公式にサポートしないため、Windows 実機でビルドする前提を採る。

## 前提ツール

Windows 機に repo を clone 済みで、次の 4 つが PATH に通っていること。

| ツール | 用途 | 備考 |
| --- | --- | --- |
| Go | Wails 本体（Go + CGO）のビルド | wails.json / go.mod と整合する版 |
| Node.js | frontend の install と build | wails が内部で `npm install` / `npm run build` を呼ぶ |
| .NET 10 SDK | C# 抽出器（Mutagen.Bethesda）の publish と実行 | ビルド機・実行機の双方で要る |
| Wails CLI | `wails build` の実行 | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |

WebView2 ランタイムは Windows 10/11 に標準搭載のため、通常は追加不要。

## ビルド

repo root で PowerShell 7 以降から次を実行する。

```
pwsh scripts/windows/build-windows.ps1
```

スクリプトは 3 段を順に行う。

- **抽出器の publish**: `dotnet publish tools/extractor -c Release -r win-x64` を `tools/extractor/bin/publish` へ出す。出力先は bootstrap が読む固定パス（`extractorDLL`）と一致させる。
- **Wails 本体のビルド**: `wails build -platform windows/amd64`。frontend の install と build は Wails が実行し、`build/bin/AITranslationEngineJp.exe` を生成する。
- **配布フォルダの組み立て**: exe と周辺データを `build/windows-dist/` へ集約する。

## 配布フォルダの構造

**exe は相対パス前提**で、作業ディレクトリ直下に周辺データを要求する。exe へ埋め込まれるのは `frontend/dist` だけで、他は実行時にファイルとして読む。よって単体 exe では動かず、次の構造を保った実行フォルダから起動する。

```
build/windows-dist/
  AITranslationEngineJp.exe
  db/migrations/                     … schema（sqlite 本体は初回起動時に生成）
  assets/                            … 感情辞書・役割語・stoplist
  tools/extractor/bin/publish/       … 抽出器 DLL（exe が dotnet で起動）
```

## 実行時の要件

- **.NET 10 ランタイム**: exe は抽出器を `dotnet <DLL>` の子プロセスとして起動するため、実行機に .NET 10 ランタイムが要る（publish は framework-dependent）。
- **参照訳と確定訳語**: 供給源は翻訳対象の Data フォルダにある Strings（english / japanese）。配布フォルダへ辞書ファイルを置く必要は無い。
- **provider 設定**: 翻訳先 LLM の接続設定は起動後に画面から行う。

## 実行ログ

配布ビルドは実行ログを exe と同じフォルダの `app.log` へ追記する。配布ビルドは console を持たないため、画面が出る前に落ちる事象は標準エラーへ出しても読めない。起動ごとに `=== start <日時> ===` の区切りを書く。

dev 起動（`npm run dev:wails:run`）は従来どおり標準エラーへ出し、`scripts/dev/run-wails.sh` が `tmp/logs/wails-dev.log` へ落とす。出し分けは build tag（`logging_production.go` と `logging_dev.go`）で行う。

## Mod Organizer から起動する場合

Mod Organizer 経由で起動すると、仮想 Data フォルダに統合された姿が見える。master の解決・load order・依存 mod がまとめて片づくため、Mod Organizer 管理下の mod を翻訳するときは経由して起動する。

**起動前に、Mod Organizer の実行ファイル ブラックリストへ `msedgewebview2.exe` を追加する。**

```
設定 → 回避策（Workarounds） → 実行ファイルのブラックリスト
  → msedgewebview2.exe を追加
```

追加しないと、ウィンドウが出ないまま終了する。エラーダイアログは出ない。原因は次のとおり。

- Mod Organizer は仮想 Data フォルダを見せるため、起動した process とその子 process へ自前の DLL を注入し、`NtCreateFile` などを差し替える。
- WebView2 は画面描画を `msedgewebview2.exe` の子 process で行い、Chromium の sandbox が同じ関数を先に差し替えている。
- 両者が重なると WebView2 のコントローラ生成が `0x8000FFFF` で失敗する。`app.log` に `[WebView2 Error] error creating controller with 8000ffff` が残る。

ブラックリストに入れると、WebView2 の子 process へは注入されなくなる。app 本体への注入は残るので、仮想 Data フォルダは従来どおり見える。子 process は画面を描くだけで Data フォルダを読まないため、除外しても抽出に影響しない。

Mod Organizer の既定のブラックリストには `Chrome.exe`・`Firefox.exe`・`Brave.exe`・`Discord.exe` が並んでいる。いずれも Chromium を土台にした app で、同じ衝突を避けるために入っている。`msedgewebview2.exe` は既定では入っていない。

## 制約

- **クロスビルド不可**: macOS から Windows 向けビルドは Wails が公式にサポートしない。Windows 実機または Windows CI ランナーでビルドする。
- **自己完結 exe ではない**: 現状の配線は相対パスで周辺データを読む dev 寄りの構造で、単体 exe 配布ではない。自己完結配布が要るなら、周辺データの埋め込みかインストーラ同梱を別 task として検討する。
