var fillTable = function (data) {
    var errorprint = document.getElementById("errorprint");
    if (data && typeof (data) === "string") {
        errorprint.innerHTML = data["Error"];
        errorprint.parentElement.style.display = "flex";
    }
    else {
        errorprint.parentElement.style.display = "none";
    }
    var tableBody = document.getElementsByTagName("tbody")[0];
    if (tableBody.innerHTML !== "") {
        tableBody.innerHTML = "";
    }
    console.log(data);
    if (!data) {
        return;
    }
    for (var i = 0; i < data.length; i++) {
        var newRow = document.createElement('tr');
        newRow.innerHTML = '<td>' + data[i]["Id"] + '</td>' +
            '<td>' + data[i]["Path"] + '</td>' +
            '<td>' + data[i]["Size"] + '</td>' +
            '<td>' + data[i]["TimeDelta"] + '</td>' +
            '<td>' + data[i]["Date"] + '</td>';
        tableBody.appendChild(newRow);
    }
};
var GetData = function (path, handler) {
    var Http = new XMLHttpRequest();
    var url = "".concat(window.location.origin, "/").concat(path);
    Http.open("GET", url, true);
    Http.send();
    Http.onreadystatechange = function (e) {
        console.log(Http.readyState, Http.status);
        if (Http.readyState == 4 && Http.status == 200) {
            console.log(Http.responseText);
            handler(JSON.parse(Http.responseText));
        }
    };
};
//@ts-ignore
var chart;
var ShowChart = function (x, y, xName, yName) {
    if (chart) {
        chart.clear();
        chart.destroy();
    }
    var ctx = document.getElementById("chart").getContext('2d');
    //@ts-ignore
    chart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: y,
            datasets: [{
                    data: x,
                    borderColor: 'blue',
                    borderWidth: 2,
                    fill: false
                }]
        },
        options: {
            legend  : {
                display: false
            },
            scales: {
                yAxes: [{
                        scaleLabel: {
                            fontSize: 36,
                            display: true,
                            labelString: yName
                        }
                    }],
                xAxes: [{
                        scaleLabel: {
                            fontSize: 36,
                            display: true,
                            labelString: xName
                        }
                    }]
            }
        }
    });
};
window.addEventListener("load", function (e) {
    GetData("table/table.php", fillTable);
    document.getElementById("show").addEventListener('click', function (e) {
        GetData("statistics/statistics.php", function (data) {
            var x = data["SizeArray"];
            var y = data["TimeDeltaArray"];
            ShowChart(y, x, "Time", "Size");
        });
    });
});
