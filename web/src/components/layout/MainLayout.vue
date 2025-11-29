<script setup lang="ts">
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import ToastContainer from '@/components/common/ToastContainer.vue'
</script>

<template>
  <div class="app-layout">
    <AppSidebar />
    <div class="app-main">
      <AppHeader />
      <main class="app-content">
        <slot />
      </main>
    </div>
  </div>
  <ToastContainer />
</template>
