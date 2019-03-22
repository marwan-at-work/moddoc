<template>
  <div class="VersionDropDown">
    <h4 @click="this.toggle">
      {{ version }}
      <span class="arrow" v-if="!this.isOpen">+</span>
      <span class="arrow" v-if="this.isOpen">-</span>
    </h4>
    <div class="list-container" v-if="this.isOpen">
      <div v-for="ver in versions" :key="ver">
        <a v-bind:href="getVerLink(ver)">{{ ver }}</a>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  naem: "VersionDropDown",
  props: {
    version: String,
    versions: Array
  },
  methods: {
    getVerLink: function(ver) {
      return `/${this.$route.params.module}/@v/${ver}`;
    },
    toggle: function() {
      this.isOpen = !this.isOpen;
    }
  },
  data() {
    return {
      // Whether the options are visible or not.
      isOpen: false
    };
  }
};
</script>

<style scoped>
.VersionDropDown h4 {
  display: inline-block;
  padding: 5px;
  border-radius: 5px;
  margin: 0;
  cursor: pointer;
  color: white;
  background: #ce3262;
  font-size: 16px;
}

.list-container {
  padding: 10px;
}
</style>
