<template>
  <div class="d-flex align-items-end">
    <img v-bind:src="imgSrc" class="rounded-circle img-fluid m-1">
    <div class="flex-1-auto mw-100">
      <div class="mw-75 p-2 rounded bg-primary text-white">{{displayMessage}}</div>
      <div class="text-muted">{{message.username}}</div>
    </div>
  </div>
</template>

<script>
  import md5 from 'md5';
  import emojione from 'emojione';

  export default {
    props: ['message'],
    computed: {
      imgSrc: function () {
        console.log(this, this.message);
        if (this.contents && this.message.email)
          return 'http://www.gravatar.com/avatar/' + md5(this.message.email);
        return '';
      },
      displayMessage: function () {
        return emojione.toImage(this.message.message);
      }
    }
  }
</script>
