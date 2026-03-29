use std::collections::BTreeMap;
use std::fs;
use std::path::Path;

use serde::{Deserialize, Serialize};

use crate::domain::xedit_export::{ImportedPluginExport, ImportedRawRecord, TranslationUnit};

pub trait XeditExportImporter {
    fn import_from_paths(&self, file_paths: &[String])
        -> Result<Vec<ImportedPluginExport>, String>;
}

pub struct FileSystemXeditExportImporter;

impl XeditExportImporter for FileSystemXeditExportImporter {
    fn import_from_paths(
        &self,
        file_paths: &[String],
    ) -> Result<Vec<ImportedPluginExport>, String> {
        if file_paths.is_empty() {
            return Err("xEdit export import requires at least one file path.".to_string());
        }

        file_paths
            .iter()
            .map(|file_path| import_single_file(Path::new(file_path)))
            .collect()
    }
}

fn import_single_file(path: &Path) -> Result<ImportedPluginExport, String> {
    let json = fs::read_to_string(path).map_err(|error| {
        format!(
            "Failed to read xEdit export JSON `{}`: {error}",
            path.display()
        )
    })?;
    let raw_export: RawPluginExport = serde_json::from_str(&json).map_err(|error| {
        format!(
            "Failed to parse xEdit export JSON `{}`: {error}. Check top-level collection fields such as `dialogue_groups`, `quests`, `items`, `magic`, `locations`, `system`, `messages`, `load_screens`, and `npcs`.",
            path.display(),
        )
    })?;

    let translation_units = raw_export
        .translation_units()
        .map_err(|error| format!("Invalid xEdit export data `{}`: {error}", path.display()))?;
    let raw_records = raw_export
        .raw_records()
        .map_err(|error| format!("Invalid xEdit export data `{}`: {error}", path.display()))?;

    ImportedPluginExport::new(
        path.display().to_string(),
        raw_export.target_plugin,
        translation_units,
        raw_records,
    )
}

#[derive(Debug, Deserialize)]
struct RawPluginExport {
    target_plugin: String,
    #[serde(default)]
    dialogue_groups: Vec<RawDialogueGroup>,
    #[serde(default)]
    quests: Vec<RawQuest>,
    #[serde(default)]
    items: Vec<RawItem>,
    #[serde(default)]
    magic: Vec<RawMagic>,
    #[serde(default)]
    locations: Vec<RawLocation>,
    #[serde(default)]
    system: Vec<RawSystemRecord>,
    #[serde(default)]
    messages: Vec<RawMessage>,
    #[serde(default)]
    load_screens: Vec<RawLoadScreen>,
    #[serde(default)]
    npcs: BTreeMap<String, RawNpc>,
}

impl RawPluginExport {
    fn translation_units(&self) -> Result<Vec<TranslationUnit>, String> {
        let mut translation_units = Vec::new();

        for dialogue_group in &self.dialogue_groups {
            dialogue_group.append_translation_units(&mut translation_units)?;
        }

        for quest in &self.quests {
            quest.append_translation_units(&mut translation_units)?;
        }

        for item in &self.items {
            item.append_translation_units(&mut translation_units)?;
        }

        for magic in &self.magic {
            magic.append_translation_units(&mut translation_units)?;
        }

        for location in &self.locations {
            location.append_translation_units(&mut translation_units)?;
        }

        for system_record in &self.system {
            system_record.append_translation_units(&mut translation_units)?;
        }

        for message in &self.messages {
            message.append_translation_units(&mut translation_units)?;
        }

        for load_screen in &self.load_screens {
            load_screen.append_translation_units(&mut translation_units)?;
        }

        for npc in self.npcs.values() {
            npc.append_translation_units(&mut translation_units)?;
        }

        Ok(translation_units)
    }

