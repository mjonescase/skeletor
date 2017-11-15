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
	title: null,
	password: null,
        joined: false // True if email and username have been filled in
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            msg.message = emojione.toImage(msg.message);
            self.messages.push(msg);

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.mobilenumber,
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text() // Strip out html
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },
        join: function () {
            if (!this.mobilenumber) {
                Materialize.toast('You must enter an mobilenumber', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
	    this.mobilenumber = $('<p>').html(this.mobilenumber).text();
	    this.firstname = $('<p>').html(this.firstname).text();
	    this.lastname = $('<p>').html(this.lastname).text();
	    this.title = $('<p>').html(this.title).text();
	    this.password = $('<p>').html(this.password).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;

	    var formData = new FormData();
	    formData.append('mobilenumber', this.mobilenumber);
	    formData.append('lastname', this.lastname);
	    formData.append('title', this.title);
	    formData.append('password', this.password);
	    formData.append('username', this.username);
	    var req = new XMLHttpRequest();

	    req.open('POST', '/register/', true);
	    req.withCredentials = true;
	    req.onload = function () {
		var data = JSON.parse(req.responseText);
		//debug here. just firing and forgetting.
	    };
	    req.send(formData);

        },
        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});
