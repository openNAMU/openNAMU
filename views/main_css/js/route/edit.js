// "use strict";

function do_stop_exit() {
    window.onbeforeunload = function() {
        do_insert_edit_data();

        let data = document.getElementById('opennamu_js_edit_textarea').value;
        let origin = document.getElementById('opennamu_js_edit_origin').value;
        if(data !== origin) {
            return '';
        }
    }
}

function do_insert_edit_data() {
    do_monaco_to_textarea();
    
    let get_data = document.getElementById('opennamu_js_edit_textarea_view').value;
    
    document.getElementById('opennamu_js_edit_textarea').value = get_data;
}

function do_stop_exit_release() {
    window.onbeforeunload = function () {}
}

function do_monaco_to_textarea() {
    try {
        document.getElementById('opennamu_js_edit_textarea_view').value = window.editor.getValue();
    } catch(e) {}
}

function do_insert_preview() {
    let s_data = new FormData();
    s_data.append('data', document.getElementById('opennamu_js_edit_textarea').value);

    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/w/Test");
    xhr.send(s_data);
    
    xhr.onreadystatechange = function() {
        if(xhr.readyState === 4 && xhr.status === 200) {
            let o_p_data = JSON.parse(xhr.responseText);
            
            console.log(o_p_data);
            
            document.getElementById('opennamu_js_preview_area').innerHTML = o_p_data['data'];
            eval(o_p_data['js_data'])
        }
    }
}


if(document.getElementById('opennamu_js_save')) {
    do_stop_exit();
    
    document.getElementById('opennamu_js_save').onclick = function() {
        do_insert_edit_data();
        do_stop_exit_release();
    };
    document.getElementById('opennamu_js_preview').onclick = function() {
        do_insert_edit_data();
        do_insert_preview();
    };
}