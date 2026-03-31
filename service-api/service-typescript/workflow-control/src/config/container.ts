import { CreateWorkflowDefinition } from "../application/create-workflow-definition.js";
import { GetWorkflowDefinitionByKey } from "../application/get-workflow-definition-by-key.js";
import { ListWorkflowDefinitions } from "../application/list-workflow-definitions.js";
import { UpdateWorkflowDefinition } from "../application/update-workflow-definition.js";
import { UpdateWorkflowDefinitionStatus } from "../application/update-workflow-definition-status.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { loadConfig } from "./env.js";
import { InMemoryWorkflowDefinitionRepository } from "../infrastructure/in-memory-workflow-definition-repository.js";
import { PostgresWorkflowDefinitionRepository } from "../infrastructure/postgres-workflow-definition-repository.js";
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

export type ReadinessDependency = {
  name: string;
  status: string;
};

const config = loadConfig();
const repository = buildRepository();

export const services = {
  createWorkflowDefinition: new CreateWorkflowDefinition(repository),
  getWorkflowDefinitionByKey: new GetWorkflowDefinitionByKey(repository),
  listWorkflowDefinitions: new ListWorkflowDefinitions(repository),
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
