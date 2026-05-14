# Service Starter Checklist

Use este checklist antes de criar ou promover um servico HTTP no ERP.

- [ ] Runtime define `ERP_ENV`, `ERP_AUTH_ENFORCEMENT`, `ERP_JWT_HS256_SECRET`, `ERP_INTERNAL_SERVICE_TOKEN`, `ERP_OPENFGA_ENFORCEMENT` e `OPENFGA_STORE_ID`.
- [ ] Auth middleware valida JWT/service account e chama OpenFGA quando habilitado.
- [ ] Mutacoes exigem correlation id.
- [ ] Tenant explicito em rotas tenant-aware.
- [ ] Erros publicos usam `code` e `message`.
- [ ] Health endpoints implementados.
- [ ] OpenAPI registra security, erros, idempotencia e rotas internas quando houver.
- [ ] Eventos usam envelope padrao e schema versionado.
- [ ] Logs nao imprimem senha, token, access link, documento fiscal ou provider secret.
- [ ] Testes cobrem unauthorized, forbidden, tenant mismatch, idempotencia e contrato publico.
