<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import Skeleton from 'primevue/skeleton'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import { statusSeverity, formatTime, formatBytes } from '../composables/useBackupUtils.js'

const route = useRoute()
const router = useRouter()
const backup = ref(null)
const loading = ref(true)
const error = ref(null)
const logs = ref(null)
const logsLoading = ref(false)
const logsError = ref(null)

onMounted(async () => {
  try {
    const res = await fetch(`/api/v1/backups/${route.params.name}`)
    if (!res.ok) {
      if (res.status === 404) throw new Error('Backup not found')
      throw new Error(`HTTP ${res.status}`)
    }
    backup.value = await res.json()
  } catch (e) {
    error.value = `Failed to load backup: ${e.message}`
  } finally {
    loading.value = false
  }
})

async function fetchLogs() {
  logsLoading.value = true
  logsError.value = null
  try {
    const res = await fetch(`/api/v1/backups/${route.params.name}/logs`)
    if (!res.ok) {
      const body = await res.text()
      throw new Error(body || `HTTP ${res.status}`)
    }
    logs.value = await res.text()
  } catch (e) {
    logsError.value = `Failed to load logs: ${e.message}`
  } finally {
    logsLoading.value = false
  }
}

function downloadPDF() {
  window.open(`/api/v1/backups/${route.params.name}/pdf`, '_blank')
}

function hasItems(arr) {
  return arr && arr.length > 0
}

function hasEntries(obj) {
  return obj && Object.keys(obj).length > 0
}
</script>

