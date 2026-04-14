<script setup>
import { ref, computed, onMounted } from 'vue'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Skeleton from 'primevue/skeleton'
import Message from 'primevue/message'
import Button from 'primevue/button'
import Chart from 'primevue/chart'
import { statusSeverity, formatTime } from '../composables/useBackupUtils.js'

const report = ref(null)
const loading = ref(true)
const error = ref(null)
const selectedHours = ref(24)
const useCustomRange = ref(false)
const rangeStart = ref('')
const rangeEnd = ref('')
const windowLabel = ref('Last 24 Hours')

onMounted(async () => {
  const now = new Date()
  const end = new Date(now.getTime() - now.getTimezoneOffset() * 60000)
  const start = new Date(end.getTime() - selectedHours.value * 60 * 60 * 1000)
  rangeEnd.value = end.toISOString().slice(0, 16)
  rangeStart.value = start.toISOString().slice(0, 16)
  await loadReport()
})

function buildQueryParams() {
  const params = new URLSearchParams()
  if (useCustomRange.value && rangeStart.value && rangeEnd.value) {
    const from = new Date(rangeStart.value)
    const to = new Date(rangeEnd.value)
    if (!Number.isNaN(from.getTime()) && !Number.isNaN(to.getTime())) {
      params.set('from', from.toISOString())
      params.set('to', to.toISOString())
    }
  } else {
    params.set('hours', String(selectedHours.value))
  }
  const query = params.toString()
  return query ? `?${query}` : ''
}

function updateWindowLabel() {
  if (useCustomRange.value && rangeStart.value && rangeEnd.value) {
    windowLabel.value = `${formatTime(new Date(rangeStart.value))} -> ${formatTime(new Date(rangeEnd.value))}`
    return
  }
  windowLabel.value = `Last ${selectedHours.value} Hours`
}

