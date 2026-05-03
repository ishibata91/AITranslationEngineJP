# agent-browser 利用規約

この文書は `agent-browser` CLI の command 例だけを扱う。
対象 URL、起動 command、証跡の置き場所は、呼び出し元の skill または task 成果物で決める。

## 環境確認

CLI とブラウザ環境は、次で確認する。

```bash
agent-browser doctor --offline --quick
```

画面を開く入口は次を使う。

```bash
agent-browser open <url>
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
  "open <url>" \
  "snapshot" \
  "console" \
  "screenshot <output-path>.png"
```

## 証跡

証跡取得は次を使う。

```bash
agent-browser console
agent-browser errors
agent-browser screenshot <output-path>.png
agent-browser screenshot --annotate --screenshot-dir <output-directory>
agent-browser network requests
```

`agent-browser screenshot` の出力先 directory が存在しない場合は、先に出力先 directory を作る。

## 終了

観測後は次を実行する。

```bash
agent-browser close
agent-browser close --all
```
