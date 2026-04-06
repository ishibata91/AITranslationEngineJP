use super::*;

use std::sync::{
    atomic::{AtomicUsize, Ordering},
    Arc, Mutex,
};

use async_trait::async_trait;
use serde_json::json;

use crate::application::dto::{
    ExecutionControlFailureCategoryDto, ProviderExecutionModeDto, ProviderRuntimeSettingsDto,
    ProviderSelectionDto,
};

struct SpyTransport {
    create_calls: Arc<AtomicUsize>,
    add_calls: Arc<AtomicUsize>,
    status_calls: Arc<AtomicUsize>,
    results_calls: Arc<AtomicUsize>,
    create_response: Result<String, XaiTransportError>,
    add_response: Result<(), XaiTransportError>,
    status_responses: Arc<Mutex<Vec<Result<String, XaiTransportError>>>>,
    results_responses: Arc<Mutex<Vec<Result<String, XaiTransportError>>>>,
    captured_pagination_tokens: Arc<Mutex<Vec<Option<String>>>>,
}

#[async_trait]
impl XaiBatchTransport for SpyTransport {
    async fn create_batch(
        &self,
        _endpoint: &str,
        _request: &XaiCreateBatchRequest,
    ) -> Result<String, XaiTransportError> {
        self.create_calls.fetch_add(1, Ordering::SeqCst);
        self.create_response.clone()
    }

    async fn add_requests(
        &self,
        _endpoint: &str,
        _request: &XaiAddRequestsRequest,
    ) -> Result<(), XaiTransportError> {
        self.add_calls.fetch_add(1, Ordering::SeqCst);
        self.add_response
    }

    async fn get_batch_status(&self, _endpoint: &str) -> Result<String, XaiTransportError> {
        self.status_calls.fetch_add(1, Ordering::SeqCst);
        let mut responses = self
            .status_responses
            .lock()
            .expect("status responses mutex");
        if responses.is_empty() {
            return Err(XaiTransportError::Response);
        }
        responses.remove(0)
    }

    async fn get_batch_results(
        &self,
        _endpoint: &str,
        _page_size: u32,
        pagination_token: Option<&str>,
    ) -> Result<String, XaiTransportError> {
        self.results_calls.fetch_add(1, Ordering::SeqCst);
        self.captured_pagination_tokens
            .lock()
            .expect("pagination tokens mutex")
            .push(pagination_token.map(ToString::to_string));

        let mut responses = self
            .results_responses
            .lock()
            .expect("results responses mutex");
        if responses.is_empty() {
            return Err(XaiTransportError::Response);
        }

        responses.remove(0)
    }
}

fn selection(provider_id: &str, execution_mode: ProviderExecutionModeDto) -> ProviderSelectionDto {
    ProviderSelectionDto {
        provider_id: provider_id.to_string(),
        execution_mode,
        runtime_settings: ProviderRuntimeSettingsDto {
            retry_limit: 1,
            max_concurrency: 1,
            pause_supported: false,
        },
    }
}

fn config_with_api_key() -> XaiConfig {
    XaiConfig {
        base_url: "https://api.x.ai".to_string(),
        api_key: Some("test-key".to_string()),
        batch_name: "test-batch".to_string(),
        model: "grok-4.20-beta-latest-non-reasoning".to_string(),
        poll_limit: 2,
        poll_interval_ms: 1,
        results_page_size: 100,
    }
}

type AdapterFixture = (
    XaiProviderRuntimeAdapter,
    Arc<AtomicUsize>,
    Arc<AtomicUsize>,
    Arc<AtomicUsize>,
    Arc<AtomicUsize>,
    Arc<Mutex<Vec<Option<String>>>>,
);

fn build_adapter(
    config: XaiConfig,
    create_response: Result<String, XaiTransportError>,
    add_response: Result<(), XaiTransportError>,
    status_responses: Vec<Result<String, XaiTransportError>>,
    results_responses: Vec<Result<String, XaiTransportError>>,
) -> AdapterFixture {
    let create_calls = Arc::new(AtomicUsize::new(0));
    let add_calls = Arc::new(AtomicUsize::new(0));
    let status_calls = Arc::new(AtomicUsize::new(0));
    let results_calls = Arc::new(AtomicUsize::new(0));
    let captured_pagination_tokens = Arc::new(Mutex::new(Vec::new()));

    let adapter = XaiProviderRuntimeAdapter::with_config_and_transport(
        config,
        Box::new(SpyTransport {
            create_calls: Arc::clone(&create_calls),
            add_calls: Arc::clone(&add_calls),
            status_calls: Arc::clone(&status_calls),
            results_calls: Arc::clone(&results_calls),
            create_response,
            add_response,
            status_responses: Arc::new(Mutex::new(status_responses)),
            results_responses: Arc::new(Mutex::new(results_responses)),
            captured_pagination_tokens: Arc::clone(&captured_pagination_tokens),
        }),
    );

    (
        adapter,
        create_calls,
        add_calls,
        status_calls,
        results_calls,
        captured_pagination_tokens,
    )
}

