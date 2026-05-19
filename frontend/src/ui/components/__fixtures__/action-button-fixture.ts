import type { ComponentProps } from "svelte"
import ActionButton from "../ActionButton.svelte"
import IconActionButton from "../IconActionButton.svelte"

type ActionButtonProps = ComponentProps<typeof ActionButton>
type IconActionButtonProps = ComponentProps<typeof IconActionButton>

const noop = (): void => {}

export const primaryActionButtonFixture: ActionButtonProps = {
  label: "保存する",
  variant: "primary",
  onClick: noop
}

export const secondaryActionButtonFixture: ActionButtonProps = {
  label: "戻る",
  variant: "secondary",
  onClick: noop
}

export const dangerActionButtonFixture: ActionButtonProps = {
  label: "削除する",
  variant: "danger",
  onClick: noop
}

export const busyActionButtonFixture: ActionButtonProps = {
  label: "処理中",
  variant: "primary",
  busy: true,
  onClick: noop
}

export const disabledActionButtonFixture: ActionButtonProps = {
  label: "入力完了後に実行できます",
  variant: "secondary",
  disabled: true,
  onClick: noop
}

export const longActionButtonFixture: ActionButtonProps = {
  label: "長い文言の操作でも折り返して対象の操作名を読める状態",
  variant: "primary",
  onClick: noop
}

export const iconActionButtonFixture: IconActionButtonProps = {
  ariaLabel: "モデル一覧を更新",
  title: "モデル一覧を更新",
  onClick: noop
}

export const iconActionButtonBusyFixture: IconActionButtonProps = {
  ariaLabel: "更新中",
  busy: true,
  onClick: noop
}

export const iconActionButtonDisabledFixture: IconActionButtonProps = {
  ariaLabel: "更新できません",
  disabled: true,
  onClick: noop
}
