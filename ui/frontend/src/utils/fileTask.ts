import { createStore } from "vuex";

export interface TransferTask {
    id: number;
    fileName: string;
    fileUrl?: string; // 下载任务需要
    file?: File; // 上传任务需要
    type: "upload" | "download"; // 任务类型
    status: "pending" | "transferring" | "completed";
    progress: number | null; // 进度：0~100，或 null 表示未开始
    completedAt?: string; // 完成时间
}

export interface State {
    transferTasks: TransferTask[];
}

export default createStore<State>({
    state: {
        transferTasks: [],
    },
    mutations: {
        ADD_TASK(state, task: TransferTask) {
            state.transferTasks.push(task);
        },
        UPDATE_TASK_PROGRESS(state, { id, progress }: { id: number; progress: number }) {
            const task = state.transferTasks.find((t) => t.id === id);
            if (task) {
                task.progress = progress;
                if (progress >= 100) {
                    task.status = "completed";
                    task.completedAt = new Date().toLocaleString();
                }
            }
        },
        START_TASK(state, id: number) {
            const task = state.transferTasks.find((t) => t.id === id);
            if (task) {
                task.status = "transferring";
                task.progress = 0;
            }
        },
    },
    actions: {
        addTransferTask({ commit }, task: TransferTask) {
            commit("ADD_TASK", { ...task, status: "pending", progress: null });
            if (task.type === "upload") {
                this.dispatch("startUpload", task.id);
            } else if (task.type === "download") {
                this.dispatch("startDownload", task.id);
            }
        },
        startDownload({ commit }, id: number) {
            commit("START_TASK", id);
            // 模拟下载过程
            let progress = 0;
            const interval = setInterval(() => {
                progress += 10;
                commit("UPDATE_TASK_PROGRESS", { id, progress });
                if (progress >= 100) clearInterval(interval);
            }, 500);
        },
        startUpload({ commit }, id: number) {
            commit("START_TASK", id);
            // 模拟上传过程
            let progress = 0;
            const interval = setInterval(() => {
                progress += 20;
                commit("UPDATE_TASK_PROGRESS", { id, progress });
                if (progress >= 100) clearInterval(interval);
            }, 500);
        },
    },
});
