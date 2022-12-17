use axum::Router;
use std::net::SocketAddr;

mod handler;

#[tokio::main]
async fn main() {
    let app = Router::new().merge(handler::routes());
    let addr = SocketAddr::from(([0, 0, 0, 0], 3005));

    println!("Agent server started on {}", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
