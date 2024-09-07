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

function opennamu_w_comment() {
    
}