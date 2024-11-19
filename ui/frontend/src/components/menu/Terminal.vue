<template>
  <div ref="terminalContainer" class="terminal-container"></div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import 'xterm/css/xterm.css';
import { useRoute } from "vue-router";
import {Exec} from "../../../wailsjs/go/client/ClientApp";

// 获取路由参数
const route = useRoute();

const id = Number(route.params.id);
const systemInfo = JSON.parse(route.query.systemInfo as string);

// 使用 ref 创建响应式引用
const terminalContainer = ref<HTMLElement | null>(null);

// 声明终端和插件变量
let terminal: Terminal | null = null;
let fitAddon: FitAddon | null = null;

// 当前路径和用户输入
const currentPath = ref(systemInfo.currentDir);
let inputBuffer = ''; // 缓存用户输入
let inputStartX: number  = 0; // 用户输入的起始位置
let inputStartY:number =0;
let flag =false;  //赋值标志位
// 格式化输出：将 \n 转为 \r\n
const formatOutput = (text: string) => text.replace(/\n/g, '\r\n');

// 输出到终端
const terminalPrint = (text: string) => {
  if (terminal) {
    terminal.write(formatOutput(text));
  }
};

// 更新提示符并记录输入起始位置
const updatePrompt = () => {
  terminalPrint(`\n${currentPath.value}> `);
  flag=false; //重置标志位
};

// 模拟命令处理
const processCommand = async (input: string): Promise<string> => {
  const [command, ...args] = input.trim().split(' ');
  if(input==""){
    return ""
  }
  console.log(input)

  var res=""
  switch (command) {
    case 'cd':
      res= await Exec(id, args[0],";" )
      return ''; // 不输出内容，仅更新路径
    default:
      res= await Exec(id, currentPath.value, input)
  }
  return res
};

// 初始化终端
const initTerminal = () => {
  if (terminalContainer.value) {
    terminal = new Terminal({
      cursorBlink: true,
      theme: {
        background: '#1d1f21',
        foreground: '#f8f8f2',
      },
    });
    fitAddon = new FitAddon();
    terminal.loadAddon(fitAddon);

    // 挂载终端
    terminal.open(terminalContainer.value);

    // 自适应终端大小
    fitAddon.fit();

    // 初始信息
    updatePrompt();
    terminalPrint(`Welcome to Vue Virtual Terminal (bash simulation)\n`);
    terminalPrint(`Operating System: ${systemInfo.os.name} ${systemInfo.os.version} (${systemInfo.os.arch})\n`);
    terminalPrint(`Current User: ${systemInfo.currentUser}\n`);
    terminalPrint(`Current Directory: ${currentPath.value}\n`);
    updatePrompt(); // 更新提示符
    // 监听用户输入
    terminal.onData(async (data) => {
      if (!flag) {
        inputStartX = terminal!.buffer.active.cursorX ;
        inputStartY = terminal!.buffer.active.cursorY ;
        flag=true;
      }
      const cursorX = terminal!.buffer.active.cursorX; // 当前光标位置
      const cursorY = terminal!.buffer.active.cursorY;

      if (data === '\r') {
        // 回车键
        const output = await processCommand(inputBuffer.trim());
        terminalPrint(`\n${output}`);
        inputBuffer = ''; // 清空输入缓存
        updatePrompt(); // 更新提示符
      } else if (data === '\x7f') {
        // 退格键
        if (cursorX > inputStartX && cursorY >= inputStartY) {
          inputBuffer = inputBuffer.slice(0, -1);
          terminalPrint('\b \b');
        }
      } else if (data === '\u001b[A') {
        // 上箭头
        if (cursorY > inputStartY) {
          terminalPrint(data); // 允许向上移动
        }
      } else if (data === '\u001b[B') {
        // 下箭头
        if (cursorY < inputStartY) {
          terminalPrint(data); // 允许向下移动
        }
      } else if (data === '\u001b[D') {
        console.log("当前："+cursorX)
        console.log("标志："+inputStartX)
        // 左箭头
        if (cursorY > inputStartY || (cursorY == inputStartY && cursorX > inputStartX)) {
          terminalPrint(data); // 允许左移
        }
      } else if (data === '\u001b[C') {
        // 右箭头
        if (cursorY > inputStartY || cursorX < inputStartX + inputBuffer.length) {
          terminalPrint(data); // 允许右移
        }
      } else {
        // 普通输入
        inputBuffer += data;
        terminalPrint(data);
      }
    });
  }
};

// 清理终端
const disposeTerminal = () => {
  terminal?.dispose();
  fitAddon?.dispose();
};

// 生命周期钩子
onMounted(() => {
  initTerminal();
  window.addEventListener('resize', () => fitAddon?.fit());
});

onBeforeUnmount(() => {
  disposeTerminal();
  window.removeEventListener('resize', () => fitAddon?.fit());
});
</script>

<style scoped>
.terminal-container {
  width: 100%;
  height: 100%;
  min-height: 300px;
  background-color: #1d1f21;
  color: white;
  font-family: monospace;
  font-size: 14px;
  padding: 10px;
}
</style>
