import type { Page } from "@playwright/test"

type BindingController = Record<string, unknown>

function pendingPromiseScript(
  controllerName: string,
  bindingName: string
): string {
  return `
    (() => {
      const controller = globalThis.go?.wails?.[${JSON.stringify(controllerName)}];
      if (!controller || typeof controller[${JSON.stringify(bindingName)}] !== "function") {
        return false;
      }
      controller[${JSON.stringify(bindingName)}] = () => new Promise(() => {});
      return true;
    })()
  `
}

export async function holdWailsBindingPending(
  page: Page,
  controllerName: string,
  bindingName: string
): Promise<void> {
  await page.waitForFunction(
    ({ controllerName: targetController, bindingName: targetBinding }) => {
      const wails = globalThis as typeof globalThis & {
        go?: { wails?: Record<string, BindingController> }
      }
      return (
        typeof wails.go?.wails?.[targetController]?.[targetBinding] ===
        "function"
      )
    },
    { controllerName, bindingName }
  )
  await page.evaluate(pendingPromiseScript(controllerName, bindingName))
}
