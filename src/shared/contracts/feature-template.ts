export type FeatureTemplateItem = {
  detail: string;
  id: string;
  status: "queued" | "running" | "ready" | "failed";
  title: string;
};

export type FeatureTemplateQuery = {
  query: string;
};

export type FeatureTemplateData = {
  items: FeatureTemplateItem[];
};

