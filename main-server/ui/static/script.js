"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
const $ = __importStar(require("jquery"));
var currentroot = "";
var pathseparator = '';
var currentsort = "asc";
const GetEntries = (root, sort, handler) => {
    $.ajax({
        url: '/fs',
        method: 'get',
        dataType: 'json',
        data: {
            root: root,
            sort: sort
        },
        success: function (data, status, XHR) {
            console.log(status, XHR);
            handler(data);
        }
    });
};
const handleRowClick = (event) => {
    console.log(event.currentTarget);
    if (event.currentTarget instanceof HTMLElement &&
        event.currentTarget.children instanceof HTMLCollection) {
        var children = event.currentTarget.children;
        if (children[0].innerText === "Folder") {
            let next_name = children[1].innerText;
            GetEntries(currentroot + pathseparator + next_name, "", (data) => {
                var res = fillTable(data);
                var backbutton = document.getElementById("backbutton");
                if (res === true) {
                    if (backbutton.style.visibility === 'hidden') {
                        backbutton.style.visibility = "visible";
                    }
                    var folders = document.getElementsByClassName("folder-name");
                    var last_folder = folders[folders.length - 1];
                    var next_folder = '<div class="path-part">' +
                        '<div class="sep">/</div>' +
                        '<div class="folder-name">' +
                        '<font color="blue">' +
                        next_name +
                        '</font>' +
                        '</div>' +
                        '</div>';
                    last_folder.insertAdjacentHTML("afterend", next_folder);
                    currentroot += pathseparator + next_name;
                    console.log(currentroot);
                }
            });
        }
    }
};
const handleBackClick = (e) => {
    var tmp = currentroot.split(pathseparator);
    if (tmp.length === 1 || (tmp.length === 2 && (currentroot.includes('') || currentroot.includes(':')))) {
        return;
    }
    var prevroot = tmp.slice(0, tmp.length - 1).join(pathseparator);
    GetEntries(prevroot, "", (data) => {
        var res = fillTable(data);
        if (res === true) {
            var folders = document.getElementsByClassName("path-part");
            console.log(folders);
            var last = folders[folders.length - 1];
            last.remove();
            currentroot = prevroot;
        }
        else {
            document.getElementById("backbutton").style.visibility = 'hidden';
        }
    });
};
function fillTable(entriesInfo) {
    const errorprint = document.getElementById("errorprint");
    if (entriesInfo && "Error" in entriesInfo) {
        errorprint.innerHTML = entriesInfo["Error"];
        errorprint.parentElement.style.display = "flex";
        return false;
    }
    else {
        errorprint.parentElement.style.display = "none";
    }
    var tableBody = document.getElementsByTagName("tbody")[0];
    if (tableBody.innerHTML !== "") {
        tableBody.innerHTML = "";
    }
    console.log(entriesInfo);
    if (!entriesInfo) {
        return;
    }
    for (var i = 0; i < entriesInfo.length; i++) {
        var newRow = document.createElement('tr');
        newRow.innerHTML = '<td>' + entriesInfo[i]["Type"] + '</td>' +
            '<td>' + entriesInfo[i]["Name"] + '</td>' +
            '<td>' + entriesInfo[i]["ConvertedSize"] + '</td>';
        tableBody.appendChild(newRow);
        newRow.addEventListener('click', handleRowClick);
    }
    return true;
}
window.addEventListener("load", (event) => {
    currentroot = "";
    var folders = document.getElementsByClassName("path-part");
    pathseparator = folders[0].getElementsByClassName("sep")[0].innerHTML;
    for (let i = 0; i < folders.length; i++) {
        if (pathseparator != folders[0].getElementsByClassName("sep")[0].innerHTML) {
            alert("ОШИБКА: УКАЗАНЫ РАЗНЫЕ РАЗДЕЛИТЕЛИ ПУТИ");
            return;
        }
        currentroot += pathseparator + folders[i].getElementsByClassName("folder-name")[0].children[0].innerHTML;
    }
    console.log("page is fully loaded");
    console.log("page is fully loaded");
    console.log(document.getElementById("backbutton"));
    GetEntries(currentroot, "", fillTable);
    document.getElementById("backbutton").addEventListener('click', handleBackClick);
    document.getElementById("sortbutton").addEventListener('click', (e) => {
        if (currentsort === "asc") {
            currentsort = "desc";
        }
        else {
            currentsort = "asc";
        }
        GetEntries(currentroot, currentsort, fillTable);
    });
});
