import { TemplateRepository } from "../domain/template-repository.js";
import { ensurePublicId } from "../domain/touchpoint.js";

export class GetTemplateByPublicId {
  constructor(private readonly repository: TemplateRepository) {}

  async execute(publicId: string) {
    return this.repository.getByPublicId(ensurePublicId(publicId, "template_public_id_invalid"));
  }
}
