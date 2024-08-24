<template>
  <div class="modal">
    <div class="modal-content">
      <span class="close" @click="$emit('close')">&times;</span>
      <h2>Edit File: {{ fileName }}</h2>
      <textarea v-model="content" rows="10" cols="50"></textarea>
      <br />
      <button @click="saveChanges">Save Changes</button>
      <button @click="$emit('close')">Cancel</button>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  props: {
    fileName: String,
    filePath: String,
  },
  data() {
    return {
      content: "",
    };
  },
  mounted() {
    this.fetchFileContent();
  },
  methods: {
    async fetchFileContent() {
      try {
        const response = await axios.get(
          `http://157.230.247.64:4567/api/files/view_file`,
          {
            params: { file: this.filePath },
          }
        );
        this.content = response.data.content;
      } catch (error) {
        console.error("Error fetching file content:", error);
      }
    },
    async saveChanges() {
      try {
        await axios.post(`http://157.230.247.64:4567/api/files/save_edit`, {
          file_path: this.filePath,
          content: this.content,
        });
        this.$emit("close");
      } catch (error) {
        console.error("Error saving file:", error);
      }
    },
  },
};
</script>

<style scoped>
.modal {
  display: block;
  position: fixed;
  z-index: 1000;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  overflow: auto;
  background-color: rgba(0, 0, 0, 0.7);
}

.modal-content {
  background-color: white;
  margin: 15% auto;
  padding: 20px;
  border: 1px solid #888;
  width: 80%;
  max-width: 600px;
  border-radius: 5px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.close {
  color: #aaa;
  float: right;
  font-size: 28px;
  font-weight: bold;
}

.close:hover,
.close:focus {
  color: black;
  text-decoration: none;
  cursor: pointer;
}
</style>
