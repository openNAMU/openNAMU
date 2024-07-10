"use strict";

function opennamu_bbs_set_post() {
    let acl_set_list = [
        "bbs_view_acl_all",
        "bbs_acl_all",
        "bbs_edit_acl_all",
        "bbs_comment_acl_all"
    ];

    for(let for_a = 0; for_a < acl_set_list.length; for_a++) {
        let post_data = new FormData();
        post_data.append('data', document.getElementById('opennamu_' + acl_set_list[for_a]).value);
        
        fetch('/api/v2/setting/' + acl_set_list[for_a], {
            method : 'PUT',
            body : post_data,
        }).then(function(res) {
            return res.json();
        }).then(function(data) {
            history.go(0);
        });
    }
}

function opennamu_bbs_set() {
    let acl_set_list = ["bbs_view_acl_all", "bbs_acl_all", "bbs_edit_acl_all", "bbs_comment_acl_all"];
    let acl_set_list_h = [2, 3, 4, 4];
    
    let lang_str = 'save reference';
    for(let for_a = 0; for_a < acl_set_list.length; for_a++) {
        lang_str += ' ' + acl_set_list[for_a];
    }

    let lang_data = new FormData();
    lang_data.append('data', lang_str);

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

            for(let for_b = 0; for_b < acl_set_list.length; for_b++) {
                acl_set_html += '<h' + acl_set_list_h[for_b] + '>' + lang[acl_set_list[for_b]] + '</h' + acl_set_list_h[for_b] + '>';
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

            document.getElementById('opennamu_bbs_set').innerHTML = renderSimpleSet('' +
                acl_set_html +
                '<hr class="main_hr">' + 
                '<button onclick="opennamu_bbs_set_post();">' + lang['save'] + '</button>' +
            '');

            for(let for_a = 0; for_a < acl_set_list.length; for_a++) {
                fetch('/api/v2/setting/' + acl_set_list[for_a]).then(function(res) {
                    return res.json();
                }).then(function(data) {
                    data = data["data"][0];

                    let select_element = document.getElementById('opennamu_' + acl_set_list[for_a]);
                    select_element.querySelector('option[value="' + data + '"]').selected = true;
                });
            }
        });
    });
}