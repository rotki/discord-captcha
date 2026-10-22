<script setup lang="ts">
import { RuiButton, RuiIcon } from '@rotki/ui-library';
import { computed, onBeforeUnmount, shallowRef } from 'vue';
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

const copied = shallowRef(false);
let copiedTimer: ReturnType<typeof setTimeout> | undefined;

async function copyLink(): Promise<void> {
  try {
    await navigator.clipboard.writeText(inviteLink);
  }
  catch {
    // Clipboard access was refused; the link is still there to select by hand.
    return;
  }
  copied.value = true;
  clearTimeout(copiedTimer);
  copiedTimer = setTimeout(() => {
    copied.value = false;
  }, 2000);
}

onBeforeUnmount(() => clearTimeout(copiedTimer));
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
    <div class="max-w-full text-body-2 text-rui-text-secondary">
      <p>{{ t('discord.invite.link') }}</p>
      <div class="flex max-w-full items-center gap-1">
        <a
          :href="inviteLink"
          target="_blank"
          rel="noopener noreferrer"
          class="min-w-0 break-all font-medium text-rui-primary underline underline-offset-2 hover:text-rui-primary-darker"
        >
          {{ inviteLink }}
        </a>
        <RuiButton
          variant="text"
          icon
          size="sm"
          class="shrink-0"
          :aria-label="copied ? t('discord.invite.copied') : t('discord.invite.copy')"
          :title="copied ? t('discord.invite.copied') : t('discord.invite.copy')"
          @click="copyLink()"
        >
          <RuiIcon
            :name="copied ? 'lu-check' : 'lu-copy'"
            size="16"
          />
        </RuiButton>
        <span
          class="sr-only"
          aria-live="polite"
        >
          {{ copied ? t('discord.invite.copied') : '' }}
        </span>
      </div>
      <i18n-t
        v-if="formattedExpiry"
        tag="p"
        class="mt-1"
        scope="global"
        keypath="discord.invite.expiry"
      >
        <template #expiry>
          <span class="font-medium text-rui-text">{{ formattedExpiry }}</span>
        </template>
      </i18n-t>
    </div>
  </div>
</template>
