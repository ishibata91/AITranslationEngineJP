---
name: codegraph-rg-explore
description: Codegraphを用いたコードベース探索用のスキル。TRIGGER WHEN　コードベースの探索時。SKIP WHEN　コードベース以外の探索・探索以外のタスク。メインエージェント使用禁止。
---
## 動作に必要な情報

- 調査対象: 人間が依頼した機能、処理経路、責務、または依存関係。

---
# ステップ1
#### `rg`で入口となるコメント・クラス・シンボルを検索する。

- 依頼文から名詞を最大5個選ぶ。抽出後，選択した語の英語と日本語を用意すること
- 選択した語でコードベース全体に向けて検索をかけ，関連するシンボルを特定する。
- 結果に現れたシンボルを確認し，調査を依頼された対象に適合すると思われるものを最大5個抽出する。

コマンド例
```sh
rg -n -i -B 1 -A 3 \
  --glob '*.go' \
  --glob '!**/*_test.go' \
  --glob '!**/.codegraph/**' \
  '((辞書|dictionary).{0,40}(置換|適用|replace|apply)|(置換|適用|replace|apply).{0,40}(辞書|dictionary))' \
  internal
```

## ステップ2
#### codegraphで実在シンボルを検索する
- codegraph CLIを用いる
- `{REPO_ROOT}/internal/`,`{REPO_ROOT}/frontend/`, `{REPO_ROOT}/tools` でindexが別々に構築されている。プロジェクトを指定する
- 指定されたフォルダがcodegraphプロジェクトでない場合に限り`rg` など別ツールでの探索を許可する。
- ステップ1で抽出したシンボルで検索を行う
- 検索結果から依頼に関係のあるものだけを選別する
- 最大呼び出し３回まで。

#### 検索に使うコマンド

- `codegraph query -p {PROJECT_PATH} "{SYMBOL}"`：実在するシンボルを名前で検索する。`-k function`のように種類を絞り，`-l 20`のように件数を増やせる。
- `codegraph explore -p {PROJECT_PATH} "{SYMBOL_OR_QUESTION}"`：関連シンボルのソースと呼び出し経路をまとめて取得する。ステップ2の選別後に，処理経路と責務を確認するために使う。
- `codegraph files -p {PROJECT_PATH} --filter {DIRECTORY}`：インデックス済みのファイル構造を絞り込んで表示する。調査対象の配置が不明な場合に使う。
- `codegraph files -p {PROJECT_PATH} --pattern "{GLOB}" --format flat`：拡張子や名前のパターンで候補ファイルを絞り込む。例は`--pattern "*Dictionary*"`である。
- `codegraph node -p {PROJECT_PATH} "{SYMBOL_OR_FILE}"`：選別済みシンボル，または候補ファイルのソースと依存関係を確認する。
- `codegraph callers -p {PROJECT_PATH} "{SYMBOL}"`：呼び出し元を確認する。
- `codegraph callees -p {PROJECT_PATH} "{SYMBOL}"`：呼び出し先を確認する。
- `codegraph status {PROJECT_PATH}`：インデックスの状態を確認する。検索結果が古い疑いがある場合に使う。

## ステップ3

- 選別した検索結果から結論を出す。
- codegraphの検索結果で足りなければ，検索結果のシンボル・ファイルで限定して`rg`でさらに詳しく探索する。


