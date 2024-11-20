<!-- ?php phpinfo()ssss ? -->
<?php

$htmlFile = 'index.html';

if (file_exists($htmlFile)) {
    header('Content-Type: text/html; charset=utf-8');
    
    $content = file_get_contents($htmlFile);
    
    // Выводим содержимое файла
    echo $content;
} else {
    echo 'Файл не найден.';
}
?>

