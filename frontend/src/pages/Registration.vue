<template>
  <div class="container">
    <h1 class="text-center">Register</h1>
    <div class="row">
      <div class="mx-sm-auto col-sm-8 mb-1">
        <form-group>
          <label>Username</label>
          <input type="text" v-model.trim="username" placeholder="username" class="form-control">
        </form-group>
      </div>
    </div>
    <div class="row">
      <div class="ml-sm-auto col-sm-4 mb-1">
        <form-group>
          <label>First Name</label>
          <input type="text" v-model.trim="firstname" placeholder="First Name" class="form-control">
        </form-group>
      </div>
      <div class="mr-sm-auto col-sm-4 mb-1">
        <form-group>
          <label>Last Name</label>
          <input type="text" v-model.trim="lastname" placeholder="Last Name" class="form-control">
        </form-group>
      </div>
    </div>
    <div class="row">
      <div class="mx-sm-auto col-sm-8 mb-1">
        <form-group>
          <label>Email</label>
          <input type="text" v-model.trim="email" placeholder="youremail@provider.com" class="form-control">
        </form-group>
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <form-group>
          <label>Mobile Number</label>
          <input type="text" v-model.trim="mobilenumber" placeholder="###-###-####" class="form-control">
        </form-group>
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <form-group>
          <label>Title</label>
          <input type="text" v-model.trim="title" placeholder="CNA" class="form-control">
        </form-group>
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <form-group>
          <label>Password</label>
          <input type="password" v-model.trim="password" class="form-control">
        </form-group>
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <button class="btn btn-primary" @click="join()">
          Join
        </button>
      </div>
    </div>
  </div>
</template>

<script>
  import FormGroup from '../components/FormGroup.vue';

  export default {
    components: { FormGroup },
    data: function () {
      return {
        mobilenumber: null,
        username: null,
        firstname: null,
        lastname: null,
        email: null,
        title: null,
        password: null,
        joined: false // True if email and username have been filled in
      }
    },
    methods: {
      join: function () {
        var self = this;
        if (!this.mobilenumber) {
          Materialize.toast('You must enter an mobilenumber', 2000);  //TODO Materialize isn't a thing - this is broken
          return
        }
        if (!this.username) {
          Materialize.toast('You must choose a username', 2000);
          return
        }
        const mobile = this.mobilenumber;
        const email = this.email;
        const firstname = this.firstname;
        const lastname = this.lastname;
        const title = this.title;
        const password = this.password;
        const username = this.username;
        this.joined = true;

        var req = new XMLHttpRequest();
        req.open('POST', '/register/', true);
        req.withCredentials = true;
        req.onload = function () {
          var data = JSON.parse(req.responseText);
          if (req.status === 200) {
            self.$router.push('/');
          }
          //debug here. just firing and forgetting.
        };
        setTimeout(function () {
          req.send(JSON.stringify({
            Firstname: firstname,
            Lastname: lastname,
            Username: username,
            Email: email,
            MobileNumber: mobile,
            Title: title,
            Password: password
          }));
        })
      }
    }
  }
</script>
