// Este arquivo inicia a API e delega o wiring para a composicao do servico.
// Regra de negocio nao deve ser implementada aqui.
using Identity.Api;
using Identity.Application;
using Identity.Infrastructure;

var builder = WebApplication.CreateBuilder(args);

builder.Services
  .AddIdentityApplication()
  .AddIdentityInfrastructure();

var app = builder.Build();

app.MapIdentityRoutes();

app.Run();

public partial class Program;
