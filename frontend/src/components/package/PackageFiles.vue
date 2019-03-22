<template>
  <div class="PackageFiles">
    <h2 id="pkg-files">Package Files</h2>
    <div class="files-container">
      <a href="#variables" v-for="file in doc.files" :key="file.Name">{{ file.Name }}</a>
    </div>
  </div>
</template>

<script>
export default {
  name: "PackageFiles",
  props: {
    doc: Object
  }
};
</script>

<style scoped>
</style>
