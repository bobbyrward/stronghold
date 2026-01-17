<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import HardcoverSearchModal from '@/components/common/HardcoverSearchModal.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import AuthorRow from '@/components/AuthorRow.vue'
import type { Author, SubscriptionScope, Notifier, HardcoverAuthorSearchResult } from '@/types/api'

const toast = useToastStore()
const authors = ref<Author[]>([])
const subscriptionScopes = ref<SubscriptionScope[]>([])
const notifiers = ref<Notifier[]>([])
const loading = ref(true)

// Search state
const searchQuery = ref('')
let searchTimeout: ReturnType<typeof setTimeout> | null = null

// Add author state
const adding = ref(false)
const newAuthor = ref({ name: '', hardcover_ref: '' })
const showHardcoverModal = ref(false)

// Edit author state
const editingId = ref<number | null>(null)

// Delete confirmation state
const deleteConfirm = ref({ show: false, id: 0, name: '' })

onMounted(async () => {
  try {
    const [authorsData, scopesData, notifiersData] = await Promise.all([
      api.authors.list(),
      api.subscriptionScopes.list(),
      api.notifiers.list()
    ])
    authors.value = authorsData
    subscriptionScopes.value = scopesData
    notifiers.value = notifiersData
  } catch (e) {
    toast.error('Failed to load data')
  } finally {
    loading.value = false
  }
})

async function loadAuthors(query?: string) {
  loading.value = true
  try {
    authors.value = await api.authors.list(query)
  } catch (e) {
    toast.error('Failed to load authors')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    loadAuthors(searchQuery.value || undefined)
  }, 300)
}

function clearSearch() {
  searchQuery.value = ''
  loadAuthors()
}

// Add author functions
function startAdd() {
  adding.value = true
  newAuthor.value = { name: '', hardcover_ref: '' }
}

function cancelAdd() {
  adding.value = false
  newAuthor.value = { name: '', hardcover_ref: '' }
}

function openHardcoverModal() {
  showHardcoverModal.value = true
}

function handleHardcoverSelect(author: HardcoverAuthorSearchResult) {
  newAuthor.value.hardcover_ref = author.slug
  if (!newAuthor.value.name) {
    newAuthor.value.name = author.name
  }
}

async function saveNew() {
  if (!newAuthor.value.name.trim()) {
    toast.error('Name is required')
    return
  }

  try {
    const created = await api.authors.create({
      name: newAuthor.value.name.trim(),
      hardcover_ref: newAuthor.value.hardcover_ref || null
    })
    authors.value.push(created)
    toast.success('Author created')
    adding.value = false
    newAuthor.value = { name: '', hardcover_ref: '' }
  } catch (e) {
    toast.error('Failed to create author')
  }
}

// Edit author functions
function startEdit(authorId: number) {
  editingId.value = authorId
}

function cancelEdit() {
  editingId.value = null
}

async function saveEdit(authorId: number, data: { name: string; hardcover_ref: string | null }) {
  if (!data.name.trim()) {
    toast.error('Name is required')
    return
  }

  try {
    const updated = await api.authors.update(authorId, {
      name: data.name.trim(),
      hardcover_ref: data.hardcover_ref
    })
    const index = authors.value.findIndex(a => a.id === authorId)
    authors.value[index] = updated
    toast.success('Author updated')
    editingId.value = null
  } catch (e) {
    toast.error('Failed to update author')
  }
}

// Delete author functions
function confirmDelete(author: Author) {
  deleteConfirm.value = { show: true, id: author.id, name: author.name }
}

async function handleDelete() {
  try {
    await api.authors.delete(deleteConfirm.value.id)
    authors.value = authors.value.filter(a => a.id !== deleteConfirm.value.id)
    toast.success('Author deleted')
  } catch (e) {
    toast.error('Failed to delete author')
  } finally {
    deleteConfirm.value = { show: false, id: 0, name: '' }
  }
}
</script>

<template>
  <div class="mt-4">
    <h2>Authors</h2>
    <p class="text-muted mb-4">Manage author subscriptions</p>

    <div class="position-relative">
      <LoadingSpinner v-if="loading" />

      <div class="row mb-3">
        <div class="col-auto">
          <button class="btn btn-primary btn-sm" @click="startAdd" :disabled="adding">
            <i class="bi bi-plus-lg me-1"></i>
            Add Author
          </button>
        </div>
        <div class="col">
          <div class="input-group input-group-sm">
            <input type="text" v-model="searchQuery" @input="handleSearch" class="form-control"
              placeholder="Search authors...">
            <button v-if="searchQuery" class="btn btn-outline-secondary" @click="clearSearch" type="button">
              <i class="bi bi-x"></i>
            </button>
          </div>
        </div>
      </div>

      <table class="table table-dark table-striped table-hover">
        <thead>
          <tr>
            <th style="width: 40px"></th>
            <th style="width: 60px">ID</th>
            <th>Name</th>
            <th>Hardcover Ref</th>
            <th style="width: 150px">Actions</th>
          </tr>
        </thead>
        <tbody>
          <!-- Add new row -->
          <tr v-if="adding">
            <td></td>
            <td>-</td>
            <td>
              <input type="text" v-model="newAuthor.name" class="form-control form-control-sm"
                placeholder="Author name">
            </td>
            <td>
              <div class="input-group input-group-sm">
                <input type="text" v-model="newAuthor.hardcover_ref" class="form-control" readonly
                  placeholder="Click search to set">
                <button class="btn btn-primary" type="button" @click="openHardcoverModal">
                  <i class="bi bi-search"></i>
                </button>
              </div>
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

          <!-- Empty state -->
          <tr v-if="authors.length === 0 && !adding && !loading">
            <td colspan="5" class="text-center text-muted py-4">
              No authors found
            </td>
          </tr>

          <!-- Author rows -->
          <AuthorRow v-for="author in authors" :key="author.id" :author="author" :is-editing="editingId === author.id"
            :subscription-scopes="subscriptionScopes" :notifiers="notifiers" @edit="startEdit(author.id)"
            @save="(data) => saveEdit(author.id, data)" @cancel="cancelEdit" @delete="confirmDelete(author)" />
        </tbody>
      </table>
    </div>

    <!-- Hardcover search modal for new author -->
    <HardcoverSearchModal :show="showHardcoverModal" :initial-query="newAuthor.name"
      @close="showHardcoverModal = false" @select="handleHardcoverSelect" />

    <!-- Delete confirmation dialog -->
    <ConfirmDialog :show="deleteConfirm.show" title="Delete Author"
      :message="`Are you sure you want to delete '${deleteConfirm.name}'?`" confirm-text="Delete" variant="danger"
      @confirm="handleDelete" @cancel="deleteConfirm.show = false" />
  </div>
</template>
