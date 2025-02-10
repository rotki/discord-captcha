<script lang="ts" setup>
interface BenefitDetails {
  title: string;
  description: string;
  icon: string;
}

const { t } = useI18n({ useScope: 'global' });

const benefits = computed<BenefitDetails[]>(() => [
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
    icon: 'lu-message-circle-question',
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
  <div class="container">
    <div
      class="text-primary py-[5rem] flex flex-col items-center self-stretch gap-[3rem]"
    >
      <div>
        <div class="text-h4">
          {{ t('discord.benefits.title') }}
        </div>
        <div class="text-body-1 mt-[1.25rem]">
          {{ t('discord.benefits.description') }}
        </div>
      </div>
      <div class="flex justify-between items-center flex-wrap gap-[3rem]">
        <BenefitDetails
          v-for="(benefit, i) in benefits"
          :key="i"
          v-bind="benefit"
        />
      </div>
    </div>
  </div>
</template>
