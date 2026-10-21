<script lang="ts" setup>
import { RuiIcon } from '@rotki/ui-library';

defineProps<{
  icon: string;
  title: string;
  description: string;
}>();
</script>

<template>
  <li class="flex items-start gap-4">
    <div
      class="flex size-12 shrink-0 items-center justify-center rounded-rui-card bg-rui-primary text-rui-dark-text"
      aria-hidden="true"
    >
      <RuiIcon :name="icon" />
    </div>
    <div>
      <h3 class="text-subtitle-1 text-rui-text">
        {{ title }}
      </h3>
      <p class="text-body-1 mt-1 text-rui-text-secondary">
        {{ description }}
      </p>
    </div>
  </li>
</template>
