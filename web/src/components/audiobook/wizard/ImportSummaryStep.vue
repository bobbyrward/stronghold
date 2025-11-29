<script setup lang="ts">
import { computed } from 'vue'
import type { ExecuteImportResponse, BookMetadata } from '@/types/api'

// Props
const props = defineProps<{
  importResult: ExecuteImportResponse | null
  metadata: BookMetadata | null
}>()

// Emits
const emit = defineEmits<{
  finish: []
}>()

// Computed
const isSuccess = computed(() => {
  return props.importResult?.success === true
})

const hasResult = computed(() => {
  return props.importResult !== null
})

const destinationPath = computed(() => {
  return props.importResult?.destination_path || 'Unknown'
})

const message = computed(() => {
  return props.importResult?.message || ''
})

// Format path for display
const formatPath = (path: string) => {
  // Split path and highlight the last directory (the book folder)
  const parts = path.split('/')
  if (parts.length > 1) {
    const bookFolder = parts[parts.length - 1]
    const parentPath = parts.slice(0, -1).join('/')
    return { parentPath, bookFolder }
  }
  return { parentPath: '', bookFolder: path }
}

const pathParts = computed(() => {
  return formatPath(destinationPath.value)
})

// Finish and close wizard
const handleFinish = () => {
  emit('finish')
}
</script>

<template>
  <div class="import-summary-step">
    <h5>Step 5: Import Summary</h5>
    <p class="text-muted mb-4">Review the import results.</p>

    <!-- No Result -->
    <div v-if="!hasResult" class="alert alert-warning">
      <i class="bi bi-exclamation-triangle"></i>
      No import result available. The import may not have been executed.
    </div>

    <!-- Import Success -->
    <div v-else-if="isSuccess" class="summary-content">
      <!-- Success Header -->
      <div class="text-center mb-4">
        <div class="success-icon">
          <i class="bi bi-check-circle-fill text-success"></i>
        </div>
        <h3 class="text-success mb-2">Import Successful!</h3>
        <p class="text-muted">The audiobook has been imported to your library.</p>
      </div>

      <!-- Book Summary Card -->
      <div v-if="metadata" class="card mb-4">
        <div class="card-body">
          <h6 class="card-title">
            <i class="bi bi-book"></i> Imported Audiobook
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
              <h5 class="mb-2">{{ metadata.title }}</h5>
              <p v-if="metadata.subtitle" class="text-muted mb-2">{{ metadata.subtitle }}</p>
              <p class="mb-1">
                <strong>Authors:</strong> {{ metadata.authors.map(a => a.name).join(', ') }}
              </p>
              <p v-if="metadata.narrators.length > 0" class="mb-1">
                <strong>Narrators:</strong> {{ metadata.narrators.map(n => n.name).join(', ') }}
              </p>
              <p v-if="metadata.seriesPrimary" class="mb-1">
                <strong>Series:</strong> {{ metadata.seriesPrimary.name }}
                <span v-if="metadata.seriesPrimary.position">
                  (Book {{ metadata.seriesPrimary.position }})
                </span>
              </p>
              <p class="mb-0">
                <strong>ASIN:</strong> <span class="font-monospace">{{ metadata.asin }}</span>
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Destination Path Card -->
      <div class="card mb-4">
        <div class="card-body">
          <h6 class="card-title">
            <i class="bi bi-folder-check"></i> Destination
          </h6>
          <div class="destination-path">
            <code>
              <span class="path-parent">{{ pathParts.parentPath }}/</span><span class="path-book">{{ pathParts.bookFolder }}</span>
            </code>
          </div>
          <p class="text-muted small mt-2 mb-0">
            <i class="bi bi-info-circle"></i>
            The audiobook files have been moved to this location and metadata files have been written.
          </p>
        </div>
      </div>

      <!-- Additional Message -->
      <div v-if="message" class="alert alert-info">
        <i class="bi bi-info-circle"></i> {{ message }}
      </div>

      <!-- Next Steps -->
      <div class="card border-success">
        <div class="card-body">
          <h6 class="card-title text-success">
            <i class="bi bi-lightbulb"></i> Next Steps
          </h6>
          <ul class="mb-0">
            <li>The audiobook is now available in your Audiobookshelf library</li>
            <li>Audiobookshelf will automatically detect the new files on next scan</li>
            <li>The torrent will continue seeding in qBittorrent</li>
            <li>You can safely close this wizard</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Import Failure -->
    <div v-else class="summary-content">
      <!-- Failure Header -->
      <div class="text-center mb-4">
        <div class="failure-icon">
          <i class="bi bi-x-circle-fill text-danger"></i>
        </div>
        <h3 class="text-danger mb-2">Import Failed</h3>
        <p class="text-muted">The audiobook import encountered an error.</p>
      </div>

      <!-- Error Details -->
      <div class="alert alert-danger">
        <h6 class="alert-heading">
          <i class="bi bi-exclamation-triangle"></i> Error Details
        </h6>
        <p class="mb-0">{{ message || 'An unknown error occurred during import.' }}</p>
      </div>

      <!-- Attempted Destination -->
      <div v-if="destinationPath !== 'Unknown'" class="card mb-4">
        <div class="card-body">
          <h6 class="card-title">
            <i class="bi bi-folder-x"></i> Attempted Destination
          </h6>
          <div class="destination-path">
            <code>{{ destinationPath }}</code>
          </div>
        </div>
      </div>

      <!-- Troubleshooting -->
      <div class="card border-warning">
        <div class="card-body">
          <h6 class="card-title text-warning">
            <i class="bi bi-tools"></i> Troubleshooting
          </h6>
          <ul class="mb-0">
            <li>Verify the destination library path exists and is writable</li>
            <li>Check that the torrent files are accessible</li>
            <li>Ensure there is sufficient disk space</li>
            <li>Review the error message above for specific details</li>
            <li>You can try the import again by closing and reopening the wizard</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Finish Button -->
    <div class="text-end mt-4">
      <button
        class="btn btn-lg"
        :class="isSuccess ? 'btn-success' : 'btn-secondary'"
        @click="handleFinish"
      >
        <i class="bi bi-check-lg"></i> Finish
      </button>
    </div>
  </div>
</template>

<style scoped>
.import-summary-step {
  min-height: 400px;
}

.summary-content {
  animation: fadeIn 0.5s ease-in;
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

.success-icon i {
  font-size: 5rem;
  animation: scaleIn 0.5s ease-out;
}

.failure-icon i {
  font-size: 5rem;
  animation: shake 0.5s ease-out;
}

@keyframes scaleIn {
  from {
    transform: scale(0);
  }
  to {
    transform: scale(1);
  }
}

@keyframes shake {
  0%, 100% {
    transform: translateX(0);
  }
  10%, 30%, 50%, 70%, 90% {
    transform: translateX(-5px);
  }
  20%, 40%, 60%, 80% {
    transform: translateX(5px);
  }
}

.destination-path {
  background-color: var(--bs-gray-900);
  border: 1px solid var(--bs-border-color);
  border-radius: 0.375rem;
  padding: 1rem;
  font-family: 'Courier New', monospace;
  word-break: break-all;
}

.destination-path code {
  color: var(--bs-success);
  font-size: 0.95rem;
}

.path-parent {
  color: var(--bs-gray-500);
}

.path-book {
  color: var(--bs-success);
  font-weight: bold;
}
</style>
