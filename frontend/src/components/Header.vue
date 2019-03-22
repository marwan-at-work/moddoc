<template>
  <header>
    <div class="logo-container">
      <img class="logo" src="../assets/Go-Logo_White.png" alt="Go Logo">
    </div>
    <div class="h1-container">
      <h1>The Go Modules Proxy Documentation</h1>
    </div>
  </header>
</template>

<script>
export default {
  name: "Header"
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
header {
  background-color: #00add8;
  background-image: linear-gradient(#00a29c, #00add8);
  height: 80px;
  display: flex;
  justify-content: start;
}

.logo-container {
  width: 10%;
  min-width: 150px;
  display: flex;
  justify-content: center;
}

.h1-container {
  width: 80%;
  display: flex;
  justify-content: center;
  min-width: 650px;
  overflow: scroll;
}

h1 {
  color: white;
  font-family: "Work Sans", sans-serif;
  min-width: 650px;
  overflow: scroll;
  text-overflow: ellipsis;
}

.logo {
  width: 150px;
  height: auto;
  position: absolute;
  top: -27px;
}
</style>
