export type CrmLeadPayload = {
  tenantSlug: string;
  name: string;
  email: string;
  source: string;
};

export type CrmLeadRecord = {
  publicId: string;
  name: string;
  email: string;
  source: string;
  status: string;
  ownerUserId: string;
};

export interface CrmGateway {
  createLead(payload: CrmLeadPayload): Promise<CrmLeadRecord>;
}

export class HttpCrmGateway implements CrmGateway {
  constructor(private readonly baseUrl: string) {}

  async createLead(payload: CrmLeadPayload): Promise<CrmLeadRecord> {
    const response = await fetch(`${this.baseUrl.replace(/\/$/, "")}/api/crm/leads`, {
      method: "POST",
      headers: {
        "content-type": "application/json"
      },
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      throw new Error("crm_gateway_error");
    }

    return (await response.json()) as CrmLeadRecord;
  }
}

export class InMemoryCrmGateway implements CrmGateway {
  async createLead(payload: CrmLeadPayload): Promise<CrmLeadRecord> {
    return {
      publicId: "00000000-0000-0000-0000-00000000e901",
      name: payload.name,
      email: payload.email,
      source: payload.source,
      status: "captured",
      ownerUserId: ""
    };
  }
}
