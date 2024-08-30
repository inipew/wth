<template>
  <div class="navigation-wrapper" v-if="currentPath">
    <nav class="navigation">
      <button class="btn btn-primary" @click="navigateTo(previousPath)">
        Go Up
      </button>
      <button class="btn btn-upload" @click="$emit('toggleUploadModal')">
        Upload File
      </button>
      <button class="btn btn-create" @click="$emit('openCreateEntityModal')">
        Make New
      </button>
    </nav>
  </div>
</template>

<script>
export default {
  name: "FileNavigation",
  props: {
    currentPath: {
      type: String,
      required: true,
    },
    previousPath: {
      type: String,
      required: true,
    },
  },
  methods: {
    navigateTo(path) {
      this.$emit("navigateTo", path);
    },
  },
};
</script>

<style scoped>
.navigation-wrapper {
  padding: 10px;
}

.navigation {
  padding: 16px;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  text-align: center;
  gap: 16px;
  background: #e9eff5; /* Light blue background for the navigation area */
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  overflow-x: auto; /* Allow horizontal scrolling if buttons overflow */
  box-sizing: border-box;
}

.navigation .btn {
  padding: 12px 16px;
  border: none;
  border-radius: 6px; /* Rounded corners for a modern look */
  cursor: pointer;
  transition: background-color 0.3s ease, box-shadow 0.3s ease,
    transform 0.3s ease;
  font-size: 1em; /* Default font size */
  display: flex;
  align-items: center; /* Vertically center text in buttons */
  text-align: center;
}

.navigation .btn-primary {
  background: linear-gradient(
    135deg,
    #007bff,
    #0056b3
  ); /* Gradient from bright blue to darker blue */
  color: #ffffff;
}

.navigation .btn-primary:hover {
  background: linear-gradient(
    135deg,
    #0056b3,
    #003d79
  ); /* Darker gradient on hover */
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
  transform: translateY(-2px); /* Lift effect */
}

.navigation .btn-primary:active {
  background: linear-gradient(
    135deg,
    #003d79,
    #002b54
  ); /* Even darker gradient on active */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  transform: translateY(0); /* Reset lift effect */
}

.navigation .btn-upload {
  background: #3b6c9b; /* Slightly muted blue */
  color: #ffffff;
}

.navigation .btn-upload:hover {
  background: #2a4a73; /* Darker muted blue on hover */
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
  transform: translateY(-2px); /* Lift effect */
}

.navigation .btn-upload:active {
  background: #1a334f; /* Even darker muted blue on active */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  transform: translateY(0); /* Reset lift effect */
}

.navigation .btn-create {
  background: linear-gradient(
    135deg,
    #1e90ff,
    #4682b4
  ); /* Gradient from bright sky blue to steel blue */
  color: #ffffff;
  font-weight: bold;
}

.navigation .btn-create:hover {
  background: linear-gradient(
    135deg,
    #4682b4,
    #4169e1
  ); /* Darker gradient on hover */
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
  transform: translateY(-2px); /* Lift effect */
}

.navigation .btn-create:active {
  background: linear-gradient(
    135deg,
    #4169e1,
    #3152a1
  ); /* Even darker gradient on active */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  transform: translateY(0); /* Reset lift effect */
}

/* Responsive adjustments */
@media (max-width: 1200px) {
  .navigation {
    padding: 16px;
    gap: 12px;
  }
  .navigation .btn {
    padding: 10px 14px;
    font-size: 0.9em; /* Slightly smaller font size for tablets */
  }
}

@media (max-width: 768px) {
  .navigation {
    flex-wrap: wrap;
    gap: 8px;
  }
  .navigation .btn {
    padding: 8px 12px;
    font-size: 0.85em;
    flex: 1 1 auto;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .navigation .btn {
    padding: 6px 10px;
    font-size: 0.8em; /* Smallest font size for very small screens */
  }
}
</style>
