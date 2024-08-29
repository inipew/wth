<template>
  <div class="edit-file-container">
    <h1>Edit File: {{ fileName }}</h1>
    <textarea v-model="fileContent" rows="15" cols="80"></textarea>
    <div class="actions">
      <button @click="saveFile">Save</button>
      <button @click="cancelEdit">Cancel</button>
    </div>
    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import axios from "axios";

export default {
  name: "EditFile",
  setup() {
    const fileName = ref("");
    const fileContent = ref("");
    const errorMessage = ref("");

    onMounted(async () => {
      const query = new URLSearchParams(window.location.search);
      const filePath = query.get("file");

      if (filePath) {
        try {
          const response = await axios.get(`/api/files/view`, {
            params: { file: decodeURIComponent(filePath) },
          });
          fileName.value = response.data.fileName;
          fileContent.value = response.data.content;
        } catch (error) {
          errorMessage.value = "Error fetching file content.";
          console.error("Error fetching file content:", error);
        }
      }
    });

    const saveFile = async () => {
      const query = new URLSearchParams(window.location.search);
      const filePath = query.get("file");

      try {
        await axios.post(`/api/files/save`, {
          path: decodeURIComponent(filePath),
          content: fileContent.value,
        });
        alert("File saved successfully!");
        window.location.href = "/"; // Redirect to home or file list
      } catch (error) {
        errorMessage.value = "Error saving file.";
        console.error("Error saving file:", error);
      }
    };

    const cancelEdit = () => {
      window.location.href = "/";
    };

    return {
      fileName,
      fileContent,
      saveFile,
      cancelEdit,
      errorMessage,
    };
  },
};
</script>

<style scoped>
.edit-file-container {
  margin: 20px;
}

textarea {
  width: 100%;
  font-family: monospace;
}

.actions {
  margin-top: 10px;
}

.button {
  margin-right: 10px;
}

.error {
  color: red;
}
</style>
