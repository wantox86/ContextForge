<script setup lang="ts">
import { ref, watch, onMounted, onActivated } from 'vue'
import { useProjectStore } from '../stores/project'
import { useExportStore } from '../stores/export'

const projectStore = useProjectStore()
const exportStore = useExportStore()

const isRefreshing = ref(false)

async function refresh() {
  const project = projectStore.currentProject
  if (!project) return
  isRefreshing.value = true
  try {
    const content = await projectStore.previewExport(project.ID, exportStore.selectedTarget)
    exportStore.setPreview(content)
  } finally {
    isRefreshing.value = false
  }
}

onMounted(refresh)
// Re-run when tab becomes active again (KeepAlive)
onActivated(refresh)
// Re-run when target changes or when a project is first analyzed
watch(() => exportStore.selectedTarget, refresh)
watch(() => projectStore.currentProject, (proj) => { if (proj) refresh() })
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Toolbar -->
    <div class="flex items-center gap-3 px-5 py-3 bg-gray-800 border-b border-gray-700 shrink-0">
      <span class="text-sm text-gray-400">Target:</span>
      <select
        v-model="exportStore.selectedTarget"
        class="bg-gray-700 border border-gray-600 rounded px-2 py-1 text-sm text-white focus:outline-none focus:border-indigo-500"
      >
        <option value="claude">CLAUDE.md</option>
        <option value="copilot">copilot-instructions.md</option>
      </select>
      <button
        @click="refresh"
        :disabled="isRefreshing"
        class="px-3 py-1 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 rounded text-xs transition-colors"
      >
        {{ isRefreshing ? 'Refreshing…' : '↻ Refresh' }}
      </button>
      <span class="ml-auto text-xs text-gray-500">
        ~{{ exportStore.tokenCount.toLocaleString() }} tokens
      </span>
      <span
        v-if="exportStore.warnings.length"
        class="text-xs text-yellow-400"
      >
        ⚠ {{ exportStore.warnings[0].Model }} limit approaching
      </span>
    </div>

    <!-- Preview content -->
    <div class="flex-1 overflow-y-auto p-5">
      <p v-if="!projectStore.currentProject" class="text-gray-500 text-sm text-center mt-20">
        Analyze a project first to see the preview.
      </p>
      <pre
        v-else-if="exportStore.previewContent"
        class="text-xs text-gray-300 font-mono whitespace-pre-wrap leading-relaxed"
      >{{ exportStore.previewContent }}</pre>
      <p v-else class="text-gray-500 text-sm text-center mt-20">
        No content yet. Click ↻ Refresh to generate preview.
      </p>
    </div>
  </div>
</template>
