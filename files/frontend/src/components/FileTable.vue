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
import axios from "axios";
import ArchiveModal from "../components/ArchiveModal.vue";
import PermissionModal from "@/components/PermissionModal.vue";
import RenameModal from "./RenameModal.vue";
import { useToast } from "vue-toastification";

export default {
  name: "FileTable",
  props: {
    files: Array,
    previousPath: String,
  },
  components: {
    ArchiveModal,
    PermissionModal,
    RenameModal,
  },
  data() {
    return {
      archiveFiles: [],
      currentPath: "",
      showArchiveModal: false,
      showPermissionModal: false,
      showRenameModal: false,
      permissionsData: {
        filepath: "",
        permissions: "",
      },
      renameData: {
        oldPath: "",
        name: "",
      },
    };
  },
  mounted() {
    this.currentPath = localStorage.getItem("currentPath") || "";
  },
  methods: {
    isArchiveFile(filename) {
      return [".zip", ".tar.gz", ".tar", ".gz"].some((ext) =>
        filename.endsWith(ext)
      );
    },
    navigateTo(path) {
      this.$emit("navigateTo", path);
    },
    async deleteFile(file) {
      if (confirm(`Are you sure you want to delete ${file.name}?`)) {
        const toast = useToast();

        try {
          const response = await axios.delete("/api/files/delete", {
            data: { path: file.path },
          });
          toast.success(response.data.message, this.$emit("getToastOptions"));
          this.$emit("fetchFiles", this.currentPath);
        } catch (error) {
          let errorMessage = "Failed to delete file.";

          // Determine specific error message
          if (error.response) {
            errorMessage = error.response.data.message || errorMessage;
          } else if (error.request) {
            errorMessage = "No response received from server.";
          } else {
            errorMessage = "An error occurred while processing your request.";
          }
          toast.error(errorMessage, this.$emit("getToastOptions"));
        }
      }
    },
    async viewArchive(path) {
      const toast = useToast();
      try {
        const response = await axios.get("/api/files/view_archive", {
          params: { path },
        });
        this.archiveFiles = response.data;
      } catch (error) {
        let errorMessage = "Failed to fetch archive contents.";

        // Determine specific error message
        if (error.response) {
          errorMessage = error.response.data.message || errorMessage;
        } else if (error.request) {
          errorMessage = "No response received from server.";
        } else {
          errorMessage = "An error occurred while processing your request.";
        }
        toast.error(errorMessage, this.$emit("getToastOptions"));
      }
    },
    async downloadFile(file) {
      if (confirm(`Are you sure you want to download ${file.name}?`)) {
        const toast = useToast();
        try {
          const url = `/api/files/download?file=${encodeURIComponent(
            file.path
          )}`;

          const response = await axios.get(url, { responseType: "blob" });

          const contentDisposition = response.headers["content-disposition"];
          const fileNameMatch = contentDisposition
            ? contentDisposition.match(/filename="(.+)"/)
            : null;
          const actualFileName = fileNameMatch ? fileNameMatch[1] : file.name;

          const blob = new Blob([response.data], {
            type:
              response.headers["content-type"] || "application/octet-stream",
          });

          const link = document.createElement("a");
          link.href = URL.createObjectURL(blob);
          link.download = actualFileName;
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
          URL.revokeObjectURL(link.href);
        } catch (error) {
          let errorMessage = "Error downloading file.";

          // Determine specific error message
          if (error.response) {
            errorMessage = error.response.data.message || errorMessage;
          } else if (error.request) {
            errorMessage = "No response received from server.";
          } else {
            errorMessage = "An error occurred while processing your request.";
          }
          toast.error(errorMessage, this.$emit("getToastOptions"));
        }
      }
    },
    openPermissionModal(file) {
      this.permissionsData = {
        filepath: file.path,
        permissions: file.permissions,
      };
      this.showPermissionModal = true;
    },
    openRenameModal(file) {
      this.renameData = {
        oldPath: file.path,
        name: file.name,
      };
      this.showRenameModal = true;
    },
    closeRenameModal() {
      this.renameData = [];
      this.showRenameModal = false;
    },
    openArchiveModal(path) {
      this.viewArchive(path);
      this.showArchiveModal = true;
    },
    closeArchiveModal() {
      this.showArchiveModal = false;
      this.archiveFiles = [];
    },
    goToEditPage(file) {
      this.$router.push({
        name: "EditItem",
        params: { filepath: encodeURIComponent(file.path) },
      });
    },
    handleRowClick(file) {
      if (file.is_dir) {
        this.navigateTo(file.path);
      } else if (this.isArchiveFile(file.name)) {
        this.openArchiveModal(file.path);
      } else if (file.is_editable) {
        this.goToEditPage(file);
      } else if (!file.is_editable && !file.is_dir) {
        this.downloadFile(file);
      }
    },
    handleError(message, error) {
      console.error(message, error);
      alert(message);
    },
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
