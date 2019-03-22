<template>
  <div class="PackageDoc">
    <div class="alldocs" v-html="rawHTML"/>
  </div>
</template>

<script>
export default {
  name: "PackageDoc",
  props: {
    rawHTML: String
  },
  mounted: function() {
    this.$nextTick(function() {
      document.querySelectorAll(".PackageDoc pre").forEach(block => {
        hljs.highlightBlock(block);
      });
    });
  }
};
</script>

<style>
.PackageDoc p {
  font-family: "Roboto", sans-serif;
}

pre {
  font-family: "Source Code Pro", monospace;
  background: #eee;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 10px;
  overflow: scroll;
}
</style>
