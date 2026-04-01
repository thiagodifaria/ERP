use std::sync::Arc;

use tokio_postgres::NoTls;

use crate::config::PostgresConfig;

pub async fn connect_postgres(
    config: &PostgresConfig,
) -> Result<Arc<tokio_postgres::Client>, tokio_postgres::Error> {
    let mut postgres_config = tokio_postgres::Config::new();
    postgres_config.host(&config.host);
    postgres_config.port(config.port);
    postgres_config.dbname(&config.db);
    postgres_config.user(&config.user);
    postgres_config.password(&config.password);
    let _ssl_mode = config.ssl_mode.to_lowercase();

    let (client, connection) = postgres_config.connect(NoTls).await?;

    tokio::spawn(async move {
        if let Err(error) = connection.await {
            eprintln!("webhook-hub postgres connection stopped: {error}");
        }
    });

    Ok(Arc::new(client))
}
