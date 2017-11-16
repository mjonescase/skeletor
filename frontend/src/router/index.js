import Vue from 'vue';
import Router from 'vue-router';
import ChatRoomPage from '../pages/ChatRoomPage.vue';

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Chat',
      component: ChatRoomPage
    }
  ]
})
