import { invoke } from "@tauri-apps/api/core";
import type { FeatureScreenGateway } from "@application/ports/gateway/feature-screen";

export function createTauriFeatureScreenGateway<TRequest extends Record<string, unknown> | undefined, TData>(
  commandName: string
): FeatureScreenGateway<TRequest, TData> {
  return {
    async load(request: TRequest): Promise<TData> {
      if (request === undefined) {
        return invoke<TData>(commandName);
      }

      return invoke<TData>(commandName, request);
    }
  };
}
