<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import Skeleton from 'primevue/skeleton'
import Message from 'primevue/message'
import { statusSeverity, formatTime } from '../composables/useBackupUtils.js'

const router = useRouter()
const backups = ref([])
const loading = ref(true)
const error = ref(null)
const filterStatus = ref(null)
const filterSchedule = ref(null)

onMounted(async () => {
  try {
    const res = await fetch('/api/v1/backups')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    backups.value = await res.json()
  } catch (e) {
    error.value = `Failed to load backups: ${e.message}`
  } finally {
    loading.value = false
  }
})

const scheduleOptions = computed(() => {
  const set = new Set(backups.value.map(b => b.scheduleName).filter(Boolean))
  return [...set].sort().map(s => ({ label: s, value: s }))
})

const statusOptions = computed(() => {
  const set = new Set(backups.value.map(b => b.status).filter(Boolean))
  return [...set].sort().map(s => ({ label: s, value: s }))
})

const filteredBackups = computed(() => {
  let result = backups.value
  if (filterStatus.value) {
    result = result.filter(b => b.status === filterStatus.value)
  }
  if (filterSchedule.value) {
    result = result.filter(b => b.scheduleName === filterSchedule.value)
  }
  return result
})

function goToDetail(name) {
  router.push({ name: 'backup-detail', params: { name } })
}
</script>

<template>
  <!-- Loading -->
  <div v-if="loading">
    <Skeleton width="8rem" height="2rem" class="mb-4" />
    <Skeleton height="20rem" />
  </div>

  <!-- Error -->
  <Message v-else-if="error" severity="error" :closable="false">{{ error }}</Message>

  <!-- Content -->
  <div v-else>
    <h2 style="margin-bottom: 1.5rem;">Backups</h2>

    <div class="filters">
      <Select
        v-model="filterStatus"
        :options="statusOptions"
        optionLabel="label"
        optionValue="value"
        placeholder="All Statuses"
        showClear
        class="filter-select"
      />
      <Select
        v-model="filterSchedule"
        :options="scheduleOptions"
        optionLabel="label"
        optionValue="value"
        placeholder="All Schedules"
        showClear
        class="filter-select"
      />
    </div>

    <DataTable
      :value="filteredBackups"
      stripedRows
      paginator
      :rows="20"
      :rowsPerPageOptions="[10, 20, 50]"
      sortField="startTimestamp"
      :sortOrder="-1"
      size="small"
      :rowHover="true"
      @rowClick="(e) => goToDetail(e.data.name)"
      style="cursor: pointer;"
    >
      <template #empty>
        <div class="empty-state">
          <i class="pi pi-inbox"></i>
          <p>No backups found</p>
        </div>
      </template>
      <Column field="name" header="Name" sortable>
        <template #body="{ data }">
          <a @click.stop="goToDetail(data.name)" class="backup-link">{{ data.name }}</a>
        </template>
      </Column>
      <Column field="scheduleName" header="Schedule" sortable>
        <template #body="{ data }">{{ data.scheduleName || '-' }}</template>
      </Column>
      <Column field="status" header="Status" sortable>
        <template #body="{ data }">
          <Tag :value="data.status" :severity="statusSeverity(data.status)" />
        </template>
      </Column>
      <Column field="startTimestamp" header="Started" sortable>
        <template #body="{ data }">{{ formatTime(data.startTimestamp) }}</template>
      </Column>
      <Column field="duration" header="Duration" sortable>
        <template #body="{ data }">{{ data.duration || '-' }}</template>
      </Column>
      <Column field="itemsBackedUp" header="Items" sortable />
      <Column field="warnings" header="Warnings" sortable>
        <template #body="{ data }">
          <span v-if="data.warnings > 0" style="color: var(--p-yellow-500); font-weight: 700;">{{ data.warnings }}</span>
          <span v-else>0</span>
        </template>
      </Column>
      <Column field="errors" header="Errors" sortable>
        <template #body="{ data }">
          <span v-if="data.errors > 0" style="color: var(--p-red-500); font-weight: 700;">{{ data.errors }}</span>
          <span v-else>0</span>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<style scoped>
.filters {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}
.filter-select {
  min-width: 12rem;
}
.backup-link {
  color: var(--p-primary-color);
  text-decoration: none;
  cursor: pointer;
  font-weight: 500;
}
.backup-link:hover {
  text-decoration: underline;
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
.mb-4 {
  margin-bottom: 1.5rem;
}
</style>
