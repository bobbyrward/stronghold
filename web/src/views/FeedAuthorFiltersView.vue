<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import DataTable, { type Column } from '@/components/common/DataTable.vue'
import type { FeedAuthorFilter, Feed, TorrentCategory, Notifier } from '@/types/api'

const toast = useToastStore()
const data = ref<FeedAuthorFilter[]>([])
const feeds = ref<Feed[]>([])
const categories = ref<TorrentCategory[]>([])
const notifiers = ref<Notifier[]>([])
const loading = ref(true)

const columns = computed<Column[]>(() => [
    { key: 'id', label: 'ID', editable: false },
    { key: 'author', label: 'Author', editable: true, type: 'text' },
    {
        key: 'feed_id',
        label: 'Feed',
        editable: true,
        type: 'select',
        displayKey: 'feed_name',
        options: feeds.value.map(f => ({ value: f.id, label: f.name }))
    },
    {
        key: 'category_id',
        label: 'Category',
        editable: true,
        type: 'select',
        displayKey: 'category_name',
        options: categories.value.map(c => ({ value: c.id, label: c.name }))
    },
    {
        key: 'notifier_id',
        label: 'Notifier',
        editable: true,
        type: 'select',
        displayKey: 'notifier_name',
        options: notifiers.value.map(n => ({ value: n.id, label: n.name }))
    }
])

onMounted(async () => {
    try {
        const [feedsData, categoriesData, notifiersData, filtersData] = await Promise.all([
            api.feeds.list(),
            api.torrentCategories.list(),
            api.notifiers.list(),
            api.feedAuthorFilters.list()
        ])
        feeds.value = feedsData
        categories.value = categoriesData
        notifiers.value = notifiersData
        data.value = filtersData
    } catch (e) {
        toast.error('Failed to load data')
    } finally {
        loading.value = false
    }
})

async function handleSave(item: FeedAuthorFilter, isNew: boolean) {
    if (!item.author || !item.feed_id || !item.category_id || !item.notifier_id) {
        toast.error('All fields are required')
        throw new Error('Validation failed')
    }

    const request = {
        author: item.author,
        feed_id: Number(item.feed_id),
        category_id: Number(item.category_id),
        notifier_id: Number(item.notifier_id)
    }

    if (isNew) {
        const created = await api.feedAuthorFilters.create(request)
        data.value.push(created)
        toast.success('Feed author filter created')
    } else {
        const updated = await api.feedAuthorFilters.update(item.id, request)
        const index = data.value.findIndex(d => d.id === item.id)
        data.value[index] = updated
        toast.success('Feed author filter updated')
    }
}

async function handleDelete(id: number) {
    await api.feedAuthorFilters.delete(id)
    data.value = data.value.filter(d => d.id !== id)
    toast.success('Feed author filter deleted')
}
</script>

<template>
    <div class="mt-4">
        <h2>Feed Author Filters</h2>
        <p class="text-muted mb-4">Manage feed author filters</p>

        <DataTable :columns="columns" :data="data" :loading="loading" :editable="true" :on-save="handleSave"
            :on-delete="handleDelete" />
    </div>
</template>
