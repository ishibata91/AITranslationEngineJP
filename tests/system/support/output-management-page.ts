import { expect, type Locator, type Page } from "@playwright/test";

import { SystemTestPageObject } from "./system-test-page-object";

export class OutputManagementPage extends SystemTestPageObject {
  constructor(page: Page) {
    super(page);
  }

  get summary(): Locator {
    return this.byTestId("output-management-output-management-summary");
  }

  get candidateList(): Locator {
    return this.byTestId("output-management-output-candidate-list");
  }

  get candidateRows(): Locator {
    return this.byTestId("output-management-output-candidate-row");
  }

  get selectedJob(): Locator {
    return this.byTestId("output-management-selected-job");
  }

  get outputActions(): Locator {
    return this.byTestId("output-management-output-actions");
  }

  get targetGameSelect(): Locator {
    return this.byTestId("output-management-target-game-select");
  }

  get outputPathInput(): Locator {
    return this.byTestId("output-management-output-path-input");
  }

  get exportButton(): Locator {
    return this.byTestId("output-management-export-button");
  }

  get reexportButton(): Locator {
    return this.byTestId("output-management-reexport-button");
  }

  get outputPathError(): Locator {
    return this.byTestId("output-management-output-path-error");
  }

  get latestResult(): Locator {
    return this.byTestId("output-management-latest-result");
  }

  get diffPreview(): Locator {
    return this.byTestId("output-management-diff-preview");
  }

  get diffRows(): Locator {
    return this.byTestId("output-management-diff-row");
  }

  async open(): Promise<void> {
    await this.openHashRoute("/#output-management");
    await this.waitFor(this.summary);
  }

  candidateRow(text: string): Locator {
    return this.candidateRows.filter({ hasText: text });
  }

  diffRow(text: string): Locator {
    return this.diffRows.filter({ hasText: text });
  }

  async expectCandidateListConnected(candidateText: string): Promise<void> {
    await expect(this.candidateList).toBeVisible();
    await expect(this.candidateRow(candidateText)).toBeVisible();
  }

  async selectCandidate(text: string): Promise<void> {
    const candidate = this.candidateRow(text);

    await this.waitFor(candidate);
    await candidate.click();
  }

  async selectTargetGame(label: string): Promise<void> {
    await this.targetGameSelect.selectOption({ label });
  }

  async fillOutputPath(outputPath: string): Promise<void> {
    await this.outputPathInput.fill(outputPath);
  }

  async exportXml(): Promise<void> {
    await this.exportButton.click();
  }

  async reexportXml(): Promise<void> {
    await this.reexportButton.click();
  }

  async selectDiffRow(text: string): Promise<void> {
    const diffRow = this.diffRow(text);

    await this.waitFor(diffRow);
    await diffRow.click();
  }
}
