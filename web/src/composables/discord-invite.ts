import { computed, type ComputedRef, readonly, type Ref, shallowRef } from 'vue';

interface DiscordInvite {
  code: string;
  expires_at: string;
}

interface UseDiscordInviteReturn {
  error: Readonly<Ref<string | undefined>>;
  expiry: ComputedRef<Date | undefined>;
  inviteLink: ComputedRef<string | undefined>;
  loading: Readonly<Ref<boolean>>;
  requestInvite: (captcha: string) => Promise<boolean>;
}

export function useDiscordInvite(): UseDiscordInviteReturn {
  const invite = shallowRef<DiscordInvite>();
  const error = shallowRef<string>();
  const loading = shallowRef(false);

  const inviteLink = computed<string | undefined>(() => {
    if (!invite.value)
      return undefined;
    return `https://discord.com/invite/${invite.value.code}`;
  });

  const expiry = computed<Date | undefined>(() => {
    if (!invite.value?.expires_at)
      return undefined;
    return new Date(invite.value.expires_at);
  });

  const requestInvite = async (captcha: string): Promise<boolean> => {
    loading.value = true;
    error.value = undefined;
    try {
      const response = await fetch('/api/discord-invite', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ captcha }),
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text.trim() || response.statusText);
      }

      invite.value = await response.json();
      return true;
    }
    catch (error_: unknown) {
      error.value = error_ instanceof Error ? error_.message : String(error_);
      invite.value = undefined;
      return false;
    }
    finally {
      loading.value = false;
    }
  };

  return {
    error: readonly(error),
    expiry,
    inviteLink,
    loading: readonly(loading),
    requestInvite,
  };
}
