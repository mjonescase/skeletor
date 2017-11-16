<template>
  <div class="d-flex align-items-end">
    <img v-bind:src="imgSrc" class="rounded-circle img-fluid m-1">
    <div class="flex-1-auto mw-100">
      <div class="mw-75 p-2 rounded bg-primary text-white">{{displayMessage}}</div>
      <div class="text-muted">{{contents.username}}</div>
    </div>
  </div>
</template>

<script>
  import md5 from 'md5';
  import emojione from 'emojione';

  export default {
    props: ['contents'],
    computed: {
      imgSrc: function () {
        console.log(this, this.contents);
        if (this.contents && this.contents.email)
          return 'http://www.gravatar.com/avatar/' + md5(this.contents.email);
        return '';
      },
      displayMessage: function () {
        return emojione.toImage(this.contents.message);
      }
    }
  }
</script>
