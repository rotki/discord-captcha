import {
  RiBook2Line,
  RiBug2Line,
  RiGithubLine,
  RiMailSendLine,
  RiNewspaperLine,
  RiQuestionnaireLine,
  RiTeamLine,
  RiTwitterLine,
  RiWechatLine,
  RuiPlugin,
} from '@rotki/ui-library';
import { defineNuxtPlugin } from '#app';
import '@fontsource/roboto/latin.css';

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(RuiPlugin, {
    mode: 'light',
    icons: [
      RiGithubLine,
      RiMailSendLine,
      RiTwitterLine,
      RiTeamLine,
      RiBug2Line,
      RiQuestionnaireLine,
      RiBook2Line,
      RiNewspaperLine,
      RiWechatLine,
    ],
  });
});
