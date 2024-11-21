import { createStore } from 'vuex'
import { InjectionKey } from 'vue'
import { Store } from 'vuex'

export interface State {
    noteContent: string
}

export const key: InjectionKey<Store<State>> = Symbol()

export const store = createStore<State>({
    state: {
        noteContent: ''
    },
    mutations: {
        SET_NOTE(state, content: string) {
            state.noteContent = content
        }
    },
    actions: {
        saveNote({ commit }, content: string) {
            commit('SET_NOTE', content)
            // 保存到 localStorage
            localStorage.setItem('userNote', content)
        },
        loadNote({ commit }) {
            const savedNote = localStorage.getItem('userNote') || ''
            commit('SET_NOTE', savedNote)
        }
    }
})

export default store