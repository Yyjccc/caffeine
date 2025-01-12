import {createMemoryHistory, createRouter, RouteRecordRaw} from "vue-router";
import SystemInfo from "../components/menu/SysemInfo.vue";
import Home from "../components/Home.vue";
import WebShell from "../components/view/WebShell.vue";
import FileManger from "../components/menu/FileManger.vue";
import Terminal from "../components/menu/Terminal.vue";
import Monitor from "../components/menu/Monitor.vue";
import Note from "../components/menu/Note.vue";
import Other from "../components/menu/Other.vue";



const routes: RouteRecordRaw[] = [
    { path: '/', component: Home },
    {
        path: '/webshell/:id',
        component: WebShell,
        children: [
            {
                path: 'terminal', // 终端页面，传递 id 参数
                component: Terminal, // 终端页面组件
                name: 'terminal'
            },
            {
                path: "files",
                component: FileManger,
                name: "files"
            },
            {
                path: "",
                component: SystemInfo,
                name: "home"
            },
            {
                path: "monitor",
                component: Monitor,
                name: "monitor"
            },
            {
                path: "note",
                component: Note,
                name: "note"
            },
            {
                path:"other",
                component:Other,
                name:"other"
            }
        ]
    },

]

const router = createRouter({
    history: createMemoryHistory(),
    routes,
})

export default router