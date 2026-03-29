import type { FeatureScreenGateway } from "@application/ports/gateway/feature-screen";
import type { FeatureTemplateData, FeatureTemplateQuery } from "@shared/contracts/feature-template";

export type FeatureTemplateGateway = FeatureScreenGateway<FeatureTemplateQuery, FeatureTemplateData>;

