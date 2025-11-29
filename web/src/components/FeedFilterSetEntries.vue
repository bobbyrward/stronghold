<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { FeedFilterSetEntry, FilterKey, FilterOperator } from '@/types/api'

const props = defineProps<{
    feedFilterSetId: number
}>()

const toast = useToastStore()
const data = ref<FeedFilterSetEntry[]>([])
const filterKeys = ref<FilterKey[]>([])
const filterOperators = ref<FilterOperator[]>([])
const loading = ref(true)

const adding = ref(false)
const newEntry = ref<Partial<FeedFilterSetEntry>>({})
const editingId = ref<number | null>(null)
const editData = ref<Partial<FeedFilterSetEntry>>({})

onMounted(async () => {
    try {
        const [entriesData, keysData, operatorsData] = await Promise.all([
            api.feedFilterSetEntries.list(props.feedFilterSetId),
            api.filterKeys.list(),
            api.filterOperators.list()
        ])
        data.value = entriesData
        filterKeys.value = keysData
        filterOperators.value = operatorsData
    } catch (e) {
        toast.error('Failed to load entries')
    } finally {
        loading.value = false
    }
})

function startAdd() {
    adding.value = true
    newEntry.value = {
        feed_filter_set_id: props.feedFilterSetId,
        key_id: filterKeys.value[0]?.id,
        operator_id: filterOperators.value[0]?.id,
        value: ''
    }
}

function cancelAdd() {
    adding.value = false
    newEntry.value = {}
}

async function saveNew() {
    if (!newEntry.value.key_id || !newEntry.value.operator_id || !newEntry.value.value) {
        toast.error('All fields are required')
        return
    }

    try {
        const created = await api.feedFilterSetEntries.create({
            feed_filter_set_id: props.feedFilterSetId,
            key_id: Number(newEntry.value.key_id),
            operator_id: Number(newEntry.value.operator_id),
            value: newEntry.value.value
        })
        data.value.push(created)
        toast.success('Entry created')
        adding.value = false
        newEntry.value = {}
    } catch (e) {
        toast.error('Failed to create entry')
    }
}

function startEdit(entry: FeedFilterSetEntry) {
    editingId.value = entry.id
    editData.value = {
        id: entry.id,
        key_id: entry.key_id,
        operator_id: entry.operator_id,
        value: entry.value
    }
}

function cancelEdit() {
    editingId.value = null
    editData.value = {}
}

async function saveEdit() {
    if (!editData.value.key_id || !editData.value.operator_id || !editData.value.value) {
        toast.error('All fields are required')
        return
    }

    try {
        const updated = await api.feedFilterSetEntries.update(editData.value.id!, {
            feed_filter_set_id: props.feedFilterSetId,
            key_id: Number(editData.value.key_id),
            operator_id: Number(editData.value.operator_id),
            value: editData.value.value
        })
        const index = data.value.findIndex(d => d.id === editData.value.id)
        data.value[index] = updated
        toast.success('Entry updated')
        editingId.value = null
        editData.value = {}
    } catch (e) {
        toast.error('Failed to update entry')
    }
}

async function deleteEntry(id: number) {
    if (confirm('Are you sure you want to delete this entry?')) {
        try {
            await api.feedFilterSetEntries.delete(id)
            data.value = data.value.filter(d => d.id !== id)
            toast.success('Entry deleted')
        } catch (e) {
            toast.error('Failed to delete entry')
        }
    }
}
</script>

<template>
    <div class="feed-filter-set-entries">
        <div class="position-relative">
            <LoadingSpinner v-if="loading" />

            <div class="d-flex justify-content-between align-items-center mb-2">
                <h6 class="mb-0 small">Entries</h6>
                <button class="btn btn-secondary btn-sm" @click="startAdd" :disabled="adding">
                    <i class="bi bi-plus-lg me-1"></i>
                    Add Entry
                </button>
            </div>

            <table class="table table-dark table-sm mb-0">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Key</th>
                        <th>Operator</th>
                        <th>Value</th>
                        <th style="width: 100px">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="adding">
                        <td>-</td>
                        <td>
                            <select v-model="newEntry.key_id" class="form-select form-select-sm">
                                <option v-for="key in filterKeys" :key="key.id" :value="key.id">
                                    {{ key.name }}
                                </option>
                            </select>
                        </td>
                        <td>
                            <select v-model="newEntry.operator_id" class="form-select form-select-sm">
                                <option v-for="op in filterOperators" :key="op.id" :value="op.id">
                                    {{ op.name }}
                                </option>
                            </select>
                        </td>
                        <td>
                            <input type="text" v-model="newEntry.value" class="form-control form-control-sm"
                                placeholder="Value" />
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
                        <td colspan="5" class="text-center text-muted py-3">
                            No entries defined
                        </td>
                    </tr>
                    <tr v-for="entry in data" :key="entry.id">
                        <td>{{ entry.id }}</td>
                        <td>
                            <select v-if="editingId === entry.id" v-model="editData.key_id"
                                class="form-select form-select-sm">
                                <option v-for="key in filterKeys" :key="key.id" :value="key.id">
                                    {{ key.name }}
                                </option>
                            </select>
                            <span v-else>{{ entry.key_name }}</span>
                        </td>
                        <td>
                            <select v-if="editingId === entry.id" v-model="editData.operator_id"
                                class="form-select form-select-sm">
                                <option v-for="op in filterOperators" :key="op.id" :value="op.id">
                                    {{ op.name }}
                                </option>
                            </select>
                            <span v-else>{{ entry.operator_name }}</span>
                        </td>
                        <td>
                            <input v-if="editingId === entry.id" type="text" v-model="editData.value"
                                class="form-control form-control-sm" />
                            <span v-else>{{ entry.value }}</span>
                        </td>
                        <td>
                            <template v-if="editingId === entry.id">
                                <button class="btn btn-success btn-sm me-1" @click="saveEdit" title="Save">
                                    <i class="bi bi-check"></i>
                                </button>
                                <button class="btn btn-secondary btn-sm" @click="cancelEdit" title="Cancel">
                                    <i class="bi bi-x"></i>
                                </button>
                            </template>
                            <template v-else>
                                <button class="btn btn-primary btn-sm me-1" @click="startEdit(entry)" title="Edit">
                                    <i class="bi bi-pencil"></i>
                                </button>
                                <button class="btn btn-danger btn-sm" @click="deleteEntry(entry.id)" title="Delete">
                                    <i class="bi bi-trash"></i>
                                </button>
                            </template>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<style scoped>
.feed-filter-set-entries {
    border-left: 3px solid var(--bs-secondary);
    padding-left: 1rem;
    background-color: rgba(0, 0, 0, 0.3);
    padding: 0.75rem;
    border-radius: 0.25rem;
}
</style>
