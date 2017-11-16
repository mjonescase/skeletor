<template>
  <div>
    <main id="app" class="container-fluid d-flex flex-column">
      <div class="row">
        <div class="col-sm-8">
          <div class="card mb-1 mt-1 flex-1-auto">
            <div id="chat-messages" class="card-content p-1">
              <message v-for="msg in messages"
                       v-bind:message="msg"></message>
            </div>
          </div>
        </div>
        <div class="col-sm-4">
          <contact v-for="contact in contacts"
                   v-bind:contact="contact"></contact>
        </div>
      </div>
      <div class="row d-flex flex-row flex-wrap">
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
    </main>
    <footer class="page-footer">
    </footer>
  </div>
</template>

<script>
  import Message from '../components/Message.vue';
  import Contact from '../components/Contact.vue';

  export default {
    name: 'chatRoomPage',
    components: {
      Message,
      Contact
    },
    data: function () {
      return {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        messages: [], // A running list of chat messages displayed on the screen
        contacts: [] // A running list of chat messages displayed on the screen
      }
    },

    created: function () {
      var self = this;
      this.ws = new WebSocket('ws://' + window.location.host + '/ws');
      this.ws.addEventListener('message', function (e) {
        var data = JSON.parse(e.data);
        var msg = data.contents;
        console.log(e, msg);
        if (data.type === 0) {
          self.messages.push(msg);
          var element = document.getElementById('chat-messages');
          element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        }
        else if (data.type === 1) {
            console.log(msg);
            self.contacts=msg;
        }
      });
    },

    methods: {
      send: function () {
        if (this.newMsg != '') {
          var user = JSON.parse(localStorage.getItem('user')) || {};
          this.ws.send(
            JSON.stringify({
              type: 0,
              contents: {
                email: user.email,
                mobilenumber: user.mobilenumber,
                username: user.username,
                message: $('<p>').html(this.newMsg).text() // Strip out html
              }
            })
          );
          this.newMsg = ''; // Reset newMsg
        }
      }
    }
  }
</script>
