<template>
  <PlaceCanvas/>
</template>

<script>
import PlaceCanvas from './components/PlaceCanvas.vue'

export default {
  name: 'App',
  components: {
    PlaceCanvas
  },
  mounted(){
  this.$components.PlaceCanvas.setup()
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
