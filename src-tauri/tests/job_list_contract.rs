use ai_translation_engine_jp_lib::application::dto::job::JobStateDto;
use ai_translation_engine_jp_lib::application::job::list::{ListJobsRepository, ListJobsUseCase};
use ai_translation_engine_jp_lib::domain::job::list::ListedJob;
use ai_translation_engine_jp_lib::domain::job_state::JobState;
use async_trait::async_trait;

#[tokio::test]
async fn given_saved_observable_jobs_when_listing_jobs_then_returns_minimal_job_id_and_state_view()
{
    let use_case = ListJobsUseCase::new(StaticListJobsRepository {
        jobs: vec![
            ListedJob::new("job-0001", JobState::Ready).expect("ready job should be listable"),
            ListedJob::new("job-0002", JobState::Completed)
                .expect("completed job should be listable"),
        ],
    });

    let result = use_case
        .execute()
        .await
        .expect("saved jobs should list successfully");

    assert_eq!(result.jobs.len(), 2);
    assert_eq!(result.jobs[0].job_id, "job-0001");
    assert_eq!(result.jobs[0].state, JobStateDto::Ready);
    assert_eq!(result.jobs[1].job_id, "job-0002");
    assert_eq!(result.jobs[1].state, JobStateDto::Completed);
}

#[tokio::test]
async fn given_no_saved_jobs_when_listing_jobs_then_returns_empty_result() {
    let use_case = ListJobsUseCase::new(StaticListJobsRepository { jobs: Vec::new() });

    let result = use_case
        .execute()
        .await
        .expect("empty persistence state should still list successfully");

    assert!(result.jobs.is_empty());
}

#[tokio::test]
async fn given_repository_failure_when_listing_jobs_then_execute_returns_the_failure() {
    let use_case = ListJobsUseCase::new(FailingListJobsRepository);

    let error = use_case
        .execute()
        .await
        .expect_err("repository failure should bubble out");

    assert_eq!(error, "job list repository read failed");
}

#[test]
fn given_list_jobs_result_dto_when_serializing_then_wire_shape_stays_minimal() {
    let result = ai_translation_engine_jp_lib::application::dto::job::ListJobsResultDto {
        jobs: vec![
            ai_translation_engine_jp_lib::application::dto::job::JobListItemDto {
                job_id: "job-0001".to_string(),
                state: JobStateDto::Ready,
            },
        ],
    };

    let serialized =
        serde_json::to_string(&result).expect("list result dto should serialize successfully");

    assert_eq!(
        serialized,
        r#"{"jobs":[{"jobId":"job-0001","state":"Ready"}]}"#
    );
}

struct StaticListJobsRepository {
    jobs: Vec<ListedJob>,
}

#[async_trait]
impl ListJobsRepository for StaticListJobsRepository {
    async fn list_jobs(&self) -> Result<Vec<ListedJob>, String> {
        Ok(self.jobs.clone())
    }
}

struct FailingListJobsRepository;

#[async_trait]
impl ListJobsRepository for FailingListJobsRepository {
    async fn list_jobs(&self) -> Result<Vec<ListedJob>, String> {
        Err("job list repository read failed".to_string())
    }
}
