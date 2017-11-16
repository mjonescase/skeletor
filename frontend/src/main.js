import Vue from 'vue'
import App from './App.vue'

new Vue({
  el: '#app',
  data: {
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
  },
  render: h => h(App)
})
