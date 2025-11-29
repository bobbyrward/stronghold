<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import DataTable, { type Column } from '@/components/common/DataTable.vue'
import type { Feed } from '@/types/api'

const toast = useToastStore()
const data = ref<Feed[]>([])
const loading = ref(true)

const columns: Column[] = [
    { key: 'id', label: 'ID', editable: false },
    { key: 'name', label: 'Name', editable: true, type: 'text' },
    { key: 'url', label: 'URL', editable: true, type: 'text' }
]

onMounted(async () => {
    try {
        data.value = await api.feeds.list()
    } catch (e) {
        toast.error('Failed to load feeds')
    } finally {
        loading.value = false
    }
})

async function handleSave(item: Feed, isNew: boolean) {
    if (!item.name || !item.url) {
        toast.error('Name and URL are required')
        throw new Error('Validation failed')
    }

    if (isNew) {
        const created = await api.feeds.create({ name: item.name, url: item.url })
        data.value.push(created)
        toast.success('Feed created')
    } else {
        const updated = await api.feeds.update(item.id, { name: item.name, url: item.url })
        const index = data.value.findIndex(d => d.id === item.id)
        data.value[index] = updated
        toast.success('Feed updated')
    }
}

async function handleDelete(id: number) {
    await api.feeds.delete(id)
    data.value = data.value.filter(d => d.id !== id)
    toast.success('Feed deleted')
}
</script>

<template>
    <div class="mt-4">
        <h2>Feeds</h2>
        <p class="text-muted mb-4">Manage RSS/torrent feeds</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="true" :on-save="handleSave"
            :on-delete="handleDelete" />
    </div>
</template>

<style scoped>
:deep(td) {
    word-break: break-all;
}
</style>
