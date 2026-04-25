import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { TokenWarning } from '../types'
import { CountTokens, CheckTokenWarnings } from '../../wailsjs/go/main/App'

export const useExportStore = defineStore('export', () => {
  const previewContent = ref('')
  const selectedTarget = ref<'claude' | 'copilot'>('claude')
  const isExporting = ref(false)
  const exportSuccess = ref(false)
  const error = ref<string | null>(null)

  const tokenCount = computed(() => {
    if (!previewContent.value) return 0
    return Math.ceil(previewContent.value.length / 4)
  })

  const warnings = ref<TokenWarning[]>([])

  async function refreshWarnings(): Promise<void> {
    if (!previewContent.value) {
      warnings.value = []
      return
    }
    try {
      warnings.value = await CheckTokenWarnings(previewContent.value)
    } catch {
      warnings.value = []
    }
  }

  function setPreview(content: string) {
    previewContent.value = content
    refreshWarnings()
  }

  function setTarget(target: 'claude' | 'copilot') {
    selectedTarget.value = target
  }

  function setExporting(val: boolean) {
    isExporting.value = val
  }

  function setSuccess(val: boolean) {
    exportSuccess.value = val
  }

  function clearError() {
    error.value = null
  }

  return {
    previewContent,
    selectedTarget,
    isExporting,
    exportSuccess,
    tokenCount,
    warnings,
    error,
    setPreview,
    setTarget,
    setExporting,
    setSuccess,
    clearError,
    refreshWarnings,
  }
})
