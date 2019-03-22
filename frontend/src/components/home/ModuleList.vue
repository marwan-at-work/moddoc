<template>
  <div class="ModuleList">
    <div class="module-item" v-for="mod in filter()" :key="mod.module">
      <a v-bind:href="getLink(mod)">{{ mod.module }}</a>
    </div>
  </div>
</template>

<script>
import FuzzySearch from "fuzzy-search";

export default {
  name: "ModuleList",
  props: {
    modules: Array,
    searchInput: String
  },
  methods: {
    getLink(mod) {
      const vers = mod.versions;
      let ver = "latest";
      if (vers.length) {
        ver = vers[vers.length - 1];
      }
      return `/${mod.module}/@v/${ver}`;
    },
    filter() {
      if (!this.searchInput) {
        return this.modules;
      }

      const searcher = new FuzzySearch(this.modules, ["module"], {
        caseSensitive: false
      });
      return searcher.search(this.searchInput);
    }
  }
};
</script>

<style scoped>
.ModuleList {
  height: 800px;
  overflow: scroll;
}

.module-item {
  font-size: 24px;
  color: #00a29c;
  padding: 5px;
}
</style>
