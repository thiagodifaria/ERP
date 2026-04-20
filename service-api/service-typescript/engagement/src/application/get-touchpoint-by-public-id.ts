import { ensurePublicId } from "../domain/touchpoint.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";

export class GetTouchpointByPublicId {
  constructor(private readonly repository: TouchpointRepository) {}

  async execute(publicId: string) {
    return this.repository.getByPublicId(ensurePublicId(publicId, "touchpoint_public_id_invalid"));
  }
}
