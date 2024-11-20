
import chart from "chart.js"

const fillTable = (data) => {
    const errorprint = document.getElementById("errorprint")!;

    if (data && "Error" in data){
        errorprint.innerHTML = data["Error"];
        errorprint.parentElement!.style.display = "flex";
    } 
    else {
        errorprint.parentElement!.style.display = "none";
    }

    var tableBody = document.getElementsByTagName("tbody")[0];
    if (tableBody.innerHTML !== "") {
        tableBody.innerHTML = ""
    }

    console.log(data);
    if (!data) {
        return
    }
    for(var i = 0; i < data.length; i++) {
        var newRow = document.createElement('tr');
        newRow.innerHTML = '<td>' + data[i]["id"]   + '</td>' + 
                           '<td>' + data[i]["path"] + '</td>' + 
                           '<td>' + data[i]["size"] + '</td>' + 
                           '<td>' + data[i]["time"] + '</td>' +
                           '<td>' + data[i]["date"] + '</td>';

        tableBody.appendChild(newRow);
    }
}

const GetStatistics = () => {
    $.ajax({
        url: '/table',
        method: 'get',
        success: function (data, status, XHR) {
            console.log(status, XHR)
            fillTable(data);
        }
    });
}
var chart: Chart = null;

const ShowChart = (x, y) => {
    if (chart) {
        chart.clear();
        chart.destroy();
    }
    var ctx = document.getElementById("chart");
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
            responsive: false, 
            scales: {
                xAxes: [{
                    display: true
                }],
                yAxes: [{
                    display: true
                }]
            }
        }
    });

}

window.addEventListener("load", (e) => {
    
    document.getElementById("show")!.addEventListener('click', (e) => {
        GetStatistics();
    });    
});