"use strict";

function opennamu_view_raw_document(render = '') {
    let name = "test";
    if(document.getElementById('opennamu_editor_doc_name')) {
        name = document.getElementById('opennamu_editor_doc_name').innerHTML;
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
                opennamu_view_raw_document_preview();
            }
        }
    });
}

function opennamu_view_raw_document_preview() {
    let name = "test";
    if(document.getElementById('opennamu_editor_doc_name')) {
        name = document.getElementById('opennamu_editor_doc_name').innerHTML;
    }

    let data = "";
    if(document.getElementById('opennamu_edit_textarea')) {
        data = document.getElementById('opennamu_edit_textarea').value;
    }

    opennamu_do_render('opennamu_preview_area', name, data);
}