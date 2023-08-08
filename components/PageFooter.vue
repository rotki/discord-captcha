<script setup lang="ts">
const year = new Date().getFullYear().toString();
</script>

<template>
  <div>
    <div class="container">
      <div class="border-t border-grey-100" />
      <div
        class="flex flex-col-reverse lg:flex-row justify-between items-center pt-4 pb-6 lg:pt-8 lg:pb-12"
      >
        <div
          class="flex flex-wrap justify-center space-x-1 text-rui-text-secondary mt-2 lg:mt-0"
        >
          <div>&copy; Rotki Solutions GmbH 2018-{{ year }}.</div>
          <div>All Rights Reserved.</div>
        </div>
        <div>
          <FooterIconLinks />
        </div>
      </div>
    </div>
  </div>
</template>
