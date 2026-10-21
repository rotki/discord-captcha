<script setup lang="ts">
import { computed } from 'vue';

const siteUrl = import.meta.env.VITE_SITE_URL;

const parsedUrl = computed<URL | undefined>(() => {
  if (!siteUrl)
    return undefined;
  try {
    return new URL(siteUrl);
  }
  catch {
    return undefined;
  }
});

const displayUrl = computed<string>(() => {
  if (!parsedUrl.value)
    return '';
  const { host, pathname } = parsedUrl.value;
  return `${host}${pathname.replace(/\/$/, '')}`;
});

// rotki.com attributes visits by these parameters, so a click from here is not counted as direct traffic.
const linkUrl = computed<string | undefined>(() => {
  if (!parsedUrl.value)
    return undefined;
  const url = new URL(parsedUrl.value);
  url.searchParams.set('utm_source', window.location.hostname);
  url.searchParams.set('utm_medium', 'referral');
  url.searchParams.set('utm_content', 'header');
  return url.toString();
});
</script>

<template>
  <header
    v-if="displayUrl"
    class="bg-rui-primary h-12 px-4 text-body-2 text-rui-dark-text flex items-center justify-center"
  >
    <i18n-t
      keypath="discord.link"
      scope="global"
    >
      <template #link>
        <a
          :href="linkUrl"
          target="_blank"
          rel="noopener"
          class="ml-1 font-medium underline underline-offset-2 hover:text-rui-grey-200"
        >
          {{ displayUrl }}
        </a>
      </template>
    </i18n-t>
  </header>
</template>
