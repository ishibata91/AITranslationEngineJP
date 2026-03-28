---

name: architecting-tests
description: AITranslationEngineJp 専用。impl / fix の実装直前に、active exec-plan と関連仕様から仕様に沿った tests / acceptance checks / fixtures / validation commands を設計し、TDD の先行テストを準備する。
---

# Architecting Tests

## Overview

実装前に、何をテストで証明するかを決める。active exec-plan と関連仕様を読んで、最小の failing tests、acceptance checks、fixtures、validation commands に落とす。

## Terms

* `failing tests`: まだ実装されていない前提で先に書く、失敗するべきテスト。証明したい振る舞いを最小単位で固定する。
* `fixtures`: テストの入力や前提を作るデータ。JSON、DB 初期データ、モック応答、サンプルファイルなどを含む。
* `acceptance checks`: 仕様を満たしたと言えるかを、ユーザー視点や業務フロー視点で確認するチェック。
* `validation commands`: そのテストやチェックを機械的に実行するためのコマンド。

## Test Style Policy

* テスト名は Given-When-Then 形式で書く
* テスト構造は Arrange-Act-Assert で書く
* 1 test では 1 つの振る舞いだけを証明する
* テスト名は「条件 / 操作 / 期待結果」を含める
* 条件分岐や境界値が複数ある場合は table-driven tests を優先する
* テスト対象は public behavior とし、内部実装には結び付けない
* mock / spy は外部境界（DB, API, clock, queue, filesystem）に限定する
* failing tests は「次に失敗させたい最小の 1 ケース」から始める
* fixture は最小の入力と前提だけを用意する
* 時刻・乱数・ID は固定する
* acceptance checks はユーザー視点の結果のみを確認し、内部実装詳細を含めない

## Workflow

1. active exec-plan と関連文書を読む。
2. 要件レベルとそれ未満の細かな仕様を分ける。
3. テストで担保する観点を、unit / integration / acceptance のどこで見るか決める。
4. failing tests ごとに観測点（戻り値 / 状態変化 / 外部出力 / エラー）を先に決める。
5. fixture、acceptance checks、validation command をその観測点に合わせて決める。
6. 仕様にない振る舞いは追加せず、必要なら docs sync に回す。
7. 必要なら active exec-plan の `Acceptance Checks` を更新する。
8. 実装へ handoff する前に、短い test brief を返すか、対象テストを直接書く。

## Lane Rules

### Impl lane

* 先に期待結果を固定し、実装はその failing tests を満たす形に寄せる
* UI / Scenario / Logic のどれをテストで証明するかを明示する

### Fix lane

* 再現手順を test に落としてから修正する
* 回帰を防ぐ最小の test case を優先する

## Few-Shot

### Impl lane

Request:

> xEdit の import 後に、1 件の翻訳可能フィールドから `TRANSLATION_UNIT` が作られることを先に保証したい。

Test brief:

- failing test: 1 件の最小 fixture を読み込んだら `TRANSLATION_UNIT` が 1 件だけ生成される
- fixture: 翻訳可能フィールド 1 件だけを持つ最小 xEdit JSON
- acceptance check: import の結果として生成物にその 1 件が反映される
- validation command: `cargo test import_creates_single_translation_unit -- --exact`

### Fix lane

Request:

> 出力 writer が空の翻訳結果で落ちるので、再現テストを先に置きたい。

Test brief:

- failing test: `dest` が空のレコードでも writer が panic せず出力を返す
- fixture: 空 `dest` を含む 1 レコードの回帰 fixture
- acceptance check: writer が失敗せず、空値の扱いが崩れない
- validation command: `cargo test writer_handles_empty_dest -- --exact`

## Rules

* 実装コードを広く直さない
* 仕様を勝手に補完しない
* test の増やし過ぎで scope を膨らませない
* 結果は短く、使える形で返す
* 1 test = 1 behavior を守る
* failing test は常に「次の 1 ケースのみ」を対象にする
* モックは外部境界に限定し、内部相互作用を固定しない
* Arrange が肥大化する場合は fixture / helper に分離する

## Reference Use

- impl lane では着手前に `../directing-implementation/references/directing-implementation.to.architecting-tests.json` を参照し、返却時は `references/architecting-tests.to.directing-implementation.json` を使う。
- fix lane では着手前に `../directing-fixes/references/directing-fixes.to.architecting-tests.json` を参照し、返却時は `references/architecting-tests.to.directing-fixes.json` を使う。
