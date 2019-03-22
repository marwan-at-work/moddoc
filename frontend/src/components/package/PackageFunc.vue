<template>
  <div class="PackageFunc">
    <h2 v-bind:id="funcID">func {{ func.Name }}</h2>
    <pre>{{ func.SignatureString }}</pre>
    <PackageDoc v-bind:rawHTML="func.Doc"/>
  </div>
</template>

<script>
import PackageDoc from "./PackageDoc";

export default {
  name: "PackageFunc",
  props: {
    func: Object,
    funcID: String
  },
  components: {
    PackageDoc
  }
};
</script>

<style scoped>
pre {
  white-space: pre-wrap;
  word-break: break-all;
  word-wrap: break-word;
}
</style>
