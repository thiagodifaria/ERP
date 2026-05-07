# supplier

Servico de fornecedores e ownership de dados de procurement.

## Escopo atual

- categorias de fornecedor por tenant
- diretório de fornecedores com dados cadastrais e financeiros básicos
- perfil de contas a pagar e termo padrão
- resumo operacional para consumo por `finance` e `analytics`

## Endpoints

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/supplier/capabilities`
- `GET /api/supplier/categories`
- `PUT /api/supplier/categories/{categoryKey}`
- `GET /api/supplier/suppliers`
- `GET /api/supplier/suppliers/export`
- `POST /api/supplier/suppliers`
- `POST /api/supplier/suppliers/bulk`
- `GET /api/supplier/suppliers/summary`
- `GET /api/supplier/suppliers/{publicId}`
- `PATCH /api/supplier/suppliers/{publicId}`
