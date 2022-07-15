"use strict";

function do_stop_exit() {
    window.onbeforeunload = function() {
        do_monaco_to_textarea();
        // section_edit_do();

        let data = document.getElementById('opennamu_js_edit_textarea').value;
        let origin = document.getElementById('opennamu_js_edit_origin').value;
        console.log([data], [origin]);
        if(data !== origin) {
            return '';
        }
    }
}

function do_stop_exit_release() {
    window.onbeforeunload = function () {}
}

function do_monaco_to_textarea() {
    try {
        document.getElementById('opennamu_js_edit_textarea_view').value = window.editor.getValue();
    } catch(e) {}
}

if(document.getElementById('opennamu_js_save')) {
    console.log('test');
    do_stop_exit();
    document.getElementById('opennamu_js_save').onclick = do_stop_exit_release;
}