"use strict";

function opennamu_w(do_type = '') {
    let name = "test";
    if(document.getElementById('opennamu_editor_doc_name')) {
        name = opennamu_xss_filter_decode(document.getElementById('opennamu_editor_doc_name').innerHTML);
    }

    fetch("/api/raw/" + opennamu_do_url_encode(name)).then(function(res) {
        return res.json();
    }).then(function(data) {
        if(data["data"]) {
            opennamu_do_render('opennamu_preview_area', data["data"], name, do_type);
        }
    });
}

function opennamu_w_page_view() {
    let name = "test";
    if(document.getElementById('opennamu_editor_doc_name')) {
        name = opennamu_xss_filter_decode(document.getElementById('opennamu_editor_doc_name').innerHTML);
    }

    fetch("/api/v2/page_view/" + opennamu_do_url_encode(name));
}

function opennamu_w_comment_post() {
    
}

function opennamu_w_comment_delete(num) {
    let name = "test";
    if(document.getElementById('opennamu_editor_doc_name')) {
        name = opennamu_xss_filter_decode(document.getElementById('opennamu_editor_doc_name').innerHTML);
    }

    fetch("/api/v2/comment/" + num + "/" + opennamu_do_url_encode(name), {
        method : 'DELETE'
    });
}

function opennamu_w_comment() {
    let name = "test";
    if(document.getElementById('opennamu_editor_doc_name')) {
        name = opennamu_xss_filter_decode(document.getElementById('opennamu_editor_doc_name').innerHTML);
    }

    fetch("/api/v2/comment/1/" + opennamu_do_url_encode(name)).then(function(res) {
        return res.json();
    }).then(function(data) {
        console.log(data);

        if(data && data.response === 'ok') {
            let data_html = '';

            data_html += '<hr>';
            data_html += '<textarea class="opennamu_textarea_100" id="opennamu_textarea"></textarea>';
            data_html += '<hr class="main_hr">';
            data_html += '<button id="opennamu_save_button" type="submit">전송</button>';

            if(data.data.length > 0) {
                data_html += '<hr>';
            }

            document.getElementById('opennamu_w_comment').innerHTML = data_html;
        }
    });
}