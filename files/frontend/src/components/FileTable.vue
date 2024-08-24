<!-- <template>
  <table class="file-table">
    <thead>
      <tr>
        <th>Name</th>
        <th>Size</th>
        <th>Last Modified</th>
        <th>Action</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td @click="$emit('navigate', previousPath)">📁 .../</td>
        <td></td>
        <td></td>
        <td></td>
      </tr>
      <tr v-for="file in files" :key="file.path">
        <td @click="file.is_dir ? $emit('navigate', file.path) : null">
          <span v-if="file.is_dir">📁 {{ file.name }}</span>
          <span v-else>📄 {{ file.name }}</span>
        </td>
        <td>{{ file.formatted_size }}</td>
        <td>{{ file.last_modified }}</td>
        <td>
          <button class="btn" @click.prevent="$emit('rename', file)">
            Rename
          </button>
          <button class="btn" @click.prevent="$emit('delete', file)">
            Delete
          </button>
          <button
            v-if="isArchiveFile(file.name)"
            class="btn"
            @click.prevent="$emit('view-archive', file.path)"
          >
            View Archive
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script>
export default {
  props: {
    files: Array,
    previousPath: String,
  },
  methods: {
    isArchiveFile(filename) {
      const archiveExtensions = [".zip", ".tar.gz", ".tar", ".gz"];
      return archiveExtensions.some((ext) => filename.endsWith(ext));
    },
  },
};
</script>

<style scoped>
.file-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

th,
td {
  border: 1px solid #ddd;
  padding: 12px;
  text-align: left;
}

th {
  background-color: #f2f2f2;
}

tr:hover {
  background-color: #f1f1f1;
}
</style> -->
<template>
  <div>
    <table class="file-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Size</th>
          <th>Last Modified</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td @click="$emit('navigate', previousPath)">📁 .../</td>
          <td></td>
          <td></td>
          <td></td>
        </tr>
        <tr v-for="file in files" :key="file.path">
          <td @click="file.is_dir ? $emit('navigate', file.path) : null">
            <span v-if="file.is_dir">📁 {{ file.name }}</span>
            <span v-else>📄 {{ file.name }}</span>
          </td>
          <td>{{ file.formatted_size }}</td>
          <td>{{ file.last_modified }}</td>
          <td>
            <button class="btn" @click.prevent="editFile(file)">Edit</button>
            <button class="btn" @click.prevent="$emit('rename', file)">
              Rename
            </button>
            <button class="btn" @click.prevent="$emit('delete', file)">
              Delete
            </button>
            <button
              v-if="isArchiveFile(file.name)"
              class="btn"
              @click.prevent="$emit('view-archive', file.path)"
            >
              View Archive
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <EditFile
      v-if="showEditModal"
      :fileName="selectedFile.name"
      :filePath="selectedFile.path"
      @close="closeEditModal"
    />
  </div>
</template>

<script>
import EditFile from "./EditFile.vue";

export default {
  components: { EditFile },
  props: {
    files: Array,
    previousPath: String,
  },
  data() {
    return {
      showEditModal: false,
      selectedFile: {},
    };
  },
  methods: {
    isArchiveFile(filename) {
      const archiveExtensions = [".zip", ".tar.gz", ".tar", ".gz"];
      return archiveExtensions.some((ext) => filename.endsWith(ext));
    },
    editFile(file) {
      this.selectedFile = file;
      this.showEditModal = true;
    },
    closeEditModal() {
      this.showEditModal = false;
      this.selectedFile = {};
    },
  },
};
</script>

<style scoped>
/* Style for the table and buttons */
.file-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

th,
td {
  border: 1px solid #ddd;
  padding: 12px;
  text-align: left;
}

th {
  background-color: #f2f2f2;
}

tr:hover {
  background-color: #f1f1f1;
}

.btn {
  margin-left: 5px;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  background-color: #42b983;
  color: white;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: #367c62;
}
</style>
