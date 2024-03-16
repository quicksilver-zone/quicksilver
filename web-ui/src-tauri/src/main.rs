use reqwest;
use serde_json::Value;
use std::collections::HashMap;
use std::process::Stdio;
use tauri::api::process::CommandEvent;
use tauri::Manager;
use tokio::io::AsyncBufReadExt;
use tokio::sync::mpsc;

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
                "https://quicksilver-rpc.polkachu.com:443,https://quicksilver-rpc.polkachu.com:443".to_string(),
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
                let mut init_command = std::process::Command::new("quicksilverd");
                init_command.arg("init").arg("test");
                let _ = init_command.output().expect("Failed to initialize chain");
            }

            // copy the genesis file to the quicksilverd config directory
            let config_dir = home_dir.join(".quicksilverd/config");
    

            let rt = tokio::runtime::Runtime::new().unwrap();
            rt.block_on(async move {
                // Create a Command to launch the quicksilverd process
                let mut quicksilverd_command = tokio::process::Command::new("quicksilverd");
                quicksilverd_command
                    .arg("start")
                    .arg("--x-crisis-skip-assert-invariants")
                    .arg("--iavl-disable-fastnode=false");

                // Set the environment variables
                for (key, value) in env_vars {
                    quicksilverd_command.env(key, value);
                }

                let (tx, mut rx) = mpsc::channel::<CommandEvent>(32); // Adjust the buffer size as needed

                // Spawn the quicksilverd process
                let mut child = quicksilverd_command
                    .stdout(Stdio::piped())
                    .spawn()
                    .expect("Failed to spawn quicksilverd");

                let stdout = child.stdout.take().expect("Failed to capture stdout");
                let tx_clone = tx.clone();
                tokio::spawn(async move {
                    let mut reader = tokio::io::BufReader::new(stdout).lines();
                    while let Ok(Some(line)) = reader.next_line().await {
                        tx_clone
                            .send(CommandEvent::Stdout(line))
                            .await
                            .expect("Failed to send event");
                    }
                });

                // Read events from the quicksilverd process
                while let Some(event) = rx.recv().await {
                    if let CommandEvent::Stdout(line) = event {
                        window
                            .emit("message", Some(format!("'{}'", line)))
                            .expect("Failed to emit event");
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
