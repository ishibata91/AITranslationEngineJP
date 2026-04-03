import type {
  FeatureTemplateData,
  FeatureTemplateQuery,
} from "@shared/contracts/feature-template";
import {
  createFeatureScreenStore,
  type FeatureScreenStore,
} from "@ui/stores/feature-screen";

export type FeatureTemplateScreenStore = FeatureScreenStore<
  FeatureTemplateData,
  string,
  FeatureTemplateQuery
>;

export function createFeatureTemplateScreenStore(
  filters: FeatureTemplateQuery,
): FeatureTemplateScreenStore {
  return createFeatureScreenStore({
    filters,
  });
}
