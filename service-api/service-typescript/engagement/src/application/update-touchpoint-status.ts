import { TouchpointRepository } from "../domain/touchpoint-repository.js";
import { ensureOptionalPublicId, ensurePublicId, ensureTouchpointStatus } from "../domain/touchpoint.js";

export class UpdateTouchpointStatus {
  constructor(private readonly repository: TouchpointRepository) {}

  async execute(publicId: string, status: string, lastWorkflowRunPublicId?: string | null) {
    return this.repository.updateStatus(
      ensurePublicId(publicId, "touchpoint_public_id_invalid"),
      ensureTouchpointStatus(status),
      ensureOptionalPublicId(lastWorkflowRunPublicId)
    );
  }
}
