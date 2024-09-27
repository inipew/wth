<template>
  <div>
    <div class="table-container">
      <table class="file-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Size</th>
            <th>Last Modified</th>
            <th>Permission</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td @click="navigateTo(previousPath)" class="directory" colspan="5">
              📁 .../
            </td>
          </tr>
          <template v-if="files && files.length">
            <tr v-for="file in files" :key="file.path">
              <td @click="handleRowClick(file)" class="directory">
                <span v-if="file.is_dir">📁 {{ file.name }}</span>
                <span v-else-if="isArchiveFile(file.name)">
                  📦 {{ file.name }}</span
                >
                <span v-else>📄 {{ file.name }}</span>
              </td>
              <td>{{ file.file_size }}</td>
              <td>{{ file.last_modified }}</td>
              <td>
                <button
                  class="btn btn-permission"
                  @click="openPermissionModal(file)"
                >
                  {{ file.permissions }}
                </button>
              </td>
              <td>
                <button
                  v-if="file.is_editable"
                  class="btn btn-edit"
                  @click="goToEditPage(file)"
                >
                  Edit
                </button>
                <button
                  v-if="isArchiveFile(file.name)"
                  class="btn btn-view"
                  @click="openArchiveModal(file.path)"
                >
                  View
                </button>
                <button class="btn btn-rename" @click="openRenameModal(file)">
                  Rename
                </button>
                <button class="btn btn-delete" @click="deleteFile(file)">
                  Delete
                </button>
              </td>
            </tr>
          </template>
          <template v-else>
            <tr>
              <td colspan="5" class="empty">Directory is empty</td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>

    <ArchiveModal
      v-if="showArchiveModal"
      :archiveFiles="archiveFiles"
      @closeArchiveModal="closeArchiveModal"
      @fetchFiles="this.$emit('fetchFiles', this.currentPath)"
      @getToastOptions="this.$emit('getToastOptions')"
    />

    <PermissionModal
      v-if="showPermissionModal"
      :show="showPermissionModal"
      :currentPath="currentPath"
      :permissionsData="permissionsData"
      @close="showPermissionModal = false"
      @fetchFiles="this.$emit('fetchFiles', this.currentPath)"
      @getToastOptions="this.$emit('getToastOptions')"
    />

    <RenameModal
      v-if="showRenameModal"
      :showModal="showRenameModal"
      :renameData="renameData"
      :currentPath="currentPath"
      @close="showRenameModal = false"
      @fetchFiles="this.$emit('fetchFiles', this.currentPath)"
      @getToastOptions="this.$emit('getToastOptions')"
    />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import axios from 'axios';
import ArchiveModal from '../components/ArchiveModal.vue';
import PermissionModal from '@/components/PermissionModal.vue';
import RenameModal from './RenameModal.vue';

export default {
  name: 'FileTable',
  components: {
    ArchiveModal,
    PermissionModal,
    RenameModal,
  },
  props: {
    files: Array,
    previousPath: String,
  },
  setup(props, { emit }) {
    const router = useRouter();
    const toast = useToast();

    const archiveFiles = ref([]);
    const currentPath = ref('');
    const showArchiveModal = ref(false);
    const showPermissionModal = ref(false);
    const showRenameModal = ref(false);
    const permissionsData = ref({ filepath: '', permissions: '' });
    const renameData = ref({ oldPath: '', name: '' });

    onMounted(() => {
      currentPath.value = localStorage.getItem('currentPath') || '';
    });

    const isArchiveFile = (filename) => {
      return ['.zip', '.tar.gz', '.tar', '.gz'].some((ext) => filename.endsWith(ext));
    };

    const navigateTo = (path) => {
      emit('navigateTo', path);
    };

    const deleteFile = async (file) => {
      if (confirm(`Are you sure you want to delete ${file.name}?`)) {
        try {
          const response = await axios.delete('/api/files/delete', { data: { path: file.path } });
          toast.success(response.data.message, emit('getToastOptions'));
          emit('fetchFiles', currentPath.value);
        } catch (error) {
          handleError('Failed to delete file.', error);
        }
      }
    };

    const viewArchive = async (path) => {
      try {
        const response = await axios.get('/api/files/view_archive', { params: { path } });
        archiveFiles.value = response.data;
      } catch (error) {
        handleError('Failed to fetch archive contents.', error);
      }
    };

    const downloadFile = async (file) => {
      if (confirm(`Are you sure you want to download ${file.name}?`)) {
        try {
          const url = `/api/files/download?file=${encodeURIComponent(file.path)}`;
          const response = await axios.get(url, { responseType: 'blob' });
          const contentDisposition = response.headers['content-disposition'];
          const fileNameMatch = contentDisposition ? contentDisposition.match(/filename="(.+)"/) : null;
          const actualFileName = fileNameMatch ? fileNameMatch[1] : file.name;
          const blob = new Blob([response.data], { type: response.headers['content-type'] || 'application/octet-stream' });
          const link = document.createElement('a');
          link.href = URL.createObjectURL(blob);
          link.download = actualFileName;
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
          URL.revokeObjectURL(link.href);
        } catch (error) {
          handleError('Error downloading file.', error);
        }
      }
    };

    const openPermissionModal = (file) => {
      permissionsData.value = { filepath: file.path, permissions: file.permissions };
      showPermissionModal.value = true;
    };

    const openRenameModal = (file) => {
      renameData.value = { oldPath: file.path, name: file.name };
      showRenameModal.value = true;
    };

    const closeRenameModal = () => {
      renameData.value = { oldPath: '', name: '' };
      showRenameModal.value = false;
    };

    const openArchiveModal = (path) => {
      viewArchive(path);
      showArchiveModal.value = true;
    };

    const closeArchiveModal = () => {
      showArchiveModal.value = false;
      archiveFiles.value = [];
    };

    const goToEditPage = (file) => {
      router.push({
        name: 'EditItem',
        params: { filepath: encodeURIComponent(file.path) },
      });
    };

    const handleRowClick = (file) => {
      if (file.is_dir) {
        navigateTo(file.path);
      } else if (isArchiveFile(file.name)) {
        openArchiveModal(file.path);
      } else if (file.is_editable) {
        goToEditPage(file);
      } else if (!file.is_editable && !file.is_dir) {
        downloadFile(file);
      }
    };

    const handleError = (message, error) => {
      console.error(message, error);
      toast.error(message, emit('getToastOptions'));
    };

    const formatDate = (dateString) => {
      const date = new Date(dateString);
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
    };
    
    return {
      archiveFiles,
      currentPath,
      showArchiveModal,
      showPermissionModal,
      showRenameModal,
      permissionsData,
      renameData,
      isArchiveFile,
      navigateTo,
      deleteFile,
      openPermissionModal,
      openRenameModal,
      closeRenameModal,
      openArchiveModal,
      closeArchiveModal,
      goToEditPage,
      handleRowClick,
      formatDate,
    };
  },
};
</script>

