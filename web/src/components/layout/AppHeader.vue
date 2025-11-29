<script setup lang="ts">
import { useUiStore } from '@/stores/ui'

const uiStore = useUiStore()
</script>

<template>
  <header class="app-header">
    <button
      class="btn btn-link text-light p-0 me-3"
      @click="uiStore.toggleSidebar"
      aria-label="Toggle sidebar"
    >
      <i class="bi bi-list fs-4"></i>
    </button>
    <h1 class="h5 mb-0">Stronghold</h1>
  </header>
</template>
