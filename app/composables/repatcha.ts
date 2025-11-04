import { set } from '@vueuse/core';
import { ref } from 'vue';

export function useRecaptcha() {
  const recaptchaPassed = ref(false);
  const recaptchaToken = ref('');
  const captchaId = ref<number>();

  const onSuccess = (token: string): void => {
    set(recaptchaToken, token);
    set(recaptchaPassed, true);
  };

  const onExpired = (): void => {
    set(recaptchaToken, '');
    set(recaptchaPassed, false);
  };

  const onError = (): void => {
    set(recaptchaPassed, false);
  };

  const resetCaptcha = (): void => {
    onExpired();
    if (import.meta.client) {
      window.grecaptcha?.reset(captchaId.value);
    }
  };

  return {
    captchaId,
    onError,
    onExpired,
    onSuccess,
    recaptchaPassed,
    recaptchaToken,
    resetCaptcha,
  };
}