<template>
  <!-- Loading -->
  <div v-if="loading">
    <Skeleton width="20rem" height="2rem" class="mb-3" />
    <Skeleton height="30rem" />
  </div>

  <!-- Error -->
  <Message v-else-if="error" severity="error" :closable="false">{{ error }}</Message>

  <!-- Content -->
  <div v-else>
    <!-- Header -->
    <div class="detail-header">
      <h2>{{ backup.name }}</h2>
      <div class="detail-actions">
        <Button label="Back" icon="pi pi-arrow-left" severity="secondary" outlined @click="router.back()" />
        <Button label="Backup Report PDF" icon="pi pi-file-pdf" severity="danger" outlined @click="downloadPDF" />
      </div>
    </div>

    <!-- Tabs -->
    <Tabs value="overview">
      <TabList>
        <Tab value="overview"><i class="pi pi-info-circle mr-1"></i>Overview</Tab>
        <Tab value="volumes"><i class="pi pi-server mr-1"></i>Volumes</Tab>
        <Tab value="logs"><i class="pi pi-code mr-1"></i>Logs</Tab>
      </TabList>

      <TabPanels>
        <!-- Overview Tab -->
        <TabPanel value="overview">
          <!-- Metadata -->
          <Card class="mb-3">
            <template #title><i class="pi pi-tag mr-2"></i>Metadata</template>
            <template #content>
              <div class="detail-grid">
                <div class="detail-col">
                  <dl class="detail-dl">
                    <dt>Name</dt><dd>{{ backup.name }}</dd>
                    <dt>Namespace</dt><dd>{{ backup.namespace }}</dd>
                    <dt>Status</dt><dd><Tag :value="backup.status" :severity="statusSeverity(backup.status)" /></dd>
                    <dt>Schedule</dt><dd>{{ backup.scheduleName || '-' }}</dd>
                  </dl>
                </div>
                <div class="detail-col">
                  <dl class="detail-dl">
                    <dt>Started</dt><dd>{{ formatTime(backup.startTimestamp) }}</dd>
                    <dt>Completed</dt><dd>{{ formatTime(backup.completionTimestamp) }}</dd>
                    <dt>Duration</dt><dd>{{ backup.duration || '-' }}</dd>
                    <dt>Expiration</dt><dd>{{ formatTime(backup.expiration) }}</dd>
                    <dt>TTL</dt><dd>{{ backup.ttl || '-' }}</dd>
                    <dt>Storage Location</dt><dd>{{ backup.storageLocation || '-' }}</dd>
                    <dt v-if="backup.formatVersion">Format Version</dt>
                    <dd v-if="backup.formatVersion">{{ backup.formatVersion }}</dd>
                  </dl>
                </div>
              </div>
            </template>
          </Card>

          <!-- Failure Reason -->
          <Message v-if="backup.failureReason" severity="error" :closable="false">
            <strong>Failure Reason:</strong> {{ backup.failureReason }}
          </Message>

          <!-- Validation Errors -->
          <Message v-if="hasItems(backup.validationErrors)" severity="warn" :closable="false">
            <strong>Validation Errors:</strong>
            <ul style="margin: 0.5rem 0 0 0; padding-left: 1.5rem;">
              <li v-for="(ve, i) in backup.validationErrors" :key="i">{{ ve }}</li>
            </ul>
          </Message>

          <!-- Configuration -->
          <Card class="mb-3">
            <template #title><i class="pi pi-cog mr-2"></i>Configuration</template>
            <template #content>
              <div class="detail-grid">
                <div class="detail-col">
                  <h4>Included Namespaces</h4>
                  <p>{{ hasItems(backup.includedNamespaces) ? backup.includedNamespaces.join(', ') : 'All (*)' }}</p>
                  <h4>Included Resources</h4>
                  <p>{{ hasItems(backup.includedResources) ? backup.includedResources.join(', ') : 'All (*)' }}</p>
                </div>
                <div class="detail-col">
                  <h4>Excluded Namespaces</h4>
                  <p>{{ hasItems(backup.excludedNamespaces) ? backup.excludedNamespaces.join(', ') : 'None' }}</p>
                  <h4>Excluded Resources</h4>
                  <p>{{ hasItems(backup.excludedResources) ? backup.excludedResources.join(', ') : 'None' }}</p>
                </div>
              </div>
            </template>
          </Card>

          <!-- Status -->
          <Card class="mb-3">
            <template #title><i class="pi pi-check-square mr-2"></i>Status</template>
            <template #content>
              <dl class="detail-dl">
                <dt>Items Backed Up</dt><dd>{{ backup.itemsBackedUp }} / {{ backup.totalItems }}</dd>
                <dt>Warnings</dt><dd>{{ backup.warnings }}</dd>
                <dt>Errors</dt><dd>{{ backup.errors }}</dd>
              </dl>
              <dl v-if="backup.hooksAttempted > 0" class="detail-dl">
                <dt>Hooks</dt>
                <dd>
                  {{ backup.hooksAttempted }} attempted<span v-if="backup.hooksFailed > 0" style="color: var(--p-red-500);">, {{ backup.hooksFailed }} failed</span>
                </dd>
              </dl>
              <dl v-if="backup.backupItemOperationsAttempted > 0" class="detail-dl">
                <dt>Async Operations</dt>
                <dd>
                  {{ backup.backupItemOperationsAttempted }} attempted, {{ backup.backupItemOperationsCompleted }} completed<span v-if="backup.backupItemOperationsFailed > 0" style="color: var(--p-red-500);">, {{ backup.backupItemOperationsFailed }} failed</span>
                </dd>
              </dl>
            </template>
          </Card>

          <!-- Labels -->
          <Card v-if="hasEntries(backup.labels)" class="mb-3">
            <template #title><i class="pi pi-tags mr-2"></i>Labels</template>
            <template #content>
              <Tag v-for="(v, k) in backup.labels" :key="k" :value="`${k}=${v}`" severity="secondary" class="mr-1 mb-1 mono-tag" />
            </template>
          </Card>

          <!-- Annotations -->
          <Card v-if="hasEntries(backup.annotations)" class="mb-3">
            <template #title><i class="pi pi-bookmark mr-2"></i>Annotations</template>
            <template #content>
              <Tag v-for="(v, k) in backup.annotations" :key="k" :value="`${k}=${v}`" severity="secondary" class="mr-1 mb-1 mono-tag" style="word-break: break-all;" />
            </template>
          </Card>
        </TabPanel>

        <!-- Volumes Tab -->
        <TabPanel value="volumes">
          <!-- Volume Snapshots -->
          <Card class="mb-3">
            <template #title><i class="pi pi-camera mr-2"></i>Volume Snapshots</template>
            <template #content>
              <dl class="detail-dl">
                <dt>Native Snapshots</dt>
                <dd>{{ backup.volumeSnapshotsCompleted }} / {{ backup.volumeSnapshotsAttempted }}</dd>
                <dt>CSI Snapshots</dt>
                <dd>{{ backup.csiVolumeSnapshotsCompleted }} / {{ backup.csiVolumeSnapshotsAttempted }}</dd>
              </dl>
            </template>
          </Card>

          <!-- File System Volume Backups -->
          <Card>
            <template #title><i class="pi pi-database mr-2"></i>File System Volume Backups</template>
            <template #content>
              <DataTable
                v-if="hasItems(backup.volumeBackups)"
                :value="backup.volumeBackups"
                stripedRows
                size="small"
              >
                <Column field="volumeName" header="Volume" />
                <Column field="podName" header="Pod">
                  <template #body="{ data }">
                    <span style="color: var(--p-text-muted-color); font-size: 0.85rem;">{{ data.podNamespace }}/</span>{{ data.podName }}
                  </template>
                </Column>
                <Column field="nodeName" header="Node">
                  <template #body="{ data }">{{ data.nodeName || '-' }}</template>
                </Column>
                <Column field="phase" header="Status">
                  <template #body="{ data }">
                    <Tag :value="data.phase" :severity="statusSeverity(data.phase)" />
                  </template>
                </Column>
                <Column field="uploaderType" header="Uploader">
                  <template #body="{ data }">{{ data.uploaderType || '-' }}</template>
                </Column>
                <Column header="Size">
                  <template #body="{ data }">{{ formatBytes(data.totalBytes) }}</template>
                </Column>
                <Column header="Progress">
                  <template #body="{ data }">{{ formatBytes(data.bytesDone) }} / {{ formatBytes(data.totalBytes) }}</template>
                </Column>
              </DataTable>
              <div v-else class="empty-state">
                <i class="pi pi-server"></i>
                <p>No file system volume backups</p>
              </div>
            </template>
          </Card>
        </TabPanel>

        <!-- Logs Tab -->
        <TabPanel value="logs">
          <div v-if="logs === null && !logsLoading && !logsError" class="empty-state">
            <i class="pi pi-code"></i>
            <Button
              v-if="backup.isTerminal"
              label="Load Logs"
              icon="pi pi-download"
              outlined
              @click="fetchLogs"
              :loading="logsLoading"
            />
            <p v-else>Logs are available after backup completes</p>
          </div>

          <div v-if="logsLoading" class="empty-state">
            <Skeleton height="15rem" />
          </div>

          <Message v-if="logsError" severity="error" :closable="false">{{ logsError }}</Message>

          <Card v-if="logs !== null">
            <template #title><i class="pi pi-code mr-2"></i>Logs</template>
            <template #content>
              <pre class="log-output">{{ logs }}</pre>
            </template>
          </Card>
        </TabPanel>
      </TabPanels>
    </Tabs>
  </div>
