<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import HardcoverSearchModal from '@/components/common/HardcoverSearchModal.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import type { Author, AuthorAlias, AuthorSubscription, AuthorSubscriptionItem, SubscriptionScope, Notifier, Library, HardcoverAuthorSearchResult } from '@/types/api'

const props = defineProps<{
  author: Author
  isEditing: boolean
  subscriptionScopes: SubscriptionScope[]
  notifiers: Notifier[]
  libraries: Library[]
}>()

const emit = defineEmits<{
  edit: []
  save: [data: { name: string; hardcover_ref: string | null }]
  cancel: []
  delete: []
}>()

const toast = useToastStore()

// Expansion state
const expanded = ref(false)
const activeTab = ref<'aliases' | 'subscription' | 'downloads'>('aliases')

// Badge counts
const aliasesCount = ref(0)
const hasSubscription = ref(false)
const downloadsCount = ref(0)
const countsLoaded = ref(false)

// Edit state for author
const editData = ref({ name: '', hardcover_ref: '' })
const showHardcoverModal = ref(false)

// Aliases state
const aliases = ref<AuthorAlias[]>([])
const aliasesLoaded = ref(false)
const newAliasName = ref('')
const editingAliasId = ref<number | null>(null)
const editAliasName = ref('')
const deleteAliasConfirm = ref({ show: false, id: 0, name: '' })

// Subscription state
const subscription = ref<AuthorSubscription | null>(null)
const subscriptionLoaded = ref(false)
const subscriptionFormMode = ref<'none' | 'create' | 'edit'>('none')
const subscriptionForm = ref({
  scope_name: '',
  notifier_id: null as number | null,
  ebook_library_name: '',
  audiobook_library_name: ''
})
const deleteSubscriptionConfirm = ref(false)

// Computed libraries filtered by book type
const ebookLibraries = computed(() =>
  props.libraries.filter(l => l.book_type_name === 'ebook')
)
const audiobookLibraries = computed(() =>
  props.libraries.filter(l => l.book_type_name === 'audiobook')
)

// Downloads state
const downloads = ref<AuthorSubscriptionItem[]>([])
const downloadsLoaded = ref(false)

// Initialize edit data when editing starts
watch(() => props.isEditing, (editing) => {
  if (editing) {
    editData.value = {
      name: props.author.name,
      hardcover_ref: props.author.hardcover_ref || ''
    }
  }
}, { immediate: true })

// Load data when tab becomes active
watch(activeTab, async (tab) => {
  if (tab === 'aliases' && !aliasesLoaded.value) {
    await loadAliases()
  }
  if (tab === 'subscription' && !subscriptionLoaded.value) {
    await loadSubscription()
  }
  if (tab === 'downloads' && !downloadsLoaded.value) {
    await loadDownloads()
  }
})

async function toggleExpand() {
  if (props.isEditing) return
  expanded.value = !expanded.value

  if (expanded.value && !countsLoaded.value) {
    await loadCounts()
  }
  if (expanded.value && activeTab.value === 'aliases' && !aliasesLoaded.value) {
    await loadAliases()
  }
  if (expanded.value && activeTab.value === 'subscription' && !subscriptionLoaded.value) {
    await loadSubscription()
  }
  if (expanded.value && activeTab.value === 'downloads' && !downloadsLoaded.value) {
    await loadDownloads()
  }
}

async function loadCounts() {
  try {
    const aliasesData = await api.authors.aliases.list(props.author.id)
    aliasesCount.value = aliasesData.length
  } catch {
    aliasesCount.value = 0
  }

  try {
    await api.authors.subscription.get(props.author.id)
    hasSubscription.value = true
    try {
      const items = await api.authors.subscription.items(props.author.id)
      downloadsCount.value = items.length
    } catch {
      downloadsCount.value = 0
    }
  } catch {
    hasSubscription.value = false
    downloadsCount.value = 0
  }

  countsLoaded.value = true
}

