# Business Operating System Control Console

Console técnico para explorar a API do projeto, ler documentação versionada, validar contratos e testar rotas em um único lugar.

## Rodar localmente

```bash
cd client-web/client-api
npm install
npm run dev
```

Abra `http://localhost:5174`.

## Conexão com backend

O modo padrão é `Local Docker via Vite proxy`. O API Explorer chama rotas como:

```txt
/__erp/crm/api/crm/leads
```

O `vite.config.ts` redireciona essa chamada para a porta local do serviço correspondente, evitando CORS durante desenvolvimento. Para testar contra URLs diretamente, use a aba `Ambientes` e altere o modo para `direct`.

## Catálogo de endpoints

O catálogo é gerado de `docs/contracts/http/*.openapi.yaml`:

```bash
npm run generate
```

O arquivo gerado fica em `src/generated/apiCatalog.ts`.
