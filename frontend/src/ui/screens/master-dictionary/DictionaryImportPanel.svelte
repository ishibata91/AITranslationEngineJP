<script lang="ts">
  import FileImportPanel from "@ui/components/FileImportPanel.svelte"
  import type { DictionaryImportPanelProps } from "./dictionary-panel-props"

  let {
    isImportRunning,
    selectedFileName,
    chooseXmlFile,
    handleXmlSelected
  }: DictionaryImportPanelProps = $props()
</script>

<FileImportPanel
  accept=".xml,text/xml,application/xml"
  eyebrow="XMLから取り込み"
  helperText="ファイルを選ぶと進行状況パネルで取り込みを開始できます。"
  inputId="xmlFileInput"
  primaryActionId="chooseXmlButton"
  primaryActionLabel="ファイルを選択"
  primaryActionDisabled={isImportRunning}
  selectedName={selectedFileName}
  testId="master-dictionary-xml-import-region"
  title="XML 取り込み"
  titleId="importHeading"
  onFileSelected={handleXmlSelected}
  onPrimaryAction={chooseXmlFile}
/>
