// Este arquivo inicia o servico e delega a inicializacao do runtime.
// Regra de negocio nao deve ser implementada aqui.
mod api;
mod config;
mod server;
mod telemetry;

#[tokio::main]
async fn main() {
    telemetry::init();

    let config = config::AppConfig::from_env();

    if let Err(error) = server::run(config).await {
        eprintln!("server stopped with error: {error}");
        std::process::exit(1);
    }
}
