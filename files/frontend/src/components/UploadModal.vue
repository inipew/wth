<template>
  <div class="modal-overlay" @click="$emit('closeUploadModal')">
    <div class="modal-content" @click.stop>
      <button class="close" @click="$emit('closeUploadModal')">&times;</button>
      <h2>Upload File</h2>
      <form @submit.prevent="$emit('uploadFile')">
        <input
          type="file"
          @change="$emit('handleFileUpload', $event)"
          required
        />
        <button type="submit" class="btn btn-primary">Upload</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: "UploadModal",
};
</script>

<style scoped>
/* Modal Overlay */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6); /* Slightly darker overlay for better focus */
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.3s ease;
}

/* Modal Content */
.modal-content {
  background: #ffffff;
  padding: 24px;
  border-radius: 12px;
  width: 90%;
  max-width: 400px; /* Max width for better consistency on larger screens */
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  transition: transform 0.3s ease;
}

.close {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 1.8em; /* Larger close button for better visibility */
  color: #888;
  border: none;
  background: transparent;
  cursor: pointer;
  transition: color 0.3s ease;
}

.close:hover {
  color: #000;
}

h2 {
  margin-top: 0;
  font-size: 1.6em; /* Slightly larger font for heading */
  color: #004080; /* Deep blue color for header */
  font-weight: 600;
}

/* Form Styling */
form {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px; /* Increased gap for better spacing */
}

/* File Input */
input[type="file"] {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px; /* Rounded corners for input */
  box-sizing: border-box;
  font-size: 1em;
}

/* Button Styling */
.btn {
  padding: 12px 20px;
  border: none;
  border-radius: 6px;
  color: #ffffff;
  font-size: 1em;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
  font-weight: 500;
}

/* Primary Button */
.btn-primary {
  background-color: #0056b3; /* Dark blue for primary action */
}

.btn-primary:hover {
  background-color: #004494; /* Darker blue on hover */
}

.btn-primary:active {
  background-color: #003366; /* Even darker blue on active */
}

/* Responsive Adjustments */
@media (max-width: 1024px) {
  .modal-content {
    width: 85%;
    max-width: 350px;
  }

  h2 {
    font-size: 1.4em;
  }

  input[type="file"] {
    padding: 10px;
    font-size: 0.9em;
  }

  .btn {
    padding: 10px 18px;
    font-size: 0.9em;
  }
}

@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    max-width: 300px;
  }

  h2 {
    font-size: 1.2em;
  }

  input[type="file"] {
    padding: 8px;
    font-size: 0.8em;
  }

  .btn {
    padding: 8px 14px;
    font-size: 0.8em;
  }
}

@media (max-width: 480px) {
  .modal-content {
    width: 100%;
    max-width: 280px;
  }

  h2 {
    font-size: 1em;
  }

  input[type="file"] {
    padding: 6px;
    font-size: 0.7em;
  }

  .btn {
    padding: 6px 12px;
    font-size: 0.7em;
  }
}
</style>
