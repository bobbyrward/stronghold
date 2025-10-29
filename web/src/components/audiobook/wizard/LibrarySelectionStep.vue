<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { api } from '@/services/api'
import type { Library, BookMetadata, ExecuteImportResponse } from '@/types/api'

// Props
const props = defineProps<{
  torrentHash: string
  metadata: BookMetadata | null
  suggestedLibrary?: string
}>()

// Emits
const emit = defineEmits<{
  next: []
  'import-complete': [result: ExecuteImportResponse]
}>()

// State
const libraries = ref<Library[]>([])
const selectedLibrary = ref<string>('')
const directoryPreview = ref<string>('')
const loadingLibraries = ref(false)
const loadingPreview = ref(false)
const importing = ref(false)
const error = ref('')
const librariesError = ref('')
const previewError = ref('')

// Computed
const canImport = computed(() => {
  return selectedLibrary.value && props.metadata && !importing.value
})

const selectedLibraryPath = computed(() => {
  const library = libraries.value.find(lib => lib.name === selectedLibrary.value)
  return library?.path || ''
})

// Load libraries on mount
const loadLibraries = async () => {
  loadingLibraries.value = true
  librariesError.value = ''

  try {
    libraries.value = await api.audiobookWizard.getLibraries()

    // Pre-select suggested library if available
    if (props.suggestedLibrary) {
      const suggested = libraries.value.find(lib => lib.name === props.suggestedLibrary)
      if (suggested) {
        selectedLibrary.value = suggested.name
      }
    }

    // If no suggestion, select first library
    if (!selectedLibrary.value && libraries.value.length > 0) {
      const firstLibrary = libraries.value[0]
      if (firstLibrary) {
        selectedLibrary.value = firstLibrary.name
      }
    }
  } catch (err) {
    librariesError.value = err instanceof Error ? err.message : 'Failed to load libraries'
  } finally {
    loadingLibraries.value = false
  }
}

// Load directory preview when metadata is available
const loadPreview = async () => {
  if (!props.metadata) {
    directoryPreview.value = ''
    return
  }

  loadingPreview.value = true
  previewError.value = ''

  try {
    const response = await api.audiobookWizard.previewDirectory({
      metadata: props.metadata
    })
    directoryPreview.value = response.directory_name
  } catch (err) {
    previewError.value = err instanceof Error ? err.message : 'Failed to load preview'
    directoryPreview.value = ''
  } finally {
    loadingPreview.value = false
  }
}

// Execute import
const executeImport = async () => {
  if (!canImport.value || !props.metadata) return

  importing.value = true
  error.value = ''

  try {
    const result = await api.audiobookWizard.executeImport({
      hash: props.torrentHash,
      metadata: props.metadata,
      library_name: selectedLibrary.value
    })

    emit('import-complete', result)
    emit('next')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Import failed'
  } finally {
    importing.value = false
  }
}

// Watch for metadata changes to load preview
watch(() => props.metadata, () => {
  loadPreview()
}, { immediate: true })

// Load libraries on mount
watch(() => props.torrentHash, () => {
  if (props.torrentHash) {
    loadLibraries()
  }
}, { immediate: true })
</script>

