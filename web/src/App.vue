<template>
  <router-view />
</template>
<script lang="ts" setup>
import { provide } from 'vue';
import * as echarts from 'echarts';
provide('echarts', echarts);
</script>
<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  background-color: #F0F2F5;
}
</style>
