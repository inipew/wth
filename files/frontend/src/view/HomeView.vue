<template>
  <FileNavigation
    :currentPath="currentPath"
    :previousPath="previousPath"
    @navigateTo="navigateTo"
    @toggleUploadModal="toggleUploadModal"
    @openCreateEntityModal="openCreateEntityModal"
    @getToastOptions="getToastOptions"
  />

  <FileTable
    :files="files"
    :previousPath="previousPath"
    @navigateTo="navigateTo"
    @fetchFiles="fetchFiles"
    @getToastOptions="getToastOptions"
  />

  <UploadModal
    v-if="showUploadModal"
    @closeUploadModal="closeUploadModal"
    @handleFileUpload="handleFileUpload"
    @uploadFile="uploadFile"
    @getToastOptions="getToastOptions"
  />

  <CreateNewModal
    v-if="showCreateEntityModal"
    :show="showCreateEntityModal"
    :currentPath="currentPath"
    @fetchFiles="fetchFiles"
    @close="closeCreateEntityModal"
    @getToastOptions="getToastOptions"
  />
</template>

<script>
import axios from "axios";
import FileNavigation from "../components/FileNavigation.vue";
import FileTable from "../components/FileTable.vue";
import UploadModal from "../components/UploadModal.vue";
import CreateNewModal from "@/components/CreateNewModal.vue";
import { useToast } from "vue-toastification";

export default {
  name: "App",
  components: {
    FileNavigation,
    FileTable,
    UploadModal,
    CreateNewModal,
  },
  data() {
    return {
      files: [],
      currentPath: "",
      previousPath: "",
      archiveFiles: [],
      showUploadModal: false,
      showCreateEntityModal: false,
      selectedFile: null,
      errorMessage: "",
    };
  },
  mounted() {
    this.currentPath = localStorage.getItem("currentPath") || "";
    this.fetchFiles(this.currentPath);
  },
  methods: {
    async fetchFiles(path) {
      const toast = useToast();
      try {
        const response = await axios.get("/api/files", { params: { path } });
        this.files = response.data.files;
        this.currentPath = response.data.current_path;
        this.previousPath = response.data.previous_path;
        localStorage.setItem("currentPath", this.currentPath);
      } catch (error) {
        let errorMessage = "Failed to fetch files.";

        // Determine specific error message
        if (error.response) {
          errorMessage = error.response.data.message || errorMessage;
        } else if (error.request) {
          errorMessage = "No response received from server.";
        } else {
          errorMessage = "An error occurred while processing your request.";
        }
        toast.error(errorMessage, this.getToastOptions());
      }
    },
    toggleUploadModal() {
      this.showUploadModal = !this.showUploadModal;
    },
    closeUploadModal() {
      this.showUploadModal = false;
      this.selectedFile = null;
    },
    openCreateEntityModal() {
      this.showCreateEntityModal = !this.showCreateEntityModal;
    },
    closeCreateEntityModal() {
      this.showCreateEntityModal = false;
    },
    handleFileUpload(event) {
      this.selectedFile = event.target.files[0];
    },
    navigateTo(path) {
      this.fetchFiles(path);

      localStorage.setItem("currentPath", path);
    },
    async uploadFile() {
      if (!this.selectedFile) {
        this.errorMessage = "Please select a file to upload.";
        return;
      }
      const toast = useToast();
      const formData = new FormData();
      formData.append("file", this.selectedFile);
      formData.append("path", this.currentPath);

      try {
        const response = await axios.post("/api/files/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        this.fetchFiles(this.currentPath);
        this.closeUploadModal();
        toast.success(response.data.message, this.getToastOptions());
      } catch (error) {
        let errorMessage = "Failed to upload a file.";

        // Determine specific error message
        if (error.response) {
          errorMessage = error.response.data.message || errorMessage;
        } else if (error.request) {
          errorMessage = "No response received from server.";
        } else {
          errorMessage = "An error occurred while processing your request.";
        }
        toast.error(errorMessage, this.getToastOptions());
      }
    },
    getToastOptions() {
      return {
        position: "top-right",
        timeout: 1989,
        closeOnClick: true,
        pauseOnFocusLoss: true,
        pauseOnHover: true,
        draggable: true,
        draggablePercent: 0.6,
        showCloseButtonOnHover: true,
        hideProgressBar: true,
        closeButton: "button",
        icon: true,
        rtl: false,
      };
    },
    handleError(message, error) {
      this.errorMessage = message;
      console.error(message, error);
      // Consider using a notification library here
      alert(message);
    },
  },
};
</script>
<style>
/* Variabel warna */
:root {
  /* Warna terang */
  --light-background: #f4f7f6;
  --light-text: #343a40;
  --light-breadcrumb-bg: #ffffff;
  --light-breadcrumb-text: #343a40;
  --light-breadcrumb-link: #007bff;
  --light-breadcrumb-link-hover: #0056b3;
  --light-button-bg: #007bff;
  --light-button-bg-hover: #0056b3;
  --light-button-shadow: rgba(0, 0, 0, 0.2);
  --light-button-border: #007bff;
  --light-modal-bg: #ffffff;
  --light-modal-header-bg: #007bff;
  --light-modal-header-text: #ffffff;
  --light-modal-body-bg: #ffffff;
  --light-modal-body-text: #343a40;

  /* Warna gelap */
  --dark-background: #222;
  --dark-text: #fff;
  --dark-breadcrumb-bg: #2c2c2c;
  --dark-breadcrumb-text: #fff;
  --dark-breadcrumb-link: #66b2ff;
  --dark-breadcrumb-link-hover: #3b8adb;
  --dark-button-bg: #555;
  --dark-button-bg-hover: #666;
  --dark-button-shadow: rgba(0, 0, 0, 0.3);
  --dark-button-border: #555;
  --dark-modal-bg: #333;
  --dark-modal-header-bg: #444;
  --dark-modal-header-text: #fff;
  --dark-modal-body-bg: #333;
  --dark-modal-body-text: #fff;

  /* Tabel */
  --dark-table-header-bg: #2c2c2c;
  --dark-table-header-text: #e1e1e1;
  --dark-table-cell-bg: #2a2a2a;
  --dark-table-cell-text: #e1e1e1;
  --dark-table-row-hover-bg: #444;

  /* Nav */
  --dark-nav-bg: #1e1e1e;
  --dark-nav-text: #e0e0e0;
  --dark-nav-link: #66b2ff;
  --dark-nav-link-hover: #3b8adb;
  --dark-nav-border: #333;
}

/* Pengaturan umum */
body {
  font-family: "Arial", sans-serif;
  color: var(--light-text);
  background-color: var(--light-background);
}

#app {
  text-align: center;
  margin-top: 20px;
}

/* Tombol */
.btn {
  margin: 0 5px;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s, box-shadow 0.3s;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-primary {
  background-color: var(--light-button-bg);
  color: #fff;
}

.btn-primary:hover {
  background-color: var(--light-button-bg-hover);
}

.btn-secondary {
  background-color: #6c757d;
  color: #fff;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

/* Responsif */
@media (max-width: 768px) {
  .btn {
    font-size: 14px;
    padding: 8px 16px;
  }
  /* table {
    font-size: 14px;
  }

  thead th,
  tbody th,
  tbody td {
    padding: 10px;
  } */
}

@media (max-width: 576px) {
  #app {
    margin-top: 10px;
  }

  /* table {
    font-size: 12px;
  }

  thead th,
  tbody th,
  tbody td {
    padding: 8px;
  } */

  .btn {
    font-size: 12px;
    padding: 6px 12px;
  }
}
</style>
