<script setup lang="ts">
// import { ref, computed } from 'vue'
import LoadingSpinner from './common/LoadingSpinner.vue'

export interface Column {
    key: string
    label: string
    editable?: boolean
    type?: 'text' | 'number' | 'select'
    options?: { value: number | string; label: string }[]
    displayKey?: string
}

const columns: Column[] = [
    { key: 'hash', label: 'Hash', editable: false },
    { key: 'name', label: 'Name', editable: false },
    { key: 'category', label: 'Category', editable: false },
    { key: 'state', label: 'State', editable: false },
    { key: 'tags', label: 'Tags', editable: false },
]

interface Props {
    data: any[]
    loading: boolean

    onChangeCategory: (hash: string) => Promise<void>
    onChangeTags: (hash: string) => Promise<void>
    onImportAudiobook?: (hash: string) => Promise<void>
}

const props = withDefaults(defineProps<Props>(), {})

async function changeTorrentCategory(hash: string) {
    try {
        await props.onChangeCategory(hash)
    } catch (error) {
        console.error('Change category failed:', error)
    }
}

async function changeTorrentTags(hash: string) {
    try {
        await props.onChangeTags(hash)
    } catch (error) {
        console.error('Change tags failed:', error)
    }
}

async function importAudiobook(hash: string) {
    if (!props.onImportAudiobook) return

    try {
        await props.onImportAudiobook(hash)
    } catch (error) {
        console.error('Import audiobook failed:', error)
    }
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
                <tr v-if="props.data.length === 0">
                    <td :colspan="columns.length + 1" class="text-center text-muted py-4">
                        No data available
                    </td>
                </tr>
                <tr v-for="item in props.data" :key="item.hash">
                    <td v-for="column in columns" :key="column.key">
                        {{ getCellValue(item, column) }}
                    </td>
                    <td>
                        <button class="btn btn-primary btn-sm me-1" @click="changeTorrentCategory(item.hash)"
                            title="Change Category">
                            <i class="bi bi-collection"></i>
                        </button>
                        <button class="btn btn-primary btn-sm me-1" @click="changeTorrentTags(item.hash)"
                            title="Change Tags">
                            <i class="bi bi-tags"></i>
                        </button>
                        <button
                            v-if="onImportAudiobook"
                            class="btn btn-success btn-sm"
                            @click="importAudiobook(item.hash)"
                            title="Import Audiobook">
                            <i class="bi bi-book"></i>
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>
