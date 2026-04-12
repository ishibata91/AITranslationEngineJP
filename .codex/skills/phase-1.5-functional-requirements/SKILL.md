---
name: phase-1.5-functional-requirements
description: 第1.5段階の機能要件固定を担当し、active exec-plan の `機能要件` section に機能要件と非機能要件を固定する。
---

# Phase 1.5 Functional Requirements

## Goal

- 第1段階の facts / constraints / gaps を前提に、機能要件と非機能要件として次工程へ渡す境界を固定する
- active exec-plan には `機能要件` section の要約だけを残し、後続の UI モック作成と前段 HITL が読める形にする
- in-scope / non_functional_requirements / out-of-scope / open_questions を混在させずに整理する

## Output

- functional_requirements_summary
- in_scope
- non_functional_requirements
- out_of_scope
- open_questions
- required_reading

## Rules

- 更新対象は active exec-plan の `機能要件` section に限定する
- `UI モック` / `Scenario テスト一覧` / `実装計画` / `review 用差分図` はまだ作らない
- 実装コード、Scenario artifact、test file を変更しない
- facts と constraints から機能要件と非機能要件を整理し、未確定事項は `open_questions` として分離する
- human 判断が必要な論点は `open_questions` に残し、確定済み要件として扱わない
- 設計判断が揺れている間は次工程へ渡さない

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-1.5-functional-requirements.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-1.5-functional-requirements.to.orchestrating-implementation.json` を返却契約として使う。
