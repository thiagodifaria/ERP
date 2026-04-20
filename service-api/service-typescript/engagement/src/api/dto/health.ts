export type HealthResponse = {
  service: string;
  status: string;
};

export type ReadinessResponse = {
  service: string;
  status: string;
  dependencies: Array<{ name: string; status: string }>;
};