async function loadAliases() {
  try {
    aliases.value = await api.authors.aliases.list(props.author.id)
    aliasesCount.value = aliases.value.length
    aliasesLoaded.value = true
  } catch (e) {
    toast.error('Failed to load aliases')
  }
}

async function loadSubscription() {
  try {
    subscription.value = await api.authors.subscription.get(props.author.id)
    hasSubscription.value = true
  } catch {
    subscription.value = null
    hasSubscription.value = false
  }
  subscriptionLoaded.value = true
}

async function loadDownloads() {
  if (!hasSubscription.value) {
    downloads.value = []
    downloadsLoaded.value = true
    return
  }

  try {
    downloads.value = await api.authors.subscription.items(props.author.id)
    downloadsCount.value = downloads.value.length
    downloadsLoaded.value = true
  } catch {
    downloads.value = []
    downloadsLoaded.value = true
  }
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleString()
}

function truncateHash(hash: string, length: number = 12): string {
  if (hash.length <= length) return hash
  return hash.substring(0, length) + '...'
}

function handleHardcoverSelect(author: HardcoverAuthorSearchResult) {
  editData.value.hardcover_ref = author.slug
}

function handleSave() {
  emit('save', {
    name: editData.value.name.trim(),
    hardcover_ref: editData.value.hardcover_ref || null
  })
}

// Alias CRUD functions
async function addAlias() {
  if (!newAliasName.value.trim()) {
    toast.error('Alias name is required')
    return
  }

  try {
    const created = await api.authors.aliases.create(props.author.id, {
      name: newAliasName.value.trim()
    })
    aliases.value.push(created)
    aliasesCount.value = aliases.value.length
    newAliasName.value = ''
    toast.success('Alias created')
  } catch (e) {
    toast.error('Failed to create alias')
  }
}

function startEditAlias(alias: AuthorAlias) {
  editingAliasId.value = alias.id
  editAliasName.value = alias.name
}

function cancelEditAlias() {
  editingAliasId.value = null
  editAliasName.value = ''
}

async function saveEditAlias(aliasId: number) {
  if (!editAliasName.value.trim()) {
    toast.error('Alias name is required')
    return
  }

  try {
    const updated = await api.authors.aliases.update(props.author.id, aliasId, {
      name: editAliasName.value.trim()
    })
    const index = aliases.value.findIndex(a => a.id === aliasId)
    aliases.value[index] = updated
    editingAliasId.value = null
    editAliasName.value = ''
    toast.success('Alias updated')
  } catch (e) {
    toast.error('Failed to update alias')
  }
}

function confirmDeleteAlias(alias: AuthorAlias) {
  deleteAliasConfirm.value = { show: true, id: alias.id, name: alias.name }
}

async function handleDeleteAlias() {
  try {
    await api.authors.aliases.delete(props.author.id, deleteAliasConfirm.value.id)
    aliases.value = aliases.value.filter(a => a.id !== deleteAliasConfirm.value.id)
    aliasesCount.value = aliases.value.length
    toast.success('Alias deleted')
  } catch (e) {
    toast.error('Failed to delete alias')
  } finally {
    deleteAliasConfirm.value = { show: false, id: 0, name: '' }
  }
}

// Subscription CRUD functions
function startCreateSubscription() {
  subscriptionFormMode.value = 'create'
  subscriptionForm.value = {
    scope_name: props.subscriptionScopes[0]?.name || '',
    notifier_id: null,
    ebook_library_name: ebookLibraries.value[0]?.name || '',
    audiobook_library_name: audiobookLibraries.value[0]?.name || ''
  }
}

function startEditSubscription() {
  if (!subscription.value) return
  subscriptionFormMode.value = 'edit'
  subscriptionForm.value = {
    scope_name: subscription.value.scope_name,
    notifier_id: subscription.value.notifier_id,
    ebook_library_name: subscription.value.ebook_library_name,
    audiobook_library_name: subscription.value.audiobook_library_name
  }
}

function cancelSubscriptionForm() {
  subscriptionFormMode.value = 'none'
  subscriptionForm.value = { scope_name: '', notifier_id: null, ebook_library_name: '', audiobook_library_name: '' }
}

