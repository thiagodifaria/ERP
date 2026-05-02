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
        critical: false,
        mode: this.config.resendApiKey.trim().length > 0 ? "configured" : "fallback",
        credentialKey: "ENGAGEMENT_RESEND_API_KEY",
        fallbackViable: true,
        status: this.config.resendApiKey.trim().length > 0 ? "ready" : "fallback",
        notes: this.config.resendApiKey.trim().length > 0 ? ["Provider ativo para email transacional."] : ["Opera em modo fallback para ambiente local e testes."],
        supportsInbound: false,
        supportsOutbound: true,
        supportsTracking: true,
        supportsCallbacks: true
      },
      {
        provider: "whatsapp_cloud",
        scope: "messaging",
        configured: this.config.whatsappAccessToken.trim().length > 0,
        critical: false,
        mode: this.config.whatsappAccessToken.trim().length > 0 ? "configured" : "fallback",
        credentialKey: "ENGAGEMENT_WHATSAPP_ACCESS_TOKEN",
        fallbackViable: true,
        status: this.config.whatsappAccessToken.trim().length > 0 ? "ready" : "fallback",
        notes: this.config.whatsappAccessToken.trim().length > 0 ? ["Canal pronto para dispatch e callback."] : ["Fallback local preserva o fluxo sem depender do provider real."],
        supportsInbound: true,
        supportsOutbound: true,
        supportsTracking: true,
        supportsCallbacks: true
      },
      {
        provider: "telegram_bot",
        scope: "messaging",
        configured: this.config.telegramBotToken.trim().length > 0,
        critical: false,
        mode: this.config.telegramBotToken.trim().length > 0 ? "configured" : "fallback",
        credentialKey: "ENGAGEMENT_TELEGRAM_BOT_TOKEN",
        fallbackViable: true,
        status: this.config.telegramBotToken.trim().length > 0 ? "ready" : "fallback",
        notes: this.config.telegramBotToken.trim().length > 0 ? ["Bot pronto para entrada e notificacao externa."] : ["Fallback local cobre exploracao do canal sem credencial real."],
        supportsInbound: true,
        supportsOutbound: true,
        supportsTracking: true,
        supportsCallbacks: true
      },
      {
        provider: "meta_ads",
        scope: "ads",
        configured: this.config.metaAdsAccessToken.trim().length > 0,
        critical: false,
        mode: this.config.metaAdsAccessToken.trim().length > 0 ? "configured" : "fallback",
        credentialKey: "ENGAGEMENT_META_ADS_ACCESS_TOKEN",
        fallbackViable: true,
        status: this.config.metaAdsAccessToken.trim().length > 0 ? "ready" : "fallback",
        notes: this.config.metaAdsAccessToken.trim().length > 0 ? ["Lead intake pronto para Meta Ads."] : ["Fallback local simula inbound de campanha para smoke e desenvolvimento."],
        supportsInbound: true,
        supportsOutbound: false,
        supportsTracking: true,
        supportsCallbacks: true
      },
      {
        provider: "manual",
        scope: "manual",
        configured: true,
        critical: false,
        mode: "manual",
        credentialKey: null,
        fallbackViable: true,
        status: "manual",
        notes: ["Modo manual sempre disponivel para operacao local e contingencia controlada."],
        supportsInbound: true,
        supportsOutbound: true,
        supportsTracking: true,
        supportsCallbacks: true
      }
    ];
  }
}
