use tauri::api::process::Command;


// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

fn main() {
  tauri::Builder::default()
    .run(tauri::generate_context!())
    .expect("error while running tauri application");

  // `new_sidecar()` expects just the filename, NOT the whole path like in JavaScript
let (mut rx, mut child) = Command::new_sidecar("quicksilverd")

// Create a HashMap to store the environment variables
let mut env_vars = HashMap::new();
env_vars.insert("SOME_ENV_VAR".to_string(), "dynamic_value1".to_string());
env_vars.insert("ANOTHER_ENV_VAR".to_string(), "dynamic_value2".to_string());
.expect("failed to create `my-sidecar` binary command")
.spawn()
.expect("Failed to spawn sidecar");

tauri::async_runtime::spawn(async move {
// read events such as stdout
while let Some(event) = rx.recv().await {
if let CommandEvent::Stdout(line) = event {
window
.emit("message", Some(format!("'{}'", line)))
.expect("failed to emit event");
// write to stdin
child.write("message from Rust\n".as_bytes()).unwrap();
}
}
});
}




