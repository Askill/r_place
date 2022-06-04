<template>
  <canvas id="main_canvas" height="1000" width="1000"> </canvas>
</template>

<script>
const wsConnection = new WebSocket('ws:localhost:8080/get', 'json');
wsConnection.onopen = (e) => {
    console.log(`wsConnection open to 127.0.0.1:8080`, e);
};
wsConnection.onerror = (e) => {
    console.error(`wsConnection error `, e);
};
wsConnection.onmessage = (e) => {
    console.log(JSON.parse(e.data));
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
canvas { border: 1px solid black; }
</style>
