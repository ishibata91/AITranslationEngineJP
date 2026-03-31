use serde::{Deserialize, Deserializer, Serialize};

use crate::application::dto::translation_unit::TranslationUnitDto;
use crate::domain::job::create::CreatedJob;
use crate::domain::job::list::ListedJob;
use crate::domain::job_state::JobState;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize)]
pub enum JobStateDto {
    Draft,
    Ready,
    Running,
    Completed,
}

impl From<JobState> for JobStateDto {
    fn from(value: JobState) -> Self {
        match value {
            JobState::Draft => Self::Draft,
            JobState::Ready => Self::Ready,
            JobState::Running => Self::Running,
            JobState::Completed => Self::Completed,
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct CreateJobRequestDto {
    pub source_groups: Vec<CreateJobSourceGroupDto>,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct CreateJobSourceGroupDto {
    pub source_json_path: String,
    pub target_plugin: String,
    pub translation_units: Vec<TranslationUnitDto>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct CreateJobRequestWireDto {
    source_groups: Vec<CreateJobSourceGroupWireDto>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct CreateJobSourceGroupWireDto {
    source_json_path: String,
    target_plugin: String,
    translation_units: Vec<TranslationUnitWireDto>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct TranslationUnitWireDto {
    source_entity_type: String,
    form_id: String,
    editor_id: String,
    record_signature: String,
    field_name: String,
    extraction_key: String,
    source_text: String,
    sort_key: String,
}

impl<'de> Deserialize<'de> for CreateJobRequestDto {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let wire = CreateJobRequestWireDto::deserialize(deserializer)?;
        Ok(Self {
            source_groups: wire
                .source_groups
                .into_iter()
                .map(|group| CreateJobSourceGroupDto {
                    source_json_path: group.source_json_path,
                    target_plugin: group.target_plugin,
                    translation_units: group
                        .translation_units
                        .into_iter()
                        .map(|unit| TranslationUnitDto {
                            source_entity_type: unit.source_entity_type,
                            form_id: unit.form_id,
                            editor_id: unit.editor_id,
                            record_signature: unit.record_signature,
                            field_name: unit.field_name,
                            extraction_key: unit.extraction_key,
                            source_text: unit.source_text,
                            sort_key: unit.sort_key,
                        })
                        .collect(),
                })
                .collect(),
        })
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct CreateJobResultDto {
    pub job_id: String,
    pub state: JobStateDto,
}

impl CreateJobResultDto {
    pub fn from_created_job(created_job: &CreatedJob) -> Self {
        Self {
            job_id: created_job.job_id.clone(),
            state: JobStateDto::from(created_job.state),
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct JobListItemDto {
    pub job_id: String,
    pub state: JobStateDto,
}

impl JobListItemDto {
    pub fn from_listed_job(listed_job: &ListedJob) -> Self {
        Self {
            job_id: listed_job.job_id.clone(),
            state: JobStateDto::from(listed_job.state),
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ListJobsResultDto {
    pub jobs: Vec<JobListItemDto>,
}

impl ListJobsResultDto {
    pub fn from_listed_jobs(listed_jobs: &[ListedJob]) -> Self {
        Self {
            jobs: listed_jobs
                .iter()
                .map(JobListItemDto::from_listed_job)
                .collect(),
        }
    }
}
