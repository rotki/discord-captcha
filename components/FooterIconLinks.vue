<script setup lang="ts">
const menus = computed(() => [
  {
    title: 'Github',
    showText: false,
    href: 'https://github.com/rotki',
    icon: 'github-line' as const,
  },
  {
    title: 'Twitter',
    showText: false,
    href: 'https://twitter.com/rotkiapp',
    icon: 'twitter-line' as const,
  },
  {
    title: 'Contact',
    showText: true,
    href: 'mailto:info@rotki.com',
    icon: 'mail-send-line' as const,
  },
]);
</script>

<template>
  <div class="flex flex-wrap items-center space-x-4">
    <div
      v-for="menu in menus"
      :key="menu.title"
    >
      <ButtonLink
        size="lg"
        :to="menu.href"
        :title="menu.title"
        external
      >
        <span class="flex space-x-2 items-center">
          <span v-if="menu.showText">{{ menu.title }}</span>
          <RuiIcon :name="menu.icon" />
        </span>
      </ButtonLink>
    </div>
  </div>
</template>
