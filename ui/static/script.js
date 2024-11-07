window.addEventListener("load", (event) => {
    console.log("page is fully loaded");
    console.log(document.getElementById("backbutton"))

    document.getElementById("pathinput").addEventListener('keyup', function (e) {
        if (e.key === 'Enter' || e.keyCode === 13) {
            //folders = document.getElementsByClassName("path")
            next = e.currentTarget.value
            // var div = document.createElement('div');
            // e.parentNode.insertBefore(div, ele);
            // div.className = 'path';
            // div.innerHTML = '<div class="mydivinside">  Text  </div>';
            next_path = '<div class="path">' + 
                            '<font color="blue">' +
                                e.currentTarget.value +
                            '</font>' + 
                        '</div>' +
                        '<div class="sep">/</div>'
    
            e.currentTarget.parentElement.insertAdjacentHTML("beforebegin", next_path);
            e.currentTarget.value = "";
        }
    });
    
    document.getElementById("backbutton").addEventListener('click', function (e) {
        folders = document.getElementsByClassName("path")
        console.log(folders)

    });    
});
  

