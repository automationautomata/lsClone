var currentroot = "/home/alexd/Projects";
var currentsort = "asc";
const pathseparator = '/';

function fillTable(entriesInfo) {
    console.log(entriesInfo);

    if ("Error" in entriesInfo){
        alert(entriesInfo["Error"])
        return false
    }
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

        tableBody.appendChild(newRow);

        newRow.addEventListener('click', (event) => {
            console.log(event.currentTarget)
            if (event.currentTarget.children[0].innerText === "Folder") {
                next_name = event.currentTarget.children[1].innerText    
                GetEntries(currentroot+pathseparator+next_name, "", (data) => {
                    var res = fillTable(data)

                    if (res === true) {
                        folders = document.getElementsByClassName("folder-name")
                        last_folder = folders[folders.length - 1]

                        next_folder = '<div class="path-part">' + 
                                            '<div class="sep">/</div>' +
                                            '<div class="folder-name">' +
                                                '<font color="blue">' + 
                                                    next_name + 
                                                '</font>' +
                                            '</div>' +
                                    '</div>';

                        last_folder.insertAdjacentHTML("afterend", next_folder);
                        console.log(next_name)
                        currentroot += pathseparator+next_name
                    }
                })
            }
        })
    }
    return true
}

function GetEntries(root, sort, handler) {
    $.ajax({
        url: '/fs',
        method: 'get',
        dataType: 'json',
        data: {
            root: root, 
            sort: sort
        },
        success: handler
    });
}

window.addEventListener("load", (event) => {
    console.log("page is fully loaded");
    console.log(document.getElementById("backbutton"))
    var folders = currentroot.split(pathseparator).filter((val, i, arr) => { return val !== "" });
    var parent = document.getElementsByClassName('path-container')[0]
    for (var i = 0; i < folders.length; i++) {
        next_folder = '<div class="path-part">' + 
                            '<div class="sep">/</div>' +
                            '<div class="folder-name">' +
                                '<font color="blue">' + 
                                    folders[i] + 
                                '</font>' +
                            '</div>' +
                      '</div>';
        parent.innerHTML += next_folder;
    }
    GetEntries(currentroot, "", fillTable)

    document.getElementById("backbutton").addEventListener('click', function (e) {
        var tmp = currentroot.split(pathseparator);
        prevroot = tmp.slice(0, tmp.length-1).join(pathseparator);

        GetEntries(prevroot, "", (data) => {
            var res = fillTable(data)

            if (res === true) {
                folders = document.getElementsByClassName("path-part");
                console.log(folders);
                
                last = folders[folders.length-1];
                last.remove();      

                currentroot = prevroot;
            }
        })
    });    

    document.getElementById("sortbutton").addEventListener('click', function (e) {
        if (currentsort === "asc") {
            currentsort = "desc";
        } else {
            currentsort = "asc";
        }
        GetEntries(currentroot, sort=currentsort, handler=fillTable);
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
