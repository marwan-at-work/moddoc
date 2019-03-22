<template>
  <div class="Index">
    <h2 id="pkg-index">Index</h2>
    <div v-if="doc.constants.length" class="index-item">
      <a href="#Constants" class="index-item-link">Constants</a>
    </div>
    <div v-if="doc.variables.length" class="index-item">
      <a href="#Variables" class="index-item-link">Variables</a>
    </div>
    <div v-for="func in doc.funcs" :key="func.SignatureString">
      <div class="index-item">
        <a v-bind:href="getFragment(func.Name)" class="index-item-link">{{ func.SignatureString }}</a>
      </div>
    </div>
    <div v-for="type in doc.types" :key="type.SignatureString">
      <div calss="index-item">
        <a v-bind:href="getFragment(type.Name)" class="index-item-link">type {{ type.Name }}</a>
      </div>
      <ul v-if="type.Funcs && type.Funcs.length">
        <li v-for="func in type.Funcs" :key="func.SignatureString">
          <a v-bind:href="getFragment(func.Name)" class="index-item-link">{{ func.SignatureString }}</a>
        </li>
      </ul>
      <ul v-if="type.Methods && type.Methods.length">
        <li v-for="func in type.Methods" :key="func.SignatureString">
          <a
            v-bind:href="getMethodFragment(type.Name, func.Name)"
            class="index-item-link"
          >{{ func.SignatureString }}</a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: "Index",
  props: {
    doc: Object
  },
  methods: {
    getFragment: name => `#${name}`,
    getMethodFragment: (typeName, methodName) => `#${typeName}.${methodName}`
  }
};
</script>

<style scoped>
.index-item {
  margin-bottom: 5px;
}
a:link {
  /* color: #00a29c; */
  color: #00758d;
}
</style>
