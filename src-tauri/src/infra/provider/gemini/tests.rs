use super::*;
use std::sync::{
    atomic::{AtomicUsize, Ordering},
    Arc, Mutex, MutexGuard, OnceLock,
};

use async_trait::async_trait;
use serde_json::{json, Value};

use crate::application::dto::{
    ExecutionControlFailureCategoryDto, ProviderExecutionModeDto, ProviderRuntimeSettingsDto,
    ProviderSelectionDto,
};

const GEMINI_API_KEY_ENV: &str = "GEMINI_API_KEY";

struct EnvVarGuard {
    _lock: MutexGuard<'static, ()>,
    key: &'static str,
    previous: Option<String>,
}

impl EnvVarGuard {
    fn lock() -> &'static Mutex<()> {
        static LOCK: OnceLock<Mutex<()>> = OnceLock::new();
        LOCK.get_or_init(|| Mutex::new(()))
    }

    fn set(key: &'static str, value: &str) -> Self {
        let lock = Self::lock()
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner());
        let previous = std::env::var(key).ok();
        std::env::set_var(key, value);
        Self {
            _lock: lock,
            key,
            previous,
        }
    }

    fn remove(key: &'static str) -> Self {
        let lock = Self::lock()
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner());
        let previous = std::env::var(key).ok();
        std::env::remove_var(key);
        Self {
            _lock: lock,
            key,
            previous,
        }
    }
}

impl Drop for EnvVarGuard {
    fn drop(&mut self) {
        if let Some(previous) = &self.previous {
            std::env::set_var(self.key, previous);
        } else {
            std::env::remove_var(self.key);
        }
    }
}

struct SpyTransport {
    calls: Arc<AtomicUsize>,
    response: Result<String, GeminiTransportError>,
    requests: Arc<Mutex<Vec<(String, Value)>>>,
}

#[async_trait]
impl GeminiTransport for SpyTransport {
    async fn execute(
        &self,
        endpoint: &str,
        request: &GeminiRequest,
    ) -> Result<String, GeminiTransportError> {
        self.calls.fetch_add(1, Ordering::SeqCst);
        self.requests.lock().expect("request capture mutex").push((
            endpoint.to_string(),
            serde_json::to_value(request).expect("request must serialize for test inspection"),
        ));

        self.response.clone()
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

fn usable_candidate_body() -> String {
    json!({
        "candidates": [{
            "content": {
                "parts": [{"text": "translated"}]
            },
            "finishReason": "STOP"
        }]
    })
    .to_string()
}

fn blocked_prompt_body() -> String {
    json!({
        "promptFeedback": {
            "blockReason": "SAFETY"
        }
    })
    .to_string()
}

fn no_candidate_body() -> String {
    json!({
        "candidates": []
    })
    .to_string()
}

#[tokio::test]
async fn given_provider_mismatch_when_run_provider_step_then_return_validation_failure_without_transport(
) {
    let _api_key = EnvVarGuard::set(GEMINI_API_KEY_ENV, "test-key");
    let calls = Arc::new(AtomicUsize::new(0));
    let adapter = GeminiProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
        calls: Arc::clone(&calls),
        response: Ok(usable_candidate_body()),
        requests: Arc::new(Mutex::new(Vec::new())),
    }));

    let error = adapter
        .run_provider_step(selection("lmstudio", ProviderExecutionModeDto::Batch))
        .await
        .expect_err("provider mismatch must fail validation");

    assert_eq!(
        error.category,
        ExecutionControlFailureCategoryDto::ValidationFailure
    );
    assert_eq!(calls.load(Ordering::SeqCst), 0);
}

#[tokio::test]
async fn given_streaming_mode_when_run_provider_step_then_return_validation_failure_without_transport(
) {
    let _api_key = EnvVarGuard::set(GEMINI_API_KEY_ENV, "test-key");
    let calls = Arc::new(AtomicUsize::new(0));
    let adapter = GeminiProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
        calls: Arc::clone(&calls),
        response: Ok(usable_candidate_body()),
        requests: Arc::new(Mutex::new(Vec::new())),
    }));

    let error = adapter
        .run_provider_step(selection("gemini", ProviderExecutionModeDto::Streaming))
        .await
        .expect_err("streaming mode must fail validation");

    assert_eq!(
        error.category,
        ExecutionControlFailureCategoryDto::ValidationFailure
    );
    assert_eq!(calls.load(Ordering::SeqCst), 0);
}

#[tokio::test]
async fn given_missing_api_key_when_run_provider_step_then_return_validation_failure_without_transport(
) {
    let _missing_api_key = EnvVarGuard::remove(GEMINI_API_KEY_ENV);
    let calls = Arc::new(AtomicUsize::new(0));
    let adapter = GeminiProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
        calls: Arc::clone(&calls),
        response: Ok(usable_candidate_body()),
        requests: Arc::new(Mutex::new(Vec::new())),
    }));

    let error = adapter
        .run_provider_step(selection("gemini", ProviderExecutionModeDto::Batch))
        .await
        .expect_err("missing API key must fail validation");

    assert_eq!(
        error.category,
        ExecutionControlFailureCategoryDto::ValidationFailure
    );
    assert_eq!(calls.load(Ordering::SeqCst), 0);
    assert!(!error.message.contains("test-key"));
}

