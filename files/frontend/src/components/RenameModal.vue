<template>
  <div v-if="showModal" class="modal-overlay">
    <div class="modal-container">
      <h2 class="modal-title">Rename File: {{ renameData.name }}</h2>
      <input v-model="newName" type="text" class="modal-input" />
      <div class="modal-buttons">
        <button @click="submit" class="btn btn-primary">Rename</button>
        <button @click="close" class="btn btn-secondary">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { useToast } from "vue-toastification";

export default {
  props: {
    showModal: {
      type: Boolean,
      required: true,
    },
    currentPath: {
      type: String,
      required: true,
    },
    renameData: {
      type: Object,
      required: true,
      default: () => ({
        name: "",
        oldPath: "",
      }),
    },
  },
  data() {
    return {
      newName: this.renameData.name,
    };
  },
  methods: {
    async submit() {
      if (this.newName && this.newName !== this.renameData.name) {
        const toast = useToast();

        try {
          const response = await axios.post("/api/files/rename", {
            oldPath: this.renameData.oldPath,
            newName: this.newName,
          });
          toast.success(response.data.message, this.$emit("getToastOptions"));
          this.$emit("fetchFiles");
          this.close();
        } catch (error) {
          let errorMessage = "Failed to rename file.";

          // Determine specific error message
          if (error.response) {
            errorMessage = error.response.data.message || errorMessage;
          } else if (error.request) {
            errorMessage = "No response received from server.";
          } else {
            errorMessage = "An error occurred while processing your request.";
          }

          // Show error toast
          toast.error(errorMessage, {
            position: "top-right",
            timeout: 1989,
            closeOnClick: true,
            pauseOnFocusLoss: true,
            pauseOnHover: true,
            draggable: true,
            draggablePercent: 0.6,
            showCloseButtonOnHover: true,
            hideProgressBar: true,
            closeButton: "button",
            icon: true,
            rtl: false,
          });
        }
      }
    },
    close() {
      this.$emit("close");
    },
  },
  watch: {
    renameData: {
      handler(newVal) {
        this.newName = newVal.name;
      },
      deep: true,
    },
  },
};
</script>

<style scoped>
/* Modal overlay styling */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  transition: opacity 0.3s ease;
}

/* Modal container styling */
.modal-container {
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  padding: 24px;
  width: 90%;
  max-width: 400px;
  animation: fadeIn 0.3s ease-out;
  box-sizing: border-box;
}

/* Title styling */
.modal-title {
  font-size: 1.5rem;
  margin-bottom: 16px;
  color: #007bff; /* Primary blue color for title */
  text-align: center;
}

/* Input styling */
.modal-input {
  width: 100%;
  padding: 12px;
  margin-bottom: 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
  box-sizing: border-box;
}

.modal-input:focus {
  border-color: #007bff;
  outline: none;
}

/* Button container styling */
.modal-buttons {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

/* Primary button styling */
.btn-primary {
  background-color: #007bff;
  color: #fff;
  border: none;
  padding: 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s ease, transform 0.2s ease;
  flex: 1;
}

.btn-primary:hover {
  background-color: #0056b3;
  transform: scale(1.05);
}

/* Secondary button styling */
.btn-secondary {
  background-color: #6c757d;
  color: #fff;
  border: none;
  padding: 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s ease, transform 0.2s ease;
  flex: 1;
}

.btn-secondary:hover {
  background-color: #5a6268;
  transform: scale(1.05);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Responsive Styles */
@media (max-width: 600px) {
  .modal-container {
    width: 95%;
    padding: 16px;
  }

  .modal-title {
    font-size: 1.25rem;
  }

  .modal-input {
    font-size: 0.9rem;
  }

  .btn-primary,
  .btn-secondary {
    font-size: 0.9rem;
    padding: 10px;
  }
}
</style>
