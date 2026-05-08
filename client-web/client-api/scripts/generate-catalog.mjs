import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(__dirname, "../../..");
const contractsDir = path.join(root, "docs/contracts/http");
const eventsDir = path.join(root, "docs/contracts/events");
const outputPath = path.join(__dirname, "../src/generated/apiCatalog.ts");

const methodNames = new Set(["get", "post", "put", "patch", "delete"]);

function serviceFromFile(fileName) {
  return fileName.replace(".openapi.yaml", "");
}

function titleCase(value) {
  return value
    .split(/[-_\s]+/)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

function parseOpenApi(fileName) {
  const service = serviceFromFile(fileName);
  const source = fs.readFileSync(path.join(contractsDir, fileName), "utf8");
  const lines = source.split(/\r?\n/);
  const title = lines.find((line) => line.trim().startsWith("title:"))?.split(":").slice(1).join(":").trim() ?? titleCase(service);
  const version = lines.find((line) => line.trim().startsWith("version:"))?.split(":").slice(1).join(":").trim() ?? "0.1.0";
  const description = lines.find((line) => line.trim().startsWith("description:"))?.split(":").slice(1).join(":").trim() ?? "";
  const endpoints = [];
  let currentPath = "";
  let currentMethod = "";

  for (const line of lines) {
    const pathMatch = line.match(/^  (\/[^:]+):\s*$/);
    if (pathMatch) {
      currentPath = pathMatch[1];
      currentMethod = "";
      continue;
    }

    const methodMatch = line.match(/^    (get|post|put|patch|delete):\s*$/);
    if (methodMatch && currentPath) {
      currentMethod = methodMatch[1].toUpperCase();
      endpoints.push({
        id: `${service}:${currentMethod}:${currentPath}`,
        service,
        method: currentMethod,
        path: currentPath,
        tag: titleCase(service),
        description: `${currentMethod} ${currentPath}`,
        summary: "",
        source: `docs/contracts/http/${fileName}`,
        hasBody: !["GET", "DELETE"].includes(currentMethod),
        pathParams: [...currentPath.matchAll(/\{([^}]+)\}/g)].map((match) => match[1])
      });
      continue;
    }

    if (currentMethod && line.match(/^      summary:/)) {
      const endpoint = endpoints[endpoints.length - 1];
      endpoint.summary = line.split(":").slice(1).join(":").trim();
      endpoint.description = endpoint.summary || endpoint.description;
    }
  }

  return {
    slug: service,
    name: title.replace(/^ERP\s+/i, "").replace(/\s+API$/i, ""),
    title,
    version,
    description,
    contractFile: `docs/contracts/http/${fileName}`,
    endpointCount: endpoints.length,
    endpoints
  };
}

function healthEndpoints(service) {
  return ["live", "ready", "details"].map((probe) => ({
    id: `${service.slug}:GET:/health/${probe}`,
    service: service.slug,
    method: "GET",
    path: `/health/${probe}`,
    tag: "Health",
    description: `Health ${probe} do serviço ${service.slug}.`,
    summary: `Health ${probe}`,
    source: "runtime",
    hasBody: false,
    pathParams: []
  }));
}

const contractFiles = fs
  .readdirSync(contractsDir)
  .filter((fileName) => fileName.endsWith(".openapi.yaml"))
  .sort();

const services = contractFiles.map(parseOpenApi);
const endpoints = services.flatMap((service) => {
  const existingIds = new Set(service.endpoints.map((endpoint) => endpoint.id));
  const probes = healthEndpoints(service).filter((endpoint) => !existingIds.has(endpoint.id));
  return [...probes, ...service.endpoints];
});
const eventSchemas = fs
  .readdirSync(eventsDir)
  .filter((fileName) => fileName.endsWith(".json"))
  .sort()
  .map((fileName) => ({
    fileName,
    name: fileName.replace(".schema.json", ""),
    source: `docs/contracts/events/${fileName}`
  }));

const content = `/* eslint-disable */\n// Gerado por scripts/generate-catalog.mjs. Nao edite manualmente.\n\nexport type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";\n\nexport type ServiceContract = {\n  slug: string;\n  name: string;\n  title: string;\n  version: string;\n  description: string;\n  contractFile: string;\n  endpointCount: number;\n};\n\nexport type EndpointContract = {\n  id: string;\n  service: string;\n  method: HttpMethod;\n  path: string;\n  tag: string;\n  description: string;\n  summary: string;\n  source: string;\n  hasBody: boolean;\n  pathParams: string[];\n};\n\nexport type EventSchemaContract = {\n  fileName: string;\n  name: string;\n  source: string;\n};\n\nexport const services: ServiceContract[] = ${JSON.stringify(
  services.map(({ endpoints: _endpoints, ...service }) => service),
  null,
  2
)};\n\nexport const endpoints: EndpointContract[] = ${JSON.stringify(endpoints, null, 2)};\n\nexport const eventSchemas: EventSchemaContract[] = ${JSON.stringify(eventSchemas, null, 2)};\n`;

fs.mkdirSync(path.dirname(outputPath), { recursive: true });
fs.writeFileSync(outputPath, content);
console.log(`generated ${endpoints.length} endpoints from ${services.length} services`);
