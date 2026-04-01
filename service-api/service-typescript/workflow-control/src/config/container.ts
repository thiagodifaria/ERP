import { CreateWorkflowDefinition } from "../application/create-workflow-definition.js";
import { GetCurrentWorkflowDefinitionVersion } from "../application/get-current-workflow-definition-version.js";
import { GetWorkflowDefinitionByKey } from "../application/get-workflow-definition-by-key.js";
import { GetWorkflowDefinitionVersionByNumber } from "../application/get-workflow-definition-version-by-number.js";
import { ListWorkflowDefinitionVersions } from "../application/list-workflow-definition-versions.js";
import { ListWorkflowDefinitions } from "../application/list-workflow-definitions.js";
import { PublishWorkflowDefinitionVersion } from "../application/publish-workflow-definition-version.js";
import { UpdateWorkflowDefinition } from "../application/update-workflow-definition.js";
import { UpdateWorkflowDefinitionStatus } from "../application/update-workflow-definition-status.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { WorkflowDefinitionVersionRepository } from "../domain/workflow-definition-version-repository.js";
import { loadConfig } from "./env.js";
import { InMemoryWorkflowDefinitionRepository } from "../infrastructure/in-memory-workflow-definition-repository.js";
import { InMemoryWorkflowDefinitionVersionRepository } from "../infrastructure/in-memory-workflow-definition-version-repository.js";
import { PostgresWorkflowDefinitionRepository } from "../infrastructure/postgres-workflow-definition-repository.js";
import { PostgresWorkflowDefinitionVersionRepository } from "../infrastructure/postgres-workflow-definition-version-repository.js";
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

export type ReadinessDependency = {
  name: string;
  status: string;
};

const config = loadConfig();
const repository = buildRepository();
const versionRepository = buildVersionRepository();

export const services = {
  createWorkflowDefinition: new CreateWorkflowDefinition(repository),
  getCurrentWorkflowDefinitionVersion: new GetCurrentWorkflowDefinitionVersion(repository, versionRepository),
  getWorkflowDefinitionByKey: new GetWorkflowDefinitionByKey(repository),
  getWorkflowDefinitionVersionByNumber: new GetWorkflowDefinitionVersionByNumber(repository, versionRepository),
  listWorkflowDefinitionVersions: new ListWorkflowDefinitionVersions(repository, versionRepository),
  listWorkflowDefinitions: new ListWorkflowDefinitions(repository),
  publishWorkflowDefinitionVersion: new PublishWorkflowDefinitionVersion(repository, versionRepository),
  updateWorkflowDefinition: new UpdateWorkflowDefinition(repository),
  updateWorkflowDefinitionStatus: new UpdateWorkflowDefinitionStatus(repository)
};

export const runtime = {
  config,
  async readinessDependencies(): Promise<ReadinessDependency[]> {
    if (config.repositoryDriver !== "postgres") {
      return [
        { name: "router", status: "ready" },
        { name: "definitions-catalog", status: "ready" }
      ];
    }

    try {
      const postgresRepository = repository as PostgresWorkflowDefinitionRepository;
      await postgresRepository.list();

      return [
        { name: "router", status: "ready" },
        { name: "postgresql", status: "ready" },
        { name: "definitions-catalog", status: "ready" }
      ];
    } catch {
      return [
        { name: "router", status: "ready" },
        { name: "postgresql", status: "not_ready" },
        { name: "definitions-catalog", status: "not_ready" }
      ];
    }
  }
};
