use crate::domain::job_state::JobState;
use crate::domain::translation_unit::TranslationUnit;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct CreateJobSourceGroup {
    pub source_json_path: String,
    pub target_plugin: String,
    pub translation_units: Vec<TranslationUnit>,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct CreatedJob {
    pub job_id: String,
    pub state: JobState,
    pub source_groups: Vec<CreateJobSourceGroup>,
}

pub fn create_ready_job(
    job_id: &str,
    source_groups: Vec<CreateJobSourceGroup>,
) -> Result<CreatedJob, String> {
    if job_id.trim().is_empty() {
        return Err("job creation requires a job_id.".to_string());
    }

    if source_groups.is_empty() {
        return Err("job creation requires at least one source group.".to_string());
    }

    let translation_unit_count = source_groups
        .iter()
        .flat_map(|group| group.translation_units.iter())
        .count();
    if translation_unit_count == 0 {
        return Err("job creation requires at least one translation unit.".to_string());
    }

    let ready_state =
        JobState::Draft
            .transition_to(JobState::Ready)
            .map_err(|transition_error| {
                format!(
                    "job creation failed to transition state from {:?} to {:?}.",
                    transition_error.from, transition_error.to
                )
            })?;

    Ok(CreatedJob {
        job_id: job_id.to_string(),
        state: ready_state,
        source_groups,
    })
}
