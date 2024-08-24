<template>
  <div class="table-container">
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
          <td @click="navigateTo(previousPath)" class="directory">📁 .../</td>
          <td></td>
          <td></td>
          <td></td>
        </tr>
        <tr v-for="file in files" :key="file.path">
          <td
            @click="file.is_dir && navigateTo(file.path)"
            :class="{ directory: file.is_dir }"
          >
            <span v-if="file.is_dir">📁 {{ file.name }}</span>
            <span v-else>📄 {{ file.name }}</span>
          </td>
          <td>{{ file.formatted_size }}</td>
          <td>{{ file.last_modified }}</td>
          <td>
            <button class="btn btn-action" @click="$emit('renameFile', file)">
              Rename
            </button>
            <button class="btn btn-action" @click="$emit('deleteFile', file)">
              Delete
            </button>
            <button
              v-if="isArchiveFile(file.name)"
              class="btn btn-action"
              @click="$emit('openArchiveTable', file.path)"
            >
              View Archive
            </button>
          </td>
        </tr>
      </tbody>
    </table>
    <p v-else>No files found.</p>
  </div>
</template>

<script>
export default {
  name: "FileTable",
  props: ["files", "currentPath", "previousPath"],
  methods: {
    navigateTo(path) {
      this.$emit("navigateTo", path);
    },
    isArchiveFile(filename) {
      const archiveExtensions = [".zip", ".tar.gz", ".tar", ".gz"];
      return archiveExtensions.some((ext) => filename.endsWith(ext));
    },
  },
};
</script>
