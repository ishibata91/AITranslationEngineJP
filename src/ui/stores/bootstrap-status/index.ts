import type { BootstrapStatusField } from "@application/ports/input/bootstrap-status";
import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";
import { createFeatureScreenStore, type FeatureScreenStore } from "@ui/stores/feature-screen";

export type BootstrapStatusScreenStore = FeatureScreenStore<
  BootstrapStatus,
  BootstrapStatusField,
  undefined
>;

export function createBootstrapStatusScreenStore(): BootstrapStatusScreenStore {
  return createFeatureScreenStore({
    filters: undefined
  });
}
