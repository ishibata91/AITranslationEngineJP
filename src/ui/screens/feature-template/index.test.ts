import { describe, expect, it } from "vitest";
import FeatureTemplateScreen from "./index";
import { FeatureTemplateView } from "@ui/views/feature-template";

describe("feature template public roots", () => {
  it("Given the screen and view roots When imported Then the template modules resolve", () => {
    expect(FeatureTemplateScreen).toBeTruthy();
    expect(FeatureTemplateView).toBeTruthy();
  });
});
