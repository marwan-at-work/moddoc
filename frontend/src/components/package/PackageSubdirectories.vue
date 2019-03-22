<template>
  <div class="PackageSubdirectories">
    <h2 id="pkg-directories">Directories</h2>
    <div class="grid-container">
      <h3>Path</h3>
      <h3>Synopsis</h3>
      <template v-for="dir in directories">
        <span :key="dir.Name+'title'">
          <a v-bind:href="pkgLink(dir.Name)">{{ dir.Name }}</a>
        </span>
        <span :key="dir.Name+'syn'">{{ dir.Synopsis }}</span>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  name: "PackageSubdirectories",
  methods: {
    pkgLink(subdir) {
      const mod = this.$route.params.module;
      const ver = this.$route.params.version;
      return `/${mod}/${subdir}/@v/${ver}`;
    }
  },
  props: {
    directories: Array
  }
};
</script>

<style scoped>
h2 {
  margin-bottom: 0px;
}

h3 {
  margin-bottom: 10px;
}

.grid-container {
  display: grid;
  grid-template-columns: 30% 70%;
  grid-template-rows: auto auto;
}

span {
  margin-bottom: 5px;
  padding-bottom: 2.5px;
  border-bottom: 1px solid #ccc;
}
</style>
