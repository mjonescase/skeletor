<template>
  <div class="h-100">
    <navbar></navbar>
    <router-view></router-view>
    <footer class="page-footer"></footer>
  </div>
</template>

<script>
  import Navbar from './components/Navbar.vue';

  export default {
    name: 'app',
    components: {
      Navbar
    }
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
