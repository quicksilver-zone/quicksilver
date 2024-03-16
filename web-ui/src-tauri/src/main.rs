use reqwest;
use serde_json::Value;
use std::collections::HashMap;
use std::fs::write;
use std::io::BufRead;
use std::process::{Command, Stdio};
use tauri::api::file::read_binary;
use tauri::Manager;

fn main() {
    // Prevents additional console window on Windows in release, DO NOT REMOVE!!
    #![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]
    tauri::Builder::default()
        .setup(|app| {
            let window = app.get_window("main").unwrap();

            // Create a HashMap to store the environment variables
            let mut env_vars = HashMap::new();
            env_vars.insert(
                "QUICKSILVERD_STATESYNC_ENABLE".to_string(),
                "true".to_string(),
            );
            env_vars.insert(
                "QUICKSILVERD_P2P_MAX_NUM_OUTBOUND_PEERS".to_string(),
                "200".to_string(),
            );
            env_vars.insert(
                "QUICKSILVERD_STATESYNC_RPC_SERVERS".to_string(),
                "https://quicksilver-rpc.polkachu.com:443,https://quicksilver-rpc.polkachu.com:443"
                    .to_string(),
            );

            // Get the latest block height and calculate the block height for state sync
            let latest_height = get_latest_block_height();
            let block_height = latest_height - 100;
            let trust_hash = get_trust_hash(block_height);

            env_vars.insert(
                "QUICKSILVERD_STATESYNC_TRUST_HEIGHT".to_string(),
                block_height.to_string(),
            );
            env_vars.insert("QUICKSILVERD_STATESYNC_TRUST_HASH".to_string(), trust_hash);

            // Initialize the chain on the first launch
            let home_dir = tauri::api::path::home_dir().unwrap();
            let genesis_path = home_dir.join(".quicksilverd/config/genesis.json");
            if !genesis_path.exists() {
                let mut init_command = Command::new("quicksilverd");
                init_command.arg("init").arg("test");
                let _ = init_command.output().expect("Failed to initialize chain");
            }

            // Copy the file from resources to $HOME/.quicksilverd/config/genesis.json
            let resource_path = "./resources/genesis.json";
            let destination_path = tauri::api::path::home_dir()
                .unwrap()
                .join(".quicksilverd/config/genesis.json");

            match read_binary(resource_path) {
                Ok(file_contents) => match write(&destination_path, &file_contents) {
                    Ok(_) => {
                        println!("File copied successfully!");
                    }
                    Err(e) => {
                        eprintln!("Error writing file: {}", e);
                    }
                },
                Err(e) => {
                    eprintln!("Error reading file: {}", e);
                }
            }

            // Spawn the quicksilverd process and capture its output
            let mut child = Command::new("quicksilverd")
                .arg("start")
                .arg("--x-crisis-skip-assert-invariants")
                .arg("--iavl-disable-fastnode=false")
                .envs(&env_vars)
                .stdout(Stdio::piped())
                .stderr(Stdio::piped())
                .spawn()
                .expect("Failed to spawn quicksilverd");

            // Show the application window
            window.show().unwrap();

            // Read the stdout of the quicksilverd process
            let stdout = child.stdout.take().unwrap();
            let stdout_reader = std::io::BufReader::new(stdout);

            // Read the stderr of the quicksilverd process
            let stderr = child.stderr.take().unwrap();
            let stderr_reader = std::io::BufReader::new(stderr);

            // Spawn a thread to read and print the stdout and stderr
            std::thread::spawn(move || {
                // Inside your thread spawn, after creating stdout_reader and stderr_reader
                let mut stdout_lines = stdout_reader.lines().peekable();
                let mut stderr_lines = stderr_reader.lines().peekable();

                loop {
                    if stdout_lines.peek().is_none() && stderr_lines.peek().is_none() {
                        break;
                    }

                    if let Some(Ok(output)) = stdout_lines.next() {
                        println!("stdout: {}", output);
                        window
                            .emit("message", Some(format!("'{}'", output)))
                            .unwrap();
                    }

                    if let Some(Ok(output)) = stderr_lines.next() {
                        eprintln!("stderr: {}", output);
                    }
                }
            });

            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("Error while running tauri application");
}

fn get_latest_block_height() -> u64 {
    let url = "https://quicksilver-rpc.polkachu.com/block";
    let response = reqwest::blocking::get(url).expect("Failed to fetch latest block height");
    let json: Value = response.json().expect("Failed to parse JSON response");
    let block_height = json["result"]["block"]["header"]["height"]
        .as_str()
        .expect("Failed to extract block height")
        .parse()
        .expect("Failed to parse block height as u64");
    block_height
}

fn get_trust_hash(block_height: u64) -> String {
    let url = format!(
        "https://quicksilver-rpc.polkachu.com/block?height={}",
        block_height
    );
    let response = reqwest::blocking::get(&url).expect("Failed to fetch trust hash");
    let json: Value = response.json().expect("Failed to parse JSON response");
    let trust_hash = json["result"]["block_id"]["hash"]
        .as_str()
        .expect("Failed to extract trust hash")
        .to_string();
    trust_hash
}
