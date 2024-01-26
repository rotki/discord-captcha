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
  createRui,
} from '@rotki/ui-library';
import '@fontsource/roboto/latin.css';

export default defineNuxtPlugin((nuxtApp) => {
  const rui = createRui({
    theme: {
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
      mode: 'light',
    },
  });
  nuxtApp.vueApp.use(rui);
});
