import { TouchpointFilters } from "../domain/touchpoint.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";

export type ConversationThread = {
  threadPublicId: string;
  tenantSlug: string;
  participantKind: string;
  participantPublicId: string;
  businessEntityType: string | null;
  businessEntityPublicId: string | null;
  channel: string;
  status: string;
  touchpoints: number;
  lastTouchpointAt: string;
};

export class ListConversations {
  constructor(private readonly repository: TouchpointRepository) {}

  async execute(filters: TouchpointFilters = {}): Promise<ConversationThread[]> {
    const touchpoints = await this.repository.list(filters);
    const grouped = new Map<string, ConversationThread>();

    for (const touchpoint of touchpoints) {
      const key = touchpoint.threadPublicId;
      const current = grouped.get(key);
      if (current === undefined) {
        grouped.set(key, {
          threadPublicId: touchpoint.threadPublicId,
          tenantSlug: touchpoint.tenantSlug,
          participantKind: touchpoint.participantKind,
          participantPublicId: touchpoint.participantPublicId,
          businessEntityType: touchpoint.businessEntityType,
          businessEntityPublicId: touchpoint.businessEntityPublicId,
          channel: touchpoint.channel,
          status: touchpoint.status,
          touchpoints: 1,
          lastTouchpointAt: touchpoint.updatedAt,
        });
        continue;
      }

      current.touchpoints += 1;
      if (touchpoint.updatedAt > current.lastTouchpointAt) {
        current.lastTouchpointAt = touchpoint.updatedAt;
        current.status = touchpoint.status;
      }
    }

    return Array.from(grouped.values()).sort((left, right) => right.lastTouchpointAt.localeCompare(left.lastTouchpointAt));
  }
}
