<!-- ?php phpinfo()ssss ? -->
<?php
// Параметры подключения к базе данных
$servername = "localhost"; // или IP-адрес сервера
$username = "username"; // имя пользователя БД
$password = "123"; // пароль БД
$dbname = "Statistics"; // имя базы данных

static $conn;
if (!isset($conn)) {
    $conn = new mysqli($servername, $username, $password, $dbname);
}
// Проверка соединения
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    // Получение данных из POST запроса
    $data = json_decode(file_get_contents("php://input"), true);
    if (isset($data['size']) && isset($data['path']) && isset($data['time'])) {

        // Подготовка и выполнение SQL-запроса для вставки данных
        $stmt = $conn->prepare("INSERT INTO FoldersInfo (path, size, timedelta) VALUES (?, ?)");
        $stmt->bind_param("sss", $data['path'], $data['size'], $data['time']);

        if ($stmt->execute()) {
            echo "Данные успешно записаны!";
        } else {
            echo "Ошибка: " . $stmt->error;
        }
        $stmt->close();
    }
}
?>
