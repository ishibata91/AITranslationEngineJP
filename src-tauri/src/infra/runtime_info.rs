pub trait RuntimeInfoProvider {
    fn backend_version(&self) -> String;
}

pub struct CargoRuntimeInfoProvider;

impl RuntimeInfoProvider for CargoRuntimeInfoProvider {
    fn backend_version(&self) -> String {
        env!("CARGO_PKG_VERSION").to_string()
    }
}
