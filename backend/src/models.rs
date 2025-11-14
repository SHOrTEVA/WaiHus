use serde::{Deserialize, Serialize};
use sqlx::FromRow;

#[derive(Serialize, Deserialize, Debug)]
pub struct VoteRequest {
    pub character_id: String,
    pub name: Option<String>,
    pub image_url: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, FromRow)]
pub struct Vote {
    pub id: i64,
    pub character_id: String,
    pub name: Option<String>,
    pub image_url: Option<String>,
    pub created_at: String,
}

#[derive(sqlx::FromRow, Serialize, Deserialize, Debug)]
pub struct ReportItem {
    pub character_id: String,
    pub name: Option<String>,
    pub image_url: Option<String>,
    pub votes: i64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct VoteResponse {
    pub id: i64,
}
