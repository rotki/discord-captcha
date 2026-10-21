<script setup lang="ts">
import { RuiButton, RuiIcon } from '@rotki/ui-library';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { inviteLink, expiry } = defineProps<{
  inviteLink: string;
  expiry?: Date;
}>();

const { t } = useI18n({ useScope: 'global' });

const formattedExpiry = computed<string | undefined>(() => {
  if (!expiry)
    return undefined;
  return new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(expiry);
});
</script>

<template>
  <div class="mt-6 flex flex-col items-start gap-3">
    <RuiButton
      tag="a"
      color="primary"
      size="lg"
      :href="inviteLink"
      target="_blank"
      rel="noopener noreferrer"
    >
      {{ t('discord.invite.join') }}
      <template #append>
        <RuiIcon
          name="lu-external-link"
          size="18"
        />
      </template>
    </RuiButton>
    <p class="text-body-2 text-rui-text-secondary">
      <i18n-t
        scope="global"
        keypath="discord.invite.link"
      >
        <template #link>
          <a
            :href="inviteLink"
            target="_blank"
            rel="noopener noreferrer"
            class="font-medium text-rui-primary underline underline-offset-2 break-all hover:text-rui-primary-darker"
          >
            {{ inviteLink }}
          </a>
        </template>
      </i18n-t>
      <template v-if="formattedExpiry">
        {{ ' ' }}
        <i18n-t
          scope="global"
          keypath="discord.invite.expiry"
        >
          <template #expiry>
            <span class="font-medium text-rui-text">{{ formattedExpiry }}</span>
          </template>
        </i18n-t>
      </template>
    </p>
  </div>
</template>
