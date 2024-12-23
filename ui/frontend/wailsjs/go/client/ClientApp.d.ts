// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {webshell} from '../models';
import {client} from '../models';
import {core} from '../models';

export function AddNewShell(arg1:{[key: string]: any}):Promise<number>;

export function CloseTerminal(arg1:number):Promise<void>;

export function CreateTerminal(arg1:number):Promise<webshell.TerminalInfo>;

export function Exec(arg1:number,arg2:string,arg3:string):Promise<string>;

export function ExecuteCommand(arg1:number,arg2:string):Promise<string>;

export function GetActiveConnections():Promise<Array<client.Connection>>;

export function GetListeningPorts():Promise<Array<client.Port>>;

export function GetLocalSystemMetrics():Promise<client.SystemMetric>;

export function GetNetworkInterfaces():Promise<Array<client.NetworkInterface>>;

export function GetNextCommand(arg1:number):Promise<string>;

export function GetPreviousCommand(arg1:number):Promise<string>;

export function GetShellID():Promise<number>;

export function GetShellList(arg1:number):Promise<Array<client.ShellEntry>>;

export function GetTerminalEnvironment(arg1:number,arg2:string):Promise<string>;

export function GetTerminalInfo(arg1:number):Promise<webshell.TerminalInfo>;

export function GetTerminalPrompt(arg1:number):Promise<string>;

export function GetTerminalWelcomeMessage(arg1:number):Promise<string>;

export function InitShell(arg1:number):Promise<core.SystemInfo>;

export function ListTerminals():Promise<Array<webshell.TerminalInfo>>;

export function SetTerminalEnvironment(arg1:number,arg2:string,arg3:string):Promise<void>;

export function TestConnect(arg1:number):Promise<boolean>;