fn create_response_body(batch_id: &str) -> String {
    json!({"batch_id": batch_id}).to_string()
}

fn status_body(
    num_requests: u32,
    num_pending: u32,
    num_success: u32,
    num_error: u32,
    num_cancelled: u32,
) -> String {
    json!({
        "state": {
            "num_requests": num_requests,
            "num_pending": num_pending,
            "num_success": num_success,
            "num_error": num_error,
            "num_cancelled": num_cancelled
        }
    })
    .to_string()
}

fn results_body(
    succeeded_content: Vec<&str>,
    failed_count: usize,
    pagination_token: Option<&str>,
) -> String {
    let succeeded = succeeded_content
        .into_iter()
        .map(|content| {
            json!({
                "response": {
                    "choices": [
                        {"message": {"content": content}}
                    ]
                }
            })
        })
        .collect::<Vec<_>>();

    let failed = (0..failed_count)
        .map(|index| {
            json!({
                "batch_request_id": format!("req-{index}"),
                "error_message": "request failed"
            })
        })
        .collect::<Vec<_>>();

    json!({
        "succeeded": succeeded,
        "failed": failed,
        "pagination_token": pagination_token
    })
    .to_string()
}

#[tokio::test]
async fn given_provider_mismatch_when_run_provider_step_then_return_validation_failure_before_transport(
) {
    let (adapter, create_calls, add_calls, status_calls, results_calls, _) = build_adapter(
        config_with_api_key(),
        Ok(create_response_body("batch_1")),
        Ok(()),
        vec![Ok(status_body(1, 0, 1, 0, 0))],
        vec![Ok(results_body(vec!["translated"], 0, None))],
    );

    let error = adapter
        .run_provider_step(selection("gemini", ProviderExecutionModeDto::Batch))
        .await
        .expect_err("provider mismatch must fail validation");

    assert_eq!(
        error.category,
        ExecutionControlFailureCategoryDto::ValidationFailure
    );
    assert_eq!(create_calls.load(Ordering::SeqCst), 0);
    assert_eq!(add_calls.load(Ordering::SeqCst), 0);
    assert_eq!(status_calls.load(Ordering::SeqCst), 0);
    assert_eq!(results_calls.load(Ordering::SeqCst), 0);
}

#[tokio::test]
async fn given_non_batch_mode_when_run_provider_step_then_return_validation_failure_before_transport(
) {
    let (adapter, create_calls, add_calls, status_calls, results_calls, _) = build_adapter(
        config_with_api_key(),
        Ok(create_response_body("batch_1")),
        Ok(()),
        vec![Ok(status_body(1, 0, 1, 0, 0))],
        vec![Ok(results_body(vec!["translated"], 0, None))],
    );

    let error = adapter
        .run_provider_step(selection("xai", ProviderExecutionModeDto::Streaming))
        .await
        .expect_err("non-batch mode must fail validation");

    assert_eq!(
        error.category,
        ExecutionControlFailureCategoryDto::ValidationFailure
    );
    assert_eq!(create_calls.load(Ordering::SeqCst), 0);
    assert_eq!(add_calls.load(Ordering::SeqCst), 0);
    assert_eq!(status_calls.load(Ordering::SeqCst), 0);
    assert_eq!(results_calls.load(Ordering::SeqCst), 0);
}

#[tokio::test]
async fn given_missing_private_config_when_run_provider_step_then_return_validation_failure_before_transport(
) {
    let mut config = config_with_api_key();
    config.api_key = None;

    let (adapter, create_calls, add_calls, status_calls, results_calls, _) = build_adapter(
        config,
        Ok(create_response_body("batch_1")),
        Ok(()),
        vec![Ok(status_body(1, 0, 1, 0, 0))],
        vec![Ok(results_body(vec!["translated"], 0, None))],
    );

    let error = adapter
        .run_provider_step(selection("xai", ProviderExecutionModeDto::Batch))
        .await
        .expect_err("missing API key must fail validation");

    assert_eq!(
        error.category,
        ExecutionControlFailureCategoryDto::ValidationFailure
    );
    assert_eq!(create_calls.load(Ordering::SeqCst), 0);
    assert_eq!(add_calls.load(Ordering::SeqCst), 0);
    assert_eq!(status_calls.load(Ordering::SeqCst), 0);
    assert_eq!(results_calls.load(Ordering::SeqCst), 0);
}

