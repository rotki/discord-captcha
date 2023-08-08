<script setup lang="ts">
const {
  public: { siteUrl },
} = useRuntimeConfig();

const displayUrl = computed(() => siteUrl.split('//')[1]);
</script>

<template>
  <header>
    <div
      class="bg-rui-primary h-[48px] text-rui-dark-text flex flex-row items-center justify-center"
    >
      <i18n-t keypath="discord.link">
        <template #link>
          <NuxtLink
            :href="siteUrl"
            target="_blank"
            class="ml-1 hover:text-rui-grey-200"
          >
            {{ displayUrl }}
          </NuxtLink>
        </template>
      </i18n-t>
    </div>
  </header>
</template>
