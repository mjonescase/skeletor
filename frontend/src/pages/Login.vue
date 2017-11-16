<template>
  <div class='ui centered card'>
    <div class='content'>
      <div class='ui form'>
        <div class='field'>
          <label>Username</label>
          <input v-model="username" type='text' ref='username' defaultValue="">
        </div>
        <div class='field'>
          <label>Password</label>
          <input v-model="passphrase" type='password' ref='passphrase' defaultValue="">
        </div>
        <div class='ui two button attached buttons'>
          <button class='ui basic blue button' v-on:click="saveForm">
            Login
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script type="text/javascript">
  export default {
    data() {
      return {
        username: '',
        passphrase: '',
      };
    },
    methods: {
      saveForm() {
        if (this.username.length > 0 && this.passphrase.length > 0) {
          const username = this.username;
          const passphrase = this.passphrase;
          var self = this;
          var req = new XMLHttpRequest();
          req.open('POST', 'http://localhost:5000/login/', true)
          req.withCredentials = true

          req.onload = function () {
            var data = JSON.parse(req.responseText);
          }

          setTimeout(function () {
            req.send(JSON.stringify({Username: username, Password: passphrase}))
          })

          this.username = '';
          this.passphrase = '';
        }
      },
    },
  };
</script>
