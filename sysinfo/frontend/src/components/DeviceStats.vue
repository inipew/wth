<!-- DeviceStats.vue -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-100 to-gray-200 py-6 px-4 sm:px-6 lg:px-8">
    <div class="max-w-7xl mx-auto">
      <h1 class="text-3xl font-bold text-gray-900 mb-8 text-center">Device Statistics Dashboard</h1>
      <div v-if="loading" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-blue-500"></div>
      </div>
      <div v-else-if="error" class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-6" role="alert">
        <p class="font-bold">Error</p>
        <p>{{ error }}</p>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div class="bg-white rounded-lg shadow-md p-6 transition duration-300 ease-in-out transform hover:scale-105">
          <h2 class="text-xl font-semibold mb-4 text-gray-800">System Info</h2>
          <div class="space-y-2">
            <p><span class="font-medium text-gray-600">OS:</span> {{ stats.goos }}</p>
            <p><span class="font-medium text-gray-600">Architecture:</span> {{ stats.goarch }}</p>
            <p><span class="font-medium text-gray-600">CPUs:</span> {{ stats.num_cpu }}</p>
            <p><span class="font-medium text-gray-600">Last Boot:</span> {{ formatDate(stats.last_boot) }}</p>
            <p><span class="font-medium text-gray-600">Uptime:</span> {{ stats.uptime }}</p>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 transition duration-300 ease-in-out transform hover:scale-105">
          <h2 class="text-xl font-semibold mb-4 text-gray-800">Memory Usage</h2>
          <div class="space-y-2">
            <p><span class="font-medium text-gray-600">Total:</span> {{ stats.mem_stats.total_mb }}</p>
            <p><span class="font-medium text-gray-600">Used:</span> {{ stats.mem_stats.used_mb }}</p>
            <p><span class="font-medium text-gray-600">Free:</span> {{ stats.mem_stats.free_mb }}</p>
          </div>
          <div class="mt-4">
            <div class="relative pt-1">
              <div class="flex mb-2 items-center justify-between">
                <div>
                  <span class="text-xs font-semibold inline-block py-1 px-2 uppercase rounded-full text-blue-600 bg-blue-200">
                    Usage
                  </span>
                </div>
                <div class="text-right">
                  <span class="text-xs font-semibold inline-block text-blue-600">
                    {{ stats.mem_stats.usage_percent.toFixed(2) }}%
                  </span>
                </div>
              </div>
              <div class="overflow-hidden h-2 mb-4 text-xs flex rounded bg-blue-200">
                <div :style="{ width: `${stats.mem_stats.usage_percent}%` }" class="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-blue-500"></div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 transition duration-300 ease-in-out transform hover:scale-105">
          <h2 class="text-xl font-semibold mb-4 text-gray-800">CPU Usage</h2>
          <div class="space-y-2">
            <p><span class="font-medium text-gray-600">Cores:</span> {{ stats.cpu_stats.cores }}</p>
          </div>
          <div class="mt-4">
            <div class="relative pt-1">
              <div class="flex mb-2 items-center justify-between">
                <div>
                  <span class="text-xs font-semibold inline-block py-1 px-2 uppercase rounded-full text-green-600 bg-green-200">
                    Usage
                  </span>
                </div>
                <div class="text-right">
                  <span class="text-xs font-semibold inline-block text-green-600">
                    {{ stats.cpu_stats.usage_percent.toFixed(2) }}%
                  </span>
                </div>
              </div>
              <div class="overflow-hidden h-2 mb-4 text-xs flex rounded bg-green-200">
                <div :style="{ width: `${stats.cpu_stats.usage_percent}%` }" class="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-green-500"></div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 col-span-full transition duration-300 ease-in-out transform hover:scale-105">
          <h2 class="text-xl font-semibold mb-4 text-gray-800">Network Adapters</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <div v-for="adapter in stats.network_stats.adapters" :key="adapter.name" class="bg-gray-50 p-4 rounded-lg">
              <p class="font-semibold text-gray-700">{{ adapter.name }}</p>
              <p class="text-sm text-gray-600">{{ adapter.ip_address }}</p>
              <p class="text-sm text-gray-600">Upload: {{ adapter.upload }}</p>
              <p class="text-sm text-gray-600">Download: {{ adapter.download }}</p>
            </div>
          </div>
        </div>
      </div>
      <div class="mt-8 text-center">
        <p class="text-lg font-semibold text-gray-800">Number of Processes: <span class="text-blue-600">{{ stats.num_processes }}</span></p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'

export default {
  name: 'DeviceStats',
  setup() {
    const stats = ref(null)
    const loading = ref(true)
    const error = ref(null)

    const fetchStats = async () => {
      try {
        const response = await fetch('/api/device-stats')
        if (!response.ok) {
          throw new Error('Failed to fetch device stats')
        }
        stats.value = await response.json()
        loading.value = false
      } catch (err) {
        error.value = err.message
        loading.value = false
      }
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleString()
    }

    onMounted(fetchStats)

    return {
      stats,
      loading,
      error,
      formatDate
    }
  }
}
</script>

<style scoped>
@import 'https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css';
</style>