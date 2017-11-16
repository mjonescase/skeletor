import Vue from 'vue';
import Router from 'vue-router';
import Login from '../pages/Login.vue';
import Registration from '../pages/Registration.vue';
import ChatRoomPage from '../pages/ChatRoomPage.vue';

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Chat',
      component: ChatRoomPage
    },
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/register',
      name: 'Register',
      component: Registration
    }
  ]
})
