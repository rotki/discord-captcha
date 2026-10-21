<script setup lang="ts">
import { RuiAlert } from '@rotki/ui-library';
import { useI18n } from 'vue-i18n';
import { useDiscordInvite } from '../composables/discord-invite';
import { useRecaptcha } from '../composables/recaptcha';
import InviteResult from './InviteResult.vue';
import Recaptcha from './Recaptcha.vue';

const { t } = useI18n({ useScope: 'global' });

const { onSuccess, onError, onExpired, captchaId, resetCaptcha } = useRecaptcha();
const { error, expiry, inviteLink, requestInvite } = useDiscordInvite();

async function onCaptchaSuccess(token: string): Promise<void> {
  onSuccess(token);
  if (!await requestInvite(token))
    resetCaptcha();
}
</script>

<template>
  <section
    class="px-4 md:px-12 xl:px-16 py-16 md:py-20 grid gap-8 md:grid-cols-[1fr_auto] md:items-center"
    aria-labelledby="discord-title"
  >
    <div>
      <p class="text-overline uppercase text-rui-primary">
        {{ t('discord.chat') }}
      </p>
      <h1
        id="discord-title"
        class="text-h4 md:text-h3 mt-1 text-rui-text"
      >
        {{ t('discord.title') }}
      </h1>
      <InviteResult
        v-if="inviteLink"
        :invite-link="inviteLink"
        :expiry="expiry"
      />
      <p
        v-else
        class="text-body-1 mt-3 max-w-prose text-rui-text-secondary"
      >
        {{ t('discord.description') }}
      </p>
      <RuiAlert
        v-if="error"
        class="mt-6 max-w-prose"
        type="error"
        :title="t('discord.error.title')"
        :description="t('discord.error.description', { message: error })"
      />
    </div>
    <Recaptcha
      v-if="!inviteLink"
      @error="onError()"
      @expired="onExpired()"
      @success="onCaptchaSuccess($event)"
      @captcha-id="captchaId = $event"
    />
  </section>
</template>
