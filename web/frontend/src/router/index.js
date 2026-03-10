import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import BackupsListView from '../views/BackupsListView.vue'
import BackupDetailView from '../views/BackupDetailView.vue'

const routes = [
  { path: '/', name: 'dashboard', component: DashboardView },
  { path: '/backups', name: 'backups', component: BackupsListView },
  { path: '/backups/:name', name: 'backup-detail', component: BackupDetailView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
