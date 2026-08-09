---
name: context-maintenance
description: メインエージェントが skill、agents、AGENTS.md などの AI コンテキストを編集する時に使う。サブエージェントは本スキルを使用しない。
---

## 入力

-　skill,agents,AGENT.mdいずれかの変更依頼

---
## 共通規約

# **望ましい結果が得られる確率を最大化する、高シグナルなトークンの最小集合を見つけること。**

- **システムプロンプト**は極めて明確でなければならず、シンプルで直接的な言葉を使い、エージェントにとって*適切な抽象度（right altitude）*で概念を提示する。
	- 適切な抽象度とは、よくある2つの失敗パターンの中間にある「ちょうどよい領域」。
	- 一方の極端では、エージェントに正確な振る舞いをさせるため、複雑で壊れやすいロジックをプロンプトにハードコードしているケースがある。
	- この方法は脆弱性を生み、時間とともに保守の複雑さを増大させる。
	- もう一方の極端では、曖昧で高レベルすぎる指示を与えるケースがある。
	- この中間を狙う。
- システムプロンプトはメモ書きではない
	- 作成時の指摘や，目的を混ぜ込まない。
	- 期待する振る舞いを完全に説明できる**最小限の情報集合**を目指す。
- few-shot prompting
	- エージェントに期待される振る舞いを効果的に示す、**多様で代表的な例の集合**を厳選すること
- just-in-time
	- 記載する内容が，常に読み込まれるべきものか，その時々に読み込まれるものなのかを考える
- 指示の粒度
	- 指示は、抽象的すぎると解釈の自由度が高くなりすぎ、詳細すぎると特定ケースにしか適用できなくなるため、再利用可能な判断基準として機能する抽象度に置く。
**抽象的すぎる**

> 必要な情報を調査してください。

**詳細すぎる**

> `src/auth/login.ts`、`src/auth/session.ts`、`tests/auth/login.test.ts`をこの順で読み、`rg "SessionId"`を実行してください。

**適切な抽象度**

> 変更対象の振る舞いを理解するため、実装・呼び出し元・関連テストを探索し、必要な範囲だけコンテキストに取り込んでください。
---
##  知識は知るべき者へ，知るべき時にだけ与えられる

- {REPO_ROOT}/protocolsはレポルートと同じフォルダ階層を持つ。
	- protocols配下に配置されたmdは，レポルートの対応したフォルダ配下のファイルや，新規ファイルの書き込み時に自動的にhookで発火してコンテキストに注入される。
	- 例
			1. protocols/internal/hoge/fuga.mdが存在する
			2. internal/hogeのファイルや，新規ファイルの書き込みを行う
			3. fuga.mdがコンテキストとして注入される。

## プロトコル

### 書き方
- メモ書きではない。
### テンプレート
- index.md
	- 地図になる。詳細を書かず，そのフォルダやファイルが担う責務のみを日本語で簡潔に，一言で記述する。
	```markdown
	フォルダ: frontend/src
	フロントのコードベース。
	- `ui/` は画面表示コードを持つ
	- `gateway/` は外部通信を担当する。

	```

- coding.md
	- そのパッケージのコーディング規約になる。
	```markdown
	# HTTP delivery

	`delivery/http/` は HTTP request と domain 操作の接続だけを担当する。

	- request の decode と response の encode は `delivery/http/` に置く。
	- HTTP status code の選択は `delivery/http/` に置く。
	- domain の業務判断は `delivery/http/` に置かない。
	- database、外部 API、認証基盤を直接呼ばない。
	- handler は application service の interface だけに依存する。
	- HTTP 固有の型を application service へ渡さない。

	# few-shots
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var request createUserRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        writeBadRequest(w, err)
        return
    }

    user, err := h.users.Create(r.Context(), app.CreateUserInput{
        Name: request.Name,
    })
    if err != nil {
        writeApplicationError(w, err)
        return
    }

    writeJSON(w, http.StatusCreated, userResponse{ID: user.ID})
}
	```
- reference.md
	- 明確に必要な場合のみ。
	- そのフォルダや，ファイル固有で利用する外部API仕様，メモなどを記載する。

---
## SKILL

スキルとは知識であり，手順書である。

- スキルが，タスクを終えられる最低限の情報を与える。
	- 最低限から始める。足りなければ補う。
	- 最初からフルで大量の情報を与えない。
- スキルが担うタスクは一つだけ。
	- スキルに複数の責務を与えない。
- descriptionはスキルをいつ呼び出すか，メイン，サブエージェントどちら向けか，のみ記載する。
	- 例 メインエージェント使用禁止。Codegraphを用いたコードベース探索用のスキル。コードベースの探索時に使う。
- 手順の明確に分かれたセクションを用いる。
---
## AGENTS

AGENTSは担当者のツールと脳を定義する設定である。

- **知識置き場ではない**
- スキルを扱う上で，適したツールとモデルを渡す。それだけ。
- モデルは安いモデルから始める。最初から高級なモデルを使わない。推論の質に問題が出た時のみ，1段ずつ高くしていく

---
## レビュー

- コンテキストを変更した後、`context_reviewer` に変更対象をレビューさせる。
- 初回を含めて最大3回まで、同じ `context_reviewer` を再開してレビューさせる。
- 否決された場合は指摘を反映し、同じ `context_reviewer` に再レビューさせる。
- 通過した場合はレビューを終了する。
- 3回目も通過しない場合は、否決理由と未解消箇所を人間へ返す。

---
