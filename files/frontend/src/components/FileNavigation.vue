<template>
  <div class="navigation" v-if="currentPath">
    <button class="btn" @click="navigateTo(previousPath)">Go Up</button>
    <button class="btn" @click="$emit('toggleUploadModal')">Upload File</button>
  </div>
</template>

<script>
export default {
  name: "FileNavigation",
  props: ["currentPath", "previousPath"],
  methods: {
    navigateTo(path) {
      this.$emit("navigateTo", path);
    },
  },
};
</script>
