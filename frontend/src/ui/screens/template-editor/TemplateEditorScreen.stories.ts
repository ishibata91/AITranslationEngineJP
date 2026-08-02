import type { Meta, StoryObj } from "@storybook/svelte-vite"
import { themes } from "storybook/theming"
import { ScreenDocsPage } from "../screen-docs"
import { screenStateDescription } from "../screen-spec"
import TemplateEditorScreen from "./TemplateEditorScreen.svelte"
import { templateEditorScreenStates } from "./template-editor-screen-specs"

// プロンプトテンプレート画面。
const meta = {
  title: "Screens/プロンプトテンプレート",
  component: TemplateEditorScreen,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen",
    docs: { page: ScreenDocsPage, theme: themes.dark }
  },
  args: {
    onFieldInput: () => {},
    onInstructionInput: () => {},
    onTabChange: () => {},
    onSave: () => {},
    onReset: () => {}
  }
} satisfies Meta<typeof TemplateEditorScreen>

export default meta
type Story = StoryObj<typeof meta>

const { baseTab, recordTab, recordTabToneDefaultEdited, recordTabDirty } =
  templateEditorScreenStates

export const BaseTab: Story = {
  name: baseTab.storyName,
  args: { ...baseTab.args },
  parameters: {
    docs: { description: { story: screenStateDescription(baseTab) } }
  }
}

export const RecordTab: Story = {
  name: recordTab.storyName,
  args: { ...recordTab.args },
  parameters: {
    docs: { description: { story: screenStateDescription(recordTab) } }
  }
}

export const RecordTabToneDefaultEdited: Story = {
  name: recordTabToneDefaultEdited.storyName,
  args: { ...recordTabToneDefaultEdited.args },
  parameters: {
    docs: {
      description: { story: screenStateDescription(recordTabToneDefaultEdited) }
    }
  }
}

export const RecordTabDirty: Story = {
  name: recordTabDirty.storyName,
  args: { ...recordTabDirty.args },
  parameters: {
    docs: { description: { story: screenStateDescription(recordTabDirty) } }
  }
}
