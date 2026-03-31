// Config centraliza o acesso a variaveis de ambiente do servico.
// Segredos e endpoints externos devem entrar por esta camada.
pub struct AppConfig {
    pub service_name: &'static str,
    pub http_address: String,
}

impl AppConfig {
    pub fn from_env() -> Self {
        let http_address =
            std::env::var("WEBHOOK_HUB_HTTP_ADDRESS").unwrap_or_else(|_| "0.0.0.0:8082".to_string());

        Self {
            service_name: "webhook-hub",
            http_address,
        }
    }
}
