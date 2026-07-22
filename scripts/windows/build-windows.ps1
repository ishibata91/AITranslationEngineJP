# build-windows.ps1 — Windows 実機で配布フォルダを組み立てるビルドスクリプト。
#
# 前提: この repo を Windows 機に clone 済みで、Go・Node.js・.NET 10 SDK・Wails CLI が入っていること。
# 実行: repo root で `pwsh scripts/windows/build-windows.ps1` （PowerShell 7 以降）。
#
# 生成物: build/windows-dist/ に exe と周辺データを同梱した実行フォルダを作る。
# 背景: exe は相対パス前提で、作業ディレクトリ直下に db/・assets/・dictionaries/・
#       tools/extractor/bin/publish/ を要求する（frontend/dist だけが exe へ埋め込まれる）。
#       単体 exe では動かないため、周辺データを 1 フォルダへ集約する。

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

# repo root（このスクリプトの 2 つ上）へ移動する。相対パス前提の exe と歩調を合わせる。
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
Set-Location $repoRoot
Write-Host "[build-windows] repo root: $repoRoot"

# 必須ツールの存在を先に確かめる。欠けていれば理由を出して止める。
function Assert-Command($name, $hint) {
  if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
    throw "$name が見つからない。$hint"
  }
}
Assert-Command "go"     "Go を入れて PATH に通す。"
Assert-Command "node"   "Node.js を入れて PATH に通す。"
Assert-Command "dotnet" ".NET 10 SDK を入れて PATH に通す。"
Assert-Command "wails"  "go install github.com/wailsapp/wails/v2/cmd/wails@latest で Wails CLI を入れる。"

# C# 抽出器を win-x64 向けに publish する。framework-dependent 出力のため、実行機に .NET 10 ランタイムが要る。
# 出力先は bootstrap の extractorDLL（tools/extractor/bin/publish/extractor.dll）と一致させる。
Write-Host "[build-windows] publishing C# extractor (win-x64)"
dotnet publish tools/extractor -c Release -r win-x64 -o tools/extractor/bin/publish

# Wails 本体を Windows amd64 向けにビルドする。frontend の install と build は wails が実行する。
# 生成物は build/bin/AITranslationEngineJp.exe。
Write-Host "[build-windows] building Wails app (windows/amd64)"
wails build -platform windows/amd64

# 配布フォルダを組み立てる。exe と、exe が実行時に読む周辺データを 1 フォルダへ集約する。
$dist = Join-Path $repoRoot "build\windows-dist"
if (Test-Path $dist) { Remove-Item $dist -Recurse -Force }
New-Item -ItemType Directory -Path $dist | Out-Null

Copy-Item "build\bin\AITranslationEngineJp.exe" $dist

# db/migrations は追跡済み。sqlite 本体は初回起動時に生成されるため同梱しない。
New-Item -ItemType Directory -Path (Join-Path $dist "db") | Out-Null
Copy-Item "db\migrations" (Join-Path $dist "db\migrations") -Recurse

Copy-Item "assets" (Join-Path $dist "assets") -Recurse

# 抽出器の publish 出力を tools/extractor/bin/publish の相対構造のまま同梱する。
$extDst = Join-Path $dist "tools\extractor\bin\publish"
New-Item -ItemType Directory -Path $extDst -Force | Out-Null
Copy-Item "tools\extractor\bin\publish\*" $extDst -Recurse

# dictionaries/xTranslatorXMLs は git 追跡外の利用者供給データ。あれば同梱、無ければ警告に留める。
if (Test-Path "dictionaries\xTranslatorXMLs") {
  New-Item -ItemType Directory -Path (Join-Path $dist "dictionaries") -Force | Out-Null
  Copy-Item "dictionaries\xTranslatorXMLs" (Join-Path $dist "dictionaries\xTranslatorXMLs") -Recurse
} else {
  Write-Warning "dictionaries\xTranslatorXMLs が無い。固有名辞書を使うなら配布フォルダへ後から置く。"
}

Write-Host "[build-windows] done. 実行フォルダ: $dist"
Write-Host "[build-windows] 実行時は .NET 10 ランタイムが要る（exe が dotnet で抽出器 DLL を起動するため）。"
