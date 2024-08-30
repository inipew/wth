<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <button class="close" @click="$emit('close')">&times;</button>
      <h2>Create New File/Folder</h2>
      <form @submit.prevent="createEntity">
        <div class="form-group">
          <label>Select Type:</label>
          <div class="radio-group">
            <div class="radio-item">
              <input
                type="radio"
                id="createFile"
                name="type"
                value="file"
                v-model="type"
              />
              <label for="createFile">File</label>
            </div>
            <div class="radio-item">
              <input
                type="radio"
                id="createDir"
                name="type"
                value="dir"
                v-model="type"
              />
              <label for="createDir">Directory</label>
            </div>
          </div>
        </div>
        <div class="form-group">
          <label for="name">Name:</label>
          <input
            type="text"
            class="form-control"
            id="name"
            v-model="name"
            required
            placeholder="Enter name"
          />
          <small class="form-text text-muted">
            Please enter a name for the new file or directory.
          </small>
        </div>
        <button type="submit" class="btn-primary">Create</button>
      </form>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { useToast } from "vue-toastification";

export default {
  props: {
    currentPath: {
      type: String,
      required: true,
    },
    show: {
      type: Boolean,
      required: true,
    },
  },
  data() {
    return {
      type: "file", // Default type
      name: "",
    };
  },
  methods: {
    async createEntity() {
      const toast = useToast();

      try {
        const response = await axios.get("/api/files/make", {
          params: {
            type: this.type,
            currentPath: this.currentPath,
            name: this.name,
          },
        });
        toast.success(response.data.message, this.$emit("getToastOptions"));
        this.clearForm();
      } catch (err) {
        let errorMessage = "Failed to make a file/dir.";

        // Determine specific error message
        if (err.response) {
          errorMessage = err.response.data.message || errorMessage;
        } else if (err.request) {
          errorMessage = "No response received from server.";
        } else {
          errorMessage = "An error occurred while processing your request.";
        }
        toast.error(errorMessage, this.$emit("getToastOptions"));
      }
    },
    clearForm() {
      this.type = "file"; // Reset type to default
      this.name = "";
      this.$emit("close");
      this.$emit("fetchFiles", this.currentPath);
    },
  },
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.3s ease;
}

.modal-content {
  background: #fff;
  padding: 16px;
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  align-items: stretch;
  position: relative;
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.close {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 1.5em;
  color: #666;
  border: none;
  background: transparent;
  cursor: pointer;
  transition: color 0.3s ease;
}

.close:hover {
  color: #333;
}

h2 {
  margin-top: 0;
  font-size: 1.5em;
  color: #222;
  font-weight: 600;
}

.form-group {
  margin-bottom: 16px;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.radio-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

input[type="radio"] {
  accent-color: #007bff; /* Custom color for radio button */
  cursor: pointer;
}

input[type="radio"]:checked + label {
  font-weight: 600;
  color: #007bff; /* Highlight color for selected radio */
}

input[type="radio"]:not(:checked) + label {
  color: #333; /* Default color for unselected radio */
}

label {
  font-size: 1em;
  cursor: pointer;
}

input[type="text"] {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-sizing: border-box;
  font-size: 1em;
}

.btn-primary {
  padding: 10px 18px;
  border: none;
  border-radius: 8px;
  color: #fff;
  background-color: #007bff;
  font-size: 1em;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.3s ease;
}

.btn-primary:hover {
  background-color: #0056b3;
  transform: translateY(-1px);
}

.btn-primary:active {
  background-color: #004085;
  transform: translateY(0);
}

.form-text {
  font-size: 0.85em;
  color: #666;
}

/* Responsive Styles */
@media (max-width: 600px) {
  .modal-content {
    padding: 12px;
    width: 95%;
    max-width: none;
  }

  .close {
    font-size: 1.3em;
  }

  h2 {
    font-size: 1.3em;
  }

  .form-group {
    margin-bottom: 14px;
  }

  .btn-primary {
    padding: 8px 16px;
    font-size: 0.9em;
  }
}

@media (max-width: 900px) {
  .modal-content {
    padding: 14px;
    width: 90%;
    max-width: 360px;
  }

  .close {
    font-size: 1.4em;
  }

  h2 {
    font-size: 1.4em;
  }

  .btn-primary {
    padding: 10px 16px;
  }
}
</style>
