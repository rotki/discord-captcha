<script setup lang="ts">
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRecaptcha } from '../composables/recaptcha';
import Recaptcha from './Recaptcha.vue';

const { t } = useI18n({ useScope: 'global' });

interface DiscordInvite {
  code: string;
  expires_at: string;
}

const invite = ref<DiscordInvite>();
const error = ref<string>();

const inviteLink = computed(() => {
  if (!invite.value)
    return null;
  return `https://discord.com/invite/${invite.value.code}`;
});

const expiry = computed(() => {
  if (!invite.value)
    return null;
  return new Date(invite.value.expires_at);
});

const { onSuccess, onError, onExpired, captchaId, resetCaptcha } = useRecaptcha();

async function onCaptchaSuccess(token: string) {
  onSuccess(token);

  try {
    const response = await fetch('/api/discord-invite', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ captcha: token }),
    });

    if (!response.ok) {
      const text = await response.text();
      throw new Error(text || response.statusText);
    }

    invite.value = await response.json();
  }
  catch (error_: any) {
    resetCaptcha();
    error.value = error_.message;
    invite.value = undefined;
  }
}
</script>

<template>
  <div class="px-4 md:px-12 xl:px-16">
    <div
      class="md:h-[294px] py-[5rem] flex flex-row items-center justify-between flex-wrap"
    >
      <div class="basis-full md:basis-auto">
        <div class="text-h6 text-rui-primary">
          {{ t('discord.chat') }}
        </div>
        <div class="text-h4 mt-[0.75rem]">
          {{ t('discord.title') }}
        </div>

        <div v-show="!inviteLink">
          {{ t('discord.description') }}
        </div>
        <div
          v-if="inviteLink"
          class="text-body-1 mt-[1.5rem]"
        >
          <i18n-t
            tag="div"
            scope="global"
            keypath="discord.invite.link"
          >
            <template #link>
              <a
                :href="inviteLink"
                target="_blank"
                rel="noreferrer"
                class="font-medium text-rui-light-primary hover:text-rui-light-primary-darker"
              >
                {{ inviteLink }}
              </a>
            </template>
          </i18n-t>
          <i18n-t
            v-if="expiry"
            tag="div"
            scope="global"
            keypath="discord.invite.expiry"
          >
            <template #expiry>
              <span class="font-medium text-rui-text-secondary">
                {{ expiry.toLocaleString() }}
              </span>
            </template>
          </i18n-t>
        </div>
      </div>
      <Recaptcha
        v-if="!inviteLink"
        class="basis-full md:basis-auto"
        @error="onError()"
        @expired="onExpired()"
        @success="onCaptchaSuccess($event)"
        @captcha-id="captchaId = $event"
      />
    </div>
  </div>
</template>
