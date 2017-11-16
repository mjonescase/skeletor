import Vue from 'vue';
import Router from 'vue-router';
import Login from '../pages/Login.vue';
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
    }
  ]
})
