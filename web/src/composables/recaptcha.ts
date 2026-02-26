import { ref } from 'vue';

export function useRecaptcha() {
  const recaptchaPassed = ref(false);
  const recaptchaToken = ref('');
  const captchaId = ref<number>();

  const onSuccess = (token: string): void => {
    recaptchaToken.value = token;
    recaptchaPassed.value = true;
  };

  const onExpired = (): void => {
    recaptchaToken.value = '';
    recaptchaPassed.value = false;
  };

  const onError = (): void => {
    recaptchaPassed.value = false;
  };

  const resetCaptcha = (): void => {
    onExpired();
    window.grecaptcha?.reset(captchaId.value);
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
