"use strict";

function opennamu_edit_move_all() {
    let lang_data = new FormData();
    lang_data.append('data', 'title_start_document title_end_document title_include_document move document_name');

    fetch('/api/lang', {
        method : 'POST',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];
    
        document.getElementById('opennamu_edit_move_all').innerHTML = '' +
            '<input placeholder="' + lang[4] + '"></input>' +
            '<hr class="main_hr"> ' +
            '<input placeholder="' + lang[4] + '"></input>' +
            '<hr class="main_hr">' +
            '<select>' +
                '<option>' + lang[0] + '</option>' +
                '<option>' + lang[1] + '</option>' +
                '<option>' + lang[2] + '</option>' +
            '</select>' +
            '<hr class="main_hr">' +
            '<button>' + lang[3] + '</button>' +
        '';
    });
}