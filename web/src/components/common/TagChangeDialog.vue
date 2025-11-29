<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  show: boolean
  currentTags?: string
}>()

const emit = defineEmits<{
  confirm: [tags: string]
  cancel: []
}>()

const tags = ref<string>(props.currentTags || '')

watch(() => props.show, (newValue) => {
  if (newValue) {
    tags.value = props.currentTags || ''
  }
})

function handleConfirm() {
  emit('confirm', tags.value)
}
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
            <h5 class="modal-title">Change Tags</h5>
            <button
              type="button"
              class="btn-close"
              aria-label="Close"
              @click="emit('cancel')"
            ></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="handleConfirm">
              <div class="mb-3">
                <label for="tags-input" class="form-label">Tags (comma-separated)</label>
                <input
                  id="tags-input"
                  v-model="tags"
                  type="text"
                  class="form-control"
                  placeholder="e.g., tag1, tag2, tag3"
                />
                <div class="form-text">Enter tags as a comma-separated list</div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-secondary"
              @click="emit('cancel')"
            >
              Cancel
            </button>
            <button
              type="button"
              class="btn btn-primary"
              @click="handleConfirm"
            >
              Save
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