<template>
  <div class="library-selection-step">
    <h5>Step 3: Library Selection & Import</h5>
    <p class="text-muted mb-4">Select the destination library and execute the import.</p>

    <!-- Error Alert -->
    <div v-if="error" class="alert alert-danger alert-dismissible fade show" role="alert">
      <i class="bi bi-exclamation-triangle"></i> <strong>Import Error:</strong> {{ error }}
      <button type="button" class="btn-close" @click="error = ''"></button>
    </div>

    <!-- Libraries Loading -->
    <div v-if="loadingLibraries" class="text-center py-4">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading libraries...</span>
      </div>
      <p class="mt-2">Loading libraries...</p>
    </div>

    <!-- Libraries Error -->
    <div v-else-if="librariesError" class="alert alert-danger">
      <i class="bi bi-exclamation-triangle"></i> {{ librariesError }}
    </div>

    <!-- No Libraries Available -->
    <div v-else-if="libraries.length === 0" class="alert alert-warning">
      <i class="bi bi-exclamation-triangle"></i>
      No Audiobookshelf libraries found. Please configure at least one library.
    </div>

    <!-- Library Selection -->
    <div v-else class="library-selection-content">
      <!-- Metadata Summary -->
      <div v-if="metadata" class="card mb-3">
        <div class="card-body">
          <h6 class="card-title">
            <i class="bi bi-book"></i> Import Summary
          </h6>
          <div class="row">
            <div class="col-md-3">
              <img
                v-if="metadata.image"
                :src="metadata.image"
                :alt="metadata.title"
                class="img-fluid rounded shadow-sm"
                style="max-width: 150px;"
              />
            </div>
            <div class="col-md-9">
              <h5 class="mb-1">{{ metadata.title }}</h5>
              <p class="text-muted mb-2">
                {{ metadata.authors.map(a => a.name).join(', ') }}
              </p>
              <p class="mb-1">
                <strong>ASIN:</strong> <span class="font-monospace">{{ metadata.asin }}</span>
              </p>
              <p v-if="metadata.seriesPrimary" class="mb-1">
                <strong>Series:</strong> {{ metadata.seriesPrimary.name }}
                <span v-if="metadata.seriesPrimary.position">
                  (Book {{ metadata.seriesPrimary.position }})
                </span>
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Library Selection Card -->
      <div class="card mb-3">
        <div class="card-body">
          <h6 class="card-title">
            <i class="bi bi-folder"></i> Select Destination Library
          </h6>

          <div class="mb-3">
            <label for="library-select" class="form-label">Library</label>
            <select
              id="library-select"
              v-model="selectedLibrary"
              class="form-select"
              :disabled="importing"
            >
              <option
                v-for="library in libraries"
                :key="library.name"
                :value="library.name"
              >
                {{ library.name }} ({{ library.path }})
              </option>
            </select>
            <div v-if="suggestedLibrary === selectedLibrary" class="form-text text-success">
              <i class="bi bi-check-circle"></i> This is the suggested library based on torrent metadata
            </div>
          </div>

          <!-- Directory Preview -->
          <div class="mb-3">
            <label class="form-label">Destination Directory</label>
            <div v-if="loadingPreview" class="text-muted">
              <span class="spinner-border spinner-border-sm me-2" role="status"></span>
              Generating preview...
            </div>
            <div v-else-if="previewError" class="alert alert-warning mb-0">
              <i class="bi bi-exclamation-triangle"></i> {{ previewError }}
            </div>
            <div v-else-if="directoryPreview" class="preview-box">
              <code>{{ selectedLibraryPath }}/{{ directoryPreview }}</code>
            </div>
            <div v-else class="text-muted">
              No preview available
            </div>
          </div>
        </div>
      </div>

      <!-- Import Button -->
      <div class="d-grid gap-2">
        <button
          class="btn btn-success btn-lg"
          :disabled="!canImport"
          @click="executeImport"
        >
          <span v-if="importing">
            <span class="spinner-border spinner-border-sm me-2" role="status"></span>
            Importing...
          </span>
          <span v-else>
            <i class="bi bi-download"></i> Execute Import
          </span>
        </button>
      </div>

      <!-- Import Info -->
      <div class="alert alert-info mt-3 mb-0">
        <i class="bi bi-info-circle"></i>
        <strong>Note:</strong> This will move the torrent files to the selected library and write metadata files.
        The torrent will remain active in qBittorrent.
      </div>
    </div>
  </div>
</template>

<style scoped>
.library-selection-step {
  min-height: 400px;
}

.preview-box {
  background-color: var(--bs-gray-900);
  border: 1px solid var(--bs-border-color);
  border-radius: 0.375rem;
  padding: 1rem;
  font-family: 'Courier New', monospace;
  word-break: break-all;
}

.preview-box code {
  color: var(--bs-success);
  font-size: 0.95rem;
}

.library-selection-content {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