    fn raw_records(&self) -> Result<Vec<ImportedRawRecord>, String> {
        let mut raw_records = Vec::new();

        for dialogue_group in &self.dialogue_groups {
            raw_records.push(raw_record("dialogue_group", dialogue_group)?);
            for response in &dialogue_group.responses {
                raw_records.push(raw_record("dialogue_response", response)?);
            }
        }

        for quest in &self.quests {
            raw_records.push(raw_record("quest", quest)?);
            for objective in &quest.objectives {
                raw_records.push(ImportedRawRecord::new(
                    "quest_objective",
                    &quest.id,
                    &quest.editor_id,
                    record_signature(&quest.record_type),
                    &serde_json::to_string(objective)
                        .map_err(|error| format!("Failed to serialize raw record: {error}"))?,
                )?);
            }
            for stage in &quest.stages {
                raw_records.push(ImportedRawRecord::new(
                    "quest_stage_log",
                    &quest.id,
                    &quest.editor_id,
                    record_signature(&quest.record_type),
                    &serde_json::to_string(stage)
                        .map_err(|error| format!("Failed to serialize raw record: {error}"))?,
                )?);
            }
        }

        for item in &self.items {
            raw_records.push(raw_record("item", item)?);
        }

        for magic in &self.magic {
            raw_records.push(raw_record("magic", magic)?);
        }

        for location in &self.locations {
            raw_records.push(raw_record("location", location)?);
        }

        for system_record in &self.system {
            raw_records.push(raw_record("system_record", system_record)?);
        }

        for message in &self.messages {
            raw_records.push(raw_record("message", message)?);
        }

        for load_screen in &self.load_screens {
            raw_records.push(raw_record("load_screen", load_screen)?);
        }

        for npc in self.npcs.values() {
            raw_records.push(raw_record("npc", npc)?);
        }

        Ok(raw_records)
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawDialogueGroup {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    player_text: String,
    #[serde(default)]
    responses: Vec<RawDialogueResponse>,
}

impl RawDialogueGroup {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        let unit_base = TranslationUnitBase {
            source_entity_type: "dialogue_group",
            form_id: &self.id,
            editor_id: &self.editor_id,
            record_type: &self.record_type,
        };

        push_translation_unit(
            translation_units,
            &unit_base,
            "player_text",
            &self.player_text,
            &format!("dialogue_group:{}:player_text", self.id),
        )?;

        for response in &self.responses {
            response.append_translation_units(translation_units)?;
        }

        Ok(())
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawDialogueResponse {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    text: String,
    #[serde(default)]
    prompt: String,
    #[serde(default)]
    topic_text: String,
    #[serde(default)]
    menu_display_text: String,
}

impl RawDialogueResponse {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        let unit_base = TranslationUnitBase {
            source_entity_type: "dialogue_response",
            form_id: &self.id,
            editor_id: &self.editor_id,
            record_type: &self.record_type,
        };

        push_translation_unit(
            translation_units,
            &unit_base,
            "text",
            &self.text,
            &format!("dialogue_response:{}:text", self.id),
        )?;
        push_translation_unit(
            translation_units,
            &unit_base,
            "prompt",
            &self.prompt,
            &format!("dialogue_response:{}:prompt", self.id),
        )?;
        push_translation_unit(
            translation_units,
            &unit_base,
            "topic_text",
            &self.topic_text,
            &format!("dialogue_response:{}:topic_text", self.id),
        )?;
        push_translation_unit(
            translation_units,
            &unit_base,
            "menu_display_text",
            &self.menu_display_text,
            &format!("dialogue_response:{}:menu_display_text", self.id),
        )
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawQuest {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    name: String,
    #[serde(default)]
    objectives: Vec<RawQuestObjective>,
    #[serde(default)]
    stages: Vec<RawQuestStageLog>,
}

impl RawQuest {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        let unit_base = TranslationUnitBase {
            source_entity_type: "quest",
            form_id: &self.id,
            editor_id: &self.editor_id,
            record_type: &self.record_type,
        };

        push_translation_unit(
            translation_units,
            &unit_base,
            "name",
            &self.name,
            &format!("quest:{}:name", self.id),
        )?;

        for objective in &self.objectives {
            objective.append_translation_units(
                translation_units,
                &self.id,
                &self.editor_id,
                &self.record_type,
            )?;
        }

        for stage in &self.stages {
            stage.append_translation_units(
                translation_units,
                &self.id,
                &self.editor_id,
                &self.record_type,
            )?;
        }

        Ok(())
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawQuestObjective {
    objective_index: String,
    #[serde(default)]
    text: String,
}

impl RawQuestObjective {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
        quest_id: &str,
        quest_editor_id: &str,
        quest_record_type: &str,
    ) -> Result<(), String> {
        if !self.text.trim().is_empty() && self.objective_index.trim().is_empty() {
            return Err(
                "xEdit export import requires `quests[].objectives[].objective_index` when objective text is present.".to_string(),
            );
        }

        let unit_base = TranslationUnitBase {
            source_entity_type: "quest_objective",
            form_id: quest_id,
            editor_id: quest_editor_id,
            record_type: quest_record_type,
        };

        push_translation_unit(
            translation_units,
            &unit_base,
            "text",
            &self.text,
            &format!("quest:{}:objective:{}:text", quest_id, self.objective_index),
        )
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawQuestStageLog {
    stage_index: i64,
    log_index: i64,
    #[serde(default)]
    text: String,
}

impl RawQuestStageLog {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
        quest_id: &str,
        quest_editor_id: &str,
        quest_record_type: &str,
    ) -> Result<(), String> {
        let unit_base = TranslationUnitBase {
            source_entity_type: "quest_stage_log",
            form_id: quest_id,
            editor_id: quest_editor_id,
            record_type: quest_record_type,
        };

        push_translation_unit(
            translation_units,
            &unit_base,
            "text",
            &self.text,
            &format!(
                "quest:{}:stage:{}:log:{}:text",
                quest_id, self.stage_index, self.log_index
            ),
        )
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawItem {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    name: String,
    #[serde(default)]
    description: String,
    #[serde(default)]
    text: String,
}

impl RawItem {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        append_named_record_translation_units(
            translation_units,
            "item",
            &self.id,
            &self.editor_id,
            &self.record_type,
            &[
                ("name", &self.name),
                ("description", &self.description),
                ("text", &self.text),
            ],
        )
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawMagic {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    name: String,
    #[serde(default)]
    description: String,
}

impl RawMagic {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        append_named_record_translation_units(
            translation_units,
            "magic",
            &self.id,
            &self.editor_id,
            &self.record_type,
            &[("name", &self.name), ("description", &self.description)],
        )
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawLocation {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    name: String,
}

impl RawLocation {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        append_named_record_translation_units(
            translation_units,
            "location",
            &self.id,
            &self.editor_id,
            &self.record_type,
            &[("name", &self.name)],
        )
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawSystemRecord {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    name: String,
    #[serde(default)]
    description: String,
}

impl RawSystemRecord {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        append_named_record_translation_units(
            translation_units,
            "system_record",
            &self.id,
            &self.editor_id,
            &self.record_type,
            &[("name", &self.name), ("description", &self.description)],
        )
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawMessage {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    text: String,
    #[serde(default)]
    title: String,
}

impl RawMessage {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        append_named_record_translation_units(
            translation_units,
            "message",
            &self.id,
            &self.editor_id,
            &self.record_type,
            &[("text", &self.text), ("title", &self.title)],
        )
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawLoadScreen {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    text: String,
}

impl RawLoadScreen {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        append_named_record_translation_units(
            translation_units,
            "load_screen",
            &self.id,
            &self.editor_id,
            &self.record_type,
            &[("text", &self.text)],
        )
    }
}

#[derive(Debug, Deserialize, Serialize)]
struct RawNpc {
    id: String,
    editor_id: String,
    #[serde(rename = "type")]
    record_type: String,
    #[serde(default)]
    name: String,
}

impl RawNpc {
    fn append_translation_units(
        &self,
        translation_units: &mut Vec<TranslationUnit>,
    ) -> Result<(), String> {
        append_named_record_translation_units(
            translation_units,
            "npc",
            &self.id,
            &self.editor_id,
            &self.record_type,
            &[("name", &self.name)],
        )
    }
}

fn append_named_record_translation_units(
    translation_units: &mut Vec<TranslationUnit>,
    source_entity_type: &str,
    form_id: &str,
    editor_id: &str,
    record_type: &str,
    fields: &[(&str, &str)],
) -> Result<(), String> {
    let unit_base = TranslationUnitBase {
        source_entity_type,
        form_id,
        editor_id,
        record_type,
    };

    for (field_name, source_text) in fields {
        push_translation_unit(
            translation_units,
            &unit_base,
            field_name,
            source_text,
            &format!("{source_entity_type}:{form_id}:{field_name}"),
        )?;
    }

    Ok(())
}

struct TranslationUnitBase<'a> {
    source_entity_type: &'a str,
    form_id: &'a str,
    editor_id: &'a str,
    record_type: &'a str,
}

fn push_translation_unit(
    translation_units: &mut Vec<TranslationUnit>,
    unit_base: &TranslationUnitBase<'_>,
    field_name: &str,
    source_text: &str,
    extraction_key: &str,
) -> Result<(), String> {
    if source_text.trim().is_empty() {
        return Ok(());
    }

    translation_units.push(TranslationUnit::new(
        unit_base.source_entity_type,
        unit_base.form_id,
        unit_base.editor_id,
        record_signature(unit_base.record_type),
        field_name,
        extraction_key,
        source_text,
    )?);

    Ok(())
}

fn record_signature(record_type: &str) -> &str {
    record_type.split_whitespace().next().unwrap_or(record_type)
}

fn raw_record<T>(source_entity_type: &str, raw: &T) -> Result<ImportedRawRecord, String>
where
    T: RawRecordIdentity + Serialize,
{
    let raw_payload = serde_json::to_string(raw)
        .map_err(|error| format!("Failed to serialize raw record: {error}"))?;

    ImportedRawRecord::new(
        source_entity_type,
        raw.form_id(),
        raw.editor_id(),
        record_signature(raw.record_type()),
        &raw_payload,
    )
}

trait RawRecordIdentity {
    fn form_id(&self) -> &str;
    fn editor_id(&self) -> &str;
    fn record_type(&self) -> &str;
}

impl RawRecordIdentity for RawDialogueGroup {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

impl RawRecordIdentity for RawDialogueResponse {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

impl RawRecordIdentity for RawQuest {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

impl RawRecordIdentity for RawItem {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

impl RawRecordIdentity for RawMagic {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

impl RawRecordIdentity for RawLocation {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

impl RawRecordIdentity for RawSystemRecord {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

impl RawRecordIdentity for RawMessage {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

impl RawRecordIdentity for RawLoadScreen {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

impl RawRecordIdentity for RawNpc {
    fn form_id(&self) -> &str {
        &self.id
    }
    fn editor_id(&self) -> &str {
        &self.editor_id
    }
    fn record_type(&self) -> &str {
        &self.record_type
    }
}

#[cfg(test)]
mod tests {
    use std::fs;
    use std::path::{Path, PathBuf};
    use std::time::{SystemTime, UNIX_EPOCH};

    use super::{FileSystemXeditExportImporter, XeditExportImporter};

    #[test]
    fn given_valid_xedit_export_json_when_importing_then_preserves_target_plugin_and_translation_units(
    ) {
        let fixture = TempJsonFixture::new(
            "valid-xedit-export",
            r#"{
  "target_plugin": "Sample.esp",
  "dialogue_groups": [
    {
      "id": "000AAA01",
      "editor_id": "SampleDialogue",
      "type": "DIAL FULL",
      "player_text": "Hello there",
      "responses": [
        {
          "id": "000AAA02",
          "editor_id": "SampleResponse",
          "type": "INFO NAM1",
          "text": "General Kenobi"
        }
      ]
    }
  ],
  "quests": [
    {
      "id": "000BBB01",
      "editor_id": "SampleQuest",
      "type": "QUST FULL",
      "name": "Quest Name",
      "objectives": [
        {
          "objective_index": "10",
          "text": "Reach the marker"
        }
      ],
      "stages": [
        {
          "stage_index": 20,
          "log_index": 0,
          "text": "Quest updated"
        }
      ]
    }
  ],
  "items": [
    {
      "id": "000CCC01",
      "editor_id": "SampleBook",
      "type": "BOOK FULL",
      "name": "Book Title",
      "text": "Book body"
    }
  ],
  "messages": [
    {
      "id": "000DDD01",
      "editor_id": "SampleMessage",
      "type": "MESG DESC",
      "text": "Message body",
      "title": "Message title"
    }
  ],
  "load_screens": [
    {
      "id": "000EEE01",
      "editor_id": "SampleLoadScreen",
      "type": "LSCR DESC",
      "text": "Load screen body"
    }
  ],
  "npcs": {
    "000FFF01": {
      "id": "000FFF01",
      "editor_id": "SampleNpc",
      "type": "NPC_ FULL",
      "name": "Npc Name"
    }
  }
}"#,
        );
        let importer = FileSystemXeditExportImporter;

        let plugin_exports = importer
            .import_from_paths(&[fixture.path().display().to_string()])
            .unwrap();

        assert_eq!(plugin_exports.len(), 1);
        assert_eq!(plugin_exports[0].target_plugin, "Sample.esp");
        assert_eq!(plugin_exports[0].translation_units.len(), 11);
        assert_eq!(
            plugin_exports[0].translation_units[0].extraction_key,
            "dialogue_group:000AAA01:player_text"
        );
        assert_eq!(
            plugin_exports[0].translation_units[0].record_signature,
            "DIAL"
        );
        assert!(plugin_exports[0]
            .translation_units
            .iter()
            .any(|unit| unit.source_entity_type == "npc" && unit.field_name == "name"));
    }

    #[test]
    fn given_invalid_xedit_export_json_when_importing_then_returns_error() {
        let fixture = TempJsonFixture::new(
            "invalid-xedit-export",
            r#"{
  "dialogue_groups": []
}"#,
        );
        let importer = FileSystemXeditExportImporter;

        let error = importer
            .import_from_paths(&[fixture.path().display().to_string()])
            .unwrap_err();

        assert!(error.contains("target_plugin"));
    }

    #[test]
    fn given_structurally_invalid_dialogue_groups_when_importing_then_returns_parse_error() {
        let fixture = TempJsonFixture::new(
            "invalid-xedit-export-structure",
            r#"{
  "target_plugin": "Sample.esp",
  "dialogue_groups": {
    "id": "000AAA01"
  }
}"#,
        );
        let importer = FileSystemXeditExportImporter;

        let error = importer
            .import_from_paths(&[fixture.path().display().to_string()])
            .unwrap_err();

        assert!(error.contains("Failed to parse xEdit export JSON"));
        assert!(error.contains(&fixture.path().display().to_string()));
        assert!(error.contains("dialogue_groups"));
    }

    #[test]
    fn given_blank_quest_objective_index_with_text_when_importing_then_returns_error() {
        let fixture = TempJsonFixture::new(
            "invalid-xedit-export-objective-index",
            r#"{
  "target_plugin": "Sample.esp",
  "quests": [
    {
      "id": "000BBB01",
      "editor_id": "SampleQuest",
      "type": "QUST FULL",
      "objectives": [
        {
          "objective_index": "",
          "text": "Reach the marker"
        }
      ]
    }
  ]
}"#,
        );
        let importer = FileSystemXeditExportImporter;

        let error = importer
            .import_from_paths(&[fixture.path().display().to_string()])
            .unwrap_err();

        assert!(error.contains(&fixture.path().display().to_string()));
        assert!(error.contains("quests[].objectives[].objective_index"));
    }

    struct TempJsonFixture {
        path: PathBuf,
    }

    impl TempJsonFixture {
        fn new(name_prefix: &str, contents: &str) -> Self {
            let timestamp = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .expect("current time should be after unix epoch")
                .as_nanos();
            let path = std::env::temp_dir().join(format!("{name_prefix}-{timestamp}.json"));

            fs::write(&path, contents).expect("fixture JSON should be writable");

            Self { path }
        }

        fn path(&self) -> &Path {
            &self.path
        }
    }

    impl Drop for TempJsonFixture {
        fn drop(&mut self) {
            let _ = fs::remove_file(&self.path);
        }
    }
}
