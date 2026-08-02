import type { Meta, StoryObj } from "@storybook/svelte-vite"
import { themes } from "storybook/theming"
import { ScreenDocsPage } from "../screen-docs"
import { screenStateDescription } from "../screen-spec"
import TranslationRunScreen from "./TranslationRunScreen.svelte"
import { translationRunScreenStates } from "./translation-run-screen-specs"

// 状態確認と手動取り込みを削除し、一つの主操作で自動処理するbatch画面の代表状態。
const meta = {
  title: "Screens/翻訳実行",
  component: TranslationRunScreen,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen",
    docs: { page: ScreenDocsPage, theme: themes.dark }
  },
  args: {
    onFieldInput: () => {},
    onLoadModels: () => {},
    onRun: () => {},
    onPagePrev: () => {},
    onPageNext: () => {},
    onUntranslatedOnlyChange: () => {},
    onProviderChange: () => {},
    onSubmit: () => {}
  }
} satisfies Meta<typeof TranslationRunScreen>

export default meta
type Story = StoryObj<typeof meta>

const { notStarted, running, paused, done, doneWithUntranslated, failed } =
  translationRunScreenStates

export const NotStarted: Story = {
  name: notStarted.storyName,
  args: { ...notStarted.args },
  parameters: {
    docs: { description: { story: screenStateDescription(notStarted) } }
  }
}

export const Running: Story = {
  name: running.storyName,
  args: { ...running.args },
  parameters: {
    docs: { description: { story: screenStateDescription(running) } }
  }
}

export const Paused: Story = {
  name: paused.storyName,
  args: { ...paused.args },
  parameters: {
    docs: { description: { story: screenStateDescription(paused) } }
  }
}

export const Done: Story = {
  name: done.storyName,
  args: { ...done.args },
  parameters: { docs: { description: { story: screenStateDescription(done) } } }
}

export const DoneWithUntranslated: Story = {
  name: doneWithUntranslated.storyName,
  args: { ...doneWithUntranslated.args },
  parameters: {
    docs: {
      description: { story: screenStateDescription(doneWithUntranslated) }
    }
  }
}

export const Failed: Story = {
  name: failed.storyName,
  args: { ...failed.args },
  parameters: {
    docs: { description: { story: screenStateDescription(failed) } }
  }
}
