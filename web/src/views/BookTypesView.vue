<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import DataTable, { type Column } from '@/components/common/DataTable.vue'
import type { BookType } from '@/types/api'

const toast = useToastStore()
const data = ref<BookType[]>([])
const loading = ref(true)

const columns: Column[] = [
    { key: 'id', label: 'ID', editable: false },
    { key: 'name', label: 'Name', editable: false }
]

onMounted(async () => {
    try {
        data.value = await api.bookTypes.list()
    } catch (e) {
        toast.error('Failed to load book types')
    } finally {
        loading.value = false
    }
})
</script>

<template>
    <div class="mt-4">
        <h2>Book Types</h2>
        <p class="text-muted mb-4">Reference data for book types (read-only)</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="false" />
    </div>
</template>
