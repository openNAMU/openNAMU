"use strict";

function opennamu_edit_move_all() {
    let lang_data = new FormData();
    lang_data.append('data', 'title_start_document title_end_document title_include_document move document_name');

    fetch('/api/v2/lang', {
        method : 'POST',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];
    
        document.getElementById('opennamu_edit_move_all').innerHTML = '' +
            '<input placeholder="' + lang['document_name'] + '"></input>' +
            '<hr class="main_hr"> ' +
            '<input placeholder="' + lang['document_name'] + '"></input>' +
            '<hr class="main_hr">' +
            '<select>' +
                '<option>' + lang['title_start_document'] + '</option>' +
                '<option>' + lang['title_end_document'] + '</option>' +
                '<option>' + lang['title_include_document'] + '</option>' +
            '</select>' +
            '<hr class="main_hr">' +
            '<button>' + lang['move'] + '</button>' +
        '';
    });
}