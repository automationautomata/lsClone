<!-- ?php phpinfo()ssss ? -->
<?php

$htmlFile = 'html/index.html';

if (file_exists($htmlFile)) {
    // header('Content-Type: text/html; charset=utf-8');
    
    $content = file_get_contents($htmlFile);
    
    // Выводим содержимое файла
    echo $content;
} else {
    echo 'Файл не найден.';
}
    // namespace html;
    // $main = new controllers\MainController();
    // header('Content-Type: text/html; charset=utf-8');
    
    // echo file_get_contents("./html/index.html");

?>

