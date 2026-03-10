<script setup>
import { ref, onMounted } from 'vue'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Skeleton from 'primevue/skeleton'
import Message from 'primevue/message'
import Button from 'primevue/button'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'
import { statusSeverity, formatTime } from '../composables/useBackupUtils.js'

const toast = useToast()
const data = ref(null)
const loading = ref(true)
const error = ref(null)
const sendingTestEmail = ref(false)

const statusCards = [
  { key: 'total', label: 'Total', icon: 'pi pi-box', severity: 'primary' },
  { key: 'completed', label: 'Completed', icon: 'pi pi-check-circle', severity: 'success' },
  { key: 'failed', label: 'Failed', icon: 'pi pi-times-circle', severity: 'danger' },
  { key: 'partiallyFailed', label: 'Partially Failed', icon: 'pi pi-exclamation-triangle', severity: 'warn' },
  { key: 'inProgress', label: 'In Progress', icon: 'pi pi-spin pi-spinner', severity: 'info' },
  { key: 'deleting', label: 'Deleting', icon: 'pi pi-trash', severity: 'secondary' },
]

onMounted(async () => {
  try {
    const res = await fetch('/api/v1/dashboard')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    data.value = await res.json()
  } catch (e) {
    error.value = `Failed to load dashboard: ${e.message}`
  } finally {
    loading.value = false
  }
})

function formatRate(rate) {
  return rate.toFixed(1) + '%'
}

async function sendTestEmail() {
  sendingTestEmail.value = true
  try {
    const res = await fetch('/api/v1/email/test', { method: 'POST' })
    const body = await res.json()
    if (!res.ok) {
      throw new Error(body.error || `HTTP ${res.status}`)
    }
    toast.add({ severity: 'success', summary: 'Email Sent', detail: body.message, life: 5000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Email Failed', detail: e.message, life: 8000 })
  } finally {
    sendingTestEmail.value = false
  }
}
</script>

<template>
  <!-- Loading skeleton -->
  <div v-if="loading">
    <Skeleton width="10rem" height="2rem" class="mb-4" />
    <div class="grid">
      <div v-for="i in 6" :key="i" class="col-12 md:col-4 lg:col-2">
        <Skeleton height="8rem" />
      </div>
    </div>
  </div>

  <!-- Error -->
  <Message v-else-if="error" severity="error" :closable="false">{{ error }}</Message>

  <!-- Content -->
  <div v-else>
    <Toast />
    <div class="dashboard-header">
      <h2>Dashboard</h2>
      <Button
        v-if="data.emailEnabled"
        label="Send Test Email"
        icon="pi pi-envelope"
        outlined
        :loading="sendingTestEmail"
        @click="sendTestEmail"
      />
    </div>

    <div class="status-cards">
      <Card v-for="card in statusCards" :key="card.key" class="status-card">
        <template #content>
          <div class="status-card-content">
            <i :class="card.icon" :style="{ color: 'var(--p-' + card.severity + '-color)', fontSize: '1.5rem' }"></i>
            <div class="status-card-count" :style="{ color: 'var(--p-' + card.severity + '-color)' }">
              {{ data.summary[card.key] }}
            </div>
            <div class="status-card-label">{{ card.label }}</div>
          </div>
        </template>
      </Card>
    </div>

    <Card style="margin-top: 1.5rem;">
      <template #title>
        <span><i class="pi pi-calendar mr-2"></i>Schedule Statistics</span>
      </template>
      <template #content>
        <DataTable
          v-if="data.schedules.length > 0"
          :value="data.schedules"
          stripedRows
          size="small"
          :sortField="'scheduleName'"
          :sortOrder="1"
        >
          <Column field="scheduleName" header="Schedule" sortable />
          <Column field="totalBackups" header="Total" sortable />
          <Column field="successfulBackups" header="Completed" sortable />
          <Column field="failedBackups" header="Failed" sortable />
          <Column field="lastBackupTime" header="Last Backup" sortable>
            <template #body="{ data }">{{ formatTime(data.lastBackupTime) }}</template>
          </Column>
          <Column field="lastBackupStatus" header="Last Status" sortable>
            <template #body="{ data }">
              <Tag v-if="data.lastBackupStatus" :value="data.lastBackupStatus" :severity="statusSeverity(data.lastBackupStatus)" />
              <span v-else>-</span>
            </template>
          </Column>
          <Column field="successRate" header="Success Rate" sortable>
            <template #body="{ data }">{{ formatRate(data.successRate) }}</template>
          </Column>
        </DataTable>
        <div v-else class="empty-state">
          <i class="pi pi-calendar-times"></i>
          <p>No schedules found</p>
        </div>
      </template>
    </Card>
  </div>
</template>

<style scoped>
.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}
.dashboard-header h2 {
  margin: 0;
}
.status-cards {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 1rem;
}
@media (max-width: 992px) {
  .status-cards { grid-template-columns: repeat(3, 1fr); }
}
@media (max-width: 576px) {
  .status-cards { grid-template-columns: repeat(2, 1fr); }
}
.status-card :deep(.p-card-body) {
  padding: 1rem;
}
.status-card :deep(.p-card-content) {
  padding: 0;
}
.status-card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}
.status-card-count {
  font-size: 2rem;
  font-weight: 700;
  line-height: 1;
}
.status-card-label {
  font-size: 0.85rem;
  color: var(--p-text-muted-color);
}
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: var(--p-text-muted-color);
}
.empty-state i {
  font-size: 2.5rem;
  margin-bottom: 0.75rem;
  display: block;
}
.mr-2 {
  margin-right: 0.5rem;
}
</style>
