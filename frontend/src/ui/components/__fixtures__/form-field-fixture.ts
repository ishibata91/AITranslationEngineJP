import type { ComponentProps } from "svelte"
import CheckboxField from "../CheckboxField.svelte"
import FormField from "../FormField.svelte"
import SelectField from "../SelectField.svelte"
import TextAreaField from "../TextAreaField.svelte"
import TextInputField from "../TextInputField.svelte"

type CheckboxFieldProps = ComponentProps<typeof CheckboxField>
type FormFieldProps = ComponentProps<typeof FormField>
type SelectFieldProps = ComponentProps<typeof SelectField>
type TextAreaFieldProps = ComponentProps<typeof TextAreaField>
type TextInputFieldProps = ComponentProps<typeof TextInputField>

const ignoreText = (): void => {}
const ignoreChecked = (): void => {}

export const formFieldFixture: FormFieldProps = {
  id: "storybook-form-field",
  label: "表示ラベル",
  help: "補助文は入力欄の下に表示します。",
  required: true
}

export const textInputFieldFixture: TextInputFieldProps = {
  id: "storybook-text-input",
  label: "検索語",
  value: "サンプル",
  placeholder: "検索語を入力",
  help: "固定 fixture の値だけを表示します。",
  required: true,
  onInput: ignoreText
}

export const textInputErrorFixture: TextInputFieldProps = {
  ...textInputFieldFixture,
  id: "storybook-text-input-error",
  value: "",
  error: "検索語を入力してください。"
}

export const textInputDisabledFixture: TextInputFieldProps = {
  ...textInputFieldFixture,
  id: "storybook-text-input-disabled",
  value: "編集できない値",
  disabled: true
}

export const textAreaFieldFixture: TextAreaFieldProps = {
  id: "storybook-text-area",
  label: "備考",
  value: "長文入力の折り返しを確認するためのサンプル文です。画面固有の保存判断は持ちません。",
  rows: 5,
  help: "複数行の入力値を表示します。",
  onInput: ignoreText
}

export const selectFieldFixture: SelectFieldProps = {
  id: "storybook-select",
  label: "分類",
  value: "sample-a",
  options: [
    { value: "sample-a", label: "分類 A" },
    { value: "sample-b", label: "分類 B" }
  ],
  help: "選択肢は synthetic value だけです。",
  onChange: ignoreText
}

export const checkboxFieldFixture: CheckboxFieldProps = {
  id: "storybook-checkbox",
  label: "確認済みにする",
  checked: true,
  help: "checked と callback だけで表示します。",
  onChange: ignoreChecked
}

export const checkboxFieldErrorFixture: CheckboxFieldProps = {
  ...checkboxFieldFixture,
  id: "storybook-checkbox-error",
  checked: false,
  error: "確認が必要です。"
}
