"use strict";

function opennamu_list_recent_discuss() {
    const url = window.location.pathname;
    const url_split = url.split('/');
    
    let set_type = '';
    let num = '';
    if(url_split.length === 2) {
        set_type = 'normal';
        num = '1';
    } else {
        set_type = url_split[3];
        num = url_split[2];
    }

    fetch('/api/v2/recent_discuss/' + set_type + '/' + num).then(function(res) {
        return res.json();
    }).then(function(data) {
        let lang = data["language"];
        let auth = data["auth"];
        data = data["data"];

        let data_html = '';

        let option_list = [
            ['normal', lang['normal']],
            ['close', lang['close_discussion']],
            ['open', lang['open_discussion']]
        ];
        for(let for_a = 0; for_a < option_list.length; for_a++) {
            data_html += '<a href="/recent_discuss/1/' + option_list[for_a][0] + '">(' + option_list[for_a][1] + ')</a> ';
        }

        data_html += '<hr class="main_hr">';

        if(auth["hidel"] === true) {
            data_html += '<a id="opennamu_list_admin_tool_button" href="javascript:void(0);">(' + lang["admin_tool"] + ')</a>';

            data_html += '<span id="opennamu_list_admin_tool" style="display: none;">';
            data_html += 'test';
            data_html += '</span>';

            data_html += '<span class="opennamu_popup_footnote" style="display: none;" id="opennamu_list_admin_tool_button_load">';
            data_html += '</span>';

            data_html += '<hr class="main_hr">';
        }

        for(let for_a = 0; for_a < data.length; for_a++) {
            let doc_name = opennamu_do_url_encode(data[for_a][0]);

            let left = '<a href="/thread/' + data[for_a][3] + '">' + opennamu_xss_filter(data[for_a][1]) + '</a> ';
            left += '<a href="/w/' + doc_name + '">(' + opennamu_xss_filter(data[for_a][0]) + ')</a> ';

            if(auth["hidel"] === true) {
                left = '<label><input type="checkbox"> ' + left + '</label>';
            }

            let right = '';
            if(data[for_a][4] === 'O') {
                right += lang['closed'] + ' | ';
            } else if(data[for_a][4] === 'S') {
                right += lang['stop'] + ' | ';
            }

            if(data[for_a][8] !== '') {
                right += lang['agreed_discussion'] + ' | ';
            }

            right += '<a href="/thread/' + data[for_a][3] + '#' + data[for_a][7] + '">#' + data[for_a][7] + '</a> | ';
            right += data[for_a][6] + ' | ';
            right += data[for_a][2];

            data_html += opennamu_make_list(left, right);
        }

        data_html += opennamu_page_control('/recent_discuss/{}/' + set_type, Number(num), data.length);

        document.getElementById('opennamu_list_recent_discuss').innerHTML = data_html;
    });
}