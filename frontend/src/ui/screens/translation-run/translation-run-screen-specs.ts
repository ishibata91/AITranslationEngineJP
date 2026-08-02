import { defineScreenState } from "../screen-spec"
import {
  OPENAI_READY_STATE,
  OPENAI_RUNNING_STATE,
  OPENAI_PAUSED_STATE,
  OPENAI_NO_UNTRANSLATED_STATE,
  OPENAI_BATCH_UNTRANSLATED_STATE,
  OPENAI_FAILED_STATE
} from "./translation-run.fixtures"

export const translationRunScreenStates = {
  notStarted: defineScreenState({
    storyName: "未開始",
    precondition: "OpenAIの接続情報が揃い、保存済みのbatch進行がない。",
    args: OPENAI_READY_STATE,
    specs: [
      {
        id: "translation-run.not-started.phase-badge",
        statement: "状態バッジに「未実行」を表示する。"
      },
      {
        id: "translation-run.not-started.main-action",
        statement: "「バッチ実行」を有効にする。"
      },
      {
        id: "translation-run.not-started.manual-actions-hidden",
        statement: "状態確認ボタンと手動取り込みボタンを表示しない。"
      }
    ]
  }),
  running: defineScreenState({
    storyName: "実行中",
    precondition: "固有名段のbatchを処理している。",
    args: OPENAI_RUNNING_STATE,
    specs: [
      {
        id: "translation-run.running.phase-badge",
        statement: "状態バッジに「実行中」を表示する。"
      },
      {
        id: "translation-run.running.main-action",
        statement: "処理中の表示を伴う「実行中…」を無効にする。"
      },
      {
        id: "translation-run.running.manual-actions-hidden",
        statement: "状態確認ボタンと手動取り込みボタンを表示しない。"
      }
    ]
  }),
  paused: defineScreenState({
    storyName: "途中停止",
    precondition: "本文段の保存済み進行があり、処理が止まっている。",
    args: OPENAI_PAUSED_STATE,
    specs: [
      {
        id: "translation-run.paused.phase-badge",
        statement: "状態バッジに「未実行」を表示する。"
      },
      {
        id: "translation-run.paused.main-action",
        statement: "「バッチ実行を再開」を有効にする。"
      },
      {
        id: "translation-run.paused.manual-actions-hidden",
        statement: "状態確認ボタンと手動取り込みボタンを表示しない。"
      }
    ]
  }),
  done: defineScreenState({
    storyName: "完了（未訳なし）",
    precondition: "固有名段と本文段が完了し、未訳が残っていない。",
    args: OPENAI_NO_UNTRANSLATED_STATE,
    specs: [
      {
        id: "translation-run.done.phase-badge",
        statement: "状態バッジに「完了」を表示する。"
      },
      {
        id: "translation-run.done.main-action",
        statement: "「完了」を無効にする。"
      },
      {
        id: "translation-run.done.manual-actions-hidden",
        statement: "状態確認ボタンと手動取り込みボタンを表示しない。"
      },
      {
        id: "translation-run.done.notice",
        statement: "本文の取り込みと翻訳の完了を案内する。"
      }
    ]
  }),
  doneWithUntranslated: defineScreenState({
    storyName: "完了（未訳あり）",
    precondition: "固有名段と本文段が完了し、未訳が3件残っている。",
    args: OPENAI_BATCH_UNTRANSLATED_STATE,
    specs: [
      {
        id: "translation-run.done-untranslated.phase-badge",
        statement: "状態バッジに「完了」を表示する。"
      },
      {
        id: "translation-run.done-untranslated.main-action",
        statement: "「未訳だけを再送信」を有効にする。"
      },
      {
        id: "translation-run.done-untranslated.manual-actions-hidden",
        statement: "状態確認ボタンと手動取り込みボタンを表示しない。"
      },
      {
        id: "translation-run.done-untranslated.notice",
        statement: "未訳件数を表示する。"
      }
    ]
  }),
  failed: defineScreenState({
    storyName: "失敗表示",
    precondition: "保存済みの本文段で外部batchが失敗し、処理が止まっている。",
    args: OPENAI_FAILED_STATE,
    specs: [
      {
        id: "translation-run.failed.phase-badge",
        statement: "状態バッジに「失敗」を表示する。"
      },
      {
        id: "translation-run.failed.main-action",
        statement: "「バッチ実行を再開」を有効にする。"
      },
      {
        id: "translation-run.failed.manual-actions-hidden",
        statement: "状態確認ボタンと手動取り込みボタンを表示しない。"
      },
      {
        id: "translation-run.failed.error",
        statement: "外部batch IDと失敗理由を表示する。"
      }
    ]
  })
} as const
