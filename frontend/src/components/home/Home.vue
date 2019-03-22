<template>
  <div class="Home">
    <Header/>
    <div class="search-container">
      <input placeholder="Search for modules..." v-model="searchInput" class="search" type="text">
    </div>
    <div class="results-container">
      <orbit-spinner
        v-if="loading"
        class="mod-spinner"
        :animation-duration="1200"
        :size="200"
        color="#00ADD8"
      />
      <ModuleList v-bind:modules="modules" v-bind:searchInput="searchInput"/>
    </div>
  </div>
</template>

<script>
import Header from "../Header";
import ModuleList from "./ModuleList";
import { OrbitSpinner } from "epic-spinners";

async function fetchCatalog() {
  const url = `/catalog`;
  const resp = await fetch(url);
  const json = await resp.json();
  return json;
}

export default {
  name: "Home",
  data() {
    return {
      searchInput: "",
      loading: true,
      modules: []
    };
  },
  mounted() {
    fetchCatalog().then(resp => {
      this.loading = false;
      this.modules = resp;
    });
  },
  methods: {
    doSearch() {}
  },
  components: {
    Header,
    ModuleList,
    OrbitSpinner
  }
};
</script>

<style scoped>
.search-container {
  margin-top: 125px;
  display: flex;
  justify-content: center;
}

.results-container {
  width: 50%;
  min-width: 680px;
  margin: 25px auto;
  display: flex;
  justify-content: center;
}

input {
  border: none;
  border-radius: 5px;
  border: 1px solid #ccc;
  width: 800px;
  padding: 7px;
  font-size: 36px;
  font-family: "Work Sans", sans-serif;
}

input:focus,
textarea {
  outline: none !important;
  border: 1px solid #00a29c;
}
</style>
