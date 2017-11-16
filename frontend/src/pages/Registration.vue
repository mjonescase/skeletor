<template>
  <div class="container">
    <h1 class="text-center">Register</h1>
    <div class="row">
      <div class="mx-sm-auto col-sm-8 mb-1">
        <input type="text" v-model.trim="mobilenumber" placeholder="Mobile Number" class="form-control">
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <input type="text" v-model.trim="username" placeholder="Username" class="form-control">
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <input type="text" v-model.trim="firstname" placeholder="First Name" class="form-control">
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <input type="text" v-model.trim="lastname" placeholder="Last Name" class="form-control">
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <input type="text" v-model.trim="email" placeholder="me@you.com" class="form-control">
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <input type="text" v-model.trim="title" placeholder="Title (ex: CNA)" class="form-control">
      </div>
      <div class="mx-sm-auto col-sm-8 mb-1">
        <input type="password" v-model.trim="password" placeholder="Password" class="form-control">
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
  export default {
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
          this.$router.push('/');
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
