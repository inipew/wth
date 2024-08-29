<template>
  <ArchiveModal
    v-if="showArchiveModal"
    :archiveFiles="archiveFiles"
    @closeArchiveModal="closeArchiveModal"
  />

  <div class="table-container">
    <table v-if="files.length" class="file-table">
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
          <td @click="$emit('navigateTo', previousPath)" class="directory">
            📁 .../
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
        </tr>
        <tr v-for="file in files" :key="file.path">
          <td @click="handleRowClick(file)" class="directory">
            <span v-if="file.is_dir">📁 {{ file.name }}</span>
            <span v-else-if="isArchiveFile(file.name)">📦 {{ file.name }}</span>
            <span v-else>📄 {{ file.name }}</span>
          </td>
          <td>{{ file.formatted_size }}</td>
          <td>{{ file.last_modified }}</td>
          <td>{{ file.permissions }}</td>
          <td>
            <button class="btn btn-action" @click="renameFile(file)">
              Rename
            </button>
            <button class="btn btn-action" @click="deleteFile(file)">
              Delete
            </button>
            <button
              v-if="isArchiveFile(file.name)"
              class="btn btn-action"
              @click="openArchiveModal(file.path)"
            >
              View
            </button>
            <!-- <button
              v-if="!file.is_dir && file.is_editable"
              class="btn btn-action"
              @click="goToEditPage(file)"
            >
              Edit
            </button> -->
          </td>
        </tr>
      </tbody>
    </table>
    <p v-else>No files found.</p>
  </div>
</template>

<script>
import axios from "axios";
import ArchiveModal from "../components/ArchiveModal.vue";

export default {
  name: "FileTable",
  props: ["files", "currentPath", "previousPath"],
  components: {
    ArchiveModal,
  },
  data() {
    return {
      archiveFiles: [],
      showArchiveModal: false,
    };
  },
  methods: {
    isArchiveFile(filename) {
      const archiveExtensions = [".zip", ".tar.gz", ".tar", ".gz"];
      return archiveExtensions.some((ext) => filename.endsWith(ext));
    },
    goToEditPage(file) {
      this.$router.push({
        name: "EditItem",
        params: { filepath: encodeURIComponent(file) },
      });
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
    openArchiveModal(path) {
      this.viewArchive(path);
      this.showArchiveModal = true;
    },
    closeArchiveModal() {
      this.showArchiveModal = false;
      this.archiveFiles = [];
    },
    handleRowClick(file) {
      if (file.is_dir) {
        this.$emit("navigateTo", file.path);
      } else if (this.isArchiveFile(file.name)) {
        this.openArchiveModal(file.path);
      } else if (file.is_editable) {
        this.goToEditPage(file.path);
      }
    },
  },
};
</script>
