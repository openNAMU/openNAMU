"use strict";

function opennamu_w(do_type = '') {
    let name = "test";
    if(document.getElementById('opennamu_editor_doc_name')) {
        name = document.getElementById('opennamu_editor_doc_name').innerHTML.replace(/&amp;/g, '&');
    }

    fetch("/api/raw/" + opennamu_do_url_encode(name)).then(function(res) {
        return res.json();
    }).then(function(data) {
        if(data["data"]) {
            opennamu_do_render('opennamu_preview_area', data["data"], name, do_type);
        }
    });
}