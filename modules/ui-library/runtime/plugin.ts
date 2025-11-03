import {
  createRui,
  LuBookText,
  LuBug,
  LuGithub,
  LuGlobe,
  LuMail,
  LuMessageCircleQuestionMark,
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
        LuMessageCircleQuestionMark,
        LuBookText,
        LuNewspaper,
        LuMessagesSquare,
      ],
      mode: 'light',
    },
  });
  nuxtApp.vueApp.use(rui);
});
