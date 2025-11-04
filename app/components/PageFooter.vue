<script setup lang="ts">
const year = new Date().getFullYear().toString();
</script>

<template>
  <div>
    <div class="container">
      <div class="border-t border-gray-100" />
      <div
        class="grid grid-cols-1 lg:grid-cols-2 gap-4 pt-4 pb-6 lg:pt-8 lg:pb-12 lg:items-center"
      >
        <div
          class="flex flex-wrap justify-center lg:justify-start space-x-1 text-rui-text-secondary order-2 lg:order-1"
        >
          <div>&copy; Rotki Solutions GmbH 2018-{{ year }}.</div>
          <div>All Rights Reserved.</div>
        </div>
        <div class="flex justify-center lg:justify-end order-1 lg:order-2">
          <FooterIconLinks />
        </div>
      </div>
    </div>
  </div>
</template>