#[tokio::test]
async fn given_recoverable_network_api_or_poll_limit_failure_when_run_provider_step_then_return_recoverable_provider_failure(
) {
    let recoverable_cases = vec![
        build_adapter(
            config_with_api_key(),
            Err(XaiTransportError::Connection),
            Ok(()),
            vec![],
            vec![],
        )
        .0,
        build_adapter(
            config_with_api_key(),
            Err(XaiTransportError::Api { status: 503 }),
            Ok(()),
            vec![],
            vec![],
        )
        .0,
        {
            let mut config = config_with_api_key();
            config.poll_limit = 2;
            build_adapter(
                config,
                Ok(create_response_body("batch_polling")),
                Ok(()),
                vec![
                    Ok(status_body(1, 1, 0, 0, 0)),
                    Ok(status_body(1, 1, 0, 0, 0)),
                ],
                vec![],
            )
            .0
        },
    ];

    for adapter in recoverable_cases {
        let error = adapter
            .run_provider_step(selection("xai", ProviderExecutionModeDto::Batch))
            .await
            .expect_err("recoverable case must fail with retryable category");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::RecoverableProviderFailure
        );
        assert!(!error.message.contains("batch_polling"));
        assert!(!error.message.contains("/v1/batches"));
        assert!(!error.message.contains("Authorization"));
        assert!(!error.message.contains("test-key"));
    }
}

#[tokio::test]
async fn given_unrecoverable_api_malformed_or_unusable_results_when_run_provider_step_then_return_unrecoverable_provider_failure(
) {
    let unrecoverable_cases = vec![
        build_adapter(
            config_with_api_key(),
            Err(XaiTransportError::Api { status: 400 }),
            Ok(()),
            vec![],
            vec![],
        )
        .0,
        build_adapter(
            config_with_api_key(),
            Ok(create_response_body("batch_1")),
            Ok(()),
            vec![Ok(status_body(1, 0, 1, 0, 0))],
            vec![Ok("not-json".to_string())],
        )
        .0,
        build_adapter(
            config_with_api_key(),
            Ok(create_response_body("batch_1")),
            Ok(()),
            vec![Ok(status_body(1, 0, 1, 0, 0))],
            vec![Ok(results_body(vec![], 1, None))],
        )
        .0,
        build_adapter(
            config_with_api_key(),
            Ok(create_response_body("batch_1")),
            Ok(()),
            vec![Ok(status_body(1, 0, 1, 0, 0))],
            vec![Ok(results_body(vec!["   "], 0, None))],
        )
        .0,
    ];

    for adapter in unrecoverable_cases {
        let error = adapter
            .run_provider_step(selection("xai", ProviderExecutionModeDto::Batch))
            .await
            .expect_err("unrecoverable case must fail");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure
        );
        assert!(!error.message.contains("batch_1"));
        assert!(!error.message.contains("/v1/batches"));
        assert!(!error.message.contains("test-key"));
    }
}

#[tokio::test]
async fn given_paginated_success_when_run_provider_step_then_fetch_all_pages_and_succeed() {
    let (adapter, create_calls, add_calls, status_calls, results_calls, captured_pagination_tokens) =
        build_adapter(
            config_with_api_key(),
            Ok(create_response_body("batch_paged")),
            Ok(()),
            vec![Ok(status_body(1, 0, 1, 0, 0))],
            vec![
                Ok(results_body(vec![""], 0, Some("next_page_token"))),
                Ok(results_body(vec!["translated"], 0, None)),
            ],
        );

    adapter
        .run_provider_step(selection("xai", ProviderExecutionModeDto::Batch))
        .await
        .expect("paginated successful results must succeed");

    assert_eq!(create_calls.load(Ordering::SeqCst), 1);
    assert_eq!(add_calls.load(Ordering::SeqCst), 1);
    assert_eq!(status_calls.load(Ordering::SeqCst), 1);
    assert_eq!(results_calls.load(Ordering::SeqCst), 2);

    let tokens = captured_pagination_tokens
        .lock()
        .expect("pagination token mutex");
    assert_eq!(
        tokens.as_slice(),
        &[None, Some("next_page_token".to_string())]
    );
}

#[tokio::test]
async fn given_nonzero_poll_interval_and_pending_then_complete_when_run_provider_step_then_poll_progresses_and_succeeds(
) {
    let (adapter, create_calls, add_calls, status_calls, results_calls, _) = build_adapter(
        config_with_api_key(),
        Ok(create_response_body("batch_pending_then_complete")),
        Ok(()),
        vec![
            Ok(status_body(1, 1, 0, 0, 0)),
            Ok(status_body(1, 0, 1, 0, 0)),
        ],
        vec![Ok(results_body(vec!["translated"], 0, None))],
    );

    adapter
        .run_provider_step(selection("xai", ProviderExecutionModeDto::Batch))
        .await
        .expect("polling should progress from pending to complete with nonzero interval");

    assert_eq!(create_calls.load(Ordering::SeqCst), 1);
    assert_eq!(add_calls.load(Ordering::SeqCst), 1);
    assert_eq!(status_calls.load(Ordering::SeqCst), 2);
    assert_eq!(results_calls.load(Ordering::SeqCst), 1);
}
