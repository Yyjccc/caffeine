import { createStore } from 'vuex'
import { InjectionKey } from 'vue'
import { Store } from 'vuex'
import {WebShellCache, WebShellSession} from "./session";

export interface State {
    noteContent: string
    webShellCache: WebShellCache;    // WebShell 会话缓存
}

export const key: InjectionKey<Store<State>> = Symbol()

export const store = createStore<State>({
    state: {
        noteContent: '',
        webShellCache: {
            sessions: {},
            activeSessionId: 0,
            userSettings: {}
        }
    },
    mutations: {
        SET_NOTE(state, content: string) {
            state.noteContent = content
        },
        // WebShell 会话操作
        SET_SESSION(state, session: WebShellSession) {
            state.webShellCache.sessions[session.sessionId] = session;
        },
        SET_ACTIVE_SESSION(state, sessionId: number | null) {
            state.webShellCache.activeSessionId = sessionId;
        },
        // 修改特定会话字段
        UPDATE_SESSION_FIELD(state, { sessionId, field, value }: { sessionId: number, field: string, value: any }) {
            const session = state.webShellCache.sessions[sessionId];
            if (session) {
                session[field] = value;
            }
        },
        SET_USER_SETTINGS(state, settings: Record<string, any>) {
            state.webShellCache.userSettings = settings;
        },
        DELETE_SESSION(state, sessionId: number) {
            // 删除指定 sessionId 的会话
            delete state.webShellCache.sessions[sessionId];
            if (state.webShellCache.activeSessionId === sessionId) {
                // 如果删除的会话是当前激活会话，则清空激活会话
                state.webShellCache.activeSessionId = 0;
            }
        },
        CLEAR_SESSIONS(state) {
            // 清空所有会话
            state.webShellCache.sessions = {};
            state.webShellCache.activeSessionId = 0;
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
        },
        // WebShell 会话管理
        createSession({ commit }, session: WebShellSession) {
            commit('SET_SESSION', session);
        },
        activateSession({ commit }, sessionId: string) {
            commit('SET_ACTIVE_SESSION', sessionId);
        },
        updateUserSettings({ commit }, settings: Record<string, any>) {
            commit('SET_USER_SETTINGS', settings);
        },
        // 获取特定会话
        getSession({ state }, sessionId: number): WebShellSession | null {
            return state.webShellCache.sessions[sessionId] || null; // 直接返回会话
        },

        // 修改会话中的某个字段
        updateSessionField({ commit }, { sessionId, field, value }: { sessionId: number, field: string, value: any }) {
            commit('UPDATE_SESSION_FIELD', { sessionId, field, value });
        },
        clearSessions({ commit }) {
            commit('CLEAR_SESSIONS');
        },
        deleteSession({ commit }, sessionId: number) {
            commit('DELETE_SESSION', sessionId);
        }
    }
})

export default store