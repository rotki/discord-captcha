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
  <a
    class="inline-flex"
    :href="to"
    :target="external ? '_blank' : '_self'"
    :rel="external ? 'noreferrer' : undefined"
  >
    <RuiButton
      v-bind="{
        variant: 'text',
        type: 'button',
        ...attrs,
      }"
    >
      <slot>
        {{ to }}
      </slot>
      <template #append>
        <slot name="append" />
      </template>
    </RuiButton>
  </a>
</template>
