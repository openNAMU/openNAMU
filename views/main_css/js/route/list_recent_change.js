"use strict";

function opennamu_list_recent_change() {
    const option_lang = function(lang_in, lang) {
        if(lang_in === 'user') {
            lang_in = lang['user_document'];
        } else if(lang[lang_in] !== undefined) {
            lang_in = lang[lang_in];
        }
    
        return lang_in;
    }

    const url = window.location.pathname;
    const url_split = url.split('/')
    
    let set_type = '';
    let num = '';
    if(url_split.length == 2) {
        set_type = 'normal';
        num = '1';
    } else {
        set_type = url_split[3];
        num = url_split[2];
    }

    fetch('/api/v2/recent_change/' + set_type + '/' + num).then(function(res) {
        return res.json();
    }).then(function(data) {
        let lang = data["language"];
        let auth = data["auth"];
        data = data["data"];

        let data_html = '';

        let option_list = ['normal', 'edit', 'move', 'delete', 'revert', 'r1', 'edit_request', 'user', 'file', 'category'];
        for(let for_a = 0; for_a < option_list.length; for_a++) {
            data_html += '<a href="/recent_change/1/' + option_list[for_a] + '">(' + option_lang(option_list[for_a], lang) + ')</a> ';
        }

        data_html += '<hr class="main_hr">'

        let date_heading = '';
        for(let for_a = 0; for_a < data.length; for_a++) {
            if(data[for_a][6] !== "" && data[for_a][1] === "") {
                data_html += openamu_make_list('----');

                continue;
            }

            let doc_name = opennamu_do_url_encode(data[for_a][1]);
            
            let left = '<a href="/w/' + doc_name + '">' + opennamu_xss_filter(data[for_a][1]) + '</a> ';

            let right = '<span id="opennamu_list_recent_change_' + String(for_a) + '_over">';

            right += '<a id="opennamu_list_recent_change_' + String(for_a) + '" href="javascript:void(0);">';
            right += '⚒️';
            right += '</a>';
            right += '<span class="opennamu_popup_footnote" id="opennamu_list_recent_change_' + String(for_a) + '_load" style="display: none;"></span>';
            right += '</span>';
            right += ' | '

            if(data[for_a][6] !== "") {
                right += '<span style="color: red;">r' + data[for_a][0] + '</span>';
            } else {
                right += 'r' + data[for_a][0];
            }
            right += ' | '
            
            if(data[for_a][5] === '0') {
                right += '<span style="color: gray;">' + data[for_a][5] + '</span>';
            } else if(data[for_a][5].match(/\+/)) {
                right += '<span style="color: green;">' + data[for_a][5] + '</span>';
            } else {
                right += '<span style="color: red;">' + data[for_a][5] + '</span>';
            }
            right += ' | ';
            
            right += data[for_a][7] + ' | ';
            
            let edit_type = 'edit';
            if(data[for_a][8] !== '') {
                edit_type = data[for_a][8];
            }

            right += option_lang(edit_type, lang) + ' | ';

            let time_split = data[for_a][2].split(' ');

            if(date_heading !== time_split[0]) {
                data_html += '<h2>' + time_split[0] + '</h2>';
                date_heading = time_split[0];
            }

            if(time_split.length > 1) {
                right += time_split[1];
            }
            
            right += '<span style="display: none;" id="opennamu_history_tool_' + String(for_a) + '">';

            right += '<a href="/raw_rev/' + data[for_a][0] + '/' + doc_name + '">' + lang['raw'] + '</a>';
            right += ' | <a href="/revert/' + data[for_a][0] + '/' + doc_name + '">' + lang['revert'] + ' (r' + data[for_a][0] + ')</a>';
            if(Number(data[for_a][0]) > 1) {
                let before_rev = String(Number(data[for_a][0]) - 1);
                right += ' | <a href="/revert/' + before_rev + '/' + doc_name + '">' + lang['revert'] + ' (r' + before_rev + ')</a>';
                right += ' | <a href="/diff/' + before_rev + '/' + data[for_a][0] + '/' + doc_name + '">' + lang['compare'] + '</a>';
            }
            right += ' | <a href="/history/' + doc_name + '">' + lang['history'] + '</a>';

            if(auth["owner"] === true || auth["hidel"] === true) {
                right += ' | <a href="/history_hidden/' + data[for_a][0] + '/' + doc_name + '">' + lang['hide'] + '</a>';
            }

            if(auth["owner"] === true) {
                right += ' | <a href="/history_delete/' + data[for_a][0] + '/' + doc_name + '">' + lang['history_delete'] + '</a>';
                right += ' | <a href="/history_send/' + data[for_a][0] + '/' + doc_name + '">' + lang['send_edit'] + '</a>';
            }

            right += '</span>';
            
            let bottom = '';
            if(data[for_a][4] !== "") {
                bottom = opennamu_send_render(opennamu_xss_filter(data[for_a][4]));
            }

            data_html += openamu_make_list(left, right, bottom);
        }

        data_html += opennamu_page_control('/recent_change/{}/' + set_type, Number(num), data.length);

        document.getElementById('opennamu_list_recent_change').innerHTML = data_html;

        for(let for_a = 0; for_a < data.length; for_a++) {
            if(data[for_a][6] !== "" && data[for_a][1] === "") {
                continue;
            }

            document.getElementById('opennamu_list_recent_change_' + String(for_a)).addEventListener("click", function() { opennamu_do_footnote_popover('opennamu_list_recent_change_' + String(for_a), '', 'opennamu_history_tool_' + String(for_a)); });
        }
    });
}