<script setup>
import { ref, computed } from 'vue'
import { RouterView, useRouter, useRoute } from 'vue-router'
import Menubar from 'primevue/menubar'
import Button from 'primevue/button'

const router = useRouter()
const route = useRoute()

const isDark = ref(document.documentElement.classList.contains('app-dark'))

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('app-dark')
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

const menuItems = computed(() => [
  {
    label: 'Dashboard',
    icon: 'pi pi-gauge',
    class: route.path === '/' ? 'active-route' : '',
    command: () => router.push('/'),
  },
  {
    label: 'Backups',
    icon: 'pi pi-list',
    class: route.path.startsWith('/backups') ? 'active-route' : '',
    command: () => router.push('/backups'),
  },
  {
    label: 'Report',
    icon: 'pi pi-file-pdf',
    class: route.path === '/report' ? 'active-route' : '',
    command: () => router.push('/report'),
  },
])
</script>

<template>
  <Menubar :model="menuItems">
    <template #start>
      <span class="font-bold text-lg cursor-pointer" @click="router.push('/')">Velero Backup Reporter</span>
    </template>
    <template #end>
      <Button
        :icon="isDark ? 'pi pi-sun' : 'pi pi-moon'"
        @click="toggleTheme"
        :title="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
        text
        rounded
      />
    </template>
  </Menubar>

  <main>
    <RouterView />
  </main>

  <footer>
    <small>Velero Backup Reporter</small>
  </footer>
</template>

<style>
:root {
  --p-font-family: 'Inter Variable', 'Inter', system-ui, -apple-system, sans-serif;
  --font-mono: 'JetBrains Mono Variable', 'JetBrains Mono', ui-monospace, 'Cascadia Code', 'Fira Code', monospace;
}
body {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  margin: 0;
  font-family: var(--p-font-family);
  font-size: 0.9375rem;
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-rendering: optimizeLegibility;
}
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
main {
  flex: 1;
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}
footer {
  text-align: center;
  padding: 1rem;
  margin-top: auto;
  border-top: 1px solid var(--p-content-border-color);
  color: var(--p-text-muted-color);
}
.active-route {
  font-weight: 600;
}
</style>
