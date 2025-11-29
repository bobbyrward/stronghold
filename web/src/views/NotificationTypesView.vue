<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import DataTable, { type Column } from '@/components/common/DataTable.vue'
import type { NotificationType } from '@/types/api'

const toast = useToastStore()
const data = ref<NotificationType[]>([])
const loading = ref(true)

const columns: Column[] = [
    { key: 'id', label: 'ID', editable: false },
    { key: 'name', label: 'Name', editable: true, type: 'text' }
]

onMounted(async () => {
    try {
        data.value = await api.notificationTypes.list()
    } catch (e) {
        toast.error('Failed to load notification types')
    } finally {
        loading.value = false
    }
})

async function handleSave(item: NotificationType, isNew: boolean) {
    if (isNew) {
        const created = await api.notificationTypes.create({ name: item.name })
        data.value.push(created)
        toast.success('Notification type created')
    } else {
        const updated = await api.notificationTypes.update(item.id, { name: item.name })
        const index = data.value.findIndex(d => d.id === item.id)
        data.value[index] = updated
        toast.success('Notification type updated')
    }
}

async function handleDelete(id: number) {
    await api.notificationTypes.delete(id)
    data.value = data.value.filter(d => d.id !== id)
    toast.success('Notification type deleted')
}
</script>

<template>
    <div class="mt-4">
        <h2>Notification Types</h2>
        <p class="text-muted mb-4">Manage notification type reference data</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="true" :on-save="handleSave"
            :on-delete="handleDelete" />
    </div>
</template>
