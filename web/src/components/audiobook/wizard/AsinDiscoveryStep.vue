<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { api } from '@/services/api'
import type { TorrentImportInfo, BookMetadata } from '@/types/api'

// Props
const props = defineProps<{
  torrentInfo: TorrentImportInfo | null
  modelValue: string
}>()

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: string]
  next: []
}>()

// State
const activeTab = ref<'metadata' | 'manual' | 'search'>('metadata')
const manualAsin = ref('')
const manualAsinError = ref('')
const searchTitle = ref('')
const searchAuthor = ref('')
const searchResults = ref<BookMetadata[]>([])
const searchLoading = ref(false)
const searchError = ref('')

// Computed
const hasExtractedAsin = computed(() => {
  return !!props.torrentInfo?.asin
})

const hasExtractedMetadata = computed(() => {
  return !!(props.torrentInfo?.title || props.torrentInfo?.author)
})

// Auto-select appropriate tab on mount
const selectDefaultTab = () => {
  if (hasExtractedAsin.value) {
    activeTab.value = 'metadata'
  } else if (hasExtractedMetadata.value) {
    activeTab.value = 'search'
    // Pre-populate search fields
    if (props.torrentInfo?.title) searchTitle.value = props.torrentInfo.title
    if (props.torrentInfo?.author) searchAuthor.value = props.torrentInfo.author
  } else {
    activeTab.value = 'manual'
  }
}

// Watch for torrentInfo changes
watch(() => props.torrentInfo, () => {
  selectDefaultTab()
}, { immediate: true })

// Validate ASIN format (10 characters, alphanumeric starting with B)
const validateAsin = (asin: string): boolean => {
  const asinPattern = /^B[0-9A-Z]{9}$/
  return asinPattern.test(asin.toUpperCase())
}

// From Metadata Tab - Use extracted ASIN
const useExtractedAsin = () => {
  if (props.torrentInfo?.asin) {
    emit('update:modelValue', props.torrentInfo.asin)
    emit('next')
  }
}

// Manual Entry Tab - Continue with manual ASIN
const continueWithManualAsin = () => {
  manualAsinError.value = ''
  const asin = manualAsin.value.trim().toUpperCase()

  if (!asin) {
    manualAsinError.value = 'Please enter an ASIN'
    return
  }

  if (!validateAsin(asin)) {
    manualAsinError.value = 'Invalid ASIN format. Must be 10 characters starting with B (e.g., B09XYZ1234)'
    return
  }

  emit('update:modelValue', asin)
  emit('next')
}

// Search Tab - Search for books
const searchForBooks = async () => {
  searchError.value = ''
  searchResults.value = []

  if (!searchTitle.value.trim()) {
    searchError.value = 'Please enter a title to search'
    return
  }

  searchLoading.value = true

  try {
    const results = await api.audiobookWizard.searchAsin({
      title: searchTitle.value.trim(),
      author: searchAuthor.value.trim()
    })

    searchResults.value = results

    if (results.length === 0) {
      searchError.value = 'No results found. Try adjusting your search terms.'
    }
  } catch (err) {
    searchError.value = err instanceof Error ? err.message : 'Search failed. Please try again.'
  } finally {
    searchLoading.value = false
  }
}

// Select a book from search results
const selectBook = (book: BookMetadata) => {
  emit('update:modelValue', book.asin)
  emit('next')
}
</script>

