<script setup lang="ts">
import { RuiButton } from '@rotki/ui-library';
import { useAttrs } from 'vue';

defineOptions({
  inheritAttrs: false,
});

const { to, external = false } = defineProps<{
  to: string;
  external?: boolean;
}>();

defineSlots<{
  default: () => void;
  append: () => void;
}>();

const attrs = useAttrs();
</script>

<template>
  <RuiButton
    v-bind="{
      variant: 'text',
      ...attrs,
    }"
    tag="a"
    :href="to"
    :target="external ? '_blank' : undefined"
    :rel="external ? 'noopener noreferrer' : undefined"
  >
    <slot>
      {{ to }}
    </slot>
    <template #append>
      <slot name="append" />
    </template>
  </RuiButton>
</template>
