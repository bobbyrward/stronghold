<script setup lang="ts">
import { useToastStore } from '@/stores/toast'

const toastStore = useToastStore()

function getToastClass(type: string) {
  switch (type) {
    case 'success':
      return 'bg-success text-white'
    case 'error':
      return 'bg-danger text-white'
    case 'warning':
      return 'bg-warning text-dark'
    case 'info':
      return 'bg-info text-white'
    default:
      return 'bg-secondary text-white'
  }
}
</script>

<template>
  <Teleport to="body">
    <div class="toast-container">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toastStore.toasts"
          :key="toast.id"
          class="toast show"
          :class="getToastClass(toast.type)"
        >
          <div class="toast-body d-flex justify-content-between align-items-center">
            <span>{{ toast.message }}</span>
            <button
              type="button"
              class="btn-close btn-close-white ms-2"
              @click="toastStore.remove(toast.id)"
              aria-label="Close"
            ></button>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 1rem;
  right: 1rem;
  z-index: 1100;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.toast {
  min-width: 250px;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
