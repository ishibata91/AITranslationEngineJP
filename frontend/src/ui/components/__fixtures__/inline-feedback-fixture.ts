import type { ComponentProps } from "svelte"
import InlineFeedback from "../InlineFeedback.svelte"

type InlineFeedbackProps = ComponentProps<typeof InlineFeedback>

const noop = (): void => {}

export const neutralInlineFeedbackFixture: InlineFeedbackProps = {
  tone: "neutral",
  title: "確認",
  message: "補助情報を表示します。"
}

export const errorInlineFeedbackFixture: InlineFeedbackProps = {
  tone: "error",
  title: "保存できません",
  message: "入力値を確認してください。",
  actionLabel: "再確認",
  onAction: noop
}

export const warningInlineFeedbackFixture: InlineFeedbackProps = {
  tone: "warning",
  title: "操作できません",
  message: "必要な入力が完了すると操作できます。"
}

export const successInlineFeedbackFixture: InlineFeedbackProps = {
  tone: "success",
  title: "完了",
  message: "処理が完了しました。"
}

export const longInlineFeedbackFixture: InlineFeedbackProps = {
  tone: "error",
  title: "長い内容",
  message:
    "長い path label や長いエラー文を想定した表示です。aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  actionLabel: "再実行",
  onAction: noop
}
