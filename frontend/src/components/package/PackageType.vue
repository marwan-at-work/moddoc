<template>
  <div class="PackageType">
    <h2 v-bind:id="type.Name">type {{ type.Name }}</h2>
    <pre>{{ type.SignatureString }}</pre>
    <PackageDoc v-bind:rawHTML="type.Doc"/>
    <PackageVars
      v-for="variable in type.Constants"
      :key="variable.SignatureString.substring(0, 10)"
      v-bind:variable="variable"
    />
    <PackageVars
      v-for="variable in type.Variables"
      :key="variable.SignatureString.substring(0, 10)"
      v-bind:variable="variable"
    />
    <PackageFunc
      v-for="func in type.Funcs"
      :key="func.Name"
      v-bind:func="func"
      v-bind:funcID="func.Name"
    />
    <PackageFunc
      v-for="func in type.Methods"
      :key="type.Name + '.' + func.Name"
      v-bind:func="func"
      v-bind:funcID="type.Name + '.' + func.Name"
    />
  </div>
</template>

<script>
import PackageDoc from "./PackageDoc";
import PackageFunc from "./PackageFunc";
import PackageVars from "./PackageVars";

export default {
  name: "PackageType",
  props: {
    type: Object
  },
  components: {
    PackageDoc,
    PackageFunc,
    PackageVars
  }
};
</script>

<style scoped>
</style>
