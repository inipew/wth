<template>
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
  />

  <p v-else>No files found.</p>

  <UploadModal
    v-if="showUploadModal"
    @closeUploadModal="closeUploadModal"
    @handleFileUpload="handleFileUpload"
    @uploadFile="uploadFile"
  />
</template>

<script>
import axios from "axios";
import FileNavigation from "../components/FileNavigation.vue";
import FileTable from "../components/FileTable.vue";
import UploadModal from "../components/UploadModal.vue";

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
<style src="../assets/styles.css"></style>
