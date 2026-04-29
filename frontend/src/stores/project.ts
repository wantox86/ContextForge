import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { DetectionResult } from '../types'
import { storage } from '../../wailsjs/go/models'
import { AnalyzeProject, ExportToFile, GetProjectByPath, PreviewExport } from '../../wailsjs/go/main/App'

export const useProjectStore = defineStore('project', () => {
  const currentProject = ref<storage.Project | null>(null)
  const detectionResult = ref<DetectionResult | null>(null)
  const isAnalyzing = ref(false)
  const error = ref<string | null>(null)

  async function analyzeProject(path: string): Promise<DetectionResult> {
    isAnalyzing.value = true
    error.value = null
    try {
      const result = await AnalyzeProject(path)
      detectionResult.value = result
      // Fetch the persisted project so we have its DB ID for preview/export
      const project = await GetProjectByPath(path)
      currentProject.value = project
      return result
    } catch (e) {
      error.value = `Failed to analyze project. Please check the path and try again.`
      throw e
    } finally {
      isAnalyzing.value = false
    }
  }

  async function previewExport(projectId: number, target: string): Promise<string> {
    try {
      return await PreviewExport(projectId, target)
    } catch (e) {
      error.value = `Failed to generate preview. Please try again.`
      throw e
    }
  }

  async function exportToFile(projectId: number, target: string, outputPath: string): Promise<void> {
    try {
      await ExportToFile(projectId, target, outputPath)
    } catch (e) {
      error.value = `Failed to export file. Check that the output path exists and is writable.`
      throw e
    }
  }

  function setCurrentProject(project: storage.Project) {
    currentProject.value = project
  }

  function clearError() {
    error.value = null
  }

  return {
    currentProject,
    detectionResult,
    isAnalyzing,
    error,
    analyzeProject,
    previewExport,
    exportToFile,
    setCurrentProject,
    clearError,
  }
})
