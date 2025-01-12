import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from "./router";
import store from './utils/store'
import { key } from './utils/store'
const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.use(store, key)
app.mount('#app')


window.addEventListener('beforeunload', () => {
    // 在浏览器关闭之前触发 Vuex 动作来清理会话
    store.dispatch('clearSessions');
});