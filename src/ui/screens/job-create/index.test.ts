import { describe, expect, it } from "vitest";
import { render } from "svelte/server";
import JobCreateScreen from "./index";
import { JobCreateView } from "@ui/views/job-create";

type RenderState = {
  error: string | null;
  isSubmitting: boolean;
  request: {
    sourceGroups: Array<{
      sourceJsonPath: string;
      targetPlugin: string;
      translationUnits: Array<{
        editorId: string;
        extractionKey: string;
        fieldName: string;
        formId: string;
        recordSignature: string;
        sortKey: string;
        sourceEntityType: string;
        sourceText: string;
      }>;
    }>;
  };
  result: {
    jobId: string;
    state: "Draft" | "Ready" | "Running" | "Completed";
  } | null;
};

function createRenderState(overrides?: Partial<RenderState>): RenderState {
  return {
    error: null,
    isSubmitting: false,
    request: {
      sourceGroups: [
        {
          sourceJsonPath: "F:/imports/sample.json",
          targetPlugin: "Example.esp",
          translationUnits: [
            {
              editorId: "Sword01",
              extractionKey: "item:0001:name",
              fieldName: "name",
              formId: "0001",
              recordSignature: "WEAP",
              sortKey: "item:0001:name",
              sourceEntityType: "item",
              sourceText: "Iron Sword"
            }
          ]
        }
      ]
    },
    result: null,
    ...overrides
  };
}

async function renderView(state: RenderState): Promise<string> {
  const module = await import("@ui/views/job-create/JobCreateView.svelte?server");
  const ServerView = module.default;
  const { body } = render(ServerView, { props: { state } });

  return body;
}

describe("job create public roots", () => {
  it("Given the screen and view roots When imported Then the job create modules resolve", () => {
    expect(JobCreateScreen).toBeTruthy();
    expect(JobCreateView).toBeTruthy();
  });

  it("Given idle create state When the view is server-rendered Then core create contract labels are present", async () => {
    const body = await renderView(createRenderState());

    expect(body).toContain("Job Create");
    expect(body).toContain("Create job");
    expect(body).toContain("Source Group");
    expect(body).toContain("Translation Unit");
  });

  it("Given success and failure state When the view is server-rendered Then result and error blocks are rendered", async () => {
    const body = await renderView(
      createRenderState({
        error: "Job creation failed. Try again.",
        result: {
          jobId: "job-0001",
          state: "Ready"
        }
      })
    );

    expect(body).toContain("Create failed");
    expect(body).toContain("Job creation failed. Try again.");
    expect(body).toContain("Created Job");
    expect(body).toContain("job-0001");
    expect(body).toContain("Observable state: Ready");
  });
});
