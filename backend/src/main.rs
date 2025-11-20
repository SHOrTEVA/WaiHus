mod db;
mod models;
mod handlers;

use actix_web::{web, App, HttpServer, middleware};
use actix_cors::Cors;
use std::sync::Arc;
use env_logger::init_from_env;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

    // Initialize database
    let pool = db::init_db()
        .await
        .expect("Failed to initialize database");

    let pool = web::Data::new(pool);

    println!("Starting server on http://0.0.0.0:3000");

    HttpServer::new(move || {
        App::new()
            .app_data(pool.clone())
            .wrap(middleware::Logger::default())
            .wrap(
                actix_web::middleware::DefaultHeaders::new()
                    .add(("Access-Control-Allow-Origin", "*"))
                    .add(("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS"))
                    .add(("Access-Control-Allow-Headers", "Content-Type"))
            )
            .wrap(
                Cors::default()
                    // For local development prefer explicitly listing the origin:
                    //.allowed_origin("http://localhost:3000")
                    // Or for quick testing only: 
                    .allow_any_origin()
                    .allowed_methods(vec!["POST", "GET", "OPTIONS"])
                    .allowed_headers(vec![header::CONTENT_TYPE, header::AUTHORIZATION, header::ACCEPT])
                    .max_age(3600)
            )
            .route("/api/vote", web::post().to(handlers::submit_vote))
            .route("/api/report", web::get().to(handlers::get_report))
            .route("/api/votes", web::get().to(handlers::get_votes))
            .route("/api/export", web::get().to(handlers::export_csv))
    })
    .bind("0.0.0.0:3000")?
    .run()
    .await
}
