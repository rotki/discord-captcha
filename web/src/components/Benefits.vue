<script lang="ts" setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import BenefitDetails from './BenefitDetails.vue';

interface BenefitDetail {
  title: string;
  description: string;
  icon: string;
}

const { t } = useI18n({ useScope: 'global' });

const benefits = computed<BenefitDetail[]>(() => [
  {
    title: t('discord.benefits.items.engagement.title'),
    description: t('discord.benefits.items.engagement.description'),
    icon: 'lu-globe',
  },
  {
    title: t('discord.benefits.items.bug_reporting.title'),
    description: t('discord.benefits.items.bug_reporting.description'),
    icon: 'lu-bug',
  },
  {
    title: t('discord.benefits.items.community_support.title'),
    description: t('discord.benefits.items.community_support.description'),
    icon: 'lu-message-circle-question-mark',
  },
  {
    title: t('discord.benefits.items.learn.title'),
    description: t('discord.benefits.items.learn.description'),
    icon: 'lu-book-text',
  },
  {
    title: t('discord.benefits.items.stay_updated.title'),
    description: t('discord.benefits.items.stay_updated.description'),
    icon: 'lu-newspaper',
  },
  {
    title: t('discord.benefits.items.networking.title'),
    description: t('discord.benefits.items.networking.description'),
    icon: 'lu-messages-square',
  },
]);
</script>

<template>
  <section
    class="px-4 md:px-12 xl:px-16 py-16 md:py-20"
    aria-labelledby="benefits-title"
  >
    <h2
      id="benefits-title"
      class="text-h5 md:text-h4 text-rui-text"
    >
      {{ t('discord.benefits.title') }}
    </h2>
    <p class="text-body-1 mt-4 max-w-prose text-rui-text-secondary">
      {{ t('discord.benefits.description') }}
    </p>
    <ul class="mt-12 grid gap-x-12 gap-y-10 md:grid-cols-2">
      <BenefitDetails
        v-for="benefit in benefits"
        :key="benefit.icon"
        v-bind="benefit"
      />
    </ul>
  </section>
</template>
