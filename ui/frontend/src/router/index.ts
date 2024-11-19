import {createMemoryHistory, createRouter} from "vue-router";
import SystemInfo from "../components/menu/SysemInfo.vue";
import Home from "../components/Home.vue";
import WebShell from "../components/view/WebShell.vue";
import FileManger from "../components/menu/FileManger.vue";
import Terminal from "../components/menu/Terminal.vue";



const routes = [
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
        ]
    },
    {path: '/webshell/info/:id',component: SystemInfo},
    {path: '/webshell/file/:id',component: FileManger},
    //{path: '/webshell/terminal/:id',component: Terminal}
]

const router = createRouter({
    history: createMemoryHistory(),
    routes,
})

export default router