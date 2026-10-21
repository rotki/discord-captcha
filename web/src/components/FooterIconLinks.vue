<script setup lang="ts">
import { RuiIcon } from '@rotki/ui-library';
import ButtonLink from './ButtonLink.vue';

interface FooterLink {
  title: string;
  showText: boolean;
  href: string;
  icon: 'lu-github' | 'lu-x-twitter' | 'lu-mail';
  external: boolean;
}

const links: FooterLink[] = [
  {
    title: 'GitHub',
    showText: false,
    href: 'https://github.com/rotki',
    icon: 'lu-github',
    external: true,
  },
  {
    title: 'X (Twitter)',
    showText: false,
    href: 'https://twitter.com/rotkiapp',
    icon: 'lu-x-twitter',
    external: true,
  },
  {
    title: 'Contact',
    showText: true,
    href: 'mailto:info@rotki.com',
    icon: 'lu-mail',
    // a new tab would open the mail client and leave an empty tab behind
    external: false,
  },
];
</script>

<template>
  <nav
    class="flex flex-wrap items-center gap-2"
    aria-label="rotki links"
  >
    <ButtonLink
      v-for="link in links"
      :key="link.href"
      size="lg"
      :to="link.href"
      :title="link.title"
      :aria-label="link.showText ? undefined : link.title"
      :external="link.external"
    >
      <span class="flex items-center gap-2">
        <span v-if="link.showText">{{ link.title }}</span>
        <RuiIcon
          :name="link.icon"
          aria-hidden="true"
        />
      </span>
    </ButtonLink>
  </nav>
</template>