<style scoped>
.table-container {
  padding: 10px;
  display: flex;
  flex-direction: column;
}

/* Table styling */
.file-table {
  width: 100%;
  border-collapse: collapse;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  background: #ffffff;
}

.file-table th,
.file-table td {
  padding: 10px;
  border: 1px solid #e0e0e0;
  text-align: left;
  color: #333;
  font-size: 0.875em;
}

.file-table thead {
  background-color: #f0f9ff; /* Light blue background */
  color: #007bff; /* Primary blue color for text */
  font-weight: 600;
}

.file-table tbody tr:nth-child(even) {
  background-color: #f8f9fa; /* Light gray for alternating rows */
}

.file-table tbody tr:hover {
  background-color: #e9ecef; /* Slightly darker gray on hover */
}

.file-table .directory {
  cursor: pointer;
  color: #007bff; /* Primary blue color */
  font-weight: 500;
}

/* Button styles */
.file-table .btn {
  padding: 8px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease, color 0.3s ease, box-shadow 0.3s ease;
  color: #ffffff;
  font-size: 0.875em;
  margin: 2px; /* Space between buttons */
}

/* Permission button */
.file-table .btn-permission {
  background-color: #f8f9fa; /* Light background for better visibility */
  color: #007bff; /* Primary blue text color */
  border: 1px solid #e0e0e0; /* Light gray border */
  min-width: 95%;
}

.file-table .btn-permission:hover {
  background-color: #e2e6ea; /* Slightly darker gray on hover */
  color: #0056b3; /* Darker blue text color on hover */
}

.file-table .btn-permission:active {
  background-color: #d6d9db; /* Even darker gray on active */
  color: #004085; /* Even darker blue text color on active */
}

/* Action buttons */
.file-table .btn-edit {
  background: linear-gradient(135deg, #1e90ff, #4682b4);
}

.file-table .btn-edit:hover {
  background: linear-gradient(135deg, #4682b4, #4169e1);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.file-table .btn-edit:active {
  background: linear-gradient(135deg, #4169e1, #3152a1);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.file-table .btn-rename {
  background: linear-gradient(135deg, #007bff, #0056b3);
}

.file-table .btn-rename:hover {
  background: linear-gradient(135deg, #0056b3, #003d79);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.file-table .btn-rename:active {
  background: linear-gradient(135deg, #003d79, #002b54);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.file-table .btn-delete {
  background: linear-gradient(135deg, #dc3545, #c82333);
}

.file-table .btn-delete:hover {
  background: linear-gradient(135deg, #c82333, #a71d2a);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.file-table .btn-delete:active {
  background: linear-gradient(135deg, #a71d2a, #721c24);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.file-table .btn-view {
  background: linear-gradient(135deg, #17a2b8, #138496);
}

.file-table .btn-view:hover {
  background: linear-gradient(135deg, #138496, #117a8b);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.file-table .btn-view:active {
  background: linear-gradient(135deg, #117a8b, #0c6e76);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

/* Responsive adjustments */
@media (max-width: 1200px) {
  .file-table {
    font-size: 1em;
  }
}

@media (max-width: 768px) {
  .file-table {
    font-size: 0.9em;
  }
  .file-table th {
    line-height: 2;
  }
  .file-table th,
  .file-table td {
    padding: 8px;
  }
  .file-table .btn {
    padding: 6px 10px;
    font-size: 0.85em;
  }
  .file-table .directory {
    line-height: 2;
  }
}

@media (max-width: 480px) {
  .file-table {
    font-size: 0.8em;
  }
  .file-table th {
    font-size: 1em;
    line-height: 2;
  }
  .file-table td {
    padding: 6px;
    font-size: 0.8em;
  }
  .file-table .btn {
    padding: 5px 8px;
    font-size: 0.8em;
    min-width: 90%;
  }
  .file-table .directory {
    line-height: 3;
  }
}
</style>
