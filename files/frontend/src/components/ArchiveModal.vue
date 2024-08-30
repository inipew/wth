<template>
  <div
    :class="{
      'modal-overlay': true,
      show: archiveFiles && archiveFiles.files && archiveFiles.files.length,
    }"
  >
    <div
      :class="{
        'modal-content': true,
        show: archiveFiles && archiveFiles.files && archiveFiles.files.length,
      }"
    >
      <header class="archive-files-header">
        <h2>Archive Files</h2>
      </header>
      <div class="table-container">
        <table
          v-if="archiveFiles && archiveFiles.files && archiveFiles.files.length"
          class="archive-files"
        >
          <thead>
            <tr>
              <th>Name</th>
              <th>Size</th>
              <th>Last Modified</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="file in archiveFiles.files" :key="file.path">
              <td>{{ file.name }}</td>
              <td>{{ file.file_size }}</td>
              <td>{{ file.last_modified }}</td>
            </tr>
          </tbody>
        </table>
        <p v-else>No files found in archive.</p>
      </div>
      <footer class="modal-footer">
        <button class="btn btn-extract" @click="extractFiles">Extract</button>
        <button class="btn btn-close" @click="$emit('closeArchiveModal')">
          Close
        </button>
      </footer>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { useToast } from "vue-toastification";
export default {
  name: "ArchiveModal",
  props: {
    archiveFiles: {
      type: Object,
      required: true,
      default: () => ({
        name: "",
        path: "",
        files: [],
      }),
    },
  },
  methods: {
    async extractFile(filePath) {
      const toast = useToast();
      try {
        const response = await axios.get("/api/files/extract", {
          params: { file: filePath },
        });
        toast.success(response.data.message, this.$emit("getToastOptions"));
        this.$emit("closeArchiveModal");
        this.$emit("fetchFiles");
      } catch (error) {
        let errorMessage = "An error occurred while extracting the file.";

        // Determine specific error message
        if (error.response) {
          errorMessage = error.response.data.message || errorMessage;
        } else if (error.request) {
          errorMessage = "No response received from server.";
        } else {
          errorMessage = "An error occurred while processing your request.";
        }

        toast.error(errorMessage, this.$emit("getToastOptions"));
      }
    },
    async extractFiles() {
      if (
        confirm(
          `Are you sure you want to extract ${this.archiveFiles.name}?`
        ) &&
        this.archiveFiles.files.length
      ) {
        try {
          this.extractFile(this.archiveFiles.path);
        } catch (error) {
          console.error("Error extracting files:", error);
          // mungkin tambahkan feedback visual untuk error
        }
      }
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
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  transition: opacity 0.4s ease, visibility 0.4s ease;
  opacity: 0;
  visibility: hidden; /* Ensure modal is hidden by default */
}

/* Modal content styling */
.modal-content {
  width: 90%;
  max-width: 600px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  padding: 20px;
  box-sizing: border-box;
  transition: transform 0.3s ease, opacity 0.3s ease;
  transform: translateY(-10px);
  opacity: 0;
}

/* Transition for showing modal */
.modal-overlay.show,
.modal-content.show {
  opacity: 1;
  visibility: visible;
}

.modal-content.show {
  transform: translateY(0);
  opacity: 1;
}

/* Header styling */
.archive-files-header {
  border-bottom: 1px solid #dcdcdc; /* Subtle border to match the footer */
  background-color: #f1f1f1; /* Light gray background */
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* Light shadow for smoothness */
}

.archive-files-header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #333; /* Dark text color for better readability */
}

/* Table container styling */
.table-container {
  overflow-x: auto;
  margin-top: 16px;
  flex: 1;
}

/* Table styling */
.archive-files {
  width: 100%;
  border-collapse: collapse;
  margin: 0;
  min-width: 300px; /* Allow table to shrink on small screens */
}

.archive-files th,
.archive-files td {
  padding: 12px;
  border-bottom: 1px solid #dcdcdc; /* Subtle border to match header and footer */
  text-align: left;
  color: #333; /* Dark text color */
  font-size: 0.875rem;
}

.archive-files th {
  background-color: #e9e9e9; /* Slightly darker gray background for headers */
  font-size: 1rem;
  font-weight: 500;
}

.archive-files td {
  word-break: break-word;
}

/* Footer styling */
.modal-footer {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
  border-top: 1px solid #dcdcdc; /* Subtle border to match header */
  background-color: #f1f1f1; /* Light gray background */
}

/* Button styling */
.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
  color: #ffffff;
  font-size: 0.875em; /* Default font size */
}

/* Extract button */
.btn-extract {
  background: linear-gradient(135deg, #28a745, #218838); /* Green gradient */
}

.btn-extract:hover {
  background: linear-gradient(
    135deg,
    #218838,
    #1e7e34
  ); /* Darker green on hover */
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.btn-extract:active {
  background: linear-gradient(
    135deg,
    #1e7e34,
    #155d27
  ); /* Even darker green on active */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

/* Close button */
.btn-close {
  background: linear-gradient(135deg, #dc3545, #c82333); /* Red gradient */
}

.btn-close:hover {
  background: linear-gradient(
    135deg,
    #c82333,
    #a71d2a
  ); /* Darker red on hover */
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.btn-close:active {
  background: linear-gradient(
    135deg,
    #a71d2a,
    #721c24
  ); /* Even darker red on active */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

/* Responsive adjustments */
@media (max-width: 600px) {
  .modal-content {
    width: 95%;
    max-width: 100%;
  }

  .archive-files-header h2 {
    font-size: 1.25rem;
  }

  .archive-files {
    min-width: 300px; /* Allow table to shrink on small screens */
  }

  .archive-files th,
  .archive-files td {
    font-size: 0.75rem;
    padding: 8px;
  }

  .btn {
    padding: 6px 12px;
    font-size: 0.85em;
  }
}
</style>
