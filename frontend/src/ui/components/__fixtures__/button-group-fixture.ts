import type { ComponentProps } from "svelte"
import ButtonGroup from "../ButtonGroup.svelte"

type ButtonGroupProps = ComponentProps<typeof ButtonGroup>

export const buttonGroupFixture: ButtonGroupProps = {
  ariaLabel: "Storybook 操作",
  align: "end"
}

export const buttonGroupStretchFixture: ButtonGroupProps = {
  ariaLabel: "Storybook 幅固定操作",
  align: "stretch"
}
