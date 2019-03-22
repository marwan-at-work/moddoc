<template>
  <nav>
    <div class="left">
      <a class="package-name" v-bind:href="packageLink">{{ packageName }}</a>:
      <span class="import-path">{{ importPath }}</span>
    </div>
    <div class="right">
      <span v-for="(link, index) in navLinks" :key="link">
        <a v-bind:href="linkFragment(index)">{{ link }}</a>
        <span class="nav-seperator" v-if="index < navLinks.length-1">|</span>
      </span>
    </div>
  </nav>
</template>

<script>
export default {
  name: "PackageNav",
  props: {
    packageName: String,
    packageLink: String,
    importPath: String,
    navLinks: Array
  },
  methods: {
    linkFragment(i) {
      return `#pkg-${this.navLinks[i].toLowerCase()}`;
    }
  }
};
</script>

<style scoped>
nav {
  padding: 10px;
  background: #dbeded;
  color: #00a29c;
  display: flex;
  justify-content: space-between;
  font-family: "Roboto", sans-serif;
  border-radius: 5px;
}

a:link {
  color: #00a29c;
  font-weight: bold;
}

.nav-seperator {
  padding: 0 5px;
}

.package-name {
  font-weight: bold;
}

.package-name:hover {
  cursor: pointer;
}
</style>
