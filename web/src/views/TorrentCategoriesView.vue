<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import DataTable, { type Column } from '@/components/common/DataTable.vue'
import type { TorrentCategory } from '@/types/api'

const toast = useToastStore()
const data = ref<TorrentCategory[]>([])
const loading = ref(true)

const columns: Column[] = [
    { key: 'id', label: 'ID', editable: false },
    { key: 'name', label: 'Name', editable: true, type: 'text' }
]

onMounted(async () => {
    try {
        data.value = await api.torrentCategories.list()
    } catch (e) {
        toast.error('Failed to load torrent categories')
    } finally {
        loading.value = false
    }
})

async function handleSave(item: TorrentCategory, isNew: boolean) {
    if (isNew) {
        const created = await api.torrentCategories.create({ name: item.name })
        data.value.push(created)
        toast.success('Torrent category created')
    } else {
        const updated = await api.torrentCategories.update(item.id, { name: item.name })
        const index = data.value.findIndex(d => d.id === item.id)
        data.value[index] = updated
        toast.success('Torrent category updated')
    }
}

async function handleDelete(id: number) {
    await api.torrentCategories.delete(id)
    data.value = data.value.filter(d => d.id !== id)
    toast.success('Torrent category deleted')
}
</script>

<template>
    <div class="mt-4">
        <h2>Torrent Categories</h2>
        <p class="text-muted mb-4">Manage torrent category reference data</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="true" :on-save="handleSave"
            :on-delete="handleDelete" />
    </div>
</template>
