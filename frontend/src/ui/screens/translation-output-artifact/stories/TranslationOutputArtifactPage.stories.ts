import type { Meta, StoryObj } from "@storybook/svelte-vite"

import { createTranslationOutputArtifactPageControllerFixture } from "../../__fixtures__/screen-page-controller-fixtures"
import TranslationOutputArtifactPage from "../TranslationOutputArtifactPage.svelte"

const meta = {
  title: "Screens/Translation Output Artifact/TranslationOutputArtifactPage",
  component: TranslationOutputArtifactPage,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    createController: createTranslationOutputArtifactPageControllerFixture()
  }
} satisfies Meta<typeof TranslationOutputArtifactPage>

export default meta

type Story = StoryObj<typeof meta>

export const Disconnected: Story = {}
