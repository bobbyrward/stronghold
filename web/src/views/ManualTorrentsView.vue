<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import TorrentTable from '@/components/TorrentTable.vue'
import CategoryChangeDialog from '@/components/common/CategoryChangeDialog.vue'
import TagChangeDialog from '@/components/common/TagChangeDialog.vue'
import AudiobookImportWizard from '@/components/audiobook/AudiobookImportWizard.vue'
import type { Torrent } from '@/types/api'

const toast = useToastStore()
const data = ref<Torrent[]>([])
const loading = ref(true)
const showCategoryModal = ref(false)
const showTagModal = ref(false)
const showImportWizard = ref(false)
const selectedHash = ref<string>('')
const currentCategory = ref<string>('')
const currentTags = ref<string>('')

onMounted(async () => {
    try {
        data.value = await api.torrents.manualIntervention()
        sortTorrents()
    } catch (e) {
        toast.error('Failed to load feeds')
    } finally {
        loading.value = false
    }
})

function sortTorrents() {
    data.value.sort((a, b) => {
        var cmp = a.state.localeCompare(b.state)
        if (cmp != 0) {
            return cmp
        }

        cmp = a.category.localeCompare(b.category)
        if (cmp != 0) {
            return cmp
        }

        return a.name.localeCompare(b.name)
    })
}

async function handleChangeCategory(hash: string) {
    const torrent = data.value.find(d => d.hash === hash)
    selectedHash.value = hash
    currentCategory.value = torrent?.category || ''
    showCategoryModal.value = true
}

async function onCategoryConfirm(category: string) {
    try {
        await api.torrents.changeCategory(selectedHash.value, category)

        data.value.filter(d => d.hash === selectedHash.value).forEach(d => {
            d.category = category
        })

        sortTorrents()
        showCategoryModal.value = false
        toast.success('Category changed')
    } catch (e) {
        toast.error('Failed to change category')
    }
}

async function onCategoryCancel() {
    showCategoryModal.value = false
    selectedHash.value = ''
    currentCategory.value = ''
}

async function handleChangeTags(hash: string) {
    const torrent = data.value.find(d => d.hash === hash)
    selectedHash.value = hash
    currentTags.value = torrent?.tags || ''
    showTagModal.value = true
}

async function onTagConfirm(tags: string) {
    try {
        await api.torrents.changeTags(selectedHash.value, tags)

        data.value.filter(d => d.hash === selectedHash.value).forEach(d => {
            d.tags = tags
        })

        showTagModal.value = false
        toast.success('Tags changed')
    } catch (e) {
        toast.error('Failed to change tags')
    }
}

async function onTagCancel() {
    showTagModal.value = false
    selectedHash.value = ''
    currentTags.value = ''
}

async function handleImportAudiobook(hash: string) {
    selectedHash.value = hash
    showImportWizard.value = true
}

async function onImportWizardClose() {
    showImportWizard.value = false
    selectedHash.value = ''
}

async function onImportWizardSuccess() {
    toast.success('Audiobook imported successfully')
    showImportWizard.value = false
    selectedHash.value = ''

    // Optionally refresh the torrent list
    try {
        data.value = await api.torrents.manualIntervention()
        sortTorrents()
    } catch (e) {
        console.error('Failed to refresh torrents:', e)
    }
}
</script>

<template>
    <div class="mt-4">
        <h2>Torrents</h2>
        <p class="text-muted mb-4">Manual Intervention Torrents</p>

        <TorrentTable :data="data" :loading="loading" :on-change-category="handleChangeCategory"
            :on-change-tags="handleChangeTags" :on-import-audiobook="handleImportAudiobook" />

        <CategoryChangeDialog :show="showCategoryModal" :current-category="currentCategory" @confirm="onCategoryConfirm"
            @cancel="onCategoryCancel" />

        <TagChangeDialog :show="showTagModal" :current-tags="currentTags" @confirm="onTagConfirm"
            @cancel="onTagCancel" />

        <AudiobookImportWizard :show="showImportWizard" :torrent-hash="selectedHash" @close="onImportWizardClose"
            @success="onImportWizardSuccess" />
    </div>
</template>

<style scoped>
:deep(td) {
    word-break: break-all;
}
</style>
