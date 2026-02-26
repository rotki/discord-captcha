<script setup lang="ts">
import PageFooter from './components/PageFooter.vue';
import PageHeader from './components/PageHeader.vue';
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <PageHeader />
    <main class="mx-auto w-full max-w-7xl flex-1 px-4">
      <RouterView />
    </main>
    <PageFooter class="mx-auto w-full max-w-7xl px-4" />
  </div>
</template>
