<template>
  <router-view></router-view>
</template>

<script>
  export default {
    name: 'app'
  }
</script>

<style>
  body {
    display: flex;
    min-height: 100vh;
    flex-direction: column;
  }

  main {
    flex: 1 0 auto;
  }

  .flex-1-auto {
    flex: 1 1 auto;
  }
  .mw-75 {
    max-width: 75%;
  }
</style>
