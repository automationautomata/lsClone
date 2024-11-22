<?php
class Container {
    public $SizeArray;
    public $TimeDeltaArray;
    public function __construct($SizeArray, $TimeDeltaArray) {
        $this->SizeArray = $SizeArray;
        $this->TimeDeltaArray = $TimeDeltaArray;
    }
}

if ($_SERVER['REQUEST_METHOD'] === 'GET') {
    $conn = null;
    try {
       $conn = new PDO("mysql:host=localhost;dbname=Statistics", "hello", "123");
       $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    } catch (PDOException $exception) {
        echo "Connection error: " . $exception->getMessage();
    }

    $query = "SELECT Size, TimeDelta FROM FoldersInfo ORDER BY TimeDelta";
    $stmt = $conn->prepare($query);
    $stmt->execute();
    
    // Fetch all results as an array of User objectsz
    $rows = $stmt->fetchAll(PDO::FETCH_CLASS);  
    $stmt->closeCursor();
    $x = [];
    $y = [];
    for($i = 0; $i < count($rows); ++$i) {
        $x[$i] = $rows[$i]->Size;
        $y[$i] = $rows[$i]->TimeDelta;
    }
    
    echo json_encode(value: new Container(SizeArray: $x, TimeDeltaArray: $y));
}
?>