#[tokio::test]
async fn given_recoverable_transport_failures_when_run_provider_step_then_return_recoverable_provider_failure(
) {
    let _api_key = EnvVarGuard::set(GEMINI_API_KEY_ENV, "test-key");

    for failure in [
        GeminiTransportError::Connection,
        GeminiTransportError::Api { status: 429 },
        GeminiTransportError::Api { status: 500 },
        GeminiTransportError::Api { status: 503 },
        GeminiTransportError::Api { status: 504 },
    ] {
        let adapter = GeminiProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
            calls: Arc::new(AtomicUsize::new(0)),
            response: Err(failure),
            requests: Arc::new(Mutex::new(Vec::new())),
        }));

        let error = adapter
            .run_provider_step(selection("gemini", ProviderExecutionModeDto::Batch))
            .await
            .expect_err("recoverable Gemini failure must not succeed");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::RecoverableProviderFailure
        );
        assert!(!error.message.contains("test-key"));
        assert!(!error.message.contains("generativelanguage.googleapis.com"));
    }
}

#[tokio::test]
async fn given_permanent_api_failures_when_run_provider_step_then_return_unrecoverable_provider_failure(
) {
    let _api_key = EnvVarGuard::set(GEMINI_API_KEY_ENV, "test-key");

    for status in [400_u16, 403_u16, 404_u16] {
        let adapter = GeminiProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
            calls: Arc::new(AtomicUsize::new(0)),
            response: Err(GeminiTransportError::Api { status }),
            requests: Arc::new(Mutex::new(Vec::new())),
        }));

        let error = adapter
            .run_provider_step(selection("gemini", ProviderExecutionModeDto::Batch))
            .await
            .expect_err("permanent Gemini API failure must not succeed");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure
        );
        assert!(!error.message.contains(&status.to_string()));
        assert!(!error.message.contains("test-key"));
    }
}

#[tokio::test]
async fn given_blocked_or_empty_candidate_body_when_run_provider_step_then_return_unrecoverable_provider_failure(
) {
    let _api_key = EnvVarGuard::set(GEMINI_API_KEY_ENV, "test-key");

    for body in [blocked_prompt_body(), no_candidate_body()] {
        let adapter = GeminiProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
            calls: Arc::new(AtomicUsize::new(0)),
            response: Ok(body.clone()),
            requests: Arc::new(Mutex::new(Vec::new())),
        }));

        let error = adapter
            .run_provider_step(selection("gemini", ProviderExecutionModeDto::Batch))
            .await
            .expect_err("blocked or empty candidate body must not succeed");

        assert_eq!(
            error.category,
            ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure
        );
        assert!(!error.message.contains("SAFETY"));
        assert!(!error.message.contains("promptFeedback"));
    }
}

#[tokio::test]
async fn given_malformed_body_when_run_provider_step_then_return_unrecoverable_provider_failure_without_raw_body_leak(
) {
    let _api_key = EnvVarGuard::set(GEMINI_API_KEY_ENV, "test-key");
    let malformed = "{\"unexpected\":true}";
    let adapter = GeminiProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
        calls: Arc::new(AtomicUsize::new(0)),
        response: Ok(malformed.to_string()),
        requests: Arc::new(Mutex::new(Vec::new())),
    }));

    let error = adapter
        .run_provider_step(selection("gemini", ProviderExecutionModeDto::Batch))
        .await
        .expect_err("malformed body must not succeed");

    assert_eq!(
        error.category,
        ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure
    );
    assert!(!error.message.contains("unexpected"));
    assert!(!error.message.contains("true"));
}

#[tokio::test]
async fn given_usable_candidate_when_run_provider_step_then_succeed_and_send_generate_content_request(
) {
    let _api_key = EnvVarGuard::set(GEMINI_API_KEY_ENV, "test-key");
    let requests = Arc::new(Mutex::new(Vec::new()));
    let adapter = GeminiProviderRuntimeAdapter::with_transport(Box::new(SpyTransport {
        calls: Arc::new(AtomicUsize::new(0)),
        response: Ok(usable_candidate_body()),
        requests: Arc::clone(&requests),
    }));

    let result = adapter
        .run_provider_step(selection("gemini", ProviderExecutionModeDto::Batch))
        .await;

    assert!(result.is_ok(), "usable candidate response must complete");

    let captured = requests.lock().expect("request capture mutex");
    let (endpoint, request) = captured.first().expect("one request must be captured");
    assert!(endpoint.contains("generativelanguage.googleapis.com/v1beta/models/"));
    assert!(endpoint.contains(":generateContent"));
    assert_eq!(
        request["contents"][0]["parts"][0]["text"],
        "Run one provider step."
    );
}
