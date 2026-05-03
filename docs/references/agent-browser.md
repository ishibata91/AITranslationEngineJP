# agent-browser 利用規約

関連文書: [`index.md`](./index.md), [`../../.codex/skills/investigate/SKILL.md`](../../.codex/skills/investigate/SKILL.md), [`../../.codex/skills/implementation-investigate/SKILL.md`](../../.codex/skills/implementation-investigate/SKILL.md)

この文書は、Codex の調査系ロールが UI 状態、console、screenshot を証跡化する時の `agent-browser` CLI 利用規約である。
Codex の UI 証跡取得は `agent-browser` CLI に統一する。

## 対象範囲

- 対象ロール: `investigator`、`implementation_investigator`
- 対象 skill: `investigate`、`implementation-investigate`
- 対象用途: UI 状態、console、screenshot、network の観測
- 証跡場所: `tmp/agent-browser/`
- 終了条件: 観測後に `agent-browser close` を実行する

## 起動確認

Wails dev server を使う場合は、先に次を実行する。

```bash
npm run dev:wails:agent-browser
```

CLI とブラウザ環境は、次で確認する。

```bash
agent-browser doctor --offline --quick
```

画面を開く入口は次を使う。

```bash
agent-browser open http://localhost:34115
agent-browser open http://localhost:34115/#dashboard
```

## 状態確認

画面状態は次で確認する。

```bash
agent-browser snapshot
agent-browser get title
agent-browser get url
agent-browser get text "#root"
agent-browser is visible "#root"
```

`@e2` のような参照値は、直前の `agent-browser snapshot` の結果から選ぶ。
CSS selector は、画面上で安定している対象にだけ使う。

## 操作

画面操作は次を使う。

```bash
agent-browser click "@e2"
agent-browser fill "#input-id" "value"
agent-browser find role button click --name "保存"
agent-browser find text "辞書" click
agent-browser press Enter
agent-browser reload
```

複数操作をまとめる場合は次を使う。

```bash
agent-browser batch --bail \
  "open http://localhost:34115" \
  "snapshot" \
  "console" \
  "screenshot tmp/agent-browser/ui-evidence.png"
```

## 証跡

証跡取得は次を使う。

```bash
agent-browser console
agent-browser errors
agent-browser screenshot tmp/agent-browser/ui-evidence.png
agent-browser screenshot --annotate --screenshot-dir tmp/agent-browser
agent-browser network requests
```

`agent-browser screenshot` の出力先 directory が存在しない場合は、先に `tmp/agent-browser/` を作る。
console、errors、screenshot、network requests は、実行コマンドと結果を完了報告入力の根拠に残す。

## 終了

観測後は次を実行する。

```bash
agent-browser close
agent-browser close --all
```

system test の Playwright runner はプロダクトテスト用の別入口として扱う。
