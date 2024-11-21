<!-- ?php phpinfo()ssss ? -->
<?php
namespace Html;

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $receive = new controllers\ReceiveController();
    $receive->run();
}
?>
