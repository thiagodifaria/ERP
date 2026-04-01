// This router starts small and grows with public workflow-control routes.
import { IncomingMessage, ServerResponse } from "node:http";
import { CreateWorkflowRunEventRequest } from "./dto/create-workflow-run-event-request.js";
import { HealthResponse, ReadinessResponse } from "./dto/health.js";
import { CreateWorkflowDefinitionRequest } from "./dto/create-workflow-definition-request.js";
import { CreateWorkflowRunRequest } from "./dto/create-workflow-run-request.js";
import { UpdateWorkflowDefinitionRequest } from "./dto/update-workflow-definition-request.js";
import { UpdateWorkflowDefinitionStatusRequest } from "./dto/update-workflow-definition-status-request.js";
import { runtime, services } from "../config/container.js";

function json(response: ServerResponse, statusCode: number, body: unknown): void {
  response.writeHead(statusCode, { "content-type": "application/json" });
  response.end(JSON.stringify(body));
}

function live(): HealthResponse {
  return {
    service: runtime.config.serviceName,
    status: "live"
  };
}

function ready(): HealthResponse {
  return {
    service: runtime.config.serviceName,
    status: "ready"
  };
}

async function details(): Promise<ReadinessResponse> {
  return {
    service: runtime.config.serviceName,
    status: "ready",
    dependencies: await runtime.readinessDependencies()
  };
}

function pathSegments(request: IncomingMessage): string[] {
  const pathname = request.url?.split("?")[0] ?? "/";
  return pathname.split("/").filter((segment) => segment.length > 0);
}