<template>
  <div class="asin-discovery-step">
    <h5>Step 1: ASIN Discovery</h5>
    <p class="text-muted mb-4">Select or search for the correct audiobook ASIN.</p>

    <!-- Tabs -->
    <ul class="nav nav-tabs mb-3" role="tablist">
      <li class="nav-item" role="presentation">
        <button
          class="nav-link"
          :class="{ active: activeTab === 'metadata' }"
          @click="activeTab = 'metadata'"
          type="button"
        >
          <i class="bi bi-file-earmark-text"></i> From Metadata
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button
          class="nav-link"
          :class="{ active: activeTab === 'manual' }"
          @click="activeTab = 'manual'"
          type="button"
        >
          <i class="bi bi-pencil"></i> Manual Entry
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button
          class="nav-link"
          :class="{ active: activeTab === 'search' }"
          @click="activeTab = 'search'"
          type="button"
        >
          <i class="bi bi-search"></i> Search
        </button>
      </li>
    </ul>

    <!-- Tab Content -->
    <div class="tab-content">
      <!-- From Metadata Tab -->
      <div v-if="activeTab === 'metadata'" class="tab-pane-content">
        <div v-if="hasExtractedAsin" class="card">
          <div class="card-body">
            <h6 class="card-title">
              <i class="bi bi-check-circle text-success"></i> ASIN Found in Metadata
            </h6>
            <p class="card-text">
              An ASIN was automatically extracted from the torrent metadata:
            </p>
            <div class="alert alert-info mb-3">
              <strong>ASIN:</strong> {{ torrentInfo?.asin }}
            </div>
            <button
              class="btn btn-primary"
              @click="useExtractedAsin"
            >
              <i class="bi bi-arrow-right-circle"></i> Use This ASIN
            </button>
          </div>
        </div>
        <div v-else class="alert alert-warning">
          <i class="bi bi-exclamation-triangle"></i>
          No ASIN found in the torrent metadata. Please use Manual Entry or Search.
        </div>
      </div>

      <!-- Manual Entry Tab -->
      <div v-if="activeTab === 'manual'" class="tab-pane-content">
        <div class="card">
          <div class="card-body">
            <h6 class="card-title">Enter ASIN Manually</h6>
            <p class="text-muted">
              Enter the Audible ASIN (10 characters, starts with B).
            </p>
            <div class="mb-3">
              <label for="manual-asin" class="form-label">ASIN</label>
              <input
                id="manual-asin"
                v-model="manualAsin"
                type="text"
                class="form-control"
                :class="{ 'is-invalid': manualAsinError }"
                placeholder="e.g., B09XYZ1234"
                maxlength="10"
                @keyup.enter="continueWithManualAsin"
              />
              <div v-if="manualAsinError" class="invalid-feedback">
                {{ manualAsinError }}
              </div>
              <div class="form-text">
                Example: B09XYZ1234
              </div>
            </div>
            <button
              class="btn btn-primary"
              @click="continueWithManualAsin"
            >
              <i class="bi bi-arrow-right-circle"></i> Continue
            </button>
          </div>
        </div>
      </div>

      <!-- Search Tab -->
      <div v-if="activeTab === 'search'" class="tab-pane-content">
        <div class="card mb-3">
          <div class="card-body">
            <h6 class="card-title">Search Audible</h6>
            <p class="text-muted">
              Search for the audiobook by title and author.
            </p>

            <div class="row mb-3">
              <div class="col-md-8">
                <label for="search-title" class="form-label">Title <span class="text-danger">*</span></label>
                <input
                  id="search-title"
                  v-model="searchTitle"
                  type="text"
                  class="form-control"
                  placeholder="Enter book title"
                  @keyup.enter="searchForBooks"
                />
              </div>
              <div class="col-md-4">
                <label for="search-author" class="form-label">Author</label>
                <input
                  id="search-author"
                  v-model="searchAuthor"
                  type="text"
                  class="form-control"
                  placeholder="Enter author name"
                  @keyup.enter="searchForBooks"
                />
              </div>
            </div>

            <button
              class="btn btn-primary"
              @click="searchForBooks"
              :disabled="searchLoading || !searchTitle.trim()"
            >
              <span v-if="searchLoading">
                <span class="spinner-border spinner-border-sm me-2" role="status"></span>
                Searching...
              </span>
              <span v-else>
                <i class="bi bi-search"></i> Search
              </span>
            </button>
          </div>
        </div>

        <!-- Search Error -->
        <div v-if="searchError" class="alert alert-danger">
          <i class="bi bi-exclamation-triangle"></i> {{ searchError }}
        </div>

        <!-- Search Results -->
        <div v-if="searchResults.length > 0" class="search-results">
          <h6 class="mb-3">Search Results ({{ searchResults.length }})</h6>
          <div class="row g-3">
            <div
              v-for="book in searchResults"
              :key="book.asin"
              class="col-md-6"
            >
              <div class="card h-100 book-result-card">
                <div class="row g-0">
                  <div v-if="book.image" class="col-4">
                    <img
                      :src="book.image"
                      class="img-fluid rounded-start book-cover"
                      :alt="book.title"
                    />
                  </div>
                  <div :class="book.image ? 'col-8' : 'col-12'">
                    <div class="card-body">
                      <h6 class="card-title">{{ book.title }}</h6>
                      <p v-if="book.subtitle" class="card-subtitle text-muted small mb-2">
                        {{ book.subtitle }}
                      </p>
                      <p class="card-text small">
                        <strong>By:</strong>
                        {{ book.authors.map(a => a.name).join(', ') }}
                      </p>
                      <p v-if="book.narrators.length > 0" class="card-text small">
                        <strong>Narrated by:</strong>
                        {{ book.narrators.map(n => n.name).join(', ') }}
                      </p>
                      <p v-if="book.seriesPrimary" class="card-text small">
                        <strong>Series:</strong>
                        {{ book.seriesPrimary.name }}
                        <span v-if="book.seriesPrimary.position">
                          (Book {{ book.seriesPrimary.position }})
                        </span>
                      </p>
                      <p class="card-text small text-muted">
                        <span v-if="book.runtimeLengthMin">
                          {{ Math.floor(book.runtimeLengthMin / 60) }}h {{ book.runtimeLengthMin % 60 }}m
                        </span>
                        <span v-if="book.releaseDate">
                          • {{ new Date(book.releaseDate).getFullYear() }}
                        </span>
                      </p>
                      <button
                        class="btn btn-sm btn-success mt-2"
                        @click="selectBook(book)"
                      >
                        <i class="bi bi-check-circle"></i> Select
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.asin-discovery-step {
  min-height: 400px;
}

.tab-pane-content {
  padding: 1rem 0;
}

.book-result-card {
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.book-result-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.book-cover {
  object-fit: cover;
  max-height: 200px;
  width: 100%;
}

.search-results {
  max-height: 500px;
  overflow-y: auto;
}
</style>