</template>

<style scoped>
.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}
.detail-actions {
  display: flex;
  gap: 0.5rem;
}
.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}
@media (max-width: 768px) {
  .detail-grid { grid-template-columns: 1fr; }
}
.detail-col h4 {
  font-size: 0.9rem;
  font-weight: 600;
  margin: 0.75rem 0 0.25rem 0;
}
.detail-col h4:first-child {
  margin-top: 0;
}
.detail-dl {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 0.35rem 1rem;
  margin: 0;
}
.detail-dl dt {
  font-weight: 600;
  color: var(--p-text-muted-color);
  font-size: 0.9rem;
}
.detail-dl dd {
  margin: 0;
}
.mb-3 {
  margin-bottom: 1rem;
}
.mr-1 {
  margin-right: 0.25rem;
}
.mb-1 {
  margin-bottom: 0.25rem;
}
.mr-2 {
  margin-right: 0.5rem;
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
.mono-tag :deep(.p-tag-label) {
  font-family: var(--font-mono);
  font-size: 0.8rem;
}
.log-output {
  max-height: 500px;
  overflow: auto;
  margin: 0;
  padding: 0.5rem;
  font-family: var(--font-mono);
  font-size: 0.8rem;
  line-height: 1.6;
  background: var(--p-surface-ground);
  border-radius: var(--p-border-radius);
}
</style>
