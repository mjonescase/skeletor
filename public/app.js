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
	        req.send(JSON.stringify({Firstname: firstname, Lastname: lastname, Username: username, Email: email, MobileNumber: mobile, Title: title, Password: password }));
            })

        },
        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});
