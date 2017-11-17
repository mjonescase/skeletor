<template>
  <div class='d-flex justify-content-center align-items-center'>
      <div class='form'>
        <h1>Login</h1>
        <form-group>
          <label>Username</label>
          <input v-model="username" type='text' ref='username' defaultValue="" class="form-control">
        </form-group>
        <form-group>
          <label>Password</label>
          <input v-model="passphrase" type='password' ref='passphrase' defaultValue="" class="form-control">
        </form-group>
        <button class='btn btn-primary btn-block' v-on:click="saveForm">
          Login
        </button>
        <router-link to="register" class="btn btn-secondary btn-block">Register</router-link>
      </div>
  </div>
</template>

<script type="text/javascript">
  import FormGroup from '../components/FormGroup.vue';
  export default {
    data() {
      return {
        username: '',
        passphrase: '',
      };
    },
    components: {
      FormGroup
    },
    methods: {
      saveForm() {
        if (this.username.length > 0 && this.passphrase.length > 0) {
          const username = this.username;
          const passphrase = this.passphrase;
          var self = this;
          var req = new XMLHttpRequest();
          req.open('POST', '/login/', true)
          req.withCredentials = true

          req.onload = function () {
            var data = JSON.parse(req.responseText);
            localStorage.setItem('user', req.responseText);
            if (req.status === 200) {
              self.$router.push('/chat');
            }
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
