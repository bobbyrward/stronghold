<script setup lang="ts">
import { ref, watch } from 'vue'
import { api } from '@/services/api'
import type {
  TorrentImportInfo,
  BookMetadata,
  ExecuteImportResponse
} from '@/types/api'
import AsinDiscoveryStep from './wizard/AsinDiscoveryStep.vue'
import MetadataConfirmationStep from './wizard/MetadataConfirmationStep.vue'
import LibrarySelectionStep from './wizard/LibrarySelectionStep.vue'
import ImportSummaryStep from './wizard/ImportSummaryStep.vue'

// Props
const props = defineProps<{
  show: boolean
  torrentHash: string
}>()

// Emits
const emit = defineEmits<{
  close: []
  success: [hash: string]
}>()

// State management
const currentStep = ref(1)
const torrentInfo = ref<TorrentImportInfo | null>(null)
const selectedAsin = ref<string>('')
const bookMetadata = ref<BookMetadata | null>(null)
const selectedLibrary = ref<string>('')
const importResult = ref<ExecuteImportResponse | null>(null)
const error = ref<string | null>(null)
const loading = ref(false)

// Step titles
const stepTitles = [
  'ASIN Discovery',
  'Metadata Confirmation',
  'Library Selection',
  'Import',
  'Summary'
]

