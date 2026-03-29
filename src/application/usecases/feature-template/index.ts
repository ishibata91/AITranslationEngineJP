import type { FeatureTemplateGateway } from "@application/ports/gateway/feature-template";
import type { FeatureTemplateScreenInput } from "@application/ports/input/feature-template";
import type { FeatureScreenStorePort } from "@application/ports/input/feature-screen";
import { createFeatureScreenUsecase } from "@application/usecases/feature-screen";
import type { FeatureTemplateData, FeatureTemplateQuery } from "@shared/contracts/feature-template";

type FeatureTemplateStorePort = FeatureScreenStorePort<
  FeatureTemplateData,
  string,
  FeatureTemplateQuery
>;

type CreateFeatureTemplateScreenUsecaseOptions = {
  gateway: FeatureTemplateGateway;
  store: FeatureTemplateStorePort;
};

export function createFeatureTemplateScreenUsecase({
  gateway,
  store
}: CreateFeatureTemplateScreenUsecaseOptions): FeatureTemplateScreenInput {
  return createFeatureScreenUsecase({
    createRequest(state) {
      return state.filters;
    },
    gateway,
    reconcileSelection: ({ currentSelection, data }) => {
      if (currentSelection === null) {
        return null;
      }

      return data.items.some((item) => item.id === currentSelection) ? currentSelection : null;
    },
    store
  });
}
