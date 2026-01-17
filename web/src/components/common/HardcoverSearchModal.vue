<script setup lang="ts">
import { ref, watch } from 'vue'
import { api } from '@/services/api'
import type { HardcoverAuthorSearchResult } from '@/types/api'

const props = defineProps<{
    show: boolean
    initialQuery?: string
}>()

const emit = defineEmits<{
    close: []
    select: [author: HardcoverAuthorSearchResult]
}>()

const searchQuery = ref('')

// Initialize search query when modal opens
watch(() => props.show, (isShowing) => {
    if (isShowing && props.initialQuery) {
        searchQuery.value = props.initialQuery
    }
})
const results = ref<HardcoverAuthorSearchResult[]>([])
const loading = ref(false)
const error = ref('')
const hasSearched = ref(false)

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
                            <div class="input-group">
                                <input v-model="searchQuery" type="text" class="form-control"
                                    placeholder="Search for author..." autofocus @keyup.enter="performSearch">
                                <button class="btn btn-primary" type="button" @click="performSearch"
                                    :disabled="loading || !searchQuery.trim()">
                                    <i class="bi bi-search"></i>
                                </button>
                            </div>
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
                            <div v-for="author in results" :key="author.slug"
                                class="list-group-item list-group-item-action">

                                <div class="row">
                                    <div class="col-md-10">
                                        <div class="fw-bold">{{ author.name }}</div>
                                        <small class="text-muted">{{ author.slug }}</small>
                                    </div>

                                    <div class="col-md-2">
                                        <div class="row">
                                            <a :href="`https://hardcover.app/authors/${author.slug}`" target="_blank"
                                                class="col-sm-6 btn btn-info btn-sm">
                                                <i class="bi bi-box-arrow-up-right"></i>
                                            </a>
                                            <a class="col-sm-6 btn btn-success btn-sm" @click="selectAuthor(author)">
                                                <i class="bi bi-check-circle-fill"></i>
                                            </a>
                                        </div>
                                    </div>
                                </div>


                            </div>
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
