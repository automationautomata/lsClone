window.addEventListener("load", (event) => {
    console.log("page is fully loaded");
    console.log(document.getElementById("backbutton"))

    
    document.getElementById("backbutton").addEventListener('click', function (e) {
        folders = document.getElementsByClassName("path")
        console.log(folders)
        last = folders[folders.length-1]
        last.parentNode.removeChild(last)
    });    
    document.querySelectorAll('td').forEach((item) => {
        item.addEventListener('click', (event) => {
            if (event.currentTarget.parentNode.children[0].innerText !== "Folder") {
                folders = document.getElementsByClassName("path")
                last_folder = folders[folders.length - 1]
    
                next_name = event.currentTarget.parentNode.children[1].innerText;
                next_folder = '<div class="sep">/</div>' +
                            '<div class="path">' + 
                                '<font color="blue">' +
                                    next_name +
                                '</font>' + 
                            '</div>';      
                console.log(event.currentTarget.parentNode.children[1].innerText, event.currentTarget.parentNode.children);
                last_folder.insertAdjacentHTML("afterend", next_folder);
            }
        });
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
