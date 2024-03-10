"use strict";

function opennamu_view_w_raw(render = '') {
    let name = "test";
    if(document.getElementById('opennamu_editor_doc_name')) {
        name = document.getElementById('opennamu_editor_doc_name').innerHTML.replace(/&amp;/g, '&');
    }

    let rev = "";
    if(document.getElementById("opennamu_editor_rev")) {
        rev = document.getElementById("opennamu_editor_rev").innerHTML;
    }

    let url = "";
    if(rev !== '') {
        url = "/api/raw_rev/" + rev + "/" + opennamu_do_url_encode(name);
    } else {
        url = "/api/raw/" + opennamu_do_url_encode(name);
    }

    fetch(url).then(function(res) {
        return res.json();
    }).then(function(data) {
        if(document.getElementById("opennamu_edit_textarea")) {
            if(data["data"]) {
                document.getElementById("opennamu_edit_textarea").value = data["data"];
            }

            if(render === 'do') {
                opennamu_view_w_raw_preview();
            }
        }
    });
}

function opennamu_view_w_raw_preview() {
    let name = "test";
    if(document.getElementById('opennamu_editor_doc_name')) {
        name = document.getElementById('opennamu_editor_doc_name').innerHTML.replace(/&amp;/g, '&');
    }

    let data = "";
    if(document.getElementById('opennamu_edit_textarea')) {
        data = document.getElementById('opennamu_edit_textarea').value;
    }
    console.log(data);

    opennamu_do_render('opennamu_preview_area', data, name);
}

function opennamu_view_w(do_type = '') {
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