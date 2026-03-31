// This router starts small and grows with public workflow-control routes.
import { IncomingMessage, ServerResponse } from "node:http";
import { HealthResponse, ReadinessResponse } from "./dto/health.js";

function json(response: ServerResponse, statusCode: number, body: unknown): void {
  response.writeHead(statusCode, { "content-type": "application/json" });
  response.end(JSON.stringify(body));
}

function live(): HealthResponse {
  return {
    service: "workflow-control",
    status: "live"
  };
}

function ready(): HealthResponse {
  return {
    service: "workflow-control",
    status: "ready"
  };
}

function details(): ReadinessResponse {
  return {
    service: "workflow-control",
    status: "ready",
    dependencies: [
      { name: "router", status: "ready" },
      { name: "definitions-catalog", status: "pending-runtime-wiring" }
    ]
  };
}

export function route(request: IncomingMessage, response: ServerResponse): void {
  if (request.url === "/health/live") {
    json(response, 200, live());
    return;
  }

  if (request.url === "/health/ready") {
    json(response, 200, ready());
    return;
  }

  if (request.url === "/health/details") {
    json(response, 200, details());
    return;
  }

  json(response, 404, {
    code: "route_not_found",
    message: "Route was not found."
  });
}
