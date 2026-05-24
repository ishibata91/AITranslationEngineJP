import type { Meta, StoryObj } from "@storybook/svelte-vite"

import { createProviderSettingsPageControllerFixture } from "../../__fixtures__/screen-page-controller-fixtures"
import ProviderSettingsPage from "../ProviderSettingsPage.svelte"

const meta = {
  title: "Screens/Provider Settings/ProviderSettingsPage",
  component: ProviderSettingsPage,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    createController: createProviderSettingsPageControllerFixture()
  }
} satisfies Meta<typeof ProviderSettingsPage>

export default meta

type Story = StoryObj<typeof meta>

export const Disconnected: Story = {}
