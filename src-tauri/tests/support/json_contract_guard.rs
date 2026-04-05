use serde_json::Value;

pub fn collect_forbidden_key_paths(value: &Value, forbidden_keys: &[&str]) -> Vec<String> {
    let mut paths = vec![];
    collect(value, "$", forbidden_keys, &mut paths);
    paths
}

fn collect(value: &Value, current_path: &str, forbidden_keys: &[&str], paths: &mut Vec<String>) {
    match value {
        Value::Object(map) => {
            for (key, nested) in map {
                let next_path = format!("{current_path}.{key}");
                if forbidden_keys.iter().any(|forbidden| forbidden == key) {
                    paths.push(next_path.clone());
                }
                collect(nested, &next_path, forbidden_keys, paths);
            }
        }
        Value::Array(items) => {
            for (index, nested) in items.iter().enumerate() {
                collect(
                    nested,
                    &format!("{current_path}[{index}]"),
                    forbidden_keys,
                    paths,
                );
            }
        }
        Value::Null | Value::Bool(_) | Value::Number(_) | Value::String(_) => {}
    }
}
