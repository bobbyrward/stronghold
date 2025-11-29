<script setup lang="ts">
import { ref } from 'vue'
import FeedFilterSets from '@/components/FeedFilterSets.vue'
import type { FeedFilter, Feed, TorrentCategory, Notifier } from '@/types/api'

const props = defineProps<{
    filter: FeedFilter
    feeds: Feed[]
    categories: TorrentCategory[]
    notifiers: Notifier[]
}>()

const emit = defineEmits<{
    save: [filter: FeedFilter]
    delete: [id: number]
}>()

const expanded = ref(false)
const editing = ref(false)
const editData = ref<Partial<FeedFilter>>({})

function toggleExpand() {
    if (!editing.value) {
        expanded.value = !expanded.value
    }
}

function startEdit(e: Event) {
    e.stopPropagation()
    editing.value = true
    editData.value = {
        id: props.filter.id,
        name: props.filter.name,
        feed_id: props.filter.feed_id,
        category_id: props.filter.category_id,
        notifier_id: props.filter.notifier_id
    }
}

function cancelEdit(e: Event) {
    e.stopPropagation()
    editing.value = false
    editData.value = {}
}

function saveEdit(e: Event) {
    e.stopPropagation()
    emit('save', editData.value as FeedFilter)
    editing.value = false
    editData.value = {}
}

function deleteFilter(e: Event) {
    e.stopPropagation()
    if (confirm('Are you sure you want to delete this filter?')) {
        emit('delete', props.filter.id)
    }
}
</script>

<template>
    <tr @click="toggleExpand" style="cursor: pointer;">
        <td>
            <i :class="expanded ? 'bi bi-chevron-down' : 'bi bi-chevron-right'"></i>
        </td>
        <td>{{ filter.id }}</td>
        <td>
            <input v-if="editing" type="text" v-model="editData.name" class="form-control form-control-sm"
                @click.stop />
            <span v-else>{{ filter.name }}</span>
        </td>
        <td>
            <select v-if="editing" v-model="editData.feed_id" class="form-select form-select-sm" @click.stop>
                <option v-for="feed in feeds" :key="feed.id" :value="feed.id">
                    {{ feed.name }}
                </option>
            </select>
            <span v-else>{{ filter.feed_name }}</span>
        </td>
        <td>
            <select v-if="editing" v-model="editData.category_id" class="form-select form-select-sm" @click.stop>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                    {{ cat.name }}
                </option>
            </select>
            <span v-else>{{ filter.category_name }}</span>
        </td>
        <td>
            <select v-if="editing" v-model="editData.notifier_id" class="form-select form-select-sm" @click.stop>
                <option v-for="notifier in notifiers" :key="notifier.id" :value="notifier.id">
                    {{ notifier.name }}
                </option>
            </select>
            <span v-else>{{ filter.notifier_name }}</span>
        </td>
        <td>
            <template v-if="editing">
                <button class="btn btn-success btn-sm me-1" @click="saveEdit" title="Save">
                    <i class="bi bi-check"></i>
                </button>
                <button class="btn btn-secondary btn-sm" @click="cancelEdit" title="Cancel">
                    <i class="bi bi-x"></i>
                </button>
            </template>
            <template v-else>
                <button class="btn btn-primary btn-sm me-1" @click="startEdit" title="Edit">
                    <i class="bi bi-pencil"></i>
                </button>
                <button class="btn btn-danger btn-sm" @click="deleteFilter" title="Delete">
                    <i class="bi bi-trash"></i>
                </button>
            </template>
        </td>
    </tr>
    <tr v-if="expanded">
        <td colspan="7">
            <div class="ps-4 py-2">
                <FeedFilterSets :feed-filter-id="filter.id" />
            </div>
        </td>
    </tr>
</template>
