import { describe, expect, it } from "vitest";
import { createFeatureTemplateScreenStore } from "./index";

describe("createFeatureTemplateScreenStore", () => {
  it("Given no data When store is created Then loading error data selection and filters start from the template defaults", () => {
    const store = createFeatureTemplateScreenStore({
      query: "dictionary",
    });

    expect(store.getState()).toEqual({
      data: null,
      error: null,
      filters: {
        query: "dictionary",
      },
      loading: false,
      selection: null,
    });
  });
});
