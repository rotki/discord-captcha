// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  ssr: true,
  modules: [
    '@nuxt/devtools',
    '@nuxtjs/i18n',
    '@nuxtjs/tailwindcss',
    '@vueuse/nuxt',
    './modules/ui-library/module.ts',
  ],
  i18n: {
    locales: [{ code: 'en', iso: 'en-US', file: 'en.json' }],
    defaultLocale: 'en',
    strategy: 'no_prefix',
    langDir: 'locales',
    lazy: true,
    vueI18n: './i18n.config.ts',
  },
  runtimeConfig: {
    recaptchaSecret: '',
    discord: {
      appId: '',
      token: '',
      publicKey: '',
      guildId: 0,
      channelId: 0,
    },
    public: {
      siteUrl: '',
      recaptcha: {
        siteKey: '',
      },
    },
  },
});
