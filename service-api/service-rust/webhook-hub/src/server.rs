// Server concentra o bootstrap HTTP e a subida do listener.
// Validacao de assinatura e normalizacao entram em camadas internas.
use crate::api::router::build_router_from_config;
use crate::config::AppConfig;

pub async fn run(config: AppConfig) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let app = build_router_from_config(&config).await?;
    let listener = tokio::net::TcpListener::bind(&config.http_address).await?;
    println!("starting {} on {}", config.service_name, config.http_address);

    axum::serve(listener, app).await?;
    Ok(())
}
