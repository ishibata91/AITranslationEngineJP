import type { Meta, StoryObj } from "@storybook/svelte-vite"
import { themes } from "storybook/theming"
import { ScreenDocsPage } from "../screen-docs"
import { screenStateDescription } from "../screen-spec"
import TargetPluginsScreen from "./TargetPluginsScreen.svelte"
import { targetPluginScreenStates } from "./target-plugins-screen-specs"

// 翻訳対象プラグイン画面。
const meta = {
  title: "Screens/翻訳対象プラグイン",
  component: TargetPluginsScreen,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen",
    docs: { page: ScreenDocsPage, theme: themes.dark }
  },
  args: {
    onSelectNewPlugin: () => {},
    onNewPluginPathInput: () => {},
    onProceedToRun: () => {},
    onOpenPlugin: () => {},
    onRequestDelete: () => {},
    onConfirmDelete: () => {},
    onCancelDelete: () => {}
  }
} satisfies Meta<typeof TargetPluginsScreen>

export default meta
type Story = StoryObj<typeof meta>

const { empty, loading, list, selected, confirmDelete, deleting, errored } =
  targetPluginScreenStates

export const Empty: Story = {
  name: empty.storyName,
  args: { ...empty.args },
  parameters: {
    docs: { description: { story: screenStateDescription(empty) } }
  }
}

export const Loading: Story = {
  name: loading.storyName,
  args: { ...loading.args },
  parameters: {
    docs: { description: { story: screenStateDescription(loading) } }
  }
}

export const List: Story = {
  name: list.storyName,
  args: { ...list.args },
  parameters: { docs: { description: { story: screenStateDescription(list) } } }
}

export const Selected: Story = {
  name: selected.storyName,
  args: { ...selected.args },
  parameters: {
    docs: { description: { story: screenStateDescription(selected) } }
  }
}

export const ConfirmDelete: Story = {
  name: confirmDelete.storyName,
  args: { ...confirmDelete.args },
  parameters: {
    docs: { description: { story: screenStateDescription(confirmDelete) } }
  }
}

export const Deleting: Story = {
  name: deleting.storyName,
  args: { ...deleting.args },
  parameters: {
    docs: { description: { story: screenStateDescription(deleting) } }
  }
}

export const Errored: Story = {
  name: errored.storyName,
  args: { ...errored.args },
  parameters: {
    docs: { description: { story: screenStateDescription(errored) } }
  }
}
