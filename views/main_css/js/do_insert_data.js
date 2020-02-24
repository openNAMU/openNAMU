// https://stackoverflow.com/questions/11076975/insert-text-into-textarea-at-cursor-position-javascript

function do_insert_data(name, data) {
    if(document.selection) {
        document.getElementById(name).focus();

        var sel = document.selection.createRange();
        sel.text = data;
    } else if(document.getElementById(name).selectionStart || document.getElementById(name).selectionStart == '0') {
        var startPos = document.getElementById(name).selectionStart;
        var endPos = document.getElementById(name).selectionEnd;
        var myPos = document.getElementById(name).value;

        document.getElementById(name).value = myPos.substring(0, startPos) + data + myPos.substring(endPos, myPos.length);
    } else {
        document.getElementById(name).value += data;
    }
}