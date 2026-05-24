import type { Meta, StoryObj } from "@storybook/svelte-vite"

import BodyTranslationPhasePanel from "../BodyTranslationPhasePanel.svelte"
import { bodyTranslationPhasePanelFixture } from "../__fixtures__/body-phase-card-fixture"

const meta = {
  title: "Screen Components/Body Translation Phase/BodyTranslationPhasePanel",
  component: BodyTranslationPhasePanel,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    viewModel: bodyTranslationPhasePanelFixture,
    onAction: () => {}
  }
} satisfies Meta<typeof BodyTranslationPhasePanel>

export default meta

type Story = StoryObj<typeof meta>

export const AiSettingsReady: Story = {}

export const AiSettingsMissing: Story = {
  args: {
    viewModel: {
      ...bodyTranslationPhasePanelFixture,
      modelLabel: "-",
      credentialRefLabel: "-"
    }
  }
}

export const RunningLocked: Story = {
  args: {
    viewModel: {
      ...bodyTranslationPhasePanelFixture,
      viewState: "running",
      phaseStateLabel: "実行中",
      statusTitle: "本文翻訳を実行中です"
    }
  }
}
