// Config centraliza o acesso a variaveis de ambiente do servico.
// Segredos e endpoints externos devem entrar por esta camada.
#[derive(Clone)]
pub struct AppConfig {
    pub service_name: &'static str,
    pub http_address: String,
    pub repository_driver: RepositoryDriver,
    pub postgres: PostgresConfig,
}

#[derive(Clone, Copy, PartialEq, Eq)]
pub enum RepositoryDriver {
    Memory,
    Postgres,
}

#[derive(Clone)]
pub struct PostgresConfig {
    pub host: String,
    pub port: u16,
    pub db: String,
    pub user: String,
    pub password: String,
    pub ssl_mode: String,
}

impl AppConfig {
    pub fn from_env() -> Self {
        let http_address =
            std::env::var("WEBHOOK_HUB_HTTP_ADDRESS").unwrap_or_else(|_| "0.0.0.0:8082".to_string());
        let repository_driver = match std::env::var("WEBHOOK_HUB_REPOSITORY_DRIVER")
            .unwrap_or_else(|_| "memory".to_string())
            .to_lowercase()
            .as_str()
        {
            "postgres" => RepositoryDriver::Postgres,
            _ => RepositoryDriver::Memory,
        };

        Self {
            service_name: "webhook-hub",
            http_address,
            repository_driver,
            postgres: PostgresConfig {
                host: std::env::var("WEBHOOK_HUB_POSTGRES_HOST")
                    .unwrap_or_else(|_| "service-postgresql".to_string()),
                port: std::env::var("WEBHOOK_HUB_POSTGRES_PORT")
                    .ok()
                    .and_then(|value| value.parse().ok())
                    .unwrap_or(5432),
                db: std::env::var("WEBHOOK_HUB_POSTGRES_DB").unwrap_or_else(|_| "erp".to_string()),
                user: std::env::var("WEBHOOK_HUB_POSTGRES_USER").unwrap_or_else(|_| "erp".to_string()),
                password: std::env::var("WEBHOOK_HUB_POSTGRES_PASSWORD")
                    .unwrap_or_else(|_| "erp".to_string()),
                ssl_mode: std::env::var("WEBHOOK_HUB_POSTGRES_SSL_MODE")
                    .unwrap_or_else(|_| "disable".to_string()),
            },
        }
    }
}