async function saveSubscription() {
  if (!subscriptionForm.value.scope_name) {
    toast.error('Scope is required')
    return
  }
  if (!subscriptionForm.value.ebook_library_name) {
    toast.error('Ebook Library is required')
    return
  }
  if (!subscriptionForm.value.audiobook_library_name) {
    toast.error('Audiobook Library is required')
    return
  }

  try {
    const requestData = {
      scope_name: subscriptionForm.value.scope_name,
      notifier_id: subscriptionForm.value.notifier_id,
      ebook_library_name: subscriptionForm.value.ebook_library_name,
      audiobook_library_name: subscriptionForm.value.audiobook_library_name
    }
    if (subscriptionFormMode.value === 'create') {
      subscription.value = await api.authors.subscription.create(props.author.id, requestData)
      hasSubscription.value = true
      toast.success('Subscription created')
    } else {
      subscription.value = await api.authors.subscription.update(props.author.id, requestData)
      toast.success('Subscription updated')
    }
    subscriptionFormMode.value = 'none'
  } catch (e) {
    toast.error('Failed to save subscription')
  }
}

async function handleDeleteSubscription() {
  try {
    await api.authors.subscription.delete(props.author.id)
    subscription.value = null
    hasSubscription.value = false
    downloadsCount.value = 0
    toast.success('Subscription deleted')
  } catch (e) {
    toast.error('Failed to delete subscription')
  } finally {
    deleteSubscriptionConfirm.value = false
  }
}
</script>

