"use strict";

function opennamu_bbs_w_set_post() {
    let acl_set_list = [
        "bbs_view_acl",
        "bbs_acl",
        "bbs_edit_acl",
        "bbs_comment_acl",
        "bbs_markup",
        "bbs_name"
    ];

    const url = window.location.pathname;
    const url_split = url.split('/');

    let set_id = url_split[3];

    for(let for_a = 0; for_a < acl_set_list.length; for_a++) {
        let post_data = new FormData();
        post_data.append('data', document.getElementById('opennamu_' + acl_set_list[for_a]).value);
        
        fetch('/api/v2/bbs/set/' + set_id + '/' + acl_set_list[for_a], {
            method : 'PUT',
            body : post_data,
        }).then(function(res) {
            return res.json();
        }).then(function(data) {
            history.go(0);
        });
    }
}

function opennamu_bbs_w_set_lang(lang, set_name) {
    if(set_name === "bbs_markup") {
        return lang["markup"];
    } else {
        return lang[set_name];
    }
}

function opennamu_bbs_w_set_select(acl_set_list, acl_set_list_h, acl_list, lang) {
    let acl_set_html = '';

    for(let for_b = 0; for_b < acl_set_list.length; for_b++) {
        acl_set_html += '<h' + acl_set_list_h[for_b] + '>' + opennamu_bbs_w_set_lang(lang, acl_set_list[for_b]) + '</h' + acl_set_list_h[for_b] + '>';
        acl_set_html += '<select id="opennamu_' + acl_set_list[for_b] + '">';
        
        let select = '';
        for(let for_a = 0; for_a < acl_list.length; for_a++) {
            let acl_list_view = acl_list[for_a];
            acl_list_view = acl_list_view === "" ? "normal" : acl_list_view;

            select += '<option value="' + acl_list[for_a] + '">' + acl_list_view + '</option>';
        }

        acl_set_html += select;
        acl_set_html += '</select>';
    }

    return acl_set_html;
}

function opennamu_bbs_w_set() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    let set_id = url_split[3];

    let acl_set_list = ["bbs_view_acl", "bbs_acl", "bbs_edit_acl", "bbs_comment_acl"];
    let acl_set_list_h = [2, 3, 4, 4];

    let markup_set_list = ["bbs_markup"];
    let markup_set_list_h = [2];
    
    let lang_str = 'save reference markup';
    for(let for_a = 0; for_a < acl_set_list.length; for_a++) {
        lang_str += ' ' + acl_set_list[for_a];
    }

    let lang_data = new FormData();
    lang_data.append('data', lang_str);

    let make_html = '';

    fetch('/api/v2/lang', {
        method : 'POST',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];

        fetch('/api/v2/list/acl/normal').then(function(res) {
            return res.json();
        }).then(function(acl_list) {
            acl_list = acl_list["data"];

            let acl_set_html = '<a href="/acl/TEST#exp">(' + lang['reference'] + ')</a>';
            acl_set_html += opennamu_bbs_w_set_select(acl_set_list, acl_set_list_h, acl_list, lang);

            make_html += acl_set_html;

            return fetch('/api/v2/list/markup');
        }).then(function(res) {
            return res.json();
        }).then(function(markup_list) {
            markup_list = markup_list["data"];

            make_html += opennamu_bbs_w_set_select(markup_set_list, markup_set_list_h, markup_list, lang);

            return;
        }).then(function() {
            document.getElementById('opennamu_bbs_w_set').innerHTML = renderSimpleSet('' +
                make_html +
                '<hr class="main_hr">' + 
                '<input id="opennamu_bbs_name">' +
                '<hr class="main_hr">' + 
                '<button onclick="opennamu_bbs_w_set_post();">' + lang['save'] + '</button>' +
            '');

            let total_set_list = [];
            total_set_list = total_set_list.concat(acl_set_list);
            total_set_list = total_set_list.concat(markup_set_list);

            for(let for_a = 0; for_a < total_set_list.length; for_a++) {
                fetch('/api/v2/bbs/set/' + set_id + '/' + total_set_list[for_a]).then(function(res) {
                    return res.json();
                }).then(function(data) {
                    data = data["data"][0][0];

                    let select_element = document.getElementById('opennamu_' + total_set_list[for_a]);
                    select_element.querySelector('option[value="' + data + '"]').selected = true;
                });
            }

            fetch('/api/v2/bbs/set/' + set_id + '/bbs_name').then(function(res) {
                return res.json();
            }).then(function(data) {
                data = data["data"][0][0];

                document.getElementById('opennamu_bbs_name').value = data;
            });
        });
    });
}