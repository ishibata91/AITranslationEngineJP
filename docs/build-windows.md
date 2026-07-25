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

## 制約

- **クロスビルド不可**: macOS から Windows 向けビルドは Wails が公式にサポートしない。Windows 実機または Windows CI ランナーでビルドする。
- **自己完結 exe ではない**: 現状の配線は相対パスで周辺データを読む dev 寄りの構造で、単体 exe 配布ではない。自己完結配布が要るなら、周辺データの埋め込みかインストーラ同梱を別 task として検討する。
