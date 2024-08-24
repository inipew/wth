<template>
  <div id="app">
    <h1>File Manager</h1>
    <FileNavigation
      :currentPath="currentPath"
      :previousPath="previousPath"
      @navigateTo="navigateTo"
      @toggleUploadModal="toggleUploadModal"
    />

    <FileTable
      v-if="files.length"
      :files="files"
      :currentPath="currentPath"
      :previousPath="previousPath"
      @navigateTo="navigateTo"
      @renameFile="renameFile"
      @deleteFile="deleteFile"
      @openArchiveTable="openArchiveTable"
    />

    <p v-else>No files found.</p>

    <div v-if="showArchiveTable" class="archive-files-container">
      <div class="archive-files-header">
        <h2>Archive Files</h2>
        <button class="btn btn-close" @click="closeArchiveTable">
          Close Archive
        </button>
      </div>
      <table v-if="archiveFiles.length" class="archive-files">
        <thead>
          <tr>
            <th>Name</th>
            <th>Size</th>
            <th>Last Modified</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="file in archiveFiles" :key="file.path">
            <td>{{ file.name }}</td>
            <td>{{ file.formatted_size }}</td>
            <td>{{ file.last_modified }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else>No files found in archive.</p>
    </div>

    <UploadModal
      v-if="showUploadModal"
      @closeUploadModal="closeUploadModal"
      @handleFileUpload="handleFileUpload"
      @uploadFile="uploadFile"
    />
  </div>
</template>

<script>
import axios from "axios";
import FileNavigation from "./components/FileNavigation.vue";
import FileTable from "./components/FileTable.vue";
import UploadModal from "./components/UploadModal.vue";

export default {
  name: "App",
  components: {
    FileNavigation,
    FileTable,
    UploadModal,
  },
  data() {
    return {
      files: [],
      currentPath: "",
      previousPath: "",
      archiveFiles: [],
      showArchiveTable: false,
      showUploadModal: false,
      selectedFile: null,
      errorMessage: "", // For user-friendly error handling
    };
  },
  mounted() {
    this.currentPath = localStorage.getItem("currentPath") || "";
    this.fetchFiles(this.currentPath);
  },
  methods: {
    async fetchFiles(path) {
      try {
        const response = await axios.get("/api/files", { params: { path } });
        this.files = response.data.files;
        this.currentPath = response.data.current_path;
        this.previousPath = response.data.previous_path;
        localStorage.setItem("currentPath", this.currentPath);
      } catch (error) {
        this.handleError("Failed to fetch files.", error);
      }
    },
    async renameFile(file) {
      const newName = prompt("Enter new name:", file.name);
      if (newName && newName !== file.name) {
        try {
          await axios.post("/api/files/rename", {
            oldPath: file.path,
            newName,
          });
          this.fetchFiles(this.currentPath);
        } catch (error) {
          this.handleError("Failed to rename file.", error);
        }
      }
    },
    async deleteFile(file) {
      if (confirm(`Are you sure you want to delete ${file.name}?`)) {
        try {
          await axios.delete("/api/files/delete", {
            data: { path: file.path },
          });
          this.fetchFiles(this.currentPath);
        } catch (error) {
          this.handleError("Failed to delete file.", error);
        }
      }
    },
    async viewArchive(path) {
      try {
        const response = await axios.get("/api/files/view_archive", {
          params: { path },
        });
        this.archiveFiles = response.data.files;
      } catch (error) {
        this.handleError("Failed to fetch archive contents.", error);
      }
    },
    openArchiveTable(path) {
      this.viewArchive(path);
      this.showArchiveTable = true;
    },
    closeArchiveTable() {
      this.showArchiveTable = false;
      this.archiveFiles = [];
    },
    toggleUploadModal() {
      this.showUploadModal = !this.showUploadModal;
    },
    closeUploadModal() {
      this.showUploadModal = false;
      this.selectedFile = null;
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

      const formData = new FormData();
      formData.append("file", this.selectedFile);
      formData.append("path", this.currentPath);

      try {
        await axios.post("/api/files/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        this.fetchFiles(this.currentPath);
        this.closeUploadModal();
      } catch (error) {
        this.handleError("Failed to upload file.", error);
      }
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

<!-- <script>
import axios from "axios";
import FileNavigation from "./components/FileNavigation.vue";
import FileTable from "./components/FileTable.vue";
import UploadModal from "./components/UploadModal.vue";

export default {
  name: "App",
  components: {
    FileNavigation,
    FileTable,
    UploadModal,
  },
  data() {
    return {
      files: [],
      currentPath: "",
      previousPath: "",
      archiveFiles: [],
      showArchiveTable: false,
      showUploadModal: false,
      selectedFile: null,
    };
  },
  mounted() {
    const savedPath = localStorage.getItem("currentPath") || "";
    this.fetchFiles(savedPath);
  },
  methods: {
    async fetchFiles(path) {
      try {
        const { data } = await axios.get("/api/files", { params: { path } });
        this.files = data.files;
        this.currentPath = data.current_path;
        this.previousPath = data.previous_path;
        localStorage.setItem("currentPath", this.currentPath);
      } catch (error) {
        console.error("Error fetching files:", error);
        alert("Failed to fetch files.");
      }
    },
    navigateTo(path) {
      this.fetchFiles(path);
      localStorage.setItem("currentPath", path);
    },
    async renameFile(file) {
      const newName = prompt("Enter new name:", file.name);
      if (newName && newName !== file.name) {
        try {
          await axios.post("/api/files/rename", {
            oldPath: file.path,
            newName,
          });
          this.fetchFiles(this.currentPath);
        } catch (error) {
          console.error("Error renaming file:", error);
          alert("Failed to rename file.");
        }
      }
    },
    async deleteFile(file) {
      if (confirm(`Are you sure you want to delete ${file.name}?`)) {
        try {
          await axios.delete("/api/files/delete", {
            data: { path: file.path },
          });
          this.fetchFiles(this.currentPath);
        } catch (error) {
          console.error("Error deleting file:", error);
          alert("Failed to delete file.");
        }
      }
    },
    async viewArchive(path) {
      try {
        const { data } = await axios.get("/api/files/view_archive", {
          params: { path },
        });
        this.archiveFiles = data.files;
      } catch (error) {
        console.error("Error fetching archive contents:", error);
        alert("Failed to fetch archive contents.");
      }
    },
    openArchiveTable(path) {
      this.viewArchive(path);
      this.showArchiveTable = true;
    },
    closeArchiveTable() {
      this.showArchiveTable = false;
      this.archiveFiles = [];
    },
    toggleUploadModal() {
      this.showUploadModal = !this.showUploadModal;
    },
    closeUploadModal() {
      this.showUploadModal = false;
      this.selectedFile = null;
    },
    handleFileUpload(event) {
      this.selectedFile = event.target.files[0];
    },
    async uploadFile() {
      if (!this.selectedFile) {
        alert("Please select a file to upload.");
        return;
      }

      const formData = new FormData();
      formData.append("file", this.selectedFile);
      formData.append("path", this.currentPath);

      try {
        await axios.post("/api/files/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        this.fetchFiles(this.currentPath);
        this.closeUploadModal();
      } catch (error) {
        console.error("Error uploading file:", error);
        alert("Failed to upload file.");
      }
    },
  },
};
</script> -->

<style src="./assets/styles.css"></style>
