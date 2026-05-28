import { expect, test } from "@playwright/test";

import { OutputManagementPage } from "./support/system-test-pages";
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks";

test.beforeEach(async ({ page }) => {
  await installScenarioWailsMocks(page);
});

test("E2E-UC-023 output management generates XML artifact", async ({
  page,
}) => {
  // 出力候補を選び出力先を入力する利用者操作で、XML 出力結果が表示されることを証明する。
  const outputManagement = new OutputManagementPage(page);
  await outputManagement.open();

  await outputManagement.expectCandidateListConnected("job #401");
  await outputManagement.selectCandidate("job #401");
  await outputManagement.selectTargetGame("Skyrim SE");
  await outputManagement.fillOutputPath("/tmp/system-test-output.xml");
  await outputManagement.exportXml();

  await expect(outputManagement.latestResult).toContainText(
    /current|出力|generate/,
  );
});

test("E2E-UC-024 output management regenerates XML artifact", async ({
  page,
}) => {
  // 既存 artifact を持つ job の再出力操作で、新しい結果が表示されることを証明する。
  const outputManagement = new OutputManagementPage(page);
  await outputManagement.open();

  await outputManagement.expectCandidateListConnected("job #402");
  await outputManagement.selectCandidate("job #402");
  await outputManagement.selectTargetGame("Skyrim SE");
  await outputManagement.fillOutputPath("/tmp/system-test-reoutput.xml");
  await expect(outputManagement.reexportButton).toBeEnabled();
  await outputManagement.reexportXml();

  await expect(outputManagement.latestResult).toContainText(
    /再出力|regenerate|current/,
  );
});

test("E2E-UC-025 output management shows diff preview rows", async ({
  page,
}) => {
  // 出力済み job の選択操作で、差分 preview の行を確認できることを証明する。
  const outputManagement = new OutputManagementPage(page);
  await outputManagement.open();

  await outputManagement.expectCandidateListConnected("job #402");
  await outputManagement.selectCandidate("job #402");

  await expect(outputManagement.diffPreview).toContainText("LydiaLine");
  await expect(outputManagement.diffRows).toHaveCount(1);
  await outputManagement.selectDiffRow("LydiaLine");
  await expect(outputManagement.selectedJob).toContainText("job id");
  await expect(outputManagement.selectedJob).toContainText("402");
});

test("E2E-UC-042 output management keeps export disabled for invalid path", async ({
  page,
}) => {
  // 不正な出力先 path では、XML 出力が成功結果へ進まないことを証明する。
  const outputManagement = new OutputManagementPage(page);
  await outputManagement.open();

  await outputManagement.expectCandidateListConnected("job #401");
  await outputManagement.selectCandidate("job #401");
  await outputManagement.fillOutputPath("invalid-path");

  await expect(outputManagement.exportButton).toBeDisabled();
  await expect(outputManagement.outputActions).toContainText(/path|xml|出力/);
  await expect(outputManagement.latestResult).toHaveCount(0);
});

test("E2E-UC-043 output management keeps reexport disabled when current artifact is already latest", async ({
  page,
}) => {
  // 再出力不要の job では、再出力操作が無効で新しい結果が表示されないことを証明する。
  const outputManagement = new OutputManagementPage(page);
  await outputManagement.open();

  await outputManagement.expectCandidateListConnected("job #403");
  await outputManagement.selectCandidate("job #403");

  await expect(outputManagement.selectedJob).toContainText(
    /0|差分なし|not ready/,
  );
  await expect(outputManagement.reexportButton).toBeDisabled();
  await expect(outputManagement.latestResult).toHaveCount(0);
});

test("E2E-UC-044 output management keeps empty diff preview when there is no diff", async ({
  page,
}) => {
  // 差分なしの job では、diff preview が空状態を維持することを証明する。
  const outputManagement = new OutputManagementPage(page);
  await outputManagement.open();

  await outputManagement.expectCandidateListConnected("job #403");
  await outputManagement.selectCandidate("job #403");

  await expect(outputManagement.diffRows).toHaveCount(0);
  await expect(outputManagement.diffPreview).toContainText(
    /未取得|row count 0/,
  );
});
