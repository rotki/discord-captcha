<script lang="ts" setup>
import { get } from '@vueuse/core';

const error = useError();
const { t } = useI18n({ useScope: 'global' });
const css = useCssModule();

const title = computed(() => get(error)?.message ?? '');
const statusCode = computed(() => {
  const err = get(error);
  return err && 'statusCode' in err ? err.statusCode : -1;
});

const handleError = () => clearError({ redirect: '/' });

useHead(() => ({
  title,
}));
</script>

<template>
  <div class="container">
    <div
      v-if="error"
      :class="css.wrapper"
    >
      <h1 :class="css.title">
        {{ statusCode }}
      </h1>
      <p :class="css.subtitle">
        {{ title }}
      </p>

      <div class="flex justify-center">
        <RuiButton
          variant="default"
          size="lg"
          filled
          color="primary"
          @click="handleError()"
        >
          {{ t('actions.go_back_home') }}
        </RuiButton>
      </div>
    </div>
  </div>
</template>

<style module lang="scss">
.wrapper {
  @apply text-center my-24;
}
.title {
  @apply block text-red-500 font-bold text-6xl;
}

.subtitle {
  @apply block font-light text-rui-text my-8 text-3xl;
}
</style>
