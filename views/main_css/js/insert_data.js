// https://stackoverflow.com/questions/11076975/insert-text-into-textarea-at-cursor-position-javascript

function insert_data(name, data) {
    if(document.selection) { 
        document.getElementById(name).focus();

        sel = document.selection.createRange();
        sel.text = data; 
    } else if(document.getElementById(name).selectionStart || document.getElementById(name).selectionStart == '0') {
        var startPos = document.getElementById(name).selectionStart;
        var endPos = document.getElementById(name).selectionEnd;

        document.getElementById(name).value = document.getElementById(name).value.substring(0, startPos) + data + document.getElementById(name).value.substring(endPos, document.getElementById(name).value.length); 
    } else {
        document.getElementById(name).value += data;
    }
}