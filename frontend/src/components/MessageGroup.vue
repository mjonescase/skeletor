<template>
  <div class="d-flex align-items-end">
    <img v-bind:src="imgSrc" class="rounded-circle img-fluid m-1">
    <div class="flex-1-auto mw-100">
      <message v-for="msg in message.messages"
               v-bind:message="msg"></message>
      <div class="text-muted">{{message.username}}</div>
    </div>
  </div>
</template>

<script>
  import Message from './Message.vue';
  import md5 from 'md5';

  export default {
    props: ['message'],
    components: { Message },
    computed: {
      imgSrc: function () {
        if (this.message && this.message.email)
          return 'http://www.gravatar.com/avatar/' + md5(this.message.email);
        return '';
      }
    }
  }
</script>
