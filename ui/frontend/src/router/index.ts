import {createMemoryHistory, createRouter} from "vue-router";
import ShellList from "../components/view/ShellList.vue";

const routes = [
    { path: '/', component: ShellList },

]

const router = createRouter({
    history: createMemoryHistory(),
    routes,
})

export default router