async function loadReport() {
  loading.value = true
  error.value = null
  try {
    updateWindowLabel()
    const res = await fetch(`/api/v1/report${buildQueryParams()}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    report.value = await res.json()
  } catch (e) {
    error.value = `Failed to load report: ${e.message}`
  } finally {
    loading.value = false
  }
}

// Summary chart data
const summaryChartData = computed(() => {
  if (!report.value) return null
  const s = report.value.summary
  return {
    labels: ['Completed', 'Failed', 'Partially Failed', 'In Progress', 'Deleting', 'Other'],
    datasets: [{
      data: [s.completed, s.failed, s.partiallyFailed, s.inProgress, s.deleting, s.total - s.completed - s.failed - s.partiallyFailed - s.inProgress - s.deleting],
      backgroundColor: ['#10b981', '#ef4444', '#f59e0b', '#3b82f6', '#a78bfa', '#6b7280'],
      borderColor: ['#10b981', '#ef4444', '#f59e0b', '#3b82f6', '#a78bfa', '#6b7280'],
      borderWidth: 2,
    }],
  }
})

const summaryChartOptions = {
  responsive: true,
  maintainAspectRatio: true,
  plugins: {
    legend: {
      labels: {
        font: { size: 12 },
        padding: 15,
      },
    },
  },
}

// Schedule success rate chart
const scheduleChartData = computed(() => {
  if (!report.value) return null
  const schedules = report.value.schedules
  return {
    labels: schedules.map(s => s.scheduleName),
    datasets: [{
      label: 'Success Rate (%)',
      data: schedules.map(s => s.successRate),
      backgroundColor: '#3b82f6',
      borderColor: '#1e40af',
      borderWidth: 2,
    }],
  }
})

const scheduleChartOptions = {
  responsive: true,
  maintainAspectRatio: true,
  indexAxis: 'y',
  plugins: {
    legend: {
      display: false,
    },
  },
  scales: {
    x: {
      beginAtZero: true,
      max: 100,
    },
  },
}

function getPeriodOrder() {
  return ['Last 24 Hours', 'Last 7 Days', 'Last 30 Days']
}

function getStatusBadge(status) {
  const severity = statusSeverity(status)
  return severity
}

function formatRate(rate) {
  return rate.toFixed(1) + '%'
}

function formatBytes(bytes) {
  if (bytes === 0) return '0'
  const sizes = ['', 'K', 'M', 'G', 'T']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return +(bytes / Math.pow(1024, i)).toFixed(2) + sizes[i]
}

function exportReport() {
  const reportData = JSON.stringify(report.value, null, 2)
  const blob = new Blob([reportData], { type: 'application/json' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `backup-report-${new Date().toISOString().split('T')[0]}.json`
  document.body.appendChild(a)
  a.click()
  window.URL.revokeObjectURL(url)
  document.body.removeChild(a)
}

function exportPDF() {
  window.open(`/api/v1/report/pdf${buildQueryParams()}`, '_blank')
}
</script>

<template>
  <!-- Loading skeleton -->
  <div v-if="loading">
    <Skeleton width="10rem" height="2rem" class="mb-4" />
    <Skeleton height="30rem" class="mb-4" />
    <Skeleton height="20rem" />
  </div>

  <!-- Error -->
  <Message v-else-if="error" severity="error" :closable="false">{{ error }}</Message>

  <!-- Content -->
  <div v-else>
    <!-- Header -->
    <div class="report-header">
      <div>
        <h2>Backup Report</h2>
        <p style="color: var(--p-text-color-secondary); margin: 0.5rem 0 0 0;">
          Generated: {{ formatTime(new Date(report.generatedAt)) }}
        </p>
        <p style="color: var(--p-text-color-secondary); margin: 0.25rem 0 0 0;">
          Window: {{ windowLabel }}
        </p>
      </div>
      <div class="report-actions">
        <div class="window-controls">
          <label class="window-toggle">
            <input v-model="useCustomRange" type="checkbox">
            Use custom range
          </label>
          <template v-if="useCustomRange">
            <input v-model="rangeStart" type="datetime-local" class="window-input">
            <input v-model="rangeEnd" type="datetime-local" class="window-input">
          </template>
          <template v-else>
            <select v-model.number="selectedHours" class="window-select">
              <option :value="6">Last 6 hours</option>
              <option :value="12">Last 12 hours</option>
              <option :value="24">Last 24 hours</option>
              <option :value="48">Last 48 hours</option>
              <option :value="72">Last 72 hours</option>
              <option :value="168">Last 7 days</option>
            </select>
          </template>
          <Button label="Apply" icon="pi pi-refresh" size="small" @click="loadReport" />
        </div>
        <Button
          label="Export PDF"
          icon="pi pi-file-pdf"
          severity="danger"
          outlined
          @click="exportPDF"
        />
        <Button
          label="Export JSON"
          icon="pi pi-download"
          severity="secondary"
          @click="exportReport"
        />
      </div>
    </div>

    <!-- Overall Summary Stats -->
    <div class="grid summary-stats">
      <Card class="col-12 md:col-6 lg:col-3">
        <template #content>
          <div class="stat-card">
            <i class="pi pi-box" style="color: var(--p-primary-color); font-size: 1.5rem;"></i>
            <div class="stat-data">
              <div class="stat-value">{{ report.summary.total }}</div>
              <div class="stat-label">Total Backups</div>
            </div>
          </div>
        </template>
      </Card>

      <Card class="col-12 md:col-6 lg:col-3">
        <template #content>
          <div class="stat-card">
            <i class="pi pi-check-circle" style="color: #10b981; font-size: 1.5rem;"></i>
            <div class="stat-data">
              <div class="stat-value">{{ report.summary.completed }}</div>
              <div class="stat-label">Completed</div>
            </div>
          </div>
        </template>
      </Card>

      <Card class="col-12 md:col-6 lg:col-3">
        <template #content>
          <div class="stat-card">
            <i class="pi pi-times-circle" style="color: #ef4444; font-size: 1.5rem;"></i>
            <div class="stat-data">
              <div class="stat-value">{{ report.summary.failed }}</div>
              <div class="stat-label">Failed</div>
            </div>
          </div>
        </template>
      </Card>

      <Card class="col-12 md:col-6 lg:col-3">
        <template #content>
          <div class="stat-card">
            <i class="pi pi-exclamation-triangle" style="color: #f59e0b; font-size: 1.5rem;"></i>
            <div class="stat-data">
              <div class="stat-value">{{ report.summary.partiallyFailed }}</div>
              <div class="stat-label">Partially Failed</div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Charts Section -->
    <div class="grid charts-section" style="margin-top: 2rem;">
      <Card class="col-12 lg:col-6">
        <template #header>
          <div style="padding: 1rem; font-weight: 600; font-size: 1rem;">Backup Status Distribution</div>
        </template>
        <template #content>
          <Chart v-if="summaryChartData" type="pie" :data="summaryChartData" :options="summaryChartOptions" />
        </template>
      </Card>

      <Card v-if="report.schedules.length > 0" class="col-12 lg:col-6">
        <template #header>
          <div style="padding: 1rem; font-weight: 600; font-size: 1rem;">Success Rate by Schedule</div>
        </template>
        <template #content>
          <Chart v-if="scheduleChartData" type="bar" :data="scheduleChartData" :options="scheduleChartOptions" />
        </template>
      </Card>
    </div>

    <!-- Time Period Summaries -->
    <div style="margin-top: 2rem;">
      <h3 style="margin-bottom: 1rem;">Performance Over Time</h3>
      <div class="grid">
        <Card v-for="period in getPeriodOrder()" :key="period" class="col-12 md:col-4">
          <template #header>
            <div style="padding: 0.75rem 1rem; font-weight: 600; background: var(--p-surface-50);">{{ period }}</div>
          </template>
          <template #content>
            <div v-if="report.periodSummaries[period]" class="period-summary">
              <div class="period-row">
                <span class="period-label">Total:</span>
                <span class="period-value">{{ report.periodSummaries[period].totalBackups }}</span>
              </div>
              <div class="period-row">
                <span class="period-label">Completed:</span>
                <Tag severity="success" :value="`${report.periodSummaries[period].completed}`" />
              </div>
              <div class="period-row">
                <span class="period-label">Failed:</span>
                <Tag severity="danger" :value="`${report.periodSummaries[period].failed}`" />
              </div>
              <div class="period-row">
                <span class="period-label">Partial:</span>
                <Tag severity="warning" :value="`${report.periodSummaries[period].partiallyFailed}`" />
              </div>
              <div class="period-row">
                <span class="period-label">Avg Duration:</span>
                <span class="period-value">{{ report.periodSummaries[period].averageDuration }}</span>
              </div>
              <div class="period-row">
                <span class="period-label">Items:</span>
                <span class="period-value">{{ report.periodSummaries[period].totalItems }}</span>
              </div>
            </div>
          </template>
        </Card>
      </div>
    </div>

    <!-- Schedules Summary -->
    <div v-if="report.schedules.length > 0" style="margin-top: 2rem;">
      <h3 style="margin-bottom: 1rem;">Schedule Summary</h3>
      <Card>
        <template #content>
          <DataTable
            :value="report.schedules"
            stripedRows
            size="small"
            responsiveLayout="scroll"
          >
            <Column field="scheduleName" header="Schedule Name" style="width: 25%;">
              <template #body="{ data }">
                <strong>{{ data.scheduleName }}</strong>
              </template>
            </Column>
            <Column field="totalBackups" header="Total" style="width: 15%;" />
            <Column field="successfulBackups" header="Successful" style="width: 15%;">
              <template #body="{ data }">
                <Tag severity="success" :value="`${data.successfulBackups}`" />
              </template>
            </Column>
            <Column field="failedBackups" header="Failed" style="width: 15%;">
              <template #body="{ data }">
                <Tag severity="danger" :value="`${data.failedBackups}`" />
              </template>
            </Column>
            <Column field="successRate" header="Success Rate" style="width: 20%;">
              <template #body="{ data }">
                <div :style="{ color: data.successRate >= 90 ? '#10b981' : data.successRate >= 70 ? '#f59e0b' : '#ef4444' }">
                  {{ formatRate(data.successRate) }}
                </div>
              </template>
            </Column>
            <Column field="lastBackupTime" header="Last Backup" style="width: 20%;">
              <template #body="{ data }">
                <div v-if="data.lastBackupTime" style="font-size: 0.85rem;">
                  {{ formatTime(new Date(data.lastBackupTime)) }}
                </div>
                <div v-else style="color: var(--p-text-color-secondary);">-</div>
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>

    <!-- Detailed Backups List -->
    <div style="margin-top: 2rem;">
      <h3 style="margin-bottom: 1rem;">All Backups (Sorted by Date - Newest First)</h3>
      <Card>
        <template #content>
          <DataTable
            :value="report.backups"
            stripedRows
            paginator
            :rows="20"
            :rowsPerPageOptions="[10, 20, 50]"
            size="small"
            responsiveLayout="scroll"
            sortField="startTimestamp"
            :sortOrder="-1"
          >
            <Column field="name" header="Backup Name" style="width: 20%;">
              <template #body="{ data }">
                <strong>{{ data.name }}</strong>
              </template>
            </Column>
            <Column field="scheduleName" header="Schedule" style="width: 15%;">
              <template #body="{ data }">
                <div v-if="data.scheduleName">{{ data.scheduleName }}</div>
                <div v-else style="color: var(--p-text-color-secondary);">Manual</div>
              </template>
            </Column>
            <Column field="status" header="Status" style="width: 12%;">
              <template #body="{ data }">
                <Tag :severity="getStatusBadge(data.status)" :value="data.status" />
              </template>
            </Column>
            <Column field="startTimestamp" header="Start Time" style="width: 15%;">
              <template #body="{ data }">
                <div v-if="data.startTimestamp" style="font-size: 0.85rem;">
                  {{ formatTime(new Date(data.startTimestamp)) }}
                </div>
                <div v-else style="color: var(--p-text-color-secondary);">-</div>
              </template>
            </Column>
            <Column field="duration" header="Duration" style="width: 12%;">
              <template #body="{ data }">
                <div v-if="data.duration" style="font-size: 0.85rem;">{{ data.duration }}</div>
                <div v-else style="color: var(--p-text-color-secondary);">-</div>
              </template>
            </Column>
            <Column field="itemsBackedUp" header="Items" style="width: 10%;">
              <template #body="{ data }">
                <div style="font-size: 0.85rem;">
                  {{ data.itemsBackedUp }} / {{ data.itemsBackedUp + (data.Warnings || 0) + (data.Errors || 0) }}
                </div>
              </template>
            </Column>
            <Column field="warnings" header="Warnings" style="width: 8%;">
              <template #body="{ data }">
                <Tag v-if="data.warnings > 0" severity="warning" :value="`${data.warnings}`" />
                <span v-else style="color: var(--p-text-color-secondary);">0</span>
              </template>
            </Column>
            <Column field="errors" header="Errors" style="width: 8%;">
              <template #body="{ data }">
                <Tag v-if="data.errors > 0" severity="danger" :value="`${data.errors}`" />
                <span v-else style="color: var(--p-text-color-secondary);">0</span>
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.report-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  gap: 1rem;
}

.report-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.window-controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.window-toggle {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.875rem;
  color: var(--p-text-color-secondary);
}

.window-select,
.window-input {
  border: 1px solid var(--p-surface-border);
  border-radius: 6px;
  padding: 0.4rem 0.55rem;
  background: var(--p-surface-0);
  color: var(--p-text-color);
  font-size: 0.875rem;
}

.summary-stats {
  margin-bottom: 2rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.stat-data {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--p-text-color);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--p-text-color-secondary);
  margin-top: 0.25rem;
}

.charts-section {
  margin-bottom: 2rem;
}

.period-summary {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.period-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  border-bottom: 1px solid var(--p-surface-border);
}

.period-row:last-child {
  border-bottom: none;
}

.period-label {
  font-weight: 500;
  color: var(--p-text-color-secondary);
  font-size: 0.875rem;
}

.period-value {
  font-weight: 600;
  color: var(--p-text-color);
}

h3 {
  margin: 0;
  color: var(--p-text-color);
  font-size: 1.125rem;
  font-weight: 600;
}

h2 {
  margin: 0;
  color: var(--p-text-color);
  font-size: 1.75rem;
  font-weight: 700;
}

p {
  margin: 0;
}

:deep(.p-card) {
  border: 1px solid var(--p-surface-border);
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

:deep(.p-datatable-thead > tr > th) {
  background: var(--p-surface-100);
  font-weight: 600;
}

:deep(.p-tag) {
  font-size: 0.75rem;
}
</style>
