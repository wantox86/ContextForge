<script setup lang="ts">
import { ref } from 'vue'
import { useProjectStore } from '../stores/project'
import { useExportStore } from '../stores/export'

const projectStore = useProjectStore()
const exportStore = useExportStore()
const outputPath = ref('')
const successMsg = ref('')

async function doExport() {
  const project = projectStore.currentProject
  if (!project || !outputPath.value.trim()) return
  exportStore.setExporting(true)
  successMsg.value = ''
  try {
    await projectStore.exportToFile(project.ID, exportStore.selectedTarget, outputPath.value.trim())
    const fileName = exportStore.selectedTarget === 'claude' ? 'CLAUDE.md' : '.github/copilot-instructions.md'
    successMsg.value = `✅ Exported to ${outputPath.value.trim()}/${fileName}`
  } catch {
    // error is set in store
  } finally {
    exportStore.setExporting(false)
  }
}
</script>

<template>
  <div class="flex flex-col items-center justify-center h-full gap-6 p-8">
    <h2 class="text-2xl font-semibold text-indigo-300">Export</h2>

    <p v-if="!projectStore.currentProject" class="text-gray-400 text-sm">
      Analyze a project first before exporting.
    </p>

    <template v-else>
      <div class="w-full max-w-lg space-y-4">
        <!-- Target selector -->
        <div class="space-y-1">
          <label class="text-xs text-gray-400">Export target</label>
          <div class="flex gap-3">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="exportStore.selectedTarget" value="claude" class="accent-indigo-500" />
              <span class="text-sm">CLAUDE.md</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="exportStore.selectedTarget" value="copilot" class="accent-indigo-500" />
              <span class="text-sm">.github/copilot-instructions.md</span>
            </label>
          </div>
        </div>

        <!-- Output path -->
        <div class="space-y-1">
          <label class="text-xs text-gray-400">Output directory</label>
          <input
            v-model="outputPath"
            type="text"
            placeholder="/Users/you/projects/my-app"
            class="w-full bg-gray-800 border border-gray-600 rounded px-4 py-2 text-sm text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500"
          />
        </div>

        <!-- Token warning -->
        <div v-if="exportStore.warnings.length" class="bg-yellow-900/30 border border-yellow-700 rounded p-3">
          <p class="text-yellow-400 text-xs font-medium">⚠ Token Limit Warning</p>
          <ul class="mt-1 space-y-0.5">
            <li v-for="w in exportStore.warnings" :key="w.Model" class="text-xs text-yellow-300">
              <strong>{{ w.Model }}</strong>: {{ w.Tokens.toLocaleString() }} tokens
              ({{ w.Exceeds ? 'exceeds max ' + w.Max.toLocaleString() : 'approaching limit ' + w.WarnAt.toLocaleString() }})
            </li>
          </ul>
        </div>

        <button
          @click="doExport"
          :disabled="exportStore.isExporting || !outputPath.trim()"
          class="w-full py-2.5 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-40 rounded text-sm font-medium transition-colors"
        >
          {{ exportStore.isExporting ? 'Exporting…' : '📤 Export File' }}
        </button>

        <p v-if="successMsg" class="text-green-400 text-sm text-center">{{ successMsg }}</p>
        <p v-if="exportStore.error" class="text-red-400 text-sm text-center">{{ exportStore.error }}</p>
      </div>
    </template>
  </div>
</template>
