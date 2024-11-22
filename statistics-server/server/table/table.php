<?php
use PDO;
use PDOException;
try  {
    $conn = null;
    try {
       $conn = new PDO("mysql:host=localhost;dbname=Statistics", "hello", "123");
       $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    } catch (PDOException $exception) {
        echo "Connection error: " . $exception->getMessage();
    }

    if ($_SERVER['REQUEST_METHOD'] === 'POST') {

        // Получение данных из POST запроса
        $data = json_decode(json: file_get_contents(filename: "php://input"), associative: true);
        if (isset($data['size']) && isset($data['path']) && isset($data['time'])) {
            // Подготовка и выполнение SQL-запроса для вставки данных
            
            $queryInsert = "INSERT INTO FoldersInfo (Path, Size, TimeDelta) VALUES (:p, :s, :td)";
            $stmt = $conn->prepare($queryInsert);

            $stmt->bindParam(':p', $data['path']);
            $stmt->bindParam(':s', $data['size']);
            $stmt->bindParam(':td', $data['time']);
            $stmt->execute();
            echo "33";

            if (true) {
                $conn->commit();
                $stmt->closeCursor();
                echo "Данные успешно записаны!";
            } else {
                $stmt->closeCursor();
                echo "Ошибка: " + $stmt->errorInfo();
            }
        }
    }
    if ($_SERVER['REQUEST_METHOD'] === 'GET') {

        $query = "SELECT Id, Path, Size, TimeDelta, Date FROM FoldersInfo;";
        $stmt = $conn->prepare(query: $query);
        $stmt->execute();
        
        $rows = $stmt->fetchAll(PDO::FETCH_CLASS);
        $stmt->closeCursor();
        echo json_encode(value: $rows);
    }
}
catch(PDOException $e) {
    echo $e->getMessage();
}

?>

