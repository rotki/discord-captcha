<script setup lang="ts">
import { computed } from 'vue';

const siteUrl = import.meta.env.VITE_SITE_URL;
const displayUrl = computed<string>(() => {
  const index = siteUrl?.indexOf('//') ?? -1;
  return index === -1 ? '' : siteUrl.slice(index + 2);
});
</script>

<template>
  <header>
    <div
      class="bg-rui-primary h-[48px] text-rui-dark-text flex flex-row items-center justify-center"
    >
      <i18n-t
        keypath="discord.link"
        scope="global"
      >
        <template #link>
          <a
            :href="siteUrl"
            target="_blank"
            class="ml-1 hover:text-rui-grey-200"
          >
            {{ displayUrl }}
          </a>
        </template>
      </i18n-t>
    </div>
  </header>
</template>
