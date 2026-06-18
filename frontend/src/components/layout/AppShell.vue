<script setup lang="ts">
/**
 * App Shell Layout
 * Wrapper for all dashboard pages
 * Includes Sidebar, Header, and main content area
 */

import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from './Sidebar.vue'
import AppHeader from './AppHeader.vue'

const route = useRoute()
const mobileSidebarOpen = ref(false)

// Close mobile sidebar on route change
watch(
  () => route.path,
  () => {
    mobileSidebarOpen.value = false
  }
)

function toggleMobileSidebar() {
  mobileSidebarOpen.value = !mobileSidebarOpen.value
}

function closeMobileSidebar() {
  mobileSidebarOpen.value = false
}
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- Sidebar -->
    <Sidebar
      :mobile-open="mobileSidebarOpen"
      @toggle="toggleMobileSidebar"
      @close="closeMobileSidebar"
    />

    <!-- Main Content Area -->
    <main class="md:ml-64 min-h-screen">
      <!-- Header -->
      <AppHeader @toggle-sidebar="toggleMobileSidebar" />

      <!-- Page Content -->
      <div class="max-w-content mx-auto px-4 md:px-8 py-8">
        <slot />
      </div>
    </main>
  </div>
</template>

<style scoped>
/* Desktop sidebar offset */
@media (min-width: 768px) {
  main {
    margin-left: 256px; /* w-64 */
  }
}
</style>