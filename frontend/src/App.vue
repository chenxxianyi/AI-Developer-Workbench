<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppShell from '@/components/layout/AppShell.vue'

const route = useRoute()

// Determine layout based on route meta
const useAppShell = computed(() => route.meta.layout === 'app')
</script>

<template>
  <!-- Landing page uses standalone layout -->
  <router-view v-if="!useAppShell" />

  <!-- Dashboard pages use AppShell wrapper -->
  <AppShell v-else>
    <router-view />
  </AppShell>
</template>
