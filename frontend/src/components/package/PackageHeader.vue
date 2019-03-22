<template>
  <div class="PackageHeader">
    <h1>package {{ packageName }}</h1>
    <h3 class="import-statement">import "{{ importPath }}"</h3>
    <VersionDropDown v-bind:versions="versions" v-bind:version="version"/>
  </div>
</template>

<script>
import VersionDropDown from "./VersionDropDown";

export default {
  name: "PackageHeader",
  props: {
    packageName: String,
    importPath: String,
    version: String,
    versions: Array
  },
  components: {
    VersionDropDown
  }
};
</script>

<style scoped>
h1 {
  margin-bottom: 5px;
}

.import-statement {
  font-family: "Source Code Pro", monospace;
}
</style>
