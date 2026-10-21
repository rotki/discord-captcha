<script setup lang="ts">
import { computed } from 'vue';

const siteUrl = import.meta.env.VITE_SITE_URL;
const displayUrl = computed<string>(() => {
  const index = siteUrl?.indexOf('//') ?? -1;
  return index === -1 ? '' : siteUrl.slice(index + 2);
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
          :href="siteUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="ml-1 font-medium underline underline-offset-2 hover:text-rui-grey-200"
        >
          {{ displayUrl }}
        </a>
      </template>
    </i18n-t>
  </header>
</template>
