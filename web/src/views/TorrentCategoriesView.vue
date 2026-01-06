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
    { key: 'name', label: 'Name', editable: false },
    { key: 'scope_name', label: 'Scope', editable: false },
    { key: 'media_type', label: 'Media Type', editable: false }
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
</script>

<template>
    <div class="mt-4">
        <h2>Torrent Categories</h2>
        <p class="text-muted mb-4">Reference data for torrent categories (read-only)</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="false" />
    </div>
</template>
