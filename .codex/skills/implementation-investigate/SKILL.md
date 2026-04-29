---
name: implementation-investigate
description: Codex implementation レーン 側の実装時調査の共通作業プロトコル。単一引き継ぎ入力 1 件内で 根拠 first に調査する判断基準を提供する。
---
# Implementation Investigate

## 目的

`implementation-investigate` は作業プロトコルである。
`implementation_investigator` agent が、`単一引き継ぎ入力` 1 件と 承認済み実装範囲 内で実装時の証拠を集める時の共通判断を提供する。

ツール権限 は [implementation_investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_investigator.toml) が持ち、引き継ぎ は skill に従う。

## 対応ロール

- `implementation_investigator` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `implementation-investigate` の出力規約で固定する。

## 入力規約

- 実装前再現、trace、再観測を行う時
- 一時観測点を add / remove する時
- 根拠 と仮説を分けて返す時
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: 単一引き継ぎ入力, 承認記録, 承認済み実装範囲, investigation_request, 検証コマンド
- 任意入力: 参照ヒント, 再現根拠, 一時観測計画
- 入力注記: {"単一引き継ぎ入力": "implementation-scope から抽出済みの 引き継ぎ 1 件だけ。implementation-scope 全文、進行中作業計画 全文、根拠成果物、後続 引き継ぎ は入力に含めない。", "参照ヒント": "implementation-investigate-reproduce、trace、観測、再観測 の参照ヒント。共通規約と完了条件は変えない。"}

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_investigator.toml) の 書き込み許可 / 実行許可 とする。
- エージェント実行定義: [implementation_investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_investigator.toml)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate-reproduce/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate-trace/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate-observe/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate-reobserve/SKILL.md

## 内部参照規約

### 拘束観点

- 根拠 first の観測
- 観測済み事実 と 仮説 の分離
- temporary observation の cleanup
- `agent-browser` CLI による UI / console / screenshot 根拠
- 重点 skill の選び方

### Browser Evidence

UI 状態、console、screenshot の観測は `execute` から `agent-browser` CLI を使う。

標準入口は次の通りである。

```bash
npm run dev:wails:agent-browser
agent-browser doctor --offline --quick
agent-browser open http://localhost:34115
agent-browser snapshot
agent-browser console
agent-browser screenshot tmp/agent-browser/ui-根拠.png
```

観測後は `agent-browser close` を実行する。
system test の Playwright runner は プロダクトテスト 用の別入口として扱う。

### 操作コマンド一覧

起動と終了:

```bash
agent-browser open http://localhost:34115
agent-browser open http://localhost:34115/#dashboard
agent-browser reload
agent-browser close
agent-browser close --all
```

状態確認:

```bash
agent-browser snapshot
agent-browser get title
agent-browser get url
agent-browser get text "#root"
agent-browser is visible "#root"
```

操作:

```bash
agent-browser click "@e2"
agent-browser fill "#input-id" "value"
agent-browser find role button click --name "保存"
agent-browser find text "辞書" click
agent-browser press Enter
```

証跡:

```bash
agent-browser console
agent-browser errors
agent-browser screenshot tmp/agent-browser/ui-根拠.png
agent-browser screenshot --annotate --screenshot-dir tmp/agent-browser
agent-browser network requests
```

複数手順をまとめる場合:

```bash
agent-browser batch --bail \
  "open http://localhost:34115" \
  "snapshot" \
  "console" \
  "screenshot tmp/agent-browser/ui-根拠.png"
```

`@e2` のような ref は直前の `スナップショット` の結果から選ぶ。
selector が安定している場合だけ CSS selector を使う。
console / errors / screenshot / network requests は 完了報告入力 の 根拠 に コマンド と結果を残す。

- 参照 型 は [investigation-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate/references/patterns/investigation-patterns.md) とする。

## 判断規約

- `単一引き継ぎ入力` 1 件と 承認済み実装範囲 を超えない
- 根拠 のない結論を固定しない
- Codex implementation レーン のブラウザ操作は `agent-browser` CLI で行う
- 一時観測点は返却前に除去する
- 恒久修正と プロダクトテスト 追加を混ぜない

- 観測条件、コマンド、結果を残す
- temporary changes と cleanup_status を返す
- recommended next step を根拠付きで返す
- active 規約 は agent 1:1。調査種別の差分は 重点 skill で扱い、出力 obligation はこの 規約 に固定する。

## 非対象規約

- 恒久修正、プロダクトテスト追加、design-time investigation は扱わない。
- 承認済み実装範囲外の調査は扱わない。
- 根拠のない結論は固定しない。
- mode 別 個別 JSON 規約は使わない。

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- 返却先: implement_lane
- 調査 focus: 承認済み実装範囲 内で何を調べたかを返す。
- 再現状態: 再現できたか、未再現か、再現不要かを返す。
- 観測事実: 観測済み事実だけを書き、仮説と混ぜない。
- 仮説: 原因候補と根拠を返す。
- 観測点: 確認した入口、経路、対象を返す。
- 一時変更: 一時観測点を使った場合だけ path と目的を返す。
- cleanup 状態: 一時観測点の除去状態を返す。
- 確認結果: 実行した 検証 と未実行理由を返す。
- 残り 不足: 未確認事項と理由を返す。
- 残留リスク: 実装判断に残る リスク を返す。
- 推奨 next step: implement、tests、戻しのどれが妥当かを根拠付きで返す。

## 完了規約

- 承認済み実装範囲 内の成果だけが返却されている。
- 検証、未実行項目、残留リスク が 根拠参照 付きで整理されている。
- 観測済み事実 と 仮説 を分けた。
- 承認済み実装範囲 内の 根拠 だけを扱った。
- temporary changes と cleanup_status を確認した。
- 必須 根拠: 承認済み実装範囲, コマンド or observation 根拠 when executed
- 完了判断材料: implement_lane が implement、tests、または implement-lane への戻しを判断できる。
- 残留リスク: 未確認事項と理由が返っている。

## 停止規約

- 恒久修正を行う時
- プロダクトテスト を追加する時
- design-time investigation を行う時
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 拒否条件: 不足 単一引き継ぎ入力
- 拒否条件: 不足 承認記録
- 拒否条件: unclear 承認済み実装範囲
- 停止条件: 一時観測点を安全に除去できない
- 停止条件: 設計判断が不足している
- 停止条件: 承認済み実装範囲 外の調査が必要である
