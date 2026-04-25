<script setup lang="ts">
import { ref } from 'vue'
import { useProjectStore } from '../stores/project'

const emit = defineEmits<{ analyzed: [] }>()
const store = useProjectStore()
const path = ref('')

async function analyze() {
  if (!path.value.trim()) return
  await store.analyzeProject(path.value.trim())
  emit('analyzed')
}
</script>

<template>
  <div class="flex flex-col items-center justify-center h-full gap-6 p-8">
    <h2 class="text-2xl font-semibold text-indigo-300">Analyze Project</h2>
    <p class="text-gray-400 text-sm max-w-md text-center">
      Enter the path to your project directory. ContextForge will auto-detect your tech stack.
    </p>

    <div class="flex gap-2 w-full max-w-lg">
      <input
        v-model="path"
        type="text"
        placeholder="/Users/you/projects/my-app"
        class="flex-1 bg-gray-800 border border-gray-600 rounded px-4 py-2 text-sm text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500"
        @keyup.enter="analyze"
      />
      <button
        @click="analyze"
        :disabled="store.isAnalyzing || !path.trim()"
        class="px-5 py-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-40 rounded text-sm font-medium transition-colors"
      >
        {{ store.isAnalyzing ? 'Scanning…' : 'Analyze' }}
      </button>
    </div>

    <p v-if="store.error" class="text-red-400 text-sm">{{ store.error }}</p>

    <div v-if="store.detectionResult" class="bg-gray-800 rounded-lg p-5 w-full max-w-lg space-y-3">
      <div class="flex items-center justify-between">
        <span class="text-sm font-medium text-gray-300">Detected Stack</span>
        <span class="text-xs text-gray-500">
          Confidence: {{ (store.detectionResult.confidence * 100).toFixed(0) }}%
        </span>
      </div>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="tag in store.detectionResult.stack"
          :key="tag"
          class="px-2 py-0.5 bg-indigo-900 text-indigo-300 rounded text-xs font-mono"
        >
          {{ tag }}
        </span>
      </div>
      <div class="text-xs text-gray-500">
        Indicators: {{ store.detectionResult.indicators.join(', ') }}
      </div>
      <button
        @click="$emit('analyzed')"
        class="w-full mt-2 py-2 bg-green-700 hover:bg-green-600 rounded text-sm font-medium transition-colors"
      >
        Continue to Rules →
      </button>
    </div>
  </div>
</template>
