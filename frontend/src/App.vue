<template>
  <div id="app">
    <Header/>
    <div v-if="isLoading" class="spinner-container">
      <circles-to-rhombuses-spinner
        class="mod-spinner"
        :animation-duration="1200"
        :circles-num="3"
        :circle-size="15"
        color="#00ADD8"
      />
    </div>
    <Package v-if="!isLoading" v-bind:doc="doc"/>
  </div>
</template>

<script>
import Header from "./components/Header";
import { CirclesToRhombusesSpinner } from "epic-spinners";
import Package from "./components/package/Package";

async function fetchDoc(mod, ver) {
  const url = `/${mod}/@v/${ver}.doc`;
  const resp = await fetch(url);
  const json = await resp.json();
  return json;
}

export default {
  name: "app",
  mounted() {
    const mod = this.$route.params.module;
    const ver = this.$route.params.version;
    fetchDoc(mod, ver).then(json => {
      this.isLoading = false;
      this.doc.packageName = json.PackageName;
      this.doc.moduleVersion = json.ModuleVersion;
      this.doc.importPath = json.ImportPath;
      this.doc.packageDoc = json.PackageDoc;
      this.doc.constants = json.Constants;
      this.doc.variables = json.Variables;
      this.doc.funcs = json.Funcs;
      this.doc.types = json.Types;
      this.doc.files = json.Files;
      this.doc.subdirs = json.Subdirs;
      this.doc.versions = json.Versions;
    });
  },

  data() {
    return {
      isLoading: true,
      doc: {
        packageName: "",
        moduleVersion: "",
        importPath: "",
        packageDoc: "",
        constants: [],
        variables: [],
        funcs: [],
        types: [],
        files: [],
        subdirs: [],
        versions: []
      }
    };
  },
  components: {
    Header,
    Package,
    CirclesToRhombusesSpinner
  }
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

#app h1,
h2,
h3,
h4,
h5 {
  font-family: "Work Sans", sans-serif;
}

.spinner-container {
  margin-top: 175px;
  display: flex;
  justify-content: center;
}
</style>
