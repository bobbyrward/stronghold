<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import DataTable, { type Column } from '@/components/common/DataTable.vue'
import type { Notifier, NotificationType } from '@/types/api'

const toast = useToastStore()
const data = ref<Notifier[]>([])
const notificationTypes = ref<NotificationType[]>([])
const loading = ref(true)

const columns = computed<Column[]>(() => [
    { key: 'id', label: 'ID', editable: false },
    { key: 'name', label: 'Name', editable: true, type: 'text' },
    {
        key: 'type_id',
        label: 'Type',
        editable: true,
        type: 'select',
        displayKey: 'type_name',
        options: notificationTypes.value.map(t => ({ value: t.id, label: t.name }))
    },
    { key: 'url', label: 'URL', editable: true, type: 'text' }
])

onMounted(async () => {
    try {
        const [notifiers, types] = await Promise.all([
            api.notifiers.list(),
            api.notificationTypes.list()
        ])
        data.value = notifiers
        notificationTypes.value = types
    } catch (e) {
        toast.error('Failed to load notifiers')
    } finally {
        loading.value = false
    }
})

async function handleSave(item: Notifier, isNew: boolean) {
    if (!item.name || !item.type_id) {
        toast.error('Name and Type are required')
        throw new Error('Validation failed')
    }

    const request = {
        name: item.name,
        type_id: item.type_id,
        url: item.url || undefined
    }

    if (isNew) {
        const created = await api.notifiers.create(request)
        data.value.push(created)
        toast.success('Notifier created')
    } else {
        const updated = await api.notifiers.update(item.id, request)
        const index = data.value.findIndex(d => d.id === item.id)
        data.value[index] = updated
        toast.success('Notifier updated')
    }
}

async function handleDelete(id: number) {
    await api.notifiers.delete(id)
    data.value = data.value.filter(d => d.id !== id)
    toast.success('Notifier deleted')
}
</script>

<template>
    <div class="mt-4">
        <h2>Notifiers</h2>
        <p class="text-muted mb-4">Manage notification channels</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="true" :on-save="handleSave"
            :on-delete="handleDelete" />
    </div>
</template>

<style scoped>
:deep(td) {
    word-break: break-all;
}
</style>
