# ERP - Control Console

Console tecnico para explorar a API do ERP, ler documentacao versionada, validar contratos e testar rotas em um unico lugar.

## Rodar localmente

```bash
cd client-web/client-api
npm install
npm run dev
```

Abra `http://localhost:5174`.

## Conexao com backend

O modo padrao e `Local Docker via Vite proxy`. O API Explorer chama rotas como:

```txt
/__erp/crm/api/crm/leads
```

O `vite.config.ts` redireciona essa chamada para a porta local do servico correspondente, evitando CORS durante desenvolvimento. Para testar contra URLs diretamente, use a aba `Ambientes` e altere o modo para `direct`.

## Catalogo de endpoints

O catalogo e gerado de `docs/contracts/http/*.openapi.yaml`:

```bash
npm run generate
```

O arquivo gerado fica em `src/generated/apiCatalog.ts`.
