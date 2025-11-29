<script setup lang="ts">
import { ref, watch } from 'vue'
import { api } from '@/services/api'
import type { TorrentCategory } from '@/types/api'

const props = defineProps<{
  show: boolean
  currentCategory?: string
}>()

const emit = defineEmits<{
  confirm: [category: string]
  cancel: []
}>()

const categories = ref<TorrentCategory[]>([])
const selectedCategory = ref<string>(props.currentCategory || '')
const loading = ref(false)

watch(() => props.show, async (newValue) => {
  if (newValue) {
    selectedCategory.value = props.currentCategory || ''
    loading.value = true
    try {
      categories.value = await api.torrentCategories.list()
    } catch (e) {
      console.error('Failed to load categories', e)
    } finally {
      loading.value = false
    }
  }
})

function handleConfirm() {
  if (!selectedCategory.value) {
    return
  }
  emit('confirm', selectedCategory.value)
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
            <h5 class="modal-title">Change Category</h5>
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
                <label for="category-select" class="form-label">Category</label>
                <select
                  id="category-select"
                  v-model="selectedCategory"
                  class="form-select"
                  :disabled="loading"
                  required
                >
                  <option value="">{{ loading ? 'Loading...' : 'Select a category' }}</option>
                  <option
                    v-for="category in categories"
                    :key="category.id"
                    :value="category.name"
                  >
                    {{ category.name }}
                  </option>
                </select>
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
              :disabled="loading || !selectedCategory"
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
