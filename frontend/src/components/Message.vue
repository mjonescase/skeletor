<template>
  <div class="mw-75 p-2 rounded bg-primary text-white my-1">{{displayMessage}}</div>
</template>

<script>
  import emojione from 'emojione';

  export default {
    props: ['message'],
    computed: {
      displayMessage: function () {
        return emojione.toImage(this.message);
      }
    }
  }
</script>
