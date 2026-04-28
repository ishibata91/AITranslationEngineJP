---
name: skill-modification
description: Codex 側の workflow / skill / agent 変更で、何を受け取り、何へ従い、どう判断し、何を返し、どこで止まるかを固定する。
---

# Skill Modification

## 目的

`skill-modification` は、`.codex` 配下の workflow、skill、agent runtime、実行権限を変更するための作業プロトコルである。

担当ロールが、skill / agent の責務境界、正本配置、入力規約、出力規約、廃止対象の削除を判断する時に使う。

## 対応ロール

- この skill は人間から直接依頼された Codex 本体だけが使う。
- サブエージェントまたは別 agent から呼び出された場合は拒否する。
- この skill は `.codex` 配下の workflow / skill / agent 変更 task を対象にする。

## 入力規約

- 入力は、変更目的、変更対象パス、非対象、残す正本、削除してよい対象を含む。
- skill / agent / 実行権限を変更する場合は、変更対象名とファイルパスを必須にする。
- 削除がある場合は、削除理由と削除後の正本確認先を必須にする。
- 入力が不足する場合は、推測で変更せず、人間へ不足項目を返す。

## 外部参照規約

- 最上位の workflow 正本は `.codex/README.md` とする。
- skill 形式の正本は `.codex/skill-template.md` とする。
- agent runtime の正本形式は `.codex/agent-template.toml` とする。
- agent binding と実行権限の配置は `.codex/README.md` の agent / skill 配置規約に従う。
- product 仕様の正本は `docs/` であり、この skill は product 仕様を変更しない。
- 外部正本が衝突する場合は、`.codex/README.md` を優先し、衝突内容を人間へ返す。

## 内部参照規約

- なし。

## 判断規約

- skill は作業プロトコルであり、手順、標準 pattern、参照タイミング一覧、知識範囲一覧を持たない。
- `.codex/skill-template.md` の `Maintenance（テンプレート外）` は、個別 skill や agent に書き込む内容ではない。
- agent は実行主体であり、agent runtime と実行権限を持つ。
- 入力規約、出力規約、完了規約、停止規約は skill 側へ統合する。
- 固定ファイル名、固定ディレクトリ名、既存 key、既存 command、既存 runtime 名以外は日本語で書く。
- 英語名を併記する必要がある場合は、先に日本語の意味を書き、括弧内に固定名を書く。
- 責務、権限、入力、出力が分かれる場合は、分岐記述ではなく skill または agent の分割を選ぶ。
- live workflow にない artifact、廃止ファイル、legacy pointer、stub、禁止文言だけの説明ファイルは残さない。
- 論理名と変更対象名は、同じ判断単位で対応が分かるように書く。

## 出力規約

- 出力は、更新する `.codex` 配下ファイル、削除するファイル、削除理由、検証方法を含む。
- skill 本体を更新する場合は、`.codex/skill-template.md` の分類順に合わせる。
- agent を更新する場合は、agent runtime、binding、実行権限の所有者を分けて示す。
- 削除対象がある場合は、削除後にどの正本を見ればよいかを示す。
- 出力に product code、product test、product 仕様 docs の変更を含めてはいけない。

## 完了規約

- `.codex/README.md`、`.codex/skill-template.md`、`.codex/agent-template.toml`、対象 `SKILL.md` との矛盾がない時に完了とする。
- skill 本体に手順、標準 pattern、参照タイミング一覧、知識範囲一覧が残っていないことを確認する。
- 実行権限、書き込み範囲、出力義務が skill 側へ混入していないことを確認する。
- 削除対象に legacy pointer、stub、禁止文言だけのファイルが残っていないことを確認する。

## 停止規約

- サブエージェントまたは別 agent から呼び出された場合は拒否する。
- product code または product test の変更が混ざる場合は停止する。
- product 仕様 docs の正本化が混ざる場合は停止する。
- 対応する実行権限ファイルが必要だが見つからない場合は停止し、不足ファイルを返す。
- 権限境界、書き込み範囲、agent 所有者が不明な場合は停止する。
- `.codex/README.md` と変更方針が衝突する場合は停止し、衝突箇所を返す。
- 削除理由または削除後の正本確認先がない削除は停止する。
