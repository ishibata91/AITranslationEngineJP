import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationResultListPanel from "../TranslationResultListPanel.svelte"
import {
  emptyTranslationResultListPanelFixture,
  translationResultListPanelFixture
} from "../__fixtures__/translation-complete-fixture"

const meta = {
  title: "Screens/Job Run/TranslationResultListPanel",
  component: TranslationResultListPanel,
  args: translationResultListPanelFixture
} satisfies Meta<typeof TranslationResultListPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Empty: Story = {
  args: emptyTranslationResultListPanelFixture
}
