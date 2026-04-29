---
name: implementation-investigate
description: Codex implementation lane 側の実装時調査の共通作業プロトコル。single_handoff_packet 1 件内で evidence first に調査する判断基準を提供する。
---
# Implementation Investigate

## 目的

`implementation-investigate` は作業プロトコルである。
`implementation_investigator` agent が、`single_handoff_packet` 1 件と owned_scope 内で実装時の証拠を集める時の共通判断を提供する。

tool policy は [implementation_investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_investigator.toml) が持ち、handoff は skill に従う。

## 対応ロール

- `implementation_investigator` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implementation-investigate` の出力規約で固定する。

## 入力規約

- 実装前再現、trace、再観測を行う時
- 一時観測点を add / remove する時
- evidence と仮説を分けて返す時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: single_handoff_packet, approval_record, owned_scope, investigation_request, validation_commands
- 任意入力: knowledge_focus, reproduction_evidence, temporary_observation_plan
- input_notes: {"single_handoff_packet": "implementation-scope から抽出済みの handoff 1 件だけ。full implementation-scope、active work plan 全文、source artifacts、後続 handoff は入力に含めない。", "knowledge_focus": "implementation-investigate-reproduce、trace、observe、reobserve の参照ヒント。共通規約と完了条件は変えない。"}

## 外部参照規約

- agent runtime と tool policy は [implementation_investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_investigator.toml) の `allowed_write_paths` / `allowed_commands` とする。
- agent runtime: [implementation_investigator.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_investigator.toml)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate-reproduce/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate-trace/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate-observe/SKILL.md, /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate-reobserve/SKILL.md

## 内部参照規約

### 拘束観点

- evidence first の観測
- observed facts と hypotheses の分離
- temporary observation の cleanup
- `agent-browser` CLI による UI / console / screenshot evidence
- focused skill の選び方

### Browser Evidence

UI state、console、screenshot の観測は `execute` から `agent-browser` CLI を使う。

標準入口は次の通りである。

```bash
npm run dev:wails:agent-browser
agent-browser doctor --offline --quick
agent-browser open http://localhost:34115
agent-browser snapshot
agent-browser console
agent-browser screenshot tmp/agent-browser/ui-evidence.png
```

観測後は `agent-browser close` を実行する。
system test の Playwright runner は product test 用の別入口として扱う。

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
agent-browser screenshot tmp/agent-browser/ui-evidence.png
agent-browser screenshot --annotate --screenshot-dir tmp/agent-browser
agent-browser network requests
```

複数手順をまとめる場合:

```bash
agent-browser batch --bail \
  "open http://localhost:34115" \
  "snapshot" \
  "console" \
  "screenshot tmp/agent-browser/ui-evidence.png"
```

`@e2` のような ref は直前の `snapshot` の結果から選ぶ。
selector が安定している場合だけ CSS selector を使う。
console / errors / screenshot / network requests は completion packet の evidence に command と結果を残す。

- 参照 pattern は [investigation-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-investigate/references/patterns/investigation-patterns.md) とする。

## 判断規約

- `single_handoff_packet` 1 件と owned_scope を超えない
- evidence のない結論を固定しない
- Codex implementation lane のブラウザ操作は `agent-browser` CLI で行う
- 一時観測点は返却前に除去する
- 恒久修正と product test 追加を混ぜない

- 観測条件、command、結果を残す
- temporary changes と cleanup_status を返す
- recommended next step を根拠付きで返す
- active 規約 は agent 1:1。調査種別の差分は focused skill で扱い、output obligation はこの 規約 に固定する。

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力に tool policy、agent runtime、product code の変更義務を含めない。
- 返却先: implement_lane
- 必須出力: investigation_focus, reproduction_status, observed_facts, hypotheses, observation_points, temporary_changes, cleanup_status, validation_results, remaining_gaps, residual_risks, recommended_next_step
- 出力 field 要件: {"observed_facts": "観測済み事実だけを書く。仮説と混ぜない", "temporary_changes": "一時観測点を使った場合だけ path と目的を返す。未使用なら none", "cleanup_status": "一時観測点の除去状態を必ず返す。未使用なら not_applicable", "recommended_next_step": "implement、tests、reroute のどれかを根拠付きで返す"}

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- observed facts と hypotheses を分けた。
- owned_scope 内の evidence だけを扱った。
- temporary changes と cleanup_status を確認した。
- 必須 evidence: owned_scope, command or observation evidence when executed
- completion signal: implement_lane が implement、tests、または implement-lane reroute を判断できる
- residual risk key: remaining_gaps

## 停止規約

- 恒久修正を行う時
- product test を追加する時
- design-time investigation を行う時
- 恒久修正を同時に行わない
- product test を追加しない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
- 恒久修正を同時に行わなかった場合は停止する。
- product test を追加しなかった場合は停止する。
- mode 別 個別 JSON 規約 を使わなかった場合は停止する。
- 拒否条件: missing single_handoff_packet
- 拒否条件: missing approval_record
- 拒否条件: unclear owned_scope
- 停止条件: 一時観測点を安全に除去できない
- 停止条件: 設計判断が不足している
- 停止条件: owned_scope 外の調査が必要である
