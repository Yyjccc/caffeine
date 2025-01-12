//webshell session 缓存

import {SystemInfo} from "../../bindings/caffeine/core";

export interface WebShellSession {
    sessionId: number;                // 会话ID
    SystemType:number                   //操作系统类型
    SystemInfo:SystemInfo
}

// 定义 Webshell 会话状态的管理结构
export interface WebShellCache {
    sessions: Record<number, WebShellSession>;  // 会话列表，键是会话ID
    activeSessionId: number | null;             // 当前活跃的会话ID
    userSettings: Record<string, any>;          // 用户偏好设置（例如：主题、字体等）
}

// 定义 Vuex store 数据结构
export interface State {
    webShellCache: WebShellCache;
}