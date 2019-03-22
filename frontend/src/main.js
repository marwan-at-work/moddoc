import Vue from 'vue';
import VueRouter from 'vue-router';
import App from './App.vue';
import Home from './components/home/Home';
import Main from './Main.vue';

Vue.use(VueRouter);

const routes = [
  { path: '/', component: Home },
  { path: '/:module(.+)/@v/:version(.+)', component: App },
]

const router = new VueRouter({
  routes,
  base: __dirname,
  mode: 'history',
});


Vue.config.productionTip = false

new Vue({
  render: h => h(Main),
  router,
}).$mount('#app')
