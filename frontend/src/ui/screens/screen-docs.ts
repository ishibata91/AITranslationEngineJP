import {
  Description,
  Primary,
  Stories,
  Subtitle,
  Title
} from "@storybook/addon-docs/blocks"
import { createElement, Fragment } from "react"

// 既定の Autodocs から画面 props の表だけを除く。
export const ScreenDocsPage = () =>
  createElement(
    Fragment,
    null,
    createElement(Title),
    createElement(Subtitle),
    createElement(Description, { of: "meta" }),
    createElement(Primary),
    createElement(Stories)
  )
