export type HealthResponse = {
  service: string;
  status: string;
};

export type DependencyHealthResponse = {
  name: string;
  status: string;
};

export type ReadinessResponse = HealthResponse & {
  dependencies: DependencyHealthResponse[];
};
