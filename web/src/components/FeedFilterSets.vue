<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import FeedFilterSetEntries from '@/components/FeedFilterSetEntries.vue'
import type { FeedFilterSet, FeedFilterSetType } from '@/types/api'

const props = defineProps<{
    feedFilterId: number
}>()

const toast = useToastStore()
const data = ref<FeedFilterSet[]>([])
const types = ref<FeedFilterSetType[]>([])
const loading = ref(true)

const expandedSetId = ref<number | null>(null)
const adding = ref(false)
const newSet = ref<Partial<FeedFilterSet>>({})
const editingId = ref<number | null>(null)
const editData = ref<Partial<FeedFilterSet>>({})

onMounted(async () => {
    try {
        const [setsData, typesData] = await Promise.all([
            api.feedFilterSets.list(props.feedFilterId),
            api.feedFilterSetTypes.list()
        ])
        data.value = setsData
        types.value = typesData
    } catch (e) {
        toast.error('Failed to load filter sets')
    } finally {
        loading.value = false
    }
})

function toggleExpand(id: number) {
    expandedSetId.value = expandedSetId.value === id ? null : id
}

function startAdd() {
    adding.value = true
    newSet.value = {
        feed_filter_id: props.feedFilterId,
        type_id: types.value[0]?.id
    }
}

function cancelAdd() {
    adding.value = false
    newSet.value = {}
}

async function saveNew() {
    if (!newSet.value.type_id) {
        toast.error('Type is required')
        return
    }

    try {
        const created = await api.feedFilterSets.create({
            feed_filter_id: props.feedFilterId,
            type_id: Number(newSet.value.type_id)
        })
        data.value.push(created)
        toast.success('Filter set created')
        adding.value = false
        newSet.value = {}
    } catch (e) {
        toast.error('Failed to create filter set')
    }
}

function startEdit(set: FeedFilterSet, e: Event) {
    e.stopPropagation()
    editingId.value = set.id
    editData.value = {
        id: set.id,
        type_id: set.type_id
    }
}

function cancelEdit(e: Event) {
    e.stopPropagation()
    editingId.value = null
    editData.value = {}
}

async function saveEdit(e: Event) {
    e.stopPropagation()
    if (!editData.value.type_id) {
        toast.error('Type is required')
        return
    }

    try {
        const updated = await api.feedFilterSets.update(editData.value.id!, {
            feed_filter_id: props.feedFilterId,
            type_id: Number(editData.value.type_id)
        })
        const index = data.value.findIndex(d => d.id === editData.value.id)
        data.value[index] = updated
        toast.success('Filter set updated')
        editingId.value = null
        editData.value = {}
    } catch (e) {
        toast.error('Failed to update filter set')
    }
}

async function deleteSet(id: number, e: Event) {
    e.stopPropagation()
    if (confirm('Are you sure you want to delete this filter set?')) {
        try {
            await api.feedFilterSets.delete(id)
            data.value = data.value.filter(d => d.id !== id)
            toast.success('Filter set deleted')
        } catch (e) {
            toast.error('Failed to delete filter set')
        }
    }
}
</script>

<template>
    <div class="feed-filter-sets">
        <div class="position-relative">
            <LoadingSpinner v-if="loading" />

            <div class="d-flex justify-content-between align-items-center mb-2">
                <h6 class="mb-0">Filter Sets</h6>
                <button class="btn btn-primary btn-sm" @click="startAdd" :disabled="adding">
                    <i class="bi bi-plus-lg me-1"></i>
                    Add Set
                </button>
            </div>

            <table class="table table-dark table-sm table-striped mb-0">
                <thead>
                    <tr>
                        <th style="width: 30px"></th>
                        <th>ID</th>
                        <th>Type</th>
                        <th style="width: 100px">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="adding">
                        <td></td>
                        <td>-</td>
                        <td>
                            <select v-model="newSet.type_id" class="form-select form-select-sm">
                                <option v-for="type in types" :key="type.id" :value="type.id">
                                    {{ type.name }}
                                </option>
                            </select>
                        </td>
                        <td>
                            <button class="btn btn-success btn-sm me-1" @click="saveNew" title="Save">
                                <i class="bi bi-check"></i>
                            </button>
                            <button class="btn btn-secondary btn-sm" @click="cancelAdd" title="Cancel">
                                <i class="bi bi-x"></i>
                            </button>
                        </td>
                    </tr>
                    <tr v-if="data.length === 0 && !adding">
                        <td colspan="4" class="text-center text-muted py-3">
                            No filter sets defined
                        </td>
                    </tr>
                    <template v-for="set in data" :key="set.id">
                        <tr @click="toggleExpand(set.id)" style="cursor: pointer;">
                            <td>
                                <i :class="expandedSetId === set.id ? 'bi bi-chevron-down' : 'bi bi-chevron-right'"></i>
                            </td>
                            <td>{{ set.id }}</td>
                            <td>
                                <select v-if="editingId === set.id" v-model="editData.type_id"
                                    class="form-select form-select-sm" @click.stop>
                                    <option v-for="type in types" :key="type.id" :value="type.id">
                                        {{ type.name }}
                                    </option>
                                </select>
                                <span v-else>{{ set.type_name }}</span>
                            </td>
                            <td>
                                <template v-if="editingId === set.id">
                                    <button class="btn btn-success btn-sm me-1" @click="saveEdit" title="Save">
                                        <i class="bi bi-check"></i>
                                    </button>
                                    <button class="btn btn-secondary btn-sm" @click="cancelEdit" title="Cancel">
                                        <i class="bi bi-x"></i>
                                    </button>
                                </template>
                                <template v-else>
                                    <button class="btn btn-primary btn-sm me-1" @click="startEdit(set, $event)"
                                        title="Edit">
                                        <i class="bi bi-pencil"></i>
                                    </button>
                                    <button class="btn btn-danger btn-sm" @click="deleteSet(set.id, $event)"
                                        title="Delete">
                                        <i class="bi bi-trash"></i>
                                    </button>
                                </template>
                            </td>
                        </tr>
                        <tr v-if="expandedSetId === set.id">
                            <td colspan="4">
                                <div class="ps-3 py-2">
                                    <FeedFilterSetEntries :feed-filter-set-id="set.id" />
                                </div>
                            </td>
                        </tr>
                    </template>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.feed-filter-sets {
    border-left: 3px solid var(--bs-primary);
    padding-left: 1rem;
    background-color: rgba(0, 0, 0, 0.2);
    padding: 1rem;
    border-radius: 0.25rem;
}
</style>
