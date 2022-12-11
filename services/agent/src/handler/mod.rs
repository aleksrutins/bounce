mod exec;
mod health;

use axum::{
    routing::{get, post},
    Router,
};

pub fn routes() -> Router {
    Router::new()
        .route("/health", get(health::handler))
        .route("/exec", post(exec::handler))
}
