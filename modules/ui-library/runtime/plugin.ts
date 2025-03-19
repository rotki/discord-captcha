import {
  createRui,
  LuBookText,
  LuBug,
  LuGithub,
  LuGlobe,
  LuMail,
  LuMessageCircleQuestion,
  LuMessagesSquare,
  LuNewspaper,
  LuXTwitter,
} from '@rotki/ui-library';
import '@fontsource/roboto/latin.css';

export default defineNuxtPlugin((nuxtApp) => {
  const rui = createRui({
    theme: {
      icons: [
        LuMail,
        LuXTwitter,
        LuGithub,
        LuGlobe,
        LuBug,
        LuMessageCircleQuestion,
        LuBookText,
        LuNewspaper,
        LuMessagesSquare,
      ],
      mode: 'light',
    },
  });
  nuxtApp.vueApp.use(rui);
});
