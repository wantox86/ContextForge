import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { storage } from '../../wailsjs/go/models'
import { GetTemplates, GetTemplatesByStack, SaveTemplate } from '../../wailsjs/go/main/App'

type Template = storage.Template

export const useTemplateStore = defineStore('template', () => {
  const templates = ref<Template[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function fetchTemplates(): Promise<void> {
    isLoading.value = true
    error.value = null
    try {
      templates.value = await GetTemplates()
    } catch (e) {
      error.value = `Failed to load templates. Please restart the app and try again.`
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function fetchByStack(stack: string[]): Promise<Template[]> {
    try {
      return await GetTemplatesByStack(stack)
    } catch (e) {
      error.value = `Failed to load templates for stack. Please try again.`
      throw e
    }
  }

  async function saveTemplate(template: Template): Promise<void> {
    try {
      await SaveTemplate(template)
      await fetchTemplates()
    } catch (e) {
      error.value = `Failed to save template. Please check your input and try again.`
      throw e
    }
  }

  function toggleRule(templateId: number, ruleId: number) {
    const tpl = templates.value.find(t => t.ID === templateId)
    if (!tpl) return
    const rule = tpl.Rules.find((r: storage.Rule) => r.ID === ruleId)
    if (rule) rule.Enabled = !rule.Enabled
  }

  function updateRuleContent(templateId: number, ruleId: number, content: string) {
    const tpl = templates.value.find(t => t.ID === templateId)
    if (!tpl) return
    const rule = tpl.Rules.find((r: storage.Rule) => r.ID === ruleId)
    if (rule) rule.Content = content
  }

  function clearError() {
    error.value = null
  }

  return {
    templates,
    isLoading,
    error,
    fetchTemplates,
    fetchByStack,
    saveTemplate,
    toggleRule,
    updateRuleContent,
    clearError,
  }
})
