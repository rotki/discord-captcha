import { createRui } from '@rotki/ui-library';
import icons from 'virtual:rotki-icons';
import { createApp } from 'vue';
import { createI18n } from 'vue-i18n';
import { createRouter, createWebHistory } from 'vue-router';
import App from './App.vue';
import en from './locales/en.json';
import '@rotki/ui-library/style.css';
import '@fontsource/roboto/latin.css';
import './assets/css/tailwind.css';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: async () => import('./pages/index.vue') },
  ],
});

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en },
});

const rui = createRui({
  theme: {
    icons,
    mode: 'light',
  },
});

const app = createApp(App);
app.use(router);
app.use(i18n);
app.use(rui);
app.mount('#app');
