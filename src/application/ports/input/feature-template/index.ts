import type {
  FeatureScreenState,
  FeatureScreenUsecase,
} from "@application/ports/input/feature-screen";
import type {
  FeatureTemplateData,
  FeatureTemplateQuery,
} from "@shared/contracts/feature-template";

export type FeatureTemplateScreenState = FeatureScreenState<
  FeatureTemplateData,
  string,
  FeatureTemplateQuery
>;
export type FeatureTemplateScreenInput = FeatureScreenUsecase<
  string,
  FeatureTemplateQuery
>;
