import { CreateWorkflowDefinition } from "../application/create-workflow-definition.js";
import { CreateWorkflowRun } from "../application/create-workflow-run.js";
import { GetCurrentWorkflowDefinitionVersion } from "../application/get-current-workflow-definition-version.js";
import { GetWorkflowDefinitionByKey } from "../application/get-workflow-definition-by-key.js";
import { GetWorkflowDefinitionVersionByNumber } from "../application/get-workflow-definition-version-by-number.js";
import { GetWorkflowRunByPublicId } from "../application/get-workflow-run-by-public-id.js";
import { GetWorkflowDefinitionVersionSummary } from "../application/get-workflow-definition-version-summary.js";
import { ListWorkflowDefinitionVersions } from "../application/list-workflow-definition-versions.js";
import { ListWorkflowDefinitions } from "../application/list-workflow-definitions.js";
import { ListWorkflowRuns } from "../application/list-workflow-runs.js";
import { PublishWorkflowDefinitionVersion } from "../application/publish-workflow-definition-version.js";
import { RestoreWorkflowDefinitionVersion } from "../application/restore-workflow-definition-version.js";
import { UpdateWorkflowDefinition } from "../application/update-workflow-definition.js";
import { UpdateWorkflowDefinitionStatus } from "../application/update-workflow-definition-status.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { WorkflowDefinitionVersionRepository } from "../domain/workflow-definition-version-repository.js";
import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";
import { loadConfig } from "./env.js";
import { InMemoryWorkflowDefinitionRepository } from "../infrastructure/in-memory-workflow-definition-repository.js";
import { InMemoryWorkflowDefinitionVersionRepository } from "../infrastructure/in-memory-workflow-definition-version-repository.js";
import { InMemoryWorkflowRunRepository } from "../infrastructure/in-memory-workflow-run-repository.js";
import { PostgresWorkflowDefinitionRepository } from "../infrastructure/postgres-workflow-definition-repository.js";
import { PostgresWorkflowDefinitionVersionRepository } from "../infrastructure/postgres-workflow-definition-version-repository.js";
import { PostgresWorkflowRunRepository } from "../infrastructure/postgres-workflow-run-repository.js";
import pg from "pg";

const { Pool } = pg;

function buildRepository(): WorkflowDefinitionRepository {
  if (config.repositoryDriver === "postgres") {
    return new PostgresWorkflowDefinitionRepository(
      new Pool({
        host: config.postgresHost,
        port: Number(config.postgresPort),
        database: config.postgresDatabase,
        user: config.postgresUser,
        password: config.postgresPassword,
        ssl: config.postgresSslMode === "disable" ? false : { rejectUnauthorized: false }
      }),
      config.bootstrapTenantSlug
    );
  }

  return new InMemoryWorkflowDefinitionRepository();
}

function buildVersionRepository(): WorkflowDefinitionVersionRepository {
  if (config.repositoryDriver === "postgres") {
    return new PostgresWorkflowDefinitionVersionRepository(
      new Pool({
        host: config.postgresHost,
        port: Number(config.postgresPort),
        database: config.postgresDatabase,
        user: config.postgresUser,
        password: config.postgresPassword,
        ssl: config.postgresSslMode === "disable" ? false : { rejectUnauthorized: false }
      }),
      config.bootstrapTenantSlug
    );
  }

  return new InMemoryWorkflowDefinitionVersionRepository();
}

function buildRunRepository(): WorkflowRunRepository {
  if (config.repositoryDriver === "postgres") {
    return new PostgresWorkflowRunRepository(
      new Pool({
        host: config.postgresHost,
        port: Number(config.postgresPort),
        database: config.postgresDatabase,
        user: config.postgresUser,
        password: config.postgresPassword,
        ssl: config.postgresSslMode === "disable" ? false : { rejectUnauthorized: false }
      }),
      config.bootstrapTenantSlug
    );
  }

  return new InMemoryWorkflowRunRepository();
}

export type ReadinessDependency = {
  name: string;
  status: string;
};

const config = loadConfig();
const repository = buildRepository();
const versionRepository = buildVersionRepository();
const runRepository = buildRunRepository();

export const repositories = {
  workflowDefinitions: repository,
  workflowDefinitionVersions: versionRepository,
  workflowRuns: runRepository
};

export const services = {
  createWorkflowDefinition: new CreateWorkflowDefinition(repository),
  createWorkflowRun: new CreateWorkflowRun(repository, versionRepository, runRepository),
  getCurrentWorkflowDefinitionVersion: new GetCurrentWorkflowDefinitionVersion(repository, versionRepository),
  getWorkflowDefinitionByKey: new GetWorkflowDefinitionByKey(repository),
  getWorkflowDefinitionVersionByNumber: new GetWorkflowDefinitionVersionByNumber(repository, versionRepository),
  getWorkflowRunByPublicId: new GetWorkflowRunByPublicId(runRepository),
  getWorkflowDefinitionVersionSummary: new GetWorkflowDefinitionVersionSummary(repository, versionRepository),
  listWorkflowDefinitionVersions: new ListWorkflowDefinitionVersions(repository, versionRepository),
  listWorkflowDefinitions: new ListWorkflowDefinitions(repository),
  listWorkflowRuns: new ListWorkflowRuns(runRepository),
  publishWorkflowDefinitionVersion: new PublishWorkflowDefinitionVersion(repository, versionRepository),
  restoreWorkflowDefinitionVersion: new RestoreWorkflowDefinitionVersion(repository, versionRepository),
  updateWorkflowDefinition: new UpdateWorkflowDefinition(repository),
  updateWorkflowDefinitionStatus: new UpdateWorkflowDefinitionStatus(repository)
};

export const runtime = {
  config,
  async readinessDependencies(): Promise<ReadinessDependency[]> {
    if (config.repositoryDriver !== "postgres") {
      return [
        { name: "router", status: "ready" },
        { name: "definitions-catalog", status: "ready" },
        { name: "workflow-runs", status: "ready" }
      ];
    }

    try {
      const postgresRepository = repository as PostgresWorkflowDefinitionRepository;
      await postgresRepository.list();
      await runRepository.list();

      return [
        { name: "router", status: "ready" },
        { name: "postgresql", status: "ready" },
        { name: "definitions-catalog", status: "ready" },
        { name: "workflow-runs", status: "ready" }
      ];
    } catch {
      return [
        { name: "router", status: "ready" },
        { name: "postgresql", status: "not_ready" },
        { name: "definitions-catalog", status: "not_ready" },
        { name: "workflow-runs", status: "not_ready" }
      ];
    }
  }
};
