<script setup lang="ts">
import { nextTick, onMounted, ref, toRef, useTemplateRef } from 'vue';

const { theme = 'light', size = 'normal', useRecaptchaNet = false } = defineProps<{
  theme?: 'light' | 'dark';
  size?: 'compact' | 'normal';
  useRecaptchaNet?: boolean;
}>();

const emit = defineEmits<{
  'error': [];
  'expired': [];
  'success': [value: string];
  'captcha-id': [value: number];
}>();

const siteKey = import.meta.env.VITE_RECAPTCHA_SITE_KEY;

const recaptchaEl = useTemplateRef<HTMLElement>('recaptchaEl');
const rendered = ref(false);

const grecaptcha = toRef(window, 'grecaptcha');

function renderCaptcha() {
  if (rendered.value || !(grecaptcha.value && recaptchaEl.value))
    return;

  const id = grecaptcha.value.render(recaptchaEl.value, {
    'sitekey': siteKey,
    'callback': (token: string) => emit('success', token),
    'expired-callback': () => emit('expired'),
    'error-callback': () => emit('error'),
    'theme': theme,
    'size': size,
  });

  emit('captcha-id', id);
  rendered.value = true;
}

window.onRecaptchaLoaded = renderCaptcha;

onMounted(() => {
  if (!grecaptcha.value) {
    const script = document.createElement('script');
    script.src = `${
      useRecaptchaNet
        ? 'https://www.recaptcha.net/recaptcha'
        : 'https://www.google.com/recaptcha'
    }/api.js?onload=onRecaptchaLoaded&render=explicit`;
    script.defer = true;
    script.async = true;
    document.head.appendChild(script);
  }
  else {
    nextTick(renderCaptcha);
  }
});
</script>

<template>
  <div class="mt-4">
    <div
      ref="recaptchaEl"
      class="h-20"
    />
  </div>
</template>
