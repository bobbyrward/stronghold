<script setup lang="ts">
import { ref, watch } from 'vue'
import { api } from '@/services/api'
import type { HardcoverAuthorSearchResult } from '@/types/api'

defineProps<{
    show: boolean
}>()

const emit = defineEmits<{
    close: []
    select: [author: HardcoverAuthorSearchResult]
}>()

const searchQuery = ref('')
const results = ref<HardcoverAuthorSearchResult[]>([])
const loading = ref(false)
const error = ref('')
const hasSearched = ref(false)

let debounceTimeout: ReturnType<typeof setTimeout> | null = null

function handleSearchInput() {
    if (debounceTimeout) {
        clearTimeout(debounceTimeout)
    }

    if (!searchQuery.value.trim()) {
        results.value = []
        hasSearched.value = false
        return
    }

    debounceTimeout = setTimeout(async () => {
        await performSearch()
    }, 300)
}

async function performSearch() {
    if (!searchQuery.value.trim()) return

    loading.value = true
    error.value = ''

    try {
        results.value = await api.hardcover.searchAuthors(searchQuery.value.trim())
        hasSearched.value = true
    } catch (e) {
        error.value = 'Failed to search Hardcover'
        results.value = []
    } finally {
        loading.value = false
    }
}

function selectAuthor(author: HardcoverAuthorSearchResult) {
    emit('select', author)
    closeModal()
}

function closeModal() {
    searchQuery.value = ''
    results.value = []
    error.value = ''
    hasSearched.value = false
    emit('close')
}

watch(() => searchQuery.value, handleSearchInput)
</script>

<template>
    <Teleport to="body">
        <div v-if="show" class="modal-backdrop fade show"></div>
        <div v-if="show" class="modal fade show d-block" tabindex="-1" @click.self="closeModal">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">Search Hardcover</h5>
                        <button type="button" class="btn-close" aria-label="Close" @click="closeModal"></button>
                    </div>
                    <div class="modal-body">
                        <div class="mb-3">
                            <input v-model="searchQuery" type="text" class="form-control"
                                placeholder="Search for author..." autofocus>
                        </div>

                        <div v-if="loading" class="text-center py-3">
                            <div class="spinner-border spinner-border-sm text-primary" role="status">
                                <span class="visually-hidden">Searching...</span>
                            </div>
                            <span class="ms-2">Searching...</span>
                        </div>

                        <div v-else-if="error" class="alert alert-danger mb-0">
                            {{ error }}
                        </div>

                        <div v-else-if="hasSearched && results.length === 0" class="text-muted text-center py-3">
                            No authors found
                        </div>

                        <div v-else-if="results.length > 0" class="list-group">
                            <button v-for="author in results" :key="author.slug" type="button"
                                class="list-group-item list-group-item-action" @click="selectAuthor(author)">
                                <div class="fw-bold">{{ author.name }}</div>
                                <small class="text-muted">{{ author.slug }}</small>
                            </button>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" @click="closeModal">
                            Cancel
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Teleport>
</template>
