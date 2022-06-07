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
  mounted(){
    var wsConnection = new WebSocket('ws://localhost:8080/get');
    var canvas = document.getElementById("main_canvas");
    var ctx = canvas.getContext("2d");

    var colorpalettte = [
        "#FFFFFF",
        "#E4E4E4",
        "#888888",
        "#222222",
        "#FFA7D1",
        "#E50000",
        "#E59500",
        "#A06A42",
        "#E5D900",
        "#94E044",
        "#02BE01",
        "#00D3DD",
        "#0083C7",
        "#0000EA",
        "#CF6EE4",
        "#820080"
        ]

    wsConnection.onopen = (e) => {
        console.log(`wsConnection open`, e);
    };
    wsConnection.onerror = (e) => {
        console.error(`wsConnection error `, e);
    };
    wsConnection.onmessage = (e) => {

        let data = JSON.parse(e.data)

        ctx.fillStyle = colorpalettte[parseInt(data["color"])];
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
body{
  background-color: grey;
}
</style>
