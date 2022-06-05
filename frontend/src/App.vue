<template>
  <HelloWorld msg="Yes"/>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'

export default {
  name: 'App',
  components: {
    HelloWorld
  },
  created(){
    var wsConnection = new WebSocket('ws://localhost:8080/get');
wsConnection.onopen = (e) => {
    console.log(`wsConnection open to 127.0.0.1:8080`, e);
};
wsConnection.onerror = (e) => {
    console.error(`wsConnection error `, e);
};
wsConnection.onmessage = (e) => {
    var canvas = document.getElementById("main_canvas");
    var ctx = canvas.getContext("2d");
    let data = JSON.parse(e.data)

    ctx.fillStyle = "rgba("+data["color"]+","+data["color"]+","+data["color"]+","+(255)+")";
    ctx.fillRect(data["y"], data["x"], 1, 1);
};
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
