use actix_web::{web, HttpResponse, Result};
use sqlx::SqlitePool;
use crate::models::{VoteRequest, Vote, ReportItem, VoteResponse};

pub async fn submit_vote(
    pool: web::Data<SqlitePool>,
    body: web::Json<VoteRequest>,
) -> Result<HttpResponse> {
    if body.character_id.is_empty() {
        return Ok(HttpResponse::BadRequest()
            .json(serde_json::json!({"error": "character_id is required"})));
    }

    let result = sqlx::query(
        r#"
        INSERT INTO votes (character_id, name, image_url)
        VALUES (?, ?, ?)
        "#,
    )
    .bind(&body.character_id)
    .bind(&body.name)
    .bind(&body.image_url)
    .execute(pool.get_ref())
    .await;

    match result {
        Ok(row) => {
            let id = row.last_insert_rowid();
            Ok(HttpResponse::Created().json(VoteResponse { id }))
        }
        Err(_) => Ok(HttpResponse::InternalServerError()
            .json(serde_json::json!({"error": "Failed to insert vote"}))),
    }
}

pub async fn get_report(pool: web::Data<SqlitePool>) -> Result<HttpResponse> {
    let rows: Vec<ReportItem> = sqlx::query_as(
        r#"
        SELECT character_id, name, image_url, COUNT(*) as votes
        FROM votes
        GROUP BY character_id
        ORDER BY votes DESC
        "#,
    )
    .fetch_all(pool.get_ref())
    .await
    .unwrap_or_default();

    Ok(HttpResponse::Ok().json(rows))
}

pub async fn get_votes(pool: web::Data<SqlitePool>) -> Result<HttpResponse> {
    let rows: Vec<Vote> = sqlx::query_as(
        r#"
        SELECT id, character_id, name, image_url, 
               strftime('%Y-%m-%d %H:%M:%S', created_at) as created_at
        FROM votes
        ORDER BY created_at DESC
        "#,
    )
    .fetch_all(pool.get_ref())
    .await
    .unwrap_or_default();

    Ok(HttpResponse::Ok().json(rows))
}

pub async fn export_csv(pool: web::Data<SqlitePool>) -> Result<HttpResponse> {
    let rows: Vec<(String, Option<String>, i64)> = sqlx::query_as(
        r#"
        SELECT character_id, name, COUNT(*) as votes
        FROM votes
        GROUP BY character_id
        ORDER BY votes DESC
        "#,
    )
    .fetch_all(pool.get_ref())
    .await
    .unwrap_or_default();

    let mut csv_data = String::from("characterId,name,votes\n");
    for (char_id, name, votes) in rows {
        let name_escaped = name.unwrap_or_default().replace("\"", "\"\"");
        csv_data.push_str(&format!(
            "{},\"{}\",{}\n",
            char_id, name_escaped, votes
        ));
    }

    Ok(HttpResponse::Ok()
        .content_type("text/csv")
        .insert_header(("Content-Disposition", "attachment; filename=\"report.csv\""))
        .body(csv_data))
}
