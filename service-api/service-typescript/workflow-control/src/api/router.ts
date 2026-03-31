// This router starts small and grows with public workflow-control routes.
import { IncomingMessage, ServerResponse } from "node:http";
import { HealthResponse, ReadinessResponse } from "./dto/health.js";
import { CreateWorkflowDefinitionRequest } from "./dto/create-workflow-definition-request.js";
import { UpdateWorkflowDefinitionStatusRequest } from "./dto/update-workflow-definition-status-request.js";
import { services } from "../config/container.js";

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

function pathSegments(request: IncomingMessage): string[] {
  const pathname = request.url?.split("?")[0] ?? "/";
  return pathname.split("/").filter((segment) => segment.length > 0);
}

async function readJson<T>(request: IncomingMessage): Promise<T> {
  const chunks: Buffer[] = [];

  for await (const chunk of request) {
    chunks.push(Buffer.from(chunk));
  }

  const rawBody = Buffer.concat(chunks).toString("utf8");

  if (rawBody.length === 0) {
    throw new Error("invalid_json");
  }

  return JSON.parse(rawBody) as T;
}

export async function route(request: IncomingMessage, response: ServerResponse): Promise<void> {
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

  if (request.method === "GET" && request.url === "/api/workflow-control/definitions") {
    json(response, 200, services.listWorkflowDefinitions.execute());
    return;
  }

  const segments = pathSegments(request);
  if (
    request.method === "GET" &&
    segments.length === 4 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "definitions"
  ) {
    const definition = services.getWorkflowDefinitionByKey.execute(segments[3]);

    if (definition === null) {
      json(response, 404, {
        code: "workflow_definition_not_found",
        message: "Workflow definition was not found."
      });
      return;
    }

    json(response, 200, definition);
    return;
  }

  if (
    request.method === "PATCH" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "definitions" &&
    segments[4] === "status"
  ) {
    try {
      const payload = await readJson<UpdateWorkflowDefinitionStatusRequest>(request);
      const updated = services.updateWorkflowDefinitionStatus.execute(segments[3], payload.status);

      json(response, 200, updated);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "invalid_json") {
        json(response, 400, {
          code: "invalid_json",
          message: "Request body must be a valid JSON object."
        });
        return;
      }

      if (code === "workflow_definition_not_found") {
        json(response, 404, {
          code,
          message: "Workflow definition was not found."
        });
        return;
      }

      if (code === "workflow_definition_status_invalid") {
        json(response, 400, {
          code,
          message: "Workflow definition status is invalid."
        });
        return;
      }

      json(response, 500, {
        code: "unexpected_error",
        message: "Unexpected error."
      });
      return;
    }
  }

  if (request.method === "POST" && request.url === "/api/workflow-control/definitions") {
    try {
      const payload = await readJson<CreateWorkflowDefinitionRequest>(request);
      const created = services.createWorkflowDefinition.execute(payload);

      json(response, 201, created);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "invalid_json") {
        json(response, 400, {
          code: "invalid_json",
          message: "Request body must be a valid JSON object."
        });
        return;
      }

      if (code === "workflow_definition_key_conflict") {
        json(response, 409, {
          code,
          message: "Workflow definition key already exists."
        });
        return;
      }

      if (code === "workflow_definition_key_required" || code === "workflow_definition_name_required" || code === "workflow_definition_trigger_required") {
        json(response, 400, {
          code,
          message: "Workflow definition payload is invalid."
        });
        return;
      }

      json(response, 500, {
        code: "unexpected_error",
        message: "Unexpected error."
      });
      return;
    }
  }

  json(response, 404, {
    code: "route_not_found",
    message: "Route was not found."
  });
}
