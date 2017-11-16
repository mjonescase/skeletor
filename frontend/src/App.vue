<template>
  <div>
    <header>
      <nav class="navbar navbar-dark bg-primary">
        <a href="/" class="navbar-brand">Simple Chat</a>
      </nav>
    </header>
    <main id="app" class="container-fluid d-flex flex-column">
      <div class="card mb-1 mt-1 flex-1-auto">
        <div id="chat-messages" class="card-content p-1">
          <message v-for="msg in messages"
                   v-bind:contents="msg.contents"></message>
        </div>
      </div>
      <div class="row d-flex flex-row flex-wrap" v-if="joined">
        <div class="col-sm-8 mb-sm">
          <input type="text" class="form-control" v-model="newMsg" @keyup.enter="send">
        </div>
        <div class="col-sm-4">
          <button class="btn btn-primary" @click="send">
            <i class="material-icons right">chat</i>
            Send
          </button>
        </div>
      </div>
      <div class="row d-flex flex-row flex-wrap" v-if="!joined">
        <div class="col-sm-8 mb-1">
          <input type="text" v-model.trim="mobilenumber" placeholder="Mobile Number" class="form-control">
        </div>
        <div class="col-sm-8 mb-1">
          <input type="text" v-model.trim="username" placeholder="Username" class="form-control">
        </div>
        <div class="col-sm-8 mb-1">
          <input type="text" v-model.trim="firstname" placeholder="First Name" class="form-control">
        </div>
        <div class="col-sm-8 mb-1">
          <input type="text" v-model.trim="lastname" placeholder="Last Name" class="form-control">
        </div>
        <div class="col-sm-8 mb-1">
          <input type="text" v-model.trim="email" placeholder="me@you.com" class="form-control">
        </div>
        <div class="col-sm-8 mb-1">
          <input type="text" v-model.trim="title" placeholder="Title (ex: CNA)" class="form-control">
        </div>
        <div class="col-sm-8 mb-1">
          <input type="password" v-model.trim="password" placeholder="Password" class="form-control">
        </div>
        <div class="col-sm-4">
          <button class="btn btn-primary btn-small" @click="join()">
            <i class="material-icons right">done</i>
            Join
          </button>
        </div>
      </div>
    </main>
    <footer class="page-footer">
    </footer>
  </div>
</template>

<script>
  import Message from './components/Message.vue';
  export default {
    name: 'app',
    components: {
      Message
    },
    data: function () {
      return {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        messages: [], // A running list of chat messages displayed on the screen
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

    created: function () {
      var self = this;
      this.ws = new WebSocket('ws://' + window.location.host + '/ws');
      this.ws.addEventListener('message', function (e) {
        var msg = JSON.parse(e.data);
        console.log(msg);
        self.messages.push(msg);
        console.log(self.messages);

        var element = document.getElementById('chat-messages');
        element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
      });
    },

    methods: {
      send: function () {
        if (this.newMsg != '') {
          this.ws.send(
            JSON.stringify({
              type: 0,
              contents: {
                email: this.mobilenumber,
                username: this.username,
                message: $('<p>').html(this.newMsg).text() // Strip out html
              }
            })
          );
          this.newMsg = ''; // Reset newMsg
        }
      },
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

