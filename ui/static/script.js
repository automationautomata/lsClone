var currentroot = "/home/alexd/Projects";
var currentsort = "asc";
const pathseparator = '\\';

function fillTable(root, sort="") {
    $.ajax({
        url: '/fs',
        method: 'get',
        dataType: 'json',
        data: {
            root: root, 
            sort: sort
        },
        success: function(entriesInfo){
            console.log(entriesInfo);
            console.log(entriesInfo, document.getElementsByTagName("tbody"));
            tableBody = document.getElementsByTagName("tbody")[0];
            if (tableBody.innerHTML !== "") {
                tableBody.innerHTML = ""
            }
            for(var i = 0; i < entriesInfo.length; i++) {
                var newRow = document.createElement('tr');
                newRow.innerHTML = '<td>' + entriesInfo[i]["Type"] + '</td>' + 
                                   '<td>' + entriesInfo[i]["Name"] + '</td>' + 
                                   '<td>' + entriesInfo[i]["ConvertedSize"] + '</td>';
                // console.log(i, entriesInfo[i]);
                tableBody.appendChild(newRow);
        
                newRow.addEventListener('click', (event) => {
                    console.log('Строка нажата:', this);
                    console.log()
                    console.log(event.currentTarget)
                    console.log(event.currentTarget.children[0].innerText);

                    if (event.currentTarget.children[0].innerText === "Folder") {
                        folders = document.getElementsByClassName("folder-name")
                        last_folder = folders[folders.length - 1]
                
                        next_name = event.currentTarget.children[1].innerText;
                        next_folder = '<div class="path-part">' + 
                                            '<div class="sep">/</div>' +
                                            '<div class="folder-name">' +
                                                '<font color="blue">' + 
                                                    next_name + 
                                                '</font>' +
                                            '</div>' +
                                      '</div>';
                        console.log(event.currentTarget.parentNode.children[1].innerText, event.currentTarget.parentNode.children);
                        last_folder.insertAdjacentHTML("afterend", next_folder);
                        currentroot += pathseparator + next_name;
                        console.log(currentroot)
                        fillTable(currentroot)
                    }
                });
            }
        }
    });
}
window.addEventListener("load", (event) => {
    console.log("page is fully loaded");
    console.log(document.getElementById("backbutton"))
    fillTable(currentroot)
    document.querySelectorAll('td').forEach((item) => {console.log(item)});

    document.getElementById("backbutton").addEventListener('click', function (e) {
        folders = document.getElementsByClassName("path-part");
        console.log(folders);
        last = folders[folders.length-1];
        last.remove();
        var tmp = currentroot.split(pathseparator);
        currentroot = tmp.slice(0, tmp.length-1).join(pathseparator);
        fillTable(currentroot);
    });    

    document.getElementById("sortbutton").addEventListener('click', function (e) {
        if (currentsort === "asc") {
            currentsort = "desc";
        } else {
            currentsort = "asc";
        }
        fillTable(currentroot, currentsort);
    });    


});


    // document.getElementById("pathinput").addEventListener('keyup', function (e) {
    //     if (e.key === 'Enter' || e.keyCode === 13) {
    //         //folders = document.getElementsByClassName("path")
    //         next = e.currentTarget.value
    //         // var div = document.createElement('div');
    //         // e.parentNode.insertBefore(div, ele);
    //         // div.className = 'path';
    //         // div.innerHTML = '<div class="mydivinside">  Text  </div>';
    //         next_path = '<div class="path">' + 
    //                         '<font color="blue">' +
    //                             e.currentTarget.value +
    //                         '</font>' + 
    //                     '</div>' +
    //                     '<div class="sep">/</div>'
    
    //         e.currentTarget.parentElement.insertAdjacentHTML("beforebegin", next_path);
    //         e.currentTarget.value = "";
    //     }
    // });
