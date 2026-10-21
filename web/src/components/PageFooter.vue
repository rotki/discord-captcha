<script setup lang="ts">
import FooterIconLinks from './FooterIconLinks.vue';

const year = new Date().getFullYear().toString();
</script>

<template>
  <footer>
    <div class="px-4 md:px-12 xl:px-16">
      <div
        class="border-t border-rui-divider grid grid-cols-1 lg:grid-cols-2 gap-4 pt-6 pb-8 lg:pt-8 lg:pb-12 lg:items-center"
      >
        <p
          class="flex flex-wrap justify-center lg:justify-start gap-x-1 text-body-2 text-rui-text-secondary order-2 lg:order-1"
        >
          <span>&copy; Rotki Solutions GmbH 2018-{{ year }}.</span>
          <span>All Rights Reserved.</span>
        </p>
        <div class="flex justify-center lg:justify-end order-1 lg:order-2">
          <FooterIconLinks />
        </div>
      </div>
    </div>
  </footer>
</template>
