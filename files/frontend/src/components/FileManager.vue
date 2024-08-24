<template>
  <div id="app">
    <h1>File Manager</h1>
    <div class="navigation" v-if="currentPath !== ''">
      <button class="btn" @click="navigateTo(previousPath)">Go Up</button>
      <button class="btn" @click="openUploadModal">Upload File</button>
    </div>
    <table v-if="files.length" class="file-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Size</th>
          <th>Last Modified</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td @click="navigateTo(previousPath)">📁 .../</td>
          <td></td>
          <td></td>
          <td></td>
        </tr>
        <tr v-for="file in files" :key="file.path">
          <td @click="file.is_dir ? navigateTo(file.path) : null">
            <span v-if="file.is_dir">📁 {{ file.name }}</span>
            <span v-else>📄 {{ file.name }}</span>
          </td>
          <td>{{ file.formatted_size }}</td>
          <td>{{ file.last_modified }}</td>
          <td>
            <button class="btn" @click.prevent="editFile(file)">Edit</button>
            <button class="btn" @click.prevent="renameFile(file)">
              Rename
            </button>
            <button class="btn" @click.prevent="deleteFile(file)">
              Delete
            </button>
            <button
              v-if="isArchiveFile(file.name)"
              class="btn"
              @click.prevent="openArchiveModal(file.path)"
            >
              View Archive
            </button>
          </td>
        </tr>
      </tbody>
    </table>
    <p v-else>No files found.</p>

    <!-- Modal for Upload Form -->
    <div v-if="showUploadModal" class="modal">
      <div class="modal-content">
        <span class="close" @click="closeUploadModal">&times;</span>
        <h2>Upload File</h2>
        <form @submit.prevent="uploadFile">
          <input type="file" @change="handleFileUpload" required />
          <button type="submit" class="btn">Upload</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "FileManager",
  data() {
    return {
      files: [],
      currentPath: "",
      previousPath: "",
      showUploadModal: false,
      selectedFile: null,
    };
  },
  mounted() {
    this.fetchFiles(".");
  },
  methods: {
    async fetchFiles(path) {
      try {
        const response = await axios.get(
          "http://157.230.247.64:4567/api/files",
          {
            params: { path },
          }
        );
        const data = response.data;
        this.files = data.files;
        this.currentPath = data.current_path;
        this.previousPath = data.previous_path;
      } catch (error) {
        console.error("Error fetching files:", error);
      }
    },
    navigateTo(path) {
      this.fetchFiles(path);
    },
    async renameFile(file) {
      const newName = prompt("Enter new name:", file.name);
      if (newName) {
        try {
          await axios.post("http://157.230.247.64:4567/api/files/rename", {
            oldPath: file.path,
            newName,
          });
          this.fetchFiles(this.currentPath);
        } catch (error) {
          console.error("Error renaming file:", error);
        }
      }
    },
    async deleteFile(file) {
      const confirmDelete = confirm(
        `Are you sure you want to delete ${file.name}?`
      );
      if (confirmDelete) {
        try {
          await axios.delete("http://157.230.247.64:4567/api/files/delete", {
            data: { path: file.path },
          });
          this.fetchFiles(this.currentPath);
        } catch (error) {
          console.error("Error deleting file:", error);
        }
      }
    },
    openUploadModal() {
      this.showUploadModal = true;
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
        await axios.post(
          "http://157.230.247.64:4567/api/files/upload",
          formData,
          {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          }
        );
        this.fetchFiles(this.currentPath);
        this.closeUploadModal();
      } catch (error) {
        console.error("Error uploading file:", error);
      }
    },
    async viewArchive(path) {
      try {
        const response = await axios.get(
          "http://157.230.247.64:4567/api/files/view_archive",
          {
            params: { path },
          }
        );
        this.archiveFiles = response.data.files;
      } catch (error) {
        console.error("Error fetching archive contents:", error);
      }
    },
    isArchiveFile(filename) {
      const archiveExtensions = [".zip", ".tar.gz", ".tar", ".gz"];
      return archiveExtensions.some((ext) => filename.endsWith(ext));
    },
    editFile(file) {
      this.$router.push({ name: "editFile", params: { filePath: file.path } });
    },
  },
};
</script>

<style scoped>
.file-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

th,
td {
  border: 1px solid #ddd;
  padding: 12px;
  text-align: left;
}

th {
  background-color: #f2f2f2;
}

tr:hover {
  background-color: #f1f1f1;
}
</style>
