<script setup lang="ts">
import { ref, computed } from 'vue'
import LoadingSpinner from './LoadingSpinner.vue'

export interface Column {
    key: string
    label: string
    editable?: boolean
    type?: 'text' | 'number' | 'select'
    options?: { value: number | string; label: string }[]
    displayKey?: string
}

interface Props {
    columns: Column[]
    data: any[]
    loading: boolean
    editable: boolean
    onSave: (item: any, isNew: boolean) => Promise<void>
    onDelete: (id: number) => Promise<void>
}

const props = withDefaults(defineProps<Props>(), {
    editable: true,
})

const editingId = ref<number | null>(null)
const editingData = ref<Record<string, any>>({})

const tableData = computed(() => {
    if (editingId.value === -1) {
        const newRow: Record<string, any> = { id: -1 }
        props.columns.forEach(col => {
            if (col.editable) {
                newRow[col.key] = col.type === 'number' ? 0 : ''
            }
        })
        return [newRow, ...props.data]
    }
    return props.data
})

function startEdit(item: any) {
    editingId.value = item.id
    editingData.value = { ...item }
}

function cancelEdit() {
    editingId.value = null
    editingData.value = {}
}

async function saveEdit() {
    const isNew = editingId.value === -1
    try {
        await props.onSave(editingData.value, isNew)
        editingId.value = null
        editingData.value = {}
    } catch (error) {
        console.error('Save failed:', error)
    }
}

async function deleteItem(id: number) {
    if (confirm('Are you sure you want to delete this item?')) {
        try {
            await props.onDelete(id)
        } catch (error) {
            console.error('Delete failed:', error)
        }
    }
}

function addNew() {
    editingId.value = -1
    editingData.value = { id: -1 }
    props.columns.forEach(col => {
        if (col.editable) {
            editingData.value[col.key] = col.type === 'number' ? 0 : ''
        }
    })
}

function getCellValue(item: any, column: Column): string {
    if (column.displayKey) {
        return item[column.displayKey] ?? ''
    }
    if (column.type === 'select' && column.options) {
        const option = column.options.find(opt => opt.value === item[column.key])
        return option?.label ?? ''
    }
    return item[column.key] ?? ''
}
</script>

<template>
    <div class="position-relative">
        <LoadingSpinner v-if="loading" />

        <div class="mb-3">
            <button class="btn btn-primary btn-sm" @click="addNew" :disabled="editingId !== null">
                <i class="bi bi-plus-lg me-1"></i>
                Add New
            </button>
        </div>

        <table class="table table-dark table-striped table-hover">
            <thead>
                <tr>
                    <th v-for="column in columns" :key="column.key">
                        {{ column.label }}
                    </th>
                    <th style="width: 120px">Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-if="tableData.length === 0">
                    <td :colspan="columns.length + 1" class="text-center text-muted py-4">
                        No data available
                    </td>
                </tr>
                <tr v-for="item in tableData" :key="item.id">
                    <td v-for="column in columns" :key="column.key">
                        <template v-if="editingId === item.id && column.editable">
                            <select v-if="column.type === 'select'" v-model="editingData[column.key]"
                                class="form-select form-select-sm">
                                <option v-for="option in column.options" :key="option.value" :value="option.value">
                                    {{ option.label }}
                                </option>
                            </select>
                            <input v-else-if="column.type === 'number'" type="number"
                                v-model.number="editingData[column.key]" class="form-control form-control-sm" />
                            <input v-else type="text" v-model="editingData[column.key]"
                                class="form-control form-control-sm" />
                        </template>
                        <template v-else>
                            {{ getCellValue(item, column) }}
                        </template>
                    </td>
                    <td>
                        <template v-if="editingId === item.id">
                            <button class="btn btn-success btn-sm me-1" @click="saveEdit" title="Save">
                                <i class="bi bi-check"></i>
                            </button>
                            <button class="btn btn-secondary btn-sm" @click="cancelEdit" title="Cancel">
                                <i class="bi bi-x"></i>
                            </button>
                        </template>
                        <template v-else-if="editable">
                            <button class="btn btn-primary btn-sm me-1" @click="startEdit(item)"
                                :disabled="editingId !== null" title="Edit">
                                <i class="bi bi-pencil"></i>
                            </button>
                            <button class="btn btn-danger btn-sm" @click="deleteItem(item.id)"
                                :disabled="editingId !== null" title="Delete">
                                <i class="bi bi-trash"></i>
                            </button>
                        </template>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>
