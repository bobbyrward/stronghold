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
    { key: 'name', label: 'Name', editable: false }
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
</script>

<template>
    <div class="mt-4">
        <h2>Filter Keys</h2>
        <p class="text-muted mb-4">Reference data for filter keys (read-only)</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="false" />
    </div>
</template>
