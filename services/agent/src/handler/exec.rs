use anyhow::anyhow;
use async_process::{Command, Stdio};
use async_std::io::WriteExt;
use axum::extract::Json;
use axum::http::StatusCode;
use serde::{Deserialize, Serialize};
use std::fs::File;
use std::io::Write;

#[derive(Deserialize)]
pub struct RequestPayload {
    id: String,
    code: String,
    input: String,
    language: String,
}

#[derive(Serialize)]
pub struct ResponsePayload {
    stdout: String,
    stderr: String,
    exit_code: i32,
}

pub async fn handler(
    Json(payload): Json<RequestPayload>,
) -> Result<Json<ResponsePayload>, StatusCode> {
    match exec_code(payload.id, payload.code, payload.input, payload.language).await {
        Ok(response) => Ok(Json(response)),
        Err(_) => Err(StatusCode::INTERNAL_SERVER_ERROR),
    }
}

async fn exec_code(
    id: String,
    code: String,
    input: String,
    language: String,
) -> Result<ResponsePayload, anyhow::Error> {
    println!("Executing code: {}", id);
    let filename = format!("/tmp/{}.{}", id, language);
    let mut file = File::create(&filename)?;
    file.write_all(code.as_bytes())?;

    let command = match language.as_str() {
        "py" => "python",
        "js" => "node",
        "sh" => "sh",
        _ => return Err(anyhow!("Unsupported language")),
    };

    let mut process = Command::new("sh")
        .arg("-c")
        .arg(format!("{} {}", command, filename))
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()?;

    if let Some(stdin) = process.stdin.as_mut() {
        stdin.write_all(input.as_bytes()).await?;
    } else {
        return Err(anyhow!("Failed to open stdin"));
    }

    let output = process.output().await?;

    Ok(ResponsePayload {
        stdout: String::from_utf8(output.stdout)?,
        stderr: String::from_utf8(output.stderr)?,
        exit_code: output.status.code().unwrap_or_default(),
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_exec_code() {
        let payload = RequestPayload {
            id: "test".to_owned(),
            code: "cat".to_owned(),
            input: "Hello, World!\n".to_owned(),
            language: "sh".to_owned(),
        };

        let response = exec_code(payload.id, payload.code, payload.input, payload.language)
            .await
            .unwrap();

        assert_eq!(response.stdout, "Hello, World!\n");
        assert_eq!(response.stderr, "");
        assert_eq!(response.exit_code, 0);
    }
}
