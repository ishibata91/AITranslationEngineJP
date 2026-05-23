# agent-browser 利用規約

この文書は `agent-browser` CLI の command 例と利用方針を扱う。  
サブエージェントの UI 証跡取得は、この文書に従う。
Codex 本体の Storybook 人間レビューコメント取得は `.codex/browser-use.md` に従う。
対象 URL、起動 command、証跡の置き場所、合否条件は、呼び出し元の skill または task 成果物で決める。

`agent-browser` は、確定的な E2E テスト実行器ではなく、AI による UI 探索、画面観測、失敗診断、証跡取得のために使う。  
CI で再現可能な合否判定が必要な場合は、Playwright などの通常テストに落とし込む。

## 基本方針

画面確認では、`snapshot` を主情報源にする。  
`screenshot` は、視覚差分、レイアウト崩れ、重なり、レスポンシブ確認のために使う。  
`get text` は補助的に使う。

画面操作では、ユーザー操作に近い command を使う。  
`@e2` のような参照値、`role`、`name`、表示テキストを優先する。  
CSS selector は、安定識別子がある場合に限定する。

`eval` は診断専用にする。  
`eval` による `click`、`submit`、入力、状態変更は原則禁止する。

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

画面状態の主確認は `snapshot` を使う。

```bash
agent-browser snapshot
agent-browser snapshot -i --compact --depth 4
agent-browser get title
agent-browser get url
agent-browser errors
```

`@e2` のような参照値は、直前の `agent-browser snapshot` の結果から選ぶ。  
AI 操作の第一選択は、`snapshot` 由来の `@ref` または `find role/text` にする。

画面全体の文字列確認が必要な場合は、次を使う。

```bash
agent-browser get text "#root"
```

ただし、`get text "#root"` は SPA 全体の文字列を取得するため、操作対象の判定には使わない。  
操作対象の有無、ラベル、状態、押下可能性は `snapshot`、`find role`、`is enabled`、`is visible` で確認する。

要素の表示状態は次で確認する。

```bash
agent-browser is visible "@e2"
agent-browser is enabled "@e2"
```

安定した識別子がある場合だけ、CSS selector を使ってよい。

```bash
agent-browser is visible "[data-testid='job-form']"
agent-browser is enabled "[data-testid='create-button']"
```

## CSS selector の利用制限

CSS selector は、画面上で安定している対象にだけ使う。  
優先順位は次の通りとする。

```text
1. snapshot 由来の @ref
2. role + accessible name
3. 表示テキスト
4. data-testid
5. id
6. 安定した aria-* 属性
7. CSS selector
```

次の selector は原則使わない。

```text
- nth-of-type
- nth-child
- 深い class chain
- DOM順序に依存する selector
- 見た目の並びだけを前提にした selector
```

避ける例:

```bash
agent-browser click "article:nth-of-type(3) .phase-actions button"
```

許容する例:

```bash
agent-browser click "[data-testid='body-translation-confirm-button']"
```

より望ましい例:

```bash
agent-browser find role button click --name "この段階の設定を確認"
```

## 操作

画面操作は、ユーザー操作に近い command を使う。

```bash
agent-browser click "@e2"
agent-browser fill "@e3" "value"
agent-browser press Enter
agent-browser reload
```

`role` と accessible name で対象を指定できる場合は、次を優先する。

```bash
agent-browser find role button click --name "保存"
agent-browser find role textbox fill --name "ジョブ名" "value"
agent-browser find role combobox click --name "翻訳サービス"
```

表示テキストで対象を指定できる場合は、次を使う。

```bash
agent-browser find text "辞書" click
agent-browser find text "本文翻訳" click
```

CSS selector を使う場合は、安定識別子に限定する。

```bash
agent-browser fill "[data-testid='job-name-input']" "value"
agent-browser click "[data-testid='create-job-button']"
```

操作後は、画面状態が変わったことを確認する。

```bash
agent-browser snapshot -i --compact --depth 4
agent-browser errors
agent-browser get text "#root"
```

## 複数操作

複数操作をまとめる場合は `batch --bail` を使う。  
`--bail` により、途中で失敗した時点で後続操作を止める。

```bash
agent-browser batch --bail \
  "open <url>" \
  "snapshot -i --compact --depth 4" \
  "errors" \
  "console" \
  "screenshot <output-path>.png"
```

操作を含む batch では、操作前後に状態確認を入れる。

```bash
agent-browser batch --bail \
  "open <url>" \
  "snapshot -i --compact --depth 4" \
  "find role button click --name 保存" \
  "snapshot -i --compact --depth 4" \
  "errors" \
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

```bash
mkdir -p <output-directory>
agent-browser screenshot <output-directory>/screen.png
```

重要操作の前後では、次の証跡を取得する。

```bash
agent-browser snapshot -i --compact --depth 4
agent-browser errors
agent-browser console
agent-browser network requests
agent-browser screenshot <output-path>.png
```

視覚確認では `screenshot` を使う。  
操作対象の有無、ラベル、状態確認では `snapshot` を優先する。  
API 失敗や非同期処理の異常確認では `network requests` と `console` を使う。

## `eval` の利用制限

`eval` は、失敗原因の診断に限定する。  
ユーザー操作の代替として `eval` を使ってはいけない。

禁止例:

```bash
agent-browser eval "document.querySelector('button').click()"
agent-browser eval "document.querySelector('form').submit()"
agent-browser eval "document.querySelector('input').value = 'value'"
agent-browser eval "window.appState.currentStep = 3"
```

許容例:

```bash
agent-browser eval "
const el = document.querySelector('[data-testid=\"create-button\"]');
const r = el.getBoundingClientRect();
({
  text: el.innerText,
  disabled: el.disabled,
  ariaDisabled: el.getAttribute('aria-disabled'),
  pointerEvents: getComputedStyle(el).pointerEvents,
  visibility: getComputedStyle(el).visibility,
  display: getComputedStyle(el).display,
  opacity: getComputedStyle(el).opacity,
  rect: { x: r.x, y: r.y, width: r.width, height: r.height },
  topElement: document.elementFromPoint(r.x + r.width / 2, r.y + r.height / 2)?.outerHTML
})
"
```

`eval` で原因を調べた後も、再操作は `click`、`fill`、`press`、`find role`、`find text` で行う。

## 失敗時の確認順序

操作が期待通りに反映されない場合は、次の順で確認する。

```text
1. snapshot で対象要素が存在するか確認する
2. is visible で対象要素が表示されているか確認する
3. is enabled で対象要素が有効か確認する
4. screenshot で重なり、モーダル、レスポンシブ崩れを確認する
5. errors で JavaScript error を確認する
6. console で警告や実行時ログを確認する
7. network requests で API 失敗を確認する
8. eval で DOM、computed style、elementFromPoint を診断する
```

失敗調査では、`eval click` で進めてはいけない。  
`eval click` で進む場合は、通常操作経路に問題があると判断し、その差分を記録する。

## viewport 確認

レスポンシブ確認では viewport を明示する。

```bash
agent-browser set viewport 390 844
agent-browser snapshot -i --compact --depth 4
agent-browser screenshot <output-path>.png
```

PC 幅と mobile 幅の両方を確認する場合は、呼び出し元の skill または task 成果物で対象 viewport を定義する。

## 終了

観測後は次を実行する。

```bash
agent-browser close
agent-browser close --all
```
