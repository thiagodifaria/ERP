import { ProviderCapability } from "../domain/provider-event.js";

export type ProviderCapabilityConfig = {
  resendApiKey: string;
  whatsappAccessToken: string;
  telegramBotToken: string;
  metaAdsAccessToken: string;
};

export class ListProviderCapabilities {
  constructor(private readonly config: ProviderCapabilityConfig) {}

  async execute(): Promise<ProviderCapability[]> {
    return [
      {
        provider: "resend",
        scope: "email",
        configured: this.config.resendApiKey.trim().length > 0,
        mode: this.config.resendApiKey.trim().length > 0 ? "configured" : "fallback",
        supportsInbound: false,
        supportsOutbound: true,
        supportsTracking: true
      },
      {
        provider: "whatsapp_cloud",
        scope: "messaging",
        configured: this.config.whatsappAccessToken.trim().length > 0,
        mode: this.config.whatsappAccessToken.trim().length > 0 ? "configured" : "fallback",
        supportsInbound: true,
        supportsOutbound: true,
        supportsTracking: true
      },
      {
        provider: "telegram_bot",
        scope: "messaging",
        configured: this.config.telegramBotToken.trim().length > 0,
        mode: this.config.telegramBotToken.trim().length > 0 ? "configured" : "fallback",
        supportsInbound: true,
        supportsOutbound: true,
        supportsTracking: true
      },
      {
        provider: "meta_ads",
        scope: "ads",
        configured: this.config.metaAdsAccessToken.trim().length > 0,
        mode: this.config.metaAdsAccessToken.trim().length > 0 ? "configured" : "fallback",
        supportsInbound: true,
        supportsOutbound: false,
        supportsTracking: true
      },
      {
        provider: "manual",
        scope: "manual",
        configured: true,
        mode: "manual",
        supportsInbound: true,
        supportsOutbound: true,
        supportsTracking: true
      }
    ];
  }
}
