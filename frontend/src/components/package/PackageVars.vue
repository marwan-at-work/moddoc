<template>
  <div class="PackageVars">
    <pre>{{ variable.SignatureString }}</pre>
    <PackageDoc v-bind:rawHTML="variable.Doc"/>
  </div>
</template>

<script>
import PackageDoc from "./PackageDoc";

export default {
  name: "PackageVars",
  props: {
    variable: Object
  },
  components: {
    PackageDoc
  }
};
</script>

<style>
</style>
