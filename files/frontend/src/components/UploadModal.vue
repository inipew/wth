<template>
  <div class="modal">
    <div class="modal-content">
      <span class="close" @click="$emit('closeUploadModal')">&times;</span>
      <h2>Upload File</h2>
      <form @submit.prevent="$emit('uploadFile')">
        <input
          type="file"
          @change="$emit('handleFileUpload', $event)"
          required
        />
        <button type="submit" class="btn">Upload</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: "UploadModal",
};
</script>
