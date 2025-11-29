<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import DataTable, { type Column } from '@/components/common/DataTable.vue'
import type { FilterKey } from '@/types/api'

const toast = useToastStore()
const data = ref<FilterKey[]>([])
const loading = ref(true)

const columns: Column[] = [
  { key: 'id', label: 'ID', editable: false },
  { key: 'name', label: 'Name', editable: true, type: 'text' }
]

onMounted(async () => {
  try {
    data.value = await api.filterKeys.list()
  } catch (e) {
    toast.error('Failed to load filter keys')
  } finally {
    loading.value = false
  }
})

async function handleSave(item: FilterKey, isNew: boolean) {
  if (isNew) {
    const created = await api.filterKeys.create({ name: item.name })
    data.value.push(created)
    toast.success('Filter key created')
  } else {
    const updated = await api.filterKeys.update(item.id, { name: item.name })
    const index = data.value.findIndex(d => d.id === item.id)
    data.value[index] = updated
    toast.success('Filter key updated')
  }
}

async function handleDelete(id: number) {
  await api.filterKeys.delete(id)
  data.value = data.value.filter(d => d.id !== id)
  toast.success('Filter key deleted')
}
</script>

<template>
  <div class="mt-4">
    <h2>Filter Keys</h2>
    <p class="text-muted mb-4">Manage filter key reference data</p>

    <DataTable :columns="columns" :data="data" :loading="loading" :editable="true" :on-save="handleSave"
      :on-delete="handleDelete" />
  </div>
</template>
