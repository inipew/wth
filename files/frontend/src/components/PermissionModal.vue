<template>
  <div v-if="show" class="modal-overlay" @click="close">
    <div class="modal-content" @click.stop>
      <h2>Edit File Permissions</h2>
      <table class="permission-table">
        <thead>
          <tr>
            <th></th>
            <th>Read</th>
            <th>Write</th>
            <th>Execute</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(role, index) in roles" :key="index">
            <td>{{ role }}</td>
            <td><input type="checkbox" v-model="permissions[role].read" /></td>
            <td><input type="checkbox" v-model="permissions[role].write" /></td>
            <td>
              <input type="checkbox" v-model="permissions[role].execute" />
            </td>
          </tr>
        </tbody>
      </table>
      <div class="button-group">
        <button class="btn btn-primary" @click="applyPermissions">Apply</button>
        <button class="btn btn-secondary" @click="close">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { useToast } from "vue-toastification";

export default {
  props: {
    show: {
      type: Boolean,
      required: true,
    },
    permissionsData: {
      type: Object,
      required: true,
      default: () => ({
        filepath: "",
        permissions: "",
      }),
    },
    currentPath: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      roles: ["owner", "group", "other"],
      permissions: this.parsePermissions(this.permissionsData.permissions),
    };
  },
  methods: {
    parsePermissions(permissionString) {
      const octal = parseInt(permissionString, 8);
      const [owner, group, other] = [octal >> 6, (octal >> 3) & 7, octal & 7];
      const toObject = (perm) => ({
        read: (perm & 4) !== 0,
        write: (perm & 2) !== 0,
        execute: (perm & 1) !== 0,
      });

      return {
        owner: toObject(owner),
        group: toObject(group),
        other: toObject(other),
      };
    },
    formatPermissions(permissions) {
      const toOctal = (perm) =>
        (perm.read ? 4 : 0) + (perm.write ? 2 : 0) + (perm.execute ? 1 : 0);
      return this.roles.map((role) => toOctal(permissions[role])).join("");
    },
    close() {
      this.permissions = "";
      this.$emit("close");
    },
    async applyPermissions() {
      const formattedPermissions = this.formatPermissions(this.permissions);
      const toast = useToast();
      try {
        const response = await axios.put("/api/files/permissions", {
          path: this.permissionsData.filepath,
          permissions: formattedPermissions,
        });
        toast.success(response.data.message, this.$emit("getToastOptions"));
        this.close();
        this.$emit("fetchFiles");
      } catch (error) {
        let errorMessage =
          "An error occurred while change permissions of the file.";

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
  },
};
</script>

<style scoped>
/* Modal overlay */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7); /* Darker overlay for better focus */
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.3s ease;
}

/* Modal content */
.modal-content {
  background: #ffffff;
  padding: 24px;
  border-radius: 12px;
  width: 90%;
  max-width: 600px; /* Maximum width for larger screens */
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.modal-content h2 {
  margin-top: 0;
  font-size: 1.6em; /* Slightly smaller for better scaling */
  color: #004080; /* Deep blue color for header */
  font-weight: 600;
}

/* Permission table styling */
.permission-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 24px;
}

.permission-table thead {
  background-color: #eaf4ff; /* Very light blue background for header */
}

.permission-table th,
.permission-table td {
  padding: 12px;
  text-align: center;
  border: 1px solid #dcdcdc;
  color: #333;
}

.permission-table th {
  background-color: #d0e3f0; /* Light blue background for header cells */
  color: #00264d; /* Very dark blue text color */
}

.permission-table tbody tr:nth-child(even) {
  background-color: #f7faff; /* Very light gray with blue tint for alternating rows */
}

.permission-table tbody tr:hover {
  background-color: #e2eafc; /* Light blue on hover */
}

.permission-table input[type="checkbox"] {
  transform: scale(1.2);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.permission-table input[type="checkbox"]:checked {
  transform: scale(1.2) translateY(-2px);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2); /* Light shadow for checked state */
}

/* Button styles */
.button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
}

.btn {
  padding: 10px 18px;
  font-size: 1em;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
  color: #ffffff;
  font-weight: 500;
}

/* Primary button */
.btn-primary {
  background-color: #0056b3; /* Primary dark blue */
}

.btn-primary:hover {
  background-color: #004494; /* Darker blue on hover */
}

.btn-primary:active {
  background-color: #003366; /* Even darker blue on active */
}

/* Secondary button */
.btn-secondary {
  background-color: #6c757d; /* Gray color */
}

.btn-secondary:hover {
  background-color: #5a6268; /* Darker gray on hover */
}

.btn-secondary:active {
  background-color: #343a40; /* Even darker gray on active */
}

/* Responsive adjustments */
@media (max-width: 1024px) {
  .modal-content {
    width: 85%;
    max-width: 500px;
  }

  .permission-table th,
  .permission-table td {
    padding: 10px;
    font-size: 0.9em;
  }

  .btn {
    padding: 8px 16px;
    font-size: 0.9em;
  }
}

@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    max-width: 400px;
  }

  .permission-table th,
  .permission-table td {
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
    max-width: 360px;
  }

  .permission-table th,
  .permission-table td {
    padding: 6px;
    font-size: 0.7em;
  }

  .btn {
    padding: 6px 12px;
    font-size: 0.7em;
  }
}
</style>
