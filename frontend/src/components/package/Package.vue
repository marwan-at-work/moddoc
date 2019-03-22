<template>
  <div class="Package">
    <PackageNav
      v-bind:packageName="doc.packageName"
      v-bind:packageLink="doc.packageLink"
      v-bind:importPath="doc.importPath"
      v-bind:navLinks="getNavLinks"
    />
    <PackageHeader
      v-bind:version="doc.moduleVersion"
      v-bind:packageName="doc.packageName"
      v-bind:importPath="doc.importPath"
      v-bind:versions="doc.versions"
    />
    <PackageDoc v-if="doc.packageDoc" v-bind:rawHTML="doc.packageDoc"/>
    <PackageIndex v-bind:doc="doc"/>
    <!-- TODO: examples index -->
    <PackageFiles v-bind:doc="doc"/>
    <h2 v-if="doc.constants && doc.constants.length" id="Constants">Constants</h2>
    <PackageVars
      v-for="variable in doc.constants"
      v-bind:variable="variable"
      :key="variable.SignatureString"
    />
    <h2 v-if="doc.variables && doc.variables.length" id="Variables">Variables</h2>
    <PackageVars
      v-for="variable in doc.variables"
      v-bind:variable="variable"
      :key="variable.SignatureString"
    />
    <PackageFunc
      v-for="func in doc.funcs"
      :key="func.Name"
      v-bind:func="func"
      v-bind:funcID="func.Name"
    />
    <PackageType v-for="type in doc.types" :key="type.Name" v-bind:type="type"/>
    <PackageSubdirectories
      v-if="doc.subdirs && doc.subdirs.length"
      v-bind:directories="doc.subdirs"
    />
  </div>
</template>

<script>
import PackageNav from "./PackageNav";
import PackageHeader from "./PackageHeader";
import PackageDoc from "./PackageDoc";
import PackageIndex from "./PackageIndex";
import PackageFiles from "./PackageFiles";
import PackageFunc from "./PackageFunc";
import PackageType from "./PackageType";
import PackageSubdirectories from "./PackageSubdirectories";
import PackageVars from "./PackageVars";

export default {
  name: "Package",
  props: {
    doc: Object // TODO: define full object
  },
  computed: {
    getNavLinks() {
      const links = ["Index"];
      if (!this.doc) {
        return links;
      }
      if (this.doc.files && this.doc.files.length) {
        links.push("Files");
      }
      if (this.doc.subdirs && this.doc.subdirs.length) {
        links.push("Directories");
      }
      return links;
    }
  },
  components: {
    PackageNav,
    PackageHeader,
    PackageDoc,
    PackageIndex,
    PackageFiles,
    PackageFunc,
    PackageType,
    PackageSubdirectories,
    PackageVars
  }
};
</script>

<style>
.Package {
  width: 50%;
  min-width: 680px;
  margin: 25px auto;
}

a:link {
  color: #00758d;
  margin-right: 5px;
}

a {
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

a:visited {
  color: #00758d;
}

.files-container {
  display: flex;
  flex-flow: row wrap;
}
</style>
