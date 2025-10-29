<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { api } from '@/services/api'
import type { BookMetadata } from '@/types/api'

// Props
const props = defineProps<{
  selectedAsin: string
}>()

// Emits
const emit = defineEmits<{
  next: []
  'metadata-loaded': [metadata: BookMetadata]
}>()

// State
const metadata = ref<BookMetadata | null>(null)
const loading = ref(false)
const error = ref('')

// Computed
const runtimeDisplay = computed(() => {
  if (!metadata.value?.runtimeLengthMin) return 'Unknown'
  const hours = Math.floor(metadata.value.runtimeLengthMin / 60)
  const minutes = metadata.value.runtimeLengthMin % 60
  return `${hours}h ${minutes}m`
})

const releaseDateDisplay = computed(() => {
  if (!metadata.value?.releaseDate) return 'Unknown'
  try {
    return new Date(metadata.value.releaseDate).toLocaleDateString()
  } catch {
    return metadata.value.releaseDate
  }
})

// Load metadata when ASIN changes
const loadMetadata = async () => {
  if (!props.selectedAsin) return

  loading.value = true
  error.value = ''
  metadata.value = null

  try {
    metadata.value = await api.audiobookWizard.getMetadataByAsin(props.selectedAsin)
    if (metadata.value) {
      emit('metadata-loaded', metadata.value)
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load metadata'
  } finally {
    loading.value = false
  }
}

// Watch for ASIN changes
watch(() => props.selectedAsin, () => {
  loadMetadata()
}, { immediate: true })

// Confirm and proceed
const confirmMetadata = () => {
  if (metadata.value) {
    emit('next')
  }
}
</script>

<template>
  <div class="metadata-confirmation-step">
    <h5>Step 2: Metadata Confirmation</h5>
    <p class="text-muted mb-4">Review the audiobook metadata before importing.</p>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading metadata...</span>
      </div>
      <p class="mt-3">Loading metadata...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="alert alert-danger">
      <i class="bi bi-exclamation-triangle"></i> {{ error }}
    </div>

    <!-- Metadata Display -->
    <div v-else-if="metadata" class="metadata-display">
      <!-- Header with Cover Image -->
      <div class="row mb-4">
        <div class="col-md-3">
          <img
            v-if="metadata.image"
            :src="metadata.image"
            :alt="metadata.title"
            class="img-fluid rounded shadow-sm cover-image"
          />
          <div v-else class="no-cover-placeholder">
            <i class="bi bi-book"></i>
            <p>No Cover</p>
          </div>
        </div>
        <div class="col-md-9">
          <h4 class="mb-2">{{ metadata.title }}</h4>
          <h6 v-if="metadata.subtitle" class="text-muted mb-3">{{ metadata.subtitle }}</h6>

          <div class="metadata-field">
            <strong>ASIN:</strong>
            <span class="ms-2 font-monospace text-muted">{{ metadata.asin }}</span>
          </div>

          <div class="metadata-field">
            <strong>Authors:</strong>
            <span class="ms-2">{{ metadata.authors.map(a => a.name).join(', ') }}</span>
          </div>

          <div v-if="metadata.narrators.length > 0" class="metadata-field">
            <strong>Narrators:</strong>
            <span class="ms-2">{{ metadata.narrators.map(n => n.name).join(', ') }}</span>
          </div>

          <div class="metadata-field">
            <strong>Publisher:</strong>
            <span class="ms-2">{{ metadata.publisherName }}</span>
          </div>

          <div class="metadata-field">
            <strong>Release Date:</strong>
            <span class="ms-2">{{ releaseDateDisplay }}</span>
          </div>

          <div class="metadata-field">
            <strong>Runtime:</strong>
            <span class="ms-2">{{ runtimeDisplay }}</span>
          </div>

          <div class="metadata-field">
            <strong>Language:</strong>
            <span class="ms-2">{{ metadata.language }}</span>
          </div>

          <div v-if="metadata.rating" class="metadata-field">
            <strong>Rating:</strong>
            <span class="ms-2">{{ metadata.rating }}</span>
          </div>

          <div v-if="metadata.isbn" class="metadata-field">
            <strong>ISBN:</strong>
            <span class="ms-2">{{ metadata.isbn }}</span>
          </div>
        </div>
      </div>

      <!-- Series Information -->
      <div v-if="metadata.seriesPrimary || metadata.seriesSecondary" class="card mb-3">
        <div class="card-body">
          <h6 class="card-title">
            <i class="bi bi-collection"></i> Series Information
          </h6>
          <div v-if="metadata.seriesPrimary" class="mb-2">
            <strong>Primary Series:</strong>
            {{ metadata.seriesPrimary.name }}
            <span v-if="metadata.seriesPrimary.position" class="text-muted">
              (Book {{ metadata.seriesPrimary.position }})
            </span>
          </div>
          <div v-if="metadata.seriesSecondary">
            <strong>Secondary Series:</strong>
            {{ metadata.seriesSecondary.name }}
            <span v-if="metadata.seriesSecondary.position" class="text-muted">
              (Book {{ metadata.seriesSecondary.position }})
            </span>
          </div>
        </div>
      </div>

      <!-- Genres -->
      <div v-if="metadata.genres.length > 0" class="card mb-3">
        <div class="card-body">
          <h6 class="card-title">
            <i class="bi bi-tags"></i> Genres
          </h6>
          <div class="d-flex flex-wrap gap-2">
            <span
              v-for="genre in metadata.genres"
              :key="genre.asin"
              class="badge bg-secondary"
            >
              {{ genre.name }}
            </span>
          </div>
        </div>
      </div>

      <!-- Description -->
      <div class="card mb-3">
        <div class="card-body">
          <h6 class="card-title">
            <i class="bi bi-card-text"></i> Description
          </h6>
          <p class="card-text description-text">{{ metadata.summary || metadata.description }}</p>
        </div>
      </div>

      <!-- Additional Information -->
      <div class="card mb-3">
        <div class="card-body">
          <h6 class="card-title">
            <i class="bi bi-info-circle"></i> Additional Information
          </h6>
          <div class="row">
            <div class="col-md-6">
              <div v-if="metadata.copyright" class="metadata-field">
                <strong>Copyright:</strong>
                <span class="ms-2">{{ metadata.copyright }}</span>
              </div>
              <div v-if="metadata.literatureType" class="metadata-field">
                <strong>Type:</strong>
                <span class="ms-2">{{ metadata.literatureType }}</span>
              </div>
            </div>
            <div class="col-md-6">
              <div class="metadata-field">
                <strong>Region:</strong>
                <span class="ms-2">{{ metadata.region }}</span>
              </div>
              <div class="metadata-field">
                <strong>Format:</strong>
                <span class="ms-2">{{ metadata.formatType }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Confirm Button -->
      <div class="text-end">
        <button
          class="btn btn-primary btn-lg"
          @click="confirmMetadata"
        >
          <i class="bi bi-check-circle"></i> Confirm Metadata
        </button>
      </div>
    </div>

    <!-- No ASIN Selected -->
    <div v-else class="alert alert-warning">
      <i class="bi bi-exclamation-triangle"></i>
      No ASIN selected. Please complete Step 1 first.
    </div>
  </div>
</template>

<style scoped>
.metadata-confirmation-step {
  min-height: 400px;
}

.cover-image {
  width: 100%;
  max-width: 250px;
  border: 1px solid #dee2e6;
}

.no-cover-placeholder {
  width: 100%;
  max-width: 250px;
  aspect-ratio: 2/3;
  background-color: #f8f9fa;
  border: 2px dashed #dee2e6;
  border-radius: 0.25rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #6c757d;
}

.no-cover-placeholder i {
  font-size: 3rem;
  margin-bottom: 0.5rem;
}

.metadata-field {
  margin-bottom: 0.75rem;
  line-height: 1.5;
}

.description-text {
  white-space: pre-wrap;
  line-height: 1.6;
}

.metadata-display {
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
