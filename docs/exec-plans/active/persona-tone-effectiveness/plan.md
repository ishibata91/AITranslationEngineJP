# Task Plan: persona-tone-effectiveness

## 状態

- 2026-08-01に最終確認v3まで完了した。
- v3は属性の対応正解率をすべて満たした。機械集計は記号的な口調18/72を失敗としたが、人間確認により対象属性に合う自然な表現まで落とす過剰判定と確定した。
- 意味、書式、自然さによる品質合格は65/72で終了条件を満たした。人間判断によりv3を採用候補とした。
- 固定した終了条件に従い、追加のAPI実験を行わずに終了した。

## 事象

- 「ぞんざい」のように強い口調指定は訳文へ現れるが、「平明・女性」では女性の話者らしい口調が訳文へ現れない。
- 話者を解決できない汎用台詞でも、条件から性別を取得できる場合は男性と女性で口調が明確に変わってほしい。
- 性別だけでなく、ペルソナ・性別・年齢・種族の組み合わせが訳文へ一貫して現れてほしい。

## 対象

- 変える対象は、翻訳 prompt へ入れる persona ごとの few-shot とする。名指し話者の「平明・男性」「平明・女性」と、汎用台詞の「男性」「女性」へ、同じ英語原文に対する性別別の日本語例を複数与える。
- 変えない対象は、基底口調の分類、性別の抽出、話者解決、台詞感情の分類、base 翻訳指示、口調 directive の共通指示、OpenAI Batch API の実装とする。persona の手掛かりだけの効果を分けて測るために変えない。
- 翻訳対象モデルは `gpt-5.6-luna` に固定する。
- 訳文を読む役割は Codex の `fresh` とし、API の判定要求は使わない。
- 追加の予備比較では、現行の説明文と1例を含む指示に対し、ペルソナ名・性別・年齢・種族名だけを渡す方式を比べる。組み合わせごとの新しい説明文とfew-shotはまだ作らない。

## 砂場

- 場所は `tmp/persona-tone-effectiveness/` とする。凍結した標本、Batch API の入力と出力、測定用の道具、`fresh` が読んだ記録を置く。
- `.gitignore` の `tmp` 規則に当たり、`git check-ignore -v tmp/persona-tone-effectiveness/probe.jsonl` で除外を確認した。砂場は commit に入れない。

## 接続先

- OpenAI API の `Batch API` を使う。endpoint は `/v1/chat/completions`、completion window は `24h` とする。
- 翻訳モデルは `gpt-5.6-luna` とする。`GET /v1/models` で利用可能であることを 2026-08-01 に確認した。
- API key は macOS Keychain の service `OPENAI_API_KEY` から実行時だけ読む。file と command 出力へ書かない。
- 標本の原文と公式日本語訳は、前の実験が凍結した `tmp/dialogue-tone-naturalness/dataset/` から読む。前の実験の prompt 例と同じ原文は除く。
- OpenAI の現在の公式資料に従い、JSONL を `purpose=batch` で upload し、`custom_id` で結果を対応付ける。

## branch 情報

- 作業 branch: `codex/persona-tone-effectiveness`
- 統合先 branch: `master`
- 分岐元 commit: `b090160aea8598ce80c3091b47f35588d3d5955b`
- worktree: `/Users/iorishibata/.codex/worktrees/persona-tone-effectiveness`

## やらないこと

- persona の判定基準や自動評価手法に関する外部調査は行わない。人間の指示により、Codex の `fresh` が標本を読む。
- `gpt-5.4-mini` の無料枠残量は確認せず、利用しない。
- 評価用標本は準備の段と探索の段では翻訳しない。達成の線を満たした回の確証だけに使う。
- 女性らしさを「だわ」「のよ」など特定の語尾1個へ固定しない。戯画的な女性語を増やすだけの変更は採らない。
- プロダクト正本と本番実装は準備の段では変更しない。
