import type { Meta, StoryObj } from "@storybook/svelte-vite"
import { themes } from "storybook/theming"
import { ScreenDocsPage } from "../screen-docs"
import { screenStateDescription } from "../screen-spec"
import PrebuiltDictionaryEditorScreen from "./PrebuiltDictionaryEditorScreen.svelte"
import { prebuiltDictionaryEditorScreenStates } from "./prebuilt-dictionary-editor-screen-specs"

const meta = {
  title: "Screens/用語辞書",
  component: PrebuiltDictionaryEditorScreen,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen",
    docs: { page: ScreenDocsPage, theme: themes.dark }
  },
  args: {
    onFilterInput: () => {},
    onCreate: () => {},
    onDelete: () => {},
    onToggleCategories: () => {},
    onConfirmChanges: () => {},
    onCancelChanges: () => {},
    onPrev: () => {},
    onNext: () => {}
  }
} satisfies Meta<typeof PrebuiltDictionaryEditorScreen>

export default meta
type Story = StoryObj<typeof meta>

const { list, edit, empty } = prebuiltDictionaryEditorScreenStates

export const List: Story = {
  name: list.storyName,
  args: list.args,
  parameters: { docs: { description: { story: screenStateDescription(list) } } }
}

export const Edit: Story = {
  name: edit.storyName,
  args: edit.args,
  parameters: { docs: { description: { story: screenStateDescription(edit) } } }
}

export const Empty: Story = {
  name: empty.storyName,
  args: empty.args,
  parameters: { docs: { description: { story: screenStateDescription(empty) } } }
}
