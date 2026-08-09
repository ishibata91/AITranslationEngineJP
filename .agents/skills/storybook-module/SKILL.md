---
name: storybook-module
description: Codex 本体が Storybook に表示できる範囲の見た目を作成または修正し、人間が承認した状態を固定するスキル。layout、文言、style、表示構造、表示用 props、story、fixture が変わる時に使う。state、API、Wails bridge、ルーティング、副作用など、表示以外の振る舞いだけが変わる時は使わない。
---

# Storybook で見た目を固定する

## 責務

Storybook で画面の状態を再現し、人間が見た目を確認できる形にする。
人間が承認した story、fixture、表示コンポーネントを、その画面状態の見た目として固定する。

## 扱う範囲

- Svelte 表示コンポーネントの template、表示用 props、表示用 script、style。
- story と表示用 fixture。
- 画面に現れる文言、配置、寸法、色、余白、装飾、表示構造。
- 見た目が異なる画面状態を再現するための mock。

次の対象は扱わない。

- store などの state 管理。
- API と Wails bridge の呼び出し。
- ルーティングと画面遷移。
- 副作用とライフサイクル処理。
- validation などの業務判断。
- backend と統合境界。

見た目の確認に必要な振る舞いは、props、fixture、mock で置き換える。
表示以外の変更が必要になった場合は、必要な変更と理由を返して止める。

## 手順

1. 承認済みの要求と、変更対象に最も近い frontend 規約と Storybook 規約を読む。
2. 人間が見た目を判断するために必要な画面状態を列挙する。
3. 画面状態を story と fixture で再現する。
4. 表示コンポーネント、story、fixture の範囲だけで見た目を実装する。
5. リポジトリ既定の方法で Storybook を起動し、対象 story をブラウザで確認する。
6. 表示崩れ、意図しない overflow、文言欠落、状態間の不整合を修正する。
7. 対象 story を人間へ提示し、見た目の承認を待つ。
8. 人間の指摘を同じ表示コンポーネント、story、fixture へ反映し、承認まで確認を繰り返す。
9. 承認後に Storybook build と関連する frontend 検証を実行する。

## 人間レビュー

人間レビューは、Storybook 上の対象 story を確認できる状態で行う。
レビュー対象には、見た目が異なる各画面状態を含める。
人間の明示的な承認がない状態を、固定済みとして扱わない。

要求を越える新しい画面要素や画面状態を独自に追加しない。
人間の指摘が要求または表示範囲を越える場合は、変更せずに差分を返す。

## 固定する対象

- 承認された story と画面状態。
- 承認時に使った fixture と mock 入力。
- 承認された文言、表示構造、layout、style。
- 表示コンポーネントが受け取る表示用 props の形。

承認後の変更で固定対象が変わる場合は、同じ Storybook 人間レビューを再実施する。

## 完了条件

- 見た目が異なる必要な画面状態を Storybook で再現できる。
- 対象 story をブラウザで確認済みである。
- 人間が対象 story の見た目を承認している。
- Storybook build と関連する frontend 検証が通過している。

## 返す結果

- 承認された story と画面状態。
- 変更した表示コンポーネント、story、fixture、関連資源。
- Storybook build と関連する frontend 検証の結果。
- 表示範囲外として残した変更と理由。