function searchParams(request: IncomingMessage): URLSearchParams {
  return new URL(request.url ?? "/", "http://workflow-control.local").searchParams;
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
    json(response, 200, await details());
    return;
  }

  if (request.method === "GET" && request.url === "/api/workflow-control/definitions") {
    json(response, 200, await services.listWorkflowDefinitions.execute());
    return;
  }

  if (request.method === "GET" && request.url === "/api/workflow-control/runs") {
    json(response, 200, await services.listWorkflowRuns.execute());
    return;
  }

  if (request.method === "GET" && request.url?.startsWith("/api/workflow-control/runs?")) {
    try {
      const params = searchParams(request);
      const workflowRuns = await services.listWorkflowRuns.execute({
        workflowDefinitionKey: params.get("workflowDefinitionKey") ?? undefined,
        status: params.get("status") ?? undefined,
        subjectType: params.get("subjectType") ?? undefined,
        initiatedBy: params.get("initiatedBy") ?? undefined
      });
      json(response, 200, workflowRuns);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_run_status_invalid") {
        json(response, 400, {
          code,
          message: "Workflow run filter is invalid."
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

  if (request.method === "GET" && request.url === "/api/workflow-control/runs/summary") {
    json(response, 200, await services.getWorkflowRunSummary.execute());
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
    const definition = await services.getWorkflowDefinitionByKey.execute(segments[3]);

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
    request.method === "GET" &&
    segments.length === 4 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "runs"
  ) {
    const workflowRun = await services.getWorkflowRunByPublicId.execute(segments[3]);

    if (workflowRun === null) {
      json(response, 404, {
        code: "workflow_run_not_found",
        message: "Workflow run was not found."
      });
      return;
    }

    json(response, 200, workflowRun);
    return;
  }

  if (
    request.method === "GET" &&
    segments.length === 6 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "runs" &&
    segments[4] === "events" &&
    segments[5] === "summary"
  ) {
    try {
      const workflowRunEventSummary = await services.getWorkflowRunEventSummary.execute(segments[3]);
      json(response, 200, workflowRunEventSummary);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_run_not_found") {
        json(response, 404, {
          code,
          message: "Workflow run was not found."
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

  if (
    request.method === "GET" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "runs" &&
    segments[4] === "events"
  ) {
    try {
      const params = searchParams(request);
      const workflowRunEvents = await services.listWorkflowRunEvents.execute(segments[3], {
        category: params.get("category") ?? undefined,
        createdBy: params.get("createdBy") ?? undefined
      });
      json(response, 200, workflowRunEvents);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_run_not_found") {
        json(response, 404, {
          code,
          message: "Workflow run was not found."
        });
        return;
      }

      if (code === "workflow_run_event_category_invalid") {
        json(response, 400, {
          code,
          message: "Workflow run event filter is invalid."
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

  if (
    request.method === "POST" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "runs" &&
    segments[4] === "events"
  ) {
    try {
      const payload = await readJson<CreateWorkflowRunEventRequest>(request);
      const workflowRunEvent = await services.createWorkflowRunNote.execute({
        workflowRunPublicId: segments[3],
        body: payload.body,
        createdBy: payload.createdBy
      });

      json(response, 201, workflowRunEvent);
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

      if (code === "workflow_run_not_found") {
        json(response, 404, {
          code,
          message: "Workflow run was not found."
        });
        return;
      }

      if (code === "workflow_run_event_body_required" || code === "workflow_run_event_created_by_required") {
        json(response, 400, {
          code,
          message: "Workflow run event payload is invalid."
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

  if (
    request.method === "POST" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "runs" &&
    segments[4] === "cancel"
  ) {
    try {
      const workflowRun = await services.cancelWorkflowRun.execute(segments[3]);
      json(response, 200, workflowRun);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_run_not_found") {
        json(response, 404, {
          code,
          message: "Workflow run was not found."
        });
        return;
      }

      if (code === "workflow_run_transition_invalid") {
        json(response, 409, {
          code,
          message: "Workflow run transition is invalid."
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

  if (
    request.method === "POST" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "runs" &&
    segments[4] === "fail"
  ) {
    try {
      const workflowRun = await services.failWorkflowRun.execute(segments[3]);
      json(response, 200, workflowRun);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_run_not_found") {
        json(response, 404, {
          code,
          message: "Workflow run was not found."
        });
        return;
      }

      if (code === "workflow_run_transition_invalid") {
        json(response, 409, {
          code,
          message: "Workflow run transition is invalid."
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

  if (
    request.method === "POST" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "runs" &&
    segments[4] === "complete"
  ) {
    try {
      const workflowRun = await services.completeWorkflowRun.execute(segments[3]);
      json(response, 200, workflowRun);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_run_not_found") {
        json(response, 404, {
          code,
          message: "Workflow run was not found."
        });
        return;
      }

      if (code === "workflow_run_transition_invalid") {
        json(response, 409, {
          code,
          message: "Workflow run transition is invalid."
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

  if (
    request.method === "POST" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "runs" &&
    segments[4] === "start"
  ) {
    try {
      const workflowRun = await services.startWorkflowRun.execute(segments[3]);
      json(response, 200, workflowRun);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_run_not_found") {
        json(response, 404, {
          code,
          message: "Workflow run was not found."
        });
        return;
      }

      if (code === "workflow_run_transition_invalid") {
        json(response, 409, {
          code,
          message: "Workflow run transition is invalid."
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

  if (
    request.method === "POST" &&
    segments.length === 7 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "definitions" &&
    segments[4] === "versions" &&
    segments[6] === "restore"
  ) {
    try {
      const versionNumber = Number(segments[5]);

      if (!Number.isInteger(versionNumber) || versionNumber <= 0) {
        json(response, 400, {
          code: "workflow_definition_version_number_invalid",
          message: "Workflow definition version number is invalid."
        });
        return;
      }

      const restoredDefinition = await services.restoreWorkflowDefinitionVersion.execute(segments[3], versionNumber);
      json(response, 200, restoredDefinition);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_definition_not_found" || code === "workflow_definition_version_not_found") {
        json(response, 404, {
          code,
          message: code === "workflow_definition_not_found"
            ? "Workflow definition was not found."
            : "Workflow definition version was not found."
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

  if (
    request.method === "GET" &&
    segments.length === 6 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "definitions" &&
    segments[4] === "versions" &&
    segments[5] === "summary"
  ) {
    try {
      const summary = await services.getWorkflowDefinitionVersionSummary.execute(segments[3]);
      json(response, 200, summary);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_definition_not_found") {
        json(response, 404, {
          code,
          message: "Workflow definition was not found."
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

  if (
    request.method === "GET" &&
    segments.length === 6 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "definitions" &&
    segments[4] === "versions" &&
    segments[5] !== "current"
  ) {
    try {
      const versionNumber = Number(segments[5]);

      if (!Number.isInteger(versionNumber) || versionNumber <= 0) {
        json(response, 400, {
          code: "workflow_definition_version_number_invalid",
          message: "Workflow definition version number is invalid."
        });
        return;
      }

      const version = await services.getWorkflowDefinitionVersionByNumber.execute(segments[3], versionNumber);

      if (version === null) {
        json(response, 404, {
          code: "workflow_definition_version_not_found",
          message: "Workflow definition version was not found."
        });
        return;
      }

      json(response, 200, version);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_definition_not_found") {
        json(response, 404, {
          code,
          message: "Workflow definition was not found."
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

  if (
    request.method === "GET" &&
    segments.length === 6 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "definitions" &&
    segments[4] === "versions" &&
    segments[5] === "current"
  ) {
    try {
      const currentVersion = await services.getCurrentWorkflowDefinitionVersion.execute(segments[3]);

      if (currentVersion === null) {
        json(response, 404, {
          code: "workflow_definition_version_not_found",
          message: "Workflow definition version was not found."
        });
        return;
      }

      json(response, 200, currentVersion);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_definition_not_found") {
        json(response, 404, {
          code,
          message: "Workflow definition was not found."
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

  if (
    request.method === "POST" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "definitions" &&
    segments[4] === "versions"
  ) {
    try {
      const version = await services.publishWorkflowDefinitionVersion.execute(segments[3]);
      json(response, 201, version);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_definition_not_found") {
        json(response, 404, {
          code,
          message: "Workflow definition was not found."
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

  if (
    request.method === "GET" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "definitions" &&
    segments[4] === "versions"
  ) {
    try {
      const versions = await services.listWorkflowDefinitionVersions.execute(segments[3]);
      json(response, 200, versions);
      return;
    } catch (error) {
      const code = error instanceof Error ? error.message : "unexpected_error";

      if (code === "workflow_definition_not_found") {
        json(response, 404, {
          code,
          message: "Workflow definition was not found."
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

  if (
    request.method === "PATCH" &&
    segments.length === 4 &&
    segments[0] === "api" &&
    segments[1] === "workflow-control" &&
    segments[2] === "definitions"
  ) {
    try {
      const payload = await readJson<UpdateWorkflowDefinitionRequest>(request);
      const updated = await services.updateWorkflowDefinition.execute(segments[3], payload);

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

      if (
        code === "workflow_definition_update_required" ||
        code === "workflow_definition_name_required" ||
        code === "workflow_definition_trigger_required"
      ) {
        json(response, 400, {
          code,
          message: "Workflow definition payload is invalid."
        });
        return;
      }

      if (code === "workflow_definition_tenant_not_found") {
        json(response, 500, {
          code,
          message: "Workflow bootstrap tenant was not found."
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
      const updated = await services.updateWorkflowDefinitionStatus.execute(segments[3], payload.status);

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

      if (code === "workflow_definition_tenant_not_found") {
        json(response, 500, {
          code,
          message: "Workflow bootstrap tenant was not found."
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
      const created = await services.createWorkflowDefinition.execute(payload);

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

      if (code === "workflow_definition_tenant_not_found") {
        json(response, 500, {
          code,
          message: "Workflow bootstrap tenant was not found."
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

  if (request.method === "POST" && request.url === "/api/workflow-control/runs") {
    try {
      const payload = await readJson<CreateWorkflowRunRequest>(request);
      const created = await services.createWorkflowRun.execute(payload);

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

      if (code === "workflow_definition_not_found" || code === "workflow_definition_version_not_found") {
        json(response, 404, {
          code,
          message: code === "workflow_definition_not_found"
            ? "Workflow definition was not found."
            : "Workflow definition version was not found."
        });
        return;
      }

      if (
        code === "workflow_run_definition_key_required" ||
        code === "workflow_run_subject_type_required" ||
        code === "workflow_run_subject_public_id_required" ||
        code === "workflow_run_initiated_by_required"
      ) {
        json(response, 400, {
          code,
          message: "Workflow run payload is invalid."
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
