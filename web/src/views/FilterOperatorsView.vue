<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import DataTable, { type Column } from '@/components/common/DataTable.vue'
import type { FilterOperator } from '@/types/api'

const toast = useToastStore()
const data = ref<FilterOperator[]>([])
const loading = ref(true)

const columns: Column[] = [
    { key: 'id', label: 'ID', editable: false },
    { key: 'name', label: 'Name', editable: true, type: 'text' }
]

onMounted(async () => {
    try {
        data.value = await api.filterOperators.list()
    } catch (e) {
        toast.error('Failed to load filter operators')
    } finally {
        loading.value = false
    }
})

async function handleSave(item: FilterOperator, isNew: boolean) {
    if (isNew) {
        const created = await api.filterOperators.create({ name: item.name })
        data.value.push(created)
        toast.success('Filter operator created')
    } else {
        const updated = await api.filterOperators.update(item.id, { name: item.name })
        const index = data.value.findIndex(d => d.id === item.id)
        data.value[index] = updated
        toast.success('Filter operator updated')
    }
}

async function handleDelete(id: number) {
    await api.filterOperators.delete(id)
    data.value = data.value.filter(d => d.id !== id)
    toast.success('Filter operator deleted')
}
</script>

<template>
    <div class="mt-4">
        <h2>Filter Operators</h2>
        <p class="text-muted mb-4">Manage filter operator reference data</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="true" :on-save="handleSave"
            :on-delete="handleDelete" />
    </div>
</template>
