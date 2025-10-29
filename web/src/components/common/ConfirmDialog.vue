<script setup lang="ts">
defineProps<{
  show: boolean
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'primary' | 'warning'
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<template>
  <Teleport to="body">
    <div v-if="show" class="modal-backdrop fade show"></div>
    <div
      v-if="show"
      class="modal fade show d-block"
      tabindex="-1"
      @click.self="emit('cancel')"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ title }}</h5>
            <button
              type="button"
              class="btn-close"
              aria-label="Close"
              @click="emit('cancel')"
            ></button>
          </div>
          <div class="modal-body">
            <p>{{ message }}</p>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-secondary"
              @click="emit('cancel')"
            >
              {{ cancelText || 'Cancel' }}
            </button>
            <button
              type="button"
              class="btn"
              :class="`btn-${variant || 'danger'}`"
              @click="emit('confirm')"
            >
              {{ confirmText || 'Confirm' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
