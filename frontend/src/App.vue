<script setup lang="ts">
import { ref } from 'vue'
import ProjectAnalyzer from './components/ProjectAnalyzer.vue'
import RuleBuilder from './components/RuleBuilder.vue'
import PreviewPanel from './components/PreviewPanel.vue'
import ExportPanel from './components/ExportPanel.vue'

type Tab = 'analyze' | 'rules' | 'preview' | 'export'
const activeTab = ref<Tab>('analyze')

const tabs: { id: Tab; label: string }[] = [
  { id: 'analyze', label: '🔍 Analyze' },
  { id: 'rules',   label: '📋 Rules' },
  { id: 'preview', label: '👁 Preview' },
  { id: 'export',  label: '📤 Export' },
]
</script>

<template>
  <div class="flex flex-col h-screen bg-gray-900 text-white">
    <!-- Header -->
    <header class="flex items-center justify-between px-6 py-3 bg-gray-800 border-b border-gray-700 shrink-0">
      <div class="flex items-center gap-3">
        <span class="text-xl font-bold text-indigo-400">⚡ ContextForge</span>
        <span class="text-xs text-gray-400">v0.1.0</span>
      </div>
      <nav class="flex gap-1">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            'px-4 py-1.5 rounded text-sm font-medium transition-colors',
            activeTab === tab.id
              ? 'bg-indigo-600 text-white'
              : 'text-gray-400 hover:text-white hover:bg-gray-700',
          ]"
        >
          {{ tab.label }}
        </button>
      </nav>
    </header>

    <!-- Main content -->
    <main class="flex-1 overflow-hidden">
      <ProjectAnalyzer v-if="activeTab === 'analyze'" @analyzed="activeTab = 'rules'" />
      <RuleBuilder     v-else-if="activeTab === 'rules'" />
      <PreviewPanel    v-else-if="activeTab === 'preview'" />
      <ExportPanel     v-else-if="activeTab === 'export'" />
    </main>
  </div>
</template>

