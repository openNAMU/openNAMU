"use strict";

function opennamu_list_recent_change(num, set_type) {
    let lang_data = new FormData();
    // user_document -> 8
    lang_data.append('data', 'tool normal edit move delete revert new_doc edit_request user_document raw revert compare history hide history_delete send_edit');

    fetch('/api/lang', {
        method : 'post',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];
        
        fetch('/api/auth_list').then(function(res) {
            return res.json();
        }).then(function(auth) {
            fetch('/api/recent_change/50/' + set_type + '/' + String(num)).then(function(res) {
                return res.json();
            }).then(function(data) {
                let data_html = '';

                let option_list = [
                    ['normal', lang[1]],
                    ['edit', lang[2]],
                    ['move', lang[3]],
                    ['delete', lang[4]],
                    ['revert', lang[5]],
                    ['r1', lang[6]],
                    ['edit_request', lang[7]],
                    ['user', lang[8]]
                ];
                for(let for_a = 0; for_a < option_list.length; for_a++) {
                    data_html += '<a href="/recent_change/1/' + option_list[for_a][0] + '">(' + option_list[for_a][1] + ')</a> ';
                }

                data_html += '<hr class="main_hr">'

                for(let for_a = 0; for_a < data.length; for_a++) {
                    if(data[for_a][6] !== "" && data[for_a][1] === "") {
                        data_html += '<div class="opennamu_recent_change">----</div>';
                        data_html += '<hr class="main_hr">';

                        continue;
                    }

                    let doc_name = opennamu_do_url_encode(data[for_a][1]);

                    data_html += '<div class="opennamu_recent_change">';
                    data_html += '<a href="/w/' + doc_name + '">' + opennamu_xss_filter(data[for_a][1]) + '</a> ';

                    data_html += '<div style="float: right;">';

                    data_html += '<span id="opennamu_list_recent_change_' + String(for_a) + '_over">';
                    data_html += '<a id="opennamu_list_recent_change_' + String(for_a) + '" href="javascript:void(0);">';
                    data_html += '⚒️';
                    data_html += '</a>';
                    data_html += '<span class="opennamu_popup_footnote" id="opennamu_list_recent_change_' + String(for_a) + '_load" style="display: none;"></span>';
                    data_html += '</span>';
                    data_html += ' | '

                    if(data[for_a][6] !== "") {
                        data_html += '<span style="color: red;">r' + data[for_a][0] + '</span>';
                    } else {
                        data_html += 'r' + data[for_a][0];
                    }
                    data_html += ' | '
                    
                    if(data[for_a][5] === '0') {
                        data_html += '<span style="color: gray;">' + data[for_a][5] + '</span>';
                    } else if(data[for_a][5].match(/\+/)) {
                        data_html += '<span style="color: green;">' + data[for_a][5] + '</span>';
                    } else {
                        data_html += '<span style="color: red;">' + data[for_a][5] + '</span>';
                    }
                    data_html += ' | ';
                    
                    data_html += data[for_a][7] + ' | ';
                    data_html += data[for_a][2];
                    
                    data_html += '<span style="display: none;" id="opennamu_history_tool_' + String(for_a) + '">';

                    data_html += '<a href="/raw_rev/' + data[for_a][0] + '/' + doc_name + '">' + lang[9] + '</a>';
                    data_html += ' | <a href="/revert/' + data[for_a][0] + '/' + doc_name + '">' + lang[10] + ' (r' + data[for_a][0] + ')</a>';
                    if(Number(data[for_a][0]) > 1) {
                        let before_rev = String(Number(data[for_a][0]) - 1);
                        data_html += ' | <a href="/revert/' + before_rev + '/' + doc_name + '">' + lang[10] + ' (r' + before_rev + ')</a>';
                        data_html += ' | <a href="/diff/' + before_rev + '/' + data[for_a][0] + '/' + doc_name + '">' + lang[11] + '</a>';
                    }
                    data_html += ' | <a href="/history/' + doc_name + '">' + lang[12] + '</a>';

                    if(auth["owner"] === true || auth["hidel"] === true) {
                        data_html += ' | <a href="/history_hidden/' + data[for_a][0] + '/' + doc_name + '">' + lang[13] + '</a>';
                    }

                    if(auth["owner"] === true) {
                        data_html += ' | <a href="/history_delete/' + data[for_a][0] + '/' + doc_name + '">' + lang[14] + '</a>';
                        data_html += ' | <a href="/history_send/' + data[for_a][0] + '/' + doc_name + '">' + lang[15] + '</a>';
                    }

                    data_html += '</span>';
                    
                    data_html += '</div>'
                    data_html += '<div style="clear: both;"></div>';

                    if(data[for_a][4] !== "") {
                        data_html += '<hr>'
                        data_html += opennamu_send_render(opennamu_xss_filter(data[for_a][4]));
                    }

                    data_html += '</div>';
                    data_html += '<hr class="main_hr">';
                }

                data_html += opennamu_page_control('/recent_change/{}/' + set_type, num, data.length);

                document.getElementById('opennamu_list_recent_change').innerHTML = data_html;

                for(let for_a = 0; for_a < data.length; for_a++) {
                    if(data[for_a][6] !== "" && data[for_a][1] === "") {
                        continue;
                    }

                    document.getElementById('opennamu_list_recent_change_' + String(for_a)).addEventListener("click", function() { opennamu_do_footnote_popover('opennamu_list_recent_change_' + String(for_a), '', 'opennamu_history_tool_' + String(for_a)); });
                }
            });
        });
    });
}