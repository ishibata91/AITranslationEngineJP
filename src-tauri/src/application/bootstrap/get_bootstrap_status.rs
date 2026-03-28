use crate::application::dto::BootstrapStatusDto;
use crate::domain::bootstrap_status::BootstrapStatus;
use crate::infra::runtime_info::RuntimeInfoProvider;

pub struct GetBootstrapStatusUseCase<P>
where
    P: RuntimeInfoProvider,
{
    provider: P,
}

impl<P> GetBootstrapStatusUseCase<P>
where
    P: RuntimeInfoProvider,
{
    pub fn new(provider: P) -> Self {
        Self { provider }
    }

    pub fn execute(&self) -> BootstrapStatusDto {
        let status = BootstrapStatus::initial(self.provider.backend_version());
        BootstrapStatusDto::from(status)
    }
}

#[cfg(test)]
mod tests {
    use super::GetBootstrapStatusUseCase;
    use crate::infra::runtime_info::RuntimeInfoProvider;

    struct StubRuntimeInfoProvider;

    impl RuntimeInfoProvider for StubRuntimeInfoProvider {
        fn backend_version(&self) -> String {
            "test-version".to_string()
        }
    }

    #[test]
    fn execute_returns_bootstrap_status_dto() {
        let use_case = GetBootstrapStatusUseCase::new(StubRuntimeInfoProvider);

        let dto = use_case.execute();

        assert_eq!(dto.backend_version, "test-version");
        assert!(dto.boundary_ready);
        assert_eq!(dto.frontend_entry, "src/main.ts");
    }
}
