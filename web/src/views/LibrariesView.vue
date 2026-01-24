<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import DataTable, { type Column } from '@/components/common/DataTable.vue'
import type { Library, BookType } from '@/types/api'

const toast = useToastStore()
const data = ref<Library[]>([])
const bookTypes = ref<BookType[]>([])
const loading = ref(true)

const columns = ref<Column[]>([
    { key: 'id', label: 'ID', editable: false },
    { key: 'name', label: 'Name', editable: true, type: 'text' },
    { key: 'path', label: 'Path', editable: true, type: 'text' },
    {
        key: 'book_type_name',
        label: 'Book Type',
        editable: true,
        type: 'select',
        displayKey: 'book_type_name',
        options: []
    }
])

onMounted(async () => {
    try {
        const [libs, types] = await Promise.all([
            api.libraries.list(),
            api.bookTypes.list()
        ])
        data.value = libs
        bookTypes.value = types

        // Populate book type options
        const typeColumn = columns.value.find(c => c.key === 'book_type_name')
        if (typeColumn) {
            typeColumn.options = types.map(t => ({
                value: t.name,
                label: t.name
            }))
        }
    } catch (e) {
        toast.error('Failed to load libraries')
    } finally {
        loading.value = false
    }
})

async function handleSave(item: Partial<Library>, isNew: boolean) {
    if (!item.name || !item.path || !item.book_type_name) {
        toast.error('Name, Path, and Book Type are required')
        throw new Error('Validation failed')
    }

    const request = {
        name: item.name,
        path: item.path,
        book_type_name: item.book_type_name
    }

    if (isNew) {
        const created = await api.libraries.create(request)
        data.value.push(created)
        toast.success('Library created')
    } else {
        const updated = await api.libraries.update(item.id!, request)
        const index = data.value.findIndex(d => d.id === item.id)
        data.value[index] = updated
        toast.success('Library updated')
    }
}

async function handleDelete(id: number) {
    await api.libraries.delete(id)
    data.value = data.value.filter(d => d.id !== id)
    toast.success('Library deleted')
}
</script>

<template>
    <div class="mt-4">
        <h2>Libraries</h2>
        <p class="text-muted mb-4">Manage book storage locations</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="true" :on-save="handleSave"
            :on-delete="handleDelete" />
    </div>
</template>

<style scoped>
:deep(td) {
    word-break: break-all;
}
</style>
