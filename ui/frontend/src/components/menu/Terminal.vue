<template>
  <div class="terminals-container" :class="{ 'split-h': isSplitH, 'split-v': isSplitV }" @contextmenu.stop="handleContextMenu">
    <div v-for="term in terminals" 
         :key="term.id" 
         class="terminal-wrapper"
         @click="activeTerminalId = term.id">
      <div class="terminal-container" :ref="el => setTerminalRef(term.id, el)"></div>
      <div v-if="isVisible && activeTerminalId === term.id" 
           class="context-menu" 
           :style="{ left: `${x}px`, top: `${y}px` }">
        <div v-for="item in menuItems" 
             :key="item.label" 
             class="menu-item"
             @click.stop="item.action(); hide()">
          <el-icon class="menu-icon"><component :is="item.icon" /></el-icon>
          {{ item.label }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue';
import type { ComponentPublicInstance } from 'vue';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import type { ITerminalOptions, ITerminalAddon } from 'xterm';
import 'xterm/css/xterm.css';
import {useRoute } from "vue-router";
import {Exec} from "../../../wailsjs/go/client/ClientApp";
import { 
  DocumentCopy, 
  DocumentAdd, 
  ArrowLeftBold, 
  ArrowUpBold,
  Delete,
  RefreshRight 
} from '@element-plus/icons-vue'

// 定义扩展类型而不是接口
type ExtendedTerminal = Terminal & {
  buffer: {
    active: {
      cursorX: number;
      cursorY: number;
    };
  };
};

interface ITerminalInfo {
  id: number;
  instance: ExtendedTerminal | null;
  fitAddon?: FitAddon;
}

interface IMenuItem {
  icon: any; // Element Plus 图标组件类型
  label: string;
  action: () => void | Promise<void>;
}

// 获取路由参数
const route = useRoute();
const id = Number(route.params.id);
const systemInfo = JSON.parse(route.query.systemInfo as string);
const currentPath = ref(systemInfo.currentDir);
const currentUser = ref(systemInfo.currentUser);

// 终端相关状态
const terminals = ref<ITerminalInfo[]>([{ id: 1, instance: null }]);
const activeTerminalId = ref<number>(1);
const terminalRefs = ref<{ [key: number]: HTMLElement | null }>({});
const isSplitH = ref<boolean>(false);
const isSplitV = ref<boolean>(false);

// 右键菜单状态
const x = ref<number>(0);
const y = ref<number>(0);
const isVisible = ref<boolean>(false);

// 用户输入状态
let inputStartX = 0;
let inputStartY = 0;
let inputBuffer = '';
let promptLength = 0; // 添加提示符长度跟踪
let flag = false; // 赋值标志位

// 格式化输出
const formatOutput = (text: string): string => text.replace(/\n/g, '\r\n');

// 输出到终端
const terminalPrint = (term: ExtendedTerminal, text: string): void => {
  const formattedText = formatOutput(text);
  term.write(formattedText);
  cursorPosition += formattedText.length;
};

// 更新提示符
const updatePrompt = (term: ExtendedTerminal): void => {
  const prompt = `\n${currentPath.value}> `;
  promptLength = prompt.length;
  terminalPrint(term, prompt);
  flag = false;
  inputBuffer = '';
  // 重置输入起始位置
  inputStartX = term.buffer.active.cursorX;
  inputStartY = term.buffer.active.cursorY;
};

// 命令处理
const processCommand = async (input: string): Promise<string> => {
  if (input === "") {
    return "";
  }

  const [command, ...args] = input.trim().split(' ');
  let res = "";
  
  switch (command) {
    case 'cd':
      res = await Exec(id, args[0], "echo ok");
      if (res.trim() === "ok") {
        currentPath.value = args[0];
      }
      return '';
    default:
      res = await Exec(id, currentPath.value, input);
      return formatOutput(res);
  }
};

// 终端事件处理
const setupTerminalEvents = (term: ExtendedTerminal): void => {
  term.onData(async (data: string) => {
    const cursorX = term.buffer.active.cursorX;
    const cursorY = term.buffer.active.cursorY;

    if (!flag) {
      inputStartX = cursorX;
      inputStartY = cursorY;
      flag = true;
    }

    // 计算当前光标相对于输入起始位置的偏移
    const currentPos = (cursorY - inputStartY) * term.cols + (cursorX - inputStartX);

    if (data === '\r') { // Enter
      const output = await processCommand(inputBuffer.trim());
      terminalPrint(term, `\n${output}`);
      inputBuffer = '';
      updatePrompt(term);
    } else if (data === '\x7f') { // Backspace
      // 只有当光标位置大于提示符长度时才允许删除
      if (inputBuffer.length > 0 && currentPos > 0) {
        inputBuffer = inputBuffer.slice(0, -1);
        term.write('\b \b');
      }
    } else if (data === '\u001b[D') { // Left arrow
      // 不允许光标移动到提示符之前
      if (currentPos > 0) {
        term.write(data);
      }
    } else if (data === '\u001b[C') { // Right arrow
      if (currentPos < inputBuffer.length) {
        term.write(data);
      }
    } else if (!data.startsWith('\u001b')) { // 普通输入
      inputBuffer += data;
      term.write(data);
    }
  });
};

// 菜单处理函数
const show = (): void => {
  isVisible.value = true;
};

const hide = (): void => {
  isVisible.value = false;
};

// 右键菜单事件处理
const handleContextMenu = (event: MouseEvent): void => {
  event.preventDefault();
  x.value = event.clientX;
  y.value = event.clientY;
  show();
};

const handleClickOutside = (event: MouseEvent) => {
  const terminalsContainer = document.querySelector('.terminals-container');
  if (terminalsContainer && !terminalsContainer.contains(event.target as Node)) {
    hide();
  }
};

// 终端引用设置
const setTerminalRef = (id: number, el: Element | ComponentPublicInstance | null): void => {
  if (el) {
    // 如果是组件实例，获取其 $el 属性
    const element = (el as ComponentPublicInstance).$el || el;
    if (element instanceof HTMLElement) {
      terminalRefs.value[id] = element;
    }
  } else {
    terminalRefs.value[id] = null;
  }
};

// 修改拆分处理函数
const handleSplitH = (): void => {
  isSplitH.value = true;
  isSplitV.value = false;
  const newId = terminals.value.length + 1;
  terminals.value.push({ id: newId, instance: null });
  nextTick(() => {
    initTerminal(newId);
    handleResize();
  });
};

const handleSplitV = (): void => {
  isSplitV.value = true;
  isSplitH.value = false;
  const newId = terminals.value.length + 1;
  terminals.value.push({ id: newId, instance: null });
  nextTick(() => {
    initTerminal(newId);
    handleResize();
  });
};

// 修改重置所有终端的函数
const resetAllTerminals = async (): Promise<void> => {
  // 清空所有终端实例
  terminals.value.forEach(term => {
    if (term.instance) {
      try {
        // 先卸载 addon
        if (term.fitAddon) {
          term.fitAddon.dispose();
        }
        // 再销毁终端实例
        term.instance.dispose();
      } catch (error) {
        console.error('Error disposing terminal:', error);
      }
    }
  });

  // 重置状态
  terminals.value = [{ id: 1, instance: null }];
  activeTerminalId.value = 1;
  isSplitH.value = false;
  isSplitV.value = false;

  // 等待 DOM 更新
  await nextTick();
  
  // 重新初始化主终端
  initTerminal(1);
};

// 菜单项配置
const menuItems: IMenuItem[] = [
  { 
    icon: DocumentCopy, 
    label: 'Copy', 
    action: () => {
      const term = terminals.value.find(t => t.id === activeTerminalId.value)?.instance;
      if (term) {
        const selection = term.getSelection();
        if (selection) {
          navigator.clipboard.writeText(selection);
        }
      }
    }
  },
  { 
    icon: DocumentAdd, 
    label: 'Paste', 
    action: async () => {
      const text = await navigator.clipboard.readText();
      const term = terminals.value.find(t => t.id === activeTerminalId.value)?.instance;
      if (term) {
        // 处理粘贴的文本，去除可能的多余换行符
        const cleanText = text.replace(/\r\n/g, '\n').replace(/\n/g, '');
        inputBuffer += cleanText;
        term.write(cleanText);
      }
    }
  },
  { icon: ArrowLeftBold, label: 'Split Horizontal', action: handleSplitH },
  { icon: ArrowUpBold, label: 'Split Vertical', action: handleSplitV },
  { 
    icon: Delete, 
    label: 'Clear Terminal', 
    action: () => {
      const term = terminals.value.find(t => t.id === activeTerminalId.value)?.instance;
      if (term) {
        term.clear();
      }
    }
  },
  { 
    icon: RefreshRight, 
    label: 'Reset Terminal', 
    action: () => {
      resetAllTerminals();
    }
  },
];

// 修改打印欢迎信息函数
const printWelcomeMessage = (term: ExtendedTerminal): void => {
  term.write('Welcome to Vue Virtual Terminal\r\n');
  term.write(`Current User: ${currentUser}\r\n`);
  term.write(`Current Directory: ${currentPath.value}\r\n`);
  updatePrompt(term);
};

// 修改初始化终端函数
const initTerminal = (terminalId: number): void => {
  const container = terminalRefs.value[terminalId];
  if (!container) return;

  const terminalOptions = {
    cursorBlink: true,
    theme: {
      background: '#1d1f21',
      foreground: '#f8f8f2',
    },
  } satisfies ITerminalOptions;

  try {
    const term = new Terminal(terminalOptions) as ExtendedTerminal;
    const fitAddon = new FitAddon();
    
    term.loadAddon(fitAddon);
    term.open(container);
    fitAddon.fit();

    const terminalInfo = terminals.value.find(t => t.id === terminalId);
    if (terminalInfo) {
      terminalInfo.instance = term;
      terminalInfo.fitAddon = fitAddon;
    }

    setupTerminalEvents(term);
    printWelcomeMessage(term);
  } catch (error) {
    console.error('Failed to initialize terminal:', error);
  }
};

// resize 处理
const handleResize = (): void => {
  terminals.value.forEach(term => {
    if (term.instance && term.fitAddon) {
      try {
        term.fitAddon.fit();
      } catch (error) {
        console.error('Failed to fit terminal:', error);
      }
    }
  });
};

// 生命周期钩子
onMounted(async () => {
  await nextTick();
  initTerminal(1);
  window.addEventListener('resize', handleResize);
  window.addEventListener('click', handleClickOutside);
});

onBeforeUnmount(() => {
  terminals.value.forEach(term => {
    if (term.instance) {
      try {
        // 先卸载 addon
        if (term.fitAddon) {
          term.fitAddon.dispose();
        }
        // 再销毁终端实例
        term.instance.dispose();
      } catch (error) {
        console.error('Error disposing terminal:', error);
      }
    }
  });
  window.removeEventListener('resize', handleResize);
  window.removeEventListener('click', handleClickOutside);
});

// 添加光标位置跟踪
let cursorPosition = 0;

// 添加光标移动函数
const moveCursor = (term: ExtendedTerminal, direction: 'left' | 'right'): void => {
  const maxPos = inputBuffer.length;
  if (direction === 'left' && cursorPosition > 0) {
    term.write('\u001b[D');
    cursorPosition--;
  } else if (direction === 'right' && cursorPosition < maxPos) {
    term.write('\u001b[C');
    cursorPosition++;
  }
};
</script>

<style scoped>
.terminals-container {
  width: 100%;
  height: 100%;
  display: flex;
  transition: all 0.3s ease;
  overflow: hidden;
}

.terminals-container.split-h {
  flex-direction: row;
}

.terminals-container.split-v {
  flex-direction: column;
}

.terminal-wrapper {
  flex: 1 1 0;
  position: relative;
  border: 1px solid #444;
  min-height: 200px;
  transition: all 0.3s ease;
  display: flex;
  overflow: hidden;
}

.terminal-container {
  width: 100%;
  height: 100%;
  background-color: #1d1f21;
  flex: 1;
  overflow: hidden;
}

.context-menu {
  position: fixed;
  background: #2d2d2d;
  border: 1px solid #444;
  border-radius: 4px;
  padding: 4px 0;
  min-width: 150px;
  z-index: 1000;
  box-shadow: 0 2px 10px rgba(0,0,0,0.2);
}

.menu-item {
  padding: 8px 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  color: #fff;
  user-select: none;
}

.menu-item:hover {
  background: #3d3d3d;
}

.menu-icon {
  margin-right: 8px;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
