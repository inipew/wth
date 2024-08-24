<template>
  <div class="modal">
    <div class="modal-content">
      <span class="close" @click="$emit('closeModal')">&times;</span>
      <h2>Archive Contents</h2>
      <table v-if="archiveFiles.length">
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
  </div>
</template>

<script>
export default {
  name: "ArchiveModal",
  props: ["archiveFiles"],
};
</script>
