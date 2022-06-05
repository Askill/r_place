<template>
  <canvas id="main_canvas" height="1000" width="1000"> </canvas>
</template>

<script>  
var wsConnection = new WebSocket('ws://localhost:8080/get');
wsConnection.onopen = (e) => {
    console.log(`wsConnection open to 127.0.0.1:8080`, e);
};
wsConnection.onerror = (e) => {
    console.error(`wsConnection error `, e);
};
wsConnection.onmessage = (e) => {
    //var canvas = document.getElementById("main_canvas");
    //var ctx = canvas.getContext("2d");
    let data = JSON.parse(e.data)
    console.log(data);
    //ctx.fillStyle = "rgba("+data["color"][0]+","+data["color"][1]+","+data["color"][2]+","+(255)+")";
    //ctx.fillRect(data["x"], data["y"], 1, 1);
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