<template>
  <!-- Main row -->
  <tr @click="toggleExpand" :style="{ cursor: isEditing ? 'default' : 'pointer' }">
    <td>
      <i :class="expanded ? 'bi bi-chevron-down' : 'bi bi-chevron-right'"
        :style="{ color: isEditing ? 'var(--bs-secondary)' : undefined }"></i>
    </td>
    <td>{{ author.id }}</td>
    <td>
      <input v-if="isEditing" type="text" v-model="editData.name" class="form-control form-control-sm" @click.stop>
      <span v-else>{{ author.name }}</span>
    </td>
    <td>
      <template v-if="isEditing">
        <div class="input-group input-group-sm" @click.stop>
          <input type="text" v-model="editData.hardcover_ref" class="form-control" readonly
            placeholder="Click search to set">
          <button class="btn btn-outline-info" type="button" @click="showHardcoverModal = true">
            <i class="bi bi-search"></i>
          </button>
        </div>
      </template>
      <template v-else>
        <span v-if="author.hardcover_ref" class="text-muted">{{ author.hardcover_ref }}</span>
        <span v-else class="text-muted fst-italic">Not set</span>
      </template>
    </td>
    <td @click.stop>
      <template v-if="isEditing">
        <button class="btn btn-success btn-sm me-1" @click="handleSave" title="Save">
          <i class="bi bi-check"></i>
        </button>
        <button class="btn btn-secondary btn-sm" @click="emit('cancel')" title="Cancel">
          <i class="bi bi-x"></i>
        </button>
      </template>
      <template v-else>
        <button class="btn btn-primary btn-sm me-1" @click="emit('edit')" title="Edit">
          <i class="bi bi-pencil"></i>
        </button>
        <button class="btn btn-danger btn-sm" @click="emit('delete')" title="Delete">
          <i class="bi bi-trash"></i>
        </button>
      </template>
    </td>
  </tr>

  <!-- Expanded row with tabs -->
  <tr v-if="expanded">
    <td colspan="5">
      <div class="ps-4 py-2">
        <ul class="nav nav-tabs">
          <li class="nav-item">
            <button class="nav-link" :class="{ active: activeTab === 'aliases' }" @click="activeTab = 'aliases'">
              Aliases
              <span class="badge bg-secondary ms-1">{{ aliasesCount }}</span>
            </button>
          </li>
          <li class="nav-item">
            <button class="nav-link" :class="{ active: activeTab === 'subscription' }"
              @click="activeTab = 'subscription'">
              Subscription
              <i v-if="hasSubscription" class="bi bi-check-circle-fill text-success ms-1"></i>
            </button>
          </li>
          <li class="nav-item">
            <button class="nav-link" :class="{ active: activeTab === 'downloads' }" @click="activeTab = 'downloads'">
              Downloads
              <span class="badge bg-secondary ms-1">{{ downloadsCount }}</span>
            </button>
          </li>
        </ul>

        <div class="tab-content p-3 border border-top-0 rounded-bottom">
          <!-- Aliases Tab -->
          <div v-if="activeTab === 'aliases'" class="tab-pane active">
            <!-- Add alias form -->
            <div class="row mb-3">
              <div class="col">
                <div class="input-group input-group-sm">
                  <input type="text" v-model="newAliasName" class="form-control" placeholder="New alias name..."
                    @keyup.enter="addAlias">
                  <button class="btn btn-primary" type="button" @click="addAlias">
                    <i class="bi bi-plus-lg me-1"></i>Add
                  </button>
                </div>
              </div>
            </div>

            <!-- Aliases table -->
            <table v-if="aliases.length > 0" class="table table-sm table-dark mb-0">
              <thead>
                <tr>
                  <th style="width: 60px">ID</th>
                  <th>Name</th>
                  <th style="width: 100px">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="alias in aliases" :key="alias.id">
                  <td>{{ alias.id }}</td>
                  <td>
                    <input v-if="editingAliasId === alias.id" type="text" v-model="editAliasName"
                      class="form-control form-control-sm" @keyup.enter="saveEditAlias(alias.id)">
                    <span v-else>{{ alias.name }}</span>
                  </td>
                  <td>
                    <template v-if="editingAliasId === alias.id">
                      <button class="btn btn-success btn-sm me-1" @click="saveEditAlias(alias.id)" title="Save">
                        <i class="bi bi-check"></i>
                      </button>
                      <button class="btn btn-secondary btn-sm" @click="cancelEditAlias" title="Cancel">
                        <i class="bi bi-x"></i>
                      </button>
                    </template>
                    <template v-else>
                      <button class="btn btn-primary btn-sm me-1" @click="startEditAlias(alias)" title="Edit">
                        <i class="bi bi-pencil"></i>
                      </button>
                      <button class="btn btn-danger btn-sm" @click="confirmDeleteAlias(alias)" title="Delete">
                        <i class="bi bi-trash"></i>
                      </button>
                    </template>
                  </td>
                </tr>
              </tbody>
            </table>

            <p v-else class="text-muted mb-0">No aliases defined</p>
          </div>

          <!-- Subscription Tab -->
          <div v-else-if="activeTab === 'subscription'" class="tab-pane active">
            <!-- No subscription state -->
            <div v-if="!subscription && subscriptionFormMode === 'none'">
              <p class="text-muted mb-3">No subscription for this author</p>
              <button class="btn btn-primary btn-sm" @click="startCreateSubscription">
                <i class="bi bi-plus-lg me-1"></i>Create Subscription
              </button>
            </div>

            <!-- Create/Edit form -->
            <div v-else-if="subscriptionFormMode !== 'none'">
              <div class="row g-3 mb-3">
                <div class="col-md-6">
                  <label class="form-label">Scope</label>
                  <select v-model="subscriptionForm.scope_name" class="form-select form-select-sm">
                    <option v-for="scope in subscriptionScopes" :key="scope.id" :value="scope.name">
                      {{ scope.name }}
                    </option>
                  </select>
                </div>
                <div class="col-md-6">
                  <label class="form-label">Notifier (optional)</label>
                  <select v-model="subscriptionForm.notifier_id" class="form-select form-select-sm">
                    <option :value="null">None</option>
                    <option v-for="notifier in notifiers" :key="notifier.id" :value="notifier.id">
                      {{ notifier.name }}
                    </option>
                  </select>
                </div>
                <div class="col-md-6">
                  <label class="form-label">Ebook Library</label>
                  <select v-model="subscriptionForm.ebook_library_name" class="form-select form-select-sm" required>
                    <option v-for="lib in ebookLibraries" :key="lib.id" :value="lib.name">
                      {{ lib.name }}
                    </option>
                  </select>
                </div>
                <div class="col-md-6">
                  <label class="form-label">Audiobook Library</label>
                  <select v-model="subscriptionForm.audiobook_library_name" class="form-select form-select-sm" required>
                    <option v-for="lib in audiobookLibraries" :key="lib.id" :value="lib.name">
                      {{ lib.name }}
                    </option>
                  </select>
                </div>
              </div>
              <div>
                <button class="btn btn-success btn-sm me-2" @click="saveSubscription">
                  <i class="bi bi-check me-1"></i>Save
                </button>
                <button class="btn btn-secondary btn-sm" @click="cancelSubscriptionForm">
                  Cancel
                </button>
              </div>
            </div>

            <!-- Subscription display -->
            <div v-else-if="subscription">
              <dl class="row mb-3">
                <dt class="col-sm-3">Scope</dt>
                <dd class="col-sm-9">{{ subscription.scope_name }}</dd>
                <dt class="col-sm-3">Notifier</dt>
                <dd class="col-sm-9">{{ subscription.notifier_name || 'None' }}</dd>
                <dt class="col-sm-3">Ebook Library</dt>
                <dd class="col-sm-9">{{ subscription.ebook_library_name }}</dd>
                <dt class="col-sm-3">Audiobook Library</dt>
                <dd class="col-sm-9">{{ subscription.audiobook_library_name }}</dd>
              </dl>
              <button class="btn btn-primary btn-sm me-2" @click="startEditSubscription">
                <i class="bi bi-pencil me-1"></i>Edit
              </button>
              <button class="btn btn-danger btn-sm" @click="deleteSubscriptionConfirm = true">
                <i class="bi bi-trash me-1"></i>Delete
              </button>
            </div>
          </div>

          <!-- Downloads Tab -->
          <div v-else-if="activeTab === 'downloads'" class="tab-pane active">
            <!-- No subscription state -->
            <div v-if="!hasSubscription">
              <p class="text-muted mb-0">No subscription for this author. Create a subscription to track downloads.</p>
            </div>

            <!-- No downloads state -->
            <div v-else-if="downloads.length === 0">
              <p class="text-muted mb-0">No downloads yet for this subscription.</p>
            </div>

            <!-- Downloads table -->
            <div v-else>
              <table class="table table-sm table-dark mb-3">
                <thead>
                  <tr>
                    <th>Downloaded At</th>
                    <th>Booksearch ID</th>
                    <th>Torrent Hash</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in downloads" :key="item.id">
                    <td>{{ formatDate(item.downloaded_at) }}</td>
                    <td>{{ item.booksearch_id }}</td>
                    <td :title="item.torrent_hash">{{ truncateHash(item.torrent_hash) }}</td>
                  </tr>
                </tbody>
              </table>
              <router-link :to="`/subscription-items?author_id=${author.id}`" class="btn btn-outline-secondary btn-sm">
                <i class="bi bi-list-ul me-1"></i>View all download history
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </td>
  </tr>

  <!-- Hardcover search modal -->
  <HardcoverSearchModal :show="showHardcoverModal" :initial-query="editData.name"
    @close="showHardcoverModal = false" @select="handleHardcoverSelect" />

  <!-- Delete alias confirmation -->
  <ConfirmDialog :show="deleteAliasConfirm.show" title="Delete Alias"
    :message="`Are you sure you want to delete alias '${deleteAliasConfirm.name}'?`" confirm-text="Delete"
    variant="danger" @confirm="handleDeleteAlias" @cancel="deleteAliasConfirm.show = false" />

  <!-- Delete subscription confirmation -->
  <ConfirmDialog :show="deleteSubscriptionConfirm" title="Delete Subscription"
    message="Are you sure you want to delete this subscription? This will also remove the download history."
    confirm-text="Delete" variant="danger" @confirm="handleDeleteSubscription"
    @cancel="deleteSubscriptionConfirm = false" />
</template>
