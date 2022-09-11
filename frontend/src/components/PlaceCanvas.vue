<template>
  <canvas id="main_canvas" height="1000" width="1000"> </canvas>
</template>

<script>  
function setup() {
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
    var background = new Image();
    background.src = "http://localhost:8080/getAll";


    wsConnection.onopen = (e) => {
        ctx.drawImage(background,0,0);   
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


export default {
  name: 'mainCanvas',
  methods: [
    setup()
  ],
  mounted(){

}
}
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
