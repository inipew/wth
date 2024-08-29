<template>
  <div class="modal-overlay" v-if="archiveFiles.length">
    <div class="modal-content">
      <div class="archive-files-header">
        <h2>Archive Files</h2>
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
      <div class="modal-footer">
        <button class="btn btn-warning" @click="extractFiles">Extract</button>
        <button class="btn btn-close" @click="$emit('closeArchiveModal')">
          Close
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ArchiveModal",
  props: ["archiveFiles"],
  methods: {
    extractFiles() {
      // Logic for extracting files goes here
      this.$emit("extractFiles");
    },
  },
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 50%;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.archive-files-header {
  display: flex;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid #ddd;
}

.archive-files-header h2 {
  margin: 0;
}

.btn {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.btn-close {
  background-color: #dc3545;
}

.btn-extract {
  background-color: #28a745;
}

.archive-files {
  width: 100%;
  border-collapse: collapse;
}

.archive-files th,
.archive-files td {
  padding: 8px 12px;
  border-bottom: 1px solid #ddd;
}

.archive-files th {
  background-color: #f8f9fa;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
  border-top: 1px solid #ddd;
}
</style>
