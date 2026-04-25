<script setup lang="ts">
import { onMounted } from 'vue'
import { useTemplateStore } from '../stores/template'
import { useProjectStore } from '../stores/project'

const templateStore = useTemplateStore()
const projectStore = useProjectStore()

onMounted(async () => {
  if (projectStore.detectionResult?.stack.length) {
    const matched = await templateStore.fetchByStack(projectStore.detectionResult.stack)
    // Merge matched templates into store if not already loaded
    if (!templateStore.templates.length) {
      await templateStore.fetchTemplates()
    }
  } else {
    await templateStore.fetchTemplates()
  }
})
</script>

<template>
  <div class="h-full overflow-y-auto p-6 space-y-6">
    <h2 class="text-xl font-semibold text-indigo-300">Rule Builder</h2>

    <p v-if="templateStore.isLoading" class="text-gray-400 text-sm">Loading templates…</p>
    <p v-else-if="templateStore.error" class="text-red-400 text-sm">{{ templateStore.error }}</p>
    <p v-else-if="!templateStore.templates.length" class="text-gray-500 text-sm">
      No templates found. Analyze a project first or add templates.
    </p>

    <div v-for="tpl in templateStore.templates" :key="tpl.ID" class="bg-gray-800 rounded-lg p-5 space-y-3">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="font-medium text-white">{{ tpl.Name }}</h3>
          <p class="text-xs text-gray-400 mt-0.5">{{ tpl.Description }}</p>
        </div>
        <div class="flex flex-wrap gap-1">
          <span
            v-for="tag in tpl.Stack"
            :key="tag"
            class="px-1.5 py-0.5 bg-gray-700 text-gray-300 rounded text-xs font-mono"
          >{{ tag }}</span>
        </div>
      </div>

      <div class="space-y-2">
        <div
          v-for="rule in tpl.Rules"
          :key="rule.ID"
          class="flex items-start gap-3 p-3 rounded bg-gray-750 border border-gray-700"
        >
          <input
            type="checkbox"
            :checked="rule.Enabled"
            @change="templateStore.toggleRule(tpl.ID, rule.ID)"
            class="mt-0.5 accent-indigo-500"
          />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <span class="text-sm font-medium text-white">{{ rule.Title }}</span>
              <span class="text-xs px-1.5 py-0.5 bg-gray-700 text-gray-400 rounded">{{ rule.Category }}</span>
            </div>
            <p class="text-xs text-gray-400 mt-1 truncate">{{ rule.Content }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