// Load torrent info when wizard opens
const loadTorrentInfo = async () => {
  if (!props.torrentHash) return

  loading.value = true
  error.value = null

  try {
    torrentInfo.value = await api.audiobookWizard.getTorrentInfo(props.torrentHash)

    // Pre-populate if we have ASIN
    if (torrentInfo.value.asin) {
      selectedAsin.value = torrentInfo.value.asin
    }

    // Pre-populate library suggestion
    if (torrentInfo.value.suggested_library) {
      selectedLibrary.value = torrentInfo.value.suggested_library
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load torrent information'
  } finally {
    loading.value = false
  }
}

// Navigation
const canGoNext = ref(true)
const canGoBack = ref(false)

const updateNavigationState = () => {
  canGoBack.value = currentStep.value > 1 && currentStep.value < 5
  canGoNext.value = currentStep.value < 5
}

const nextStep = () => {
  if (currentStep.value < 5) {
    error.value = null
    currentStep.value++

    // Skip Step 4 (import progress) since it's now handled in Step 3
    if (currentStep.value === 4) {
      currentStep.value = 5
    }

    updateNavigationState()
  }
}

const previousStep = () => {
  if (currentStep.value > 1) {
    error.value = null
    currentStep.value--

    // Skip Step 4 (import progress) since it's now handled in Step 3
    if (currentStep.value === 4) {
      currentStep.value = 3
    }

    updateNavigationState()
  }
}

const closeWizard = () => {
  // Reset state
  currentStep.value = 1
  torrentInfo.value = null
  selectedAsin.value = ''
  bookMetadata.value = null
  selectedLibrary.value = ''
  importResult.value = null
  error.value = null
  loading.value = false

  emit('close')
}

const handleSuccess = () => {
  emit('success', props.torrentHash)
  closeWizard()
}

// Watch for show prop changes to load data
watch(() => props.show, (newVal) => {
  if (newVal) {
    loadTorrentInfo()
    updateNavigationState()
  }
})

// Helper to check if step is completed
const isStepCompleted = (step: number) => {
  return currentStep.value > step
}
</script>

<template>
  <Teleport to="body">
    <div v-if="show" class="modal-backdrop fade show"></div>
    <div
      v-if="show"
      class="modal fade show d-block"
      tabindex="-1"
    >
      <div class="modal-dialog modal-xl">
        <div class="modal-content">
          <!-- Modal Header -->
          <div class="modal-header">
            <h5 class="modal-title">Audiobook Import Wizard</h5>
            <button
              type="button"
              class="btn-close"
              aria-label="Close"
              @click="closeWizard"
            ></button>
          </div>

          <!-- Modal Body -->
          <div class="modal-body">
            <!-- Error Alert -->
            <div v-if="error" class="alert alert-danger alert-dismissible fade show" role="alert">
              <strong>Error:</strong> {{ error }}
              <button type="button" class="btn-close" @click="error = null"></button>
            </div>

            <!-- Loading Spinner -->
            <div v-if="loading" class="text-center py-5">
              <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
              <p class="mt-3">Loading...</p>
            </div>

            <!-- Wizard Content -->
            <div v-else>
              <!-- Step Indicator -->
              <div class="wizard-steps mb-4">
                <div class="d-flex justify-content-between align-items-center">
                  <div
                    v-for="(title, index) in stepTitles"
                    :key="index"
                    class="wizard-step flex-grow-1"
                    :class="{
                      'active': currentStep === index + 1,
                      'completed': isStepCompleted(index + 1)
                    }"
                  >
                    <div class="step-circle-container">
                      <div class="step-circle">
                        <span v-if="!isStepCompleted(index + 1)">{{ index + 1 }}</span>
                        <i v-else class="bi bi-check-lg"></i>
                      </div>
                      <div class="step-title">{{ title }}</div>
                    </div>
                    <div v-if="index < stepTitles.length - 1" class="step-line"></div>
                  </div>
                </div>
              </div>

              <!-- Torrent Info Card -->
              <div v-if="torrentInfo" class="card mb-3">
                <div class="card-body">
                  <h6 class="card-title">Torrent Information</h6>
                  <div class="row">
                    <div class="col-md-6">
                      <strong>Name:</strong> {{ torrentInfo.name }}
                    </div>
                    <div class="col-md-3">
                      <strong>Category:</strong> {{ torrentInfo.category }}
                    </div>
                    <div class="col-md-3">
                      <strong>Tags:</strong> {{ torrentInfo.tags }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Step Content -->
              <div class="step-content">
                <!-- Step 1: ASIN Discovery -->
                <div v-if="currentStep === 1" class="step-container">
                  <AsinDiscoveryStep
                    :torrent-info="torrentInfo"
                    v-model="selectedAsin"
                    @next="nextStep"
                  />
                </div>

                <!-- Step 2: Metadata Confirmation -->
                <div v-if="currentStep === 2" class="step-container">
                  <MetadataConfirmationStep
                    :selected-asin="selectedAsin"
                    @metadata-loaded="(meta) => bookMetadata = meta"
                    @next="nextStep"
                  />
                </div>

                <!-- Step 3: Library Selection & Import -->
                <div v-if="currentStep === 3" class="step-container">
                  <LibrarySelectionStep
                    :torrent-hash="torrentHash"
                    :metadata="bookMetadata"
                    :suggested-library="torrentInfo?.suggested_library"
                    @import-complete="(result) => importResult = result"
                    @next="nextStep"
                  />
                </div>

                <!-- Step 4: Import Progress (merged into Step 3) -->
                <div v-if="currentStep === 4" class="step-container">
                  <h5>Step 4: Processing...</h5>
                  <p class="text-muted">Import in progress.</p>
                  <div class="alert alert-info">
                    This step is now handled by Step 3
                  </div>
                </div>

                <!-- Step 5: Import Summary -->
                <div v-if="currentStep === 5" class="step-container">
                  <ImportSummaryStep
                    :import-result="importResult"
                    :metadata="bookMetadata"
                    @finish="handleSuccess"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Modal Footer -->
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-secondary"
              @click="closeWizard"
            >
              Cancel
            </button>
            <button
              v-if="canGoBack"
              type="button"
              class="btn btn-outline-primary"
              @click="previousStep"
              :disabled="loading"
            >
              <i class="bi bi-arrow-left"></i> Back
            </button>
            <button
              v-if="canGoNext && currentStep < 5"
              type="button"
              class="btn btn-primary"
              @click="nextStep"
              :disabled="loading"
            >
              Next <i class="bi bi-arrow-right"></i>
            </button>
            <button
              v-if="currentStep === 5"
              type="button"
              class="btn btn-success"
              @click="handleSuccess"
            >
              <i class="bi bi-check-lg"></i> Finish
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* Wizard Steps Styling */
.wizard-steps {
  padding: 20px 0;
}

.wizard-step {
  position: relative;
  display: flex;
  align-items: center;
}

.step-circle-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  z-index: 2;
}

.step-circle {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-color: var(--bs-secondary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 1.2rem;
  transition: all 0.3s ease;
  border: 3px solid var(--bs-secondary);
}

.wizard-step.active .step-circle {
  background-color: var(--bs-primary);
  border-color: var(--bs-primary);
  box-shadow: 0 0 0 0.2rem rgba(var(--bs-primary-rgb), 0.25);
}

.wizard-step.completed .step-circle {
  background-color: var(--bs-success);
  border-color: var(--bs-success);
}

.step-title {
  margin-top: 10px;
  font-size: 0.85rem;
  text-align: center;
  color: var(--bs-secondary);
  font-weight: 500;
}

.wizard-step.active .step-title {
  color: var(--bs-primary);
  font-weight: bold;
}

.wizard-step.completed .step-title {
  color: var(--bs-success);
}

.step-line {
  flex-grow: 1;
  height: 3px;
  background-color: var(--bs-secondary);
  margin: 0 10px;
  position: relative;
  top: -20px;
}

.wizard-step.completed + .wizard-step .step-line {
  background-color: var(--bs-success);
}

.step-container {
  min-height: 300px;
  padding: 20px 0;
}

/* Modal adjustments */
.modal-xl {
  max-width: 900px;
}
</style>
