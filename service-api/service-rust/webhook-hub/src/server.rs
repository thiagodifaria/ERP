// Server concentra o bootstrap HTTP e a subida do listener.
// Validacao de assinatura e normalizacao entram em camadas internas.
use crate::api::router::build_router;
use crate::config::AppConfig;

pub async fn run(config: AppConfig) -> Result<(), std::io::Error> {
    let listener = tokio::net::TcpListener::bind(&config.http_address).await?;
    println!("starting {} on {}", config.service_name, config.http_address);

    axum::serve(listener, build_router()).await
}
