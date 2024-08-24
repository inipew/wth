<template>
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
</template>
<script>
export default {
  name: "ArchiveTable",
  props: ["archiveFiles"],
};
</script>
