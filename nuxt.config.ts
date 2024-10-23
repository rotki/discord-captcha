// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  app: {
    head: {
      htmlAttrs: {
        lang: 'en',
      },
      link: [
        {
          href: '/apple-touch-icon.png',
          rel: 'apple-touch-icon',
          sizes: '180x180',
        },
        { href: '/favicon.ico', rel: 'icon', type: 'image/x-icon' },
        {
          href: '/favicon-32x32.png',
          rel: 'icon',
          sizes: '32x32',
          type: 'image/png',
        },
        {
          href: '/favicon-16x16.png',
          rel: 'icon',
          sizes: '16x16',
          type: 'image/png',
        },
        {
          crossorigin: 'use-credentials',
          href: '/site.webmanifest',
          rel: 'manifest',
        },
        {
          color: '#5bbad5',
          href: '/safari-pinned-tab.svg',
          rel: 'mask-icon',
        },
      ],
      meta: [
        { charset: 'utf-8' },
        { content: 'width=device-width, initial-scale=1', name: 'viewport' },
        { content: '#00aba9', name: 'msapplication-TileColor' },
        { content: '#ffffff', name: 'theme-color' },
      ],
      title: 'discord.rotki.com',
    },
  },

  compatibilityDate: '2024-08-13',

  devtools: { enabled: true },

  i18n: {
    defaultLocale: 'en',
    langDir: 'locales',
    lazy: true,
    locales: [{ code: 'en', file: 'en.json', language: 'en-US' }],
    strategy: 'no_prefix',
    vueI18n: './i18n.config.ts',
  },

  modules: [
    '@nuxt/devtools',
    '@nuxtjs/i18n',
    '@nuxtjs/tailwindcss',
    '@nuxtjs/robots',
    '@vueuse/nuxt',
    '@nuxt/test-utils/module',
    './modules/ui-library/module.ts',
  ],

  robots: {
    groups: [
      {
        disallow: ['/health'],
        userAgent: '*',
      },
    ],
  },

  runtimeConfig: {
    discord: {
      appId: '',
      channelId: 0,
      guildId: 0,
      publicKey: '',
      roleId: 0,
      token: '',
    },
    public: {
      logLevel: 3, // defaults at info level
      recaptcha: {
        siteKey: '',
      },
      siteUrl: '',
    },
    recaptchaSecret: '',
    redis: {
      host: '',
      password: '',
    },
  },

  ssr: true,

  vite: {
    css: {
      preprocessorOptions: {
        scss: {
          api: 'modern',
        },
      },
    },
    optimizeDeps: {
      exclude: ['fsevents', 'zlib-sync'],
    },
  },
});
