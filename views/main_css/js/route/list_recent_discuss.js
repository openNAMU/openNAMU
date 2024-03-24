"use strict";

function opennamu_list_recent_discuss(tool = 'normal') {
    let lang_data = new FormData();
    lang_data.append('data', 'tool normal close_discussion open_discussion_list closed')

    fetch('/api/lang', {
        method : 'post',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];

        fetch('/api/recent_discuss/' + tool + '/50').then(function(res) {
            return res.json();
        }).then(function(data) {
            let data_html = '';

            let option_list = [
                ['normal', lang[1]],
                ['close', lang[2]],
                ['open', lang[3]]
            ];
            for(let for_a = 0; for_a < option_list.length; for_a++) {
                data_html += '<a href="/recent_discuss/' + option_list[for_a][0] + '">(' + option_list[for_a][1] + ')</a> ';
            }

            data_html += '<hr class="main_hr">'

            for(let for_a = 0; for_a < data.length; for_a++) {
                let doc_name = opennamu_do_url_encode(data[for_a][0]);

                data_html += '<div class="opennamu_recent_change">';
                data_html += '<a href="/thread/' + data[for_a][3] + '">' + opennamu_xss_filter(data[for_a][1]) + '</a> ';
                data_html += '<a href="/w/' + doc_name + '">(' + opennamu_xss_filter(data[for_a][0]) + ')</a> ';

                data_html += '<div style="float: right;">';

                if(data[for_a][4] === 'O') {
                    data_html += lang[4] + ' | ';
                }

                data_html += '<a href="/thread/' + data[for_a][3] + '#' + data[for_a][7] + '">#' + data[for_a][7] + '</a> | ';
                data_html += data[for_a][6] + ' | ';
                data_html += data[for_a][2];

                data_html += '</div>';
                data_html += '<div style="clear: both;"></div>';

                data_html += '</div>';
                data_html += '<hr class="main_hr">';
            }

            document.getElementById('opennamu_list_recent_discuss').innerHTML = data_html;
        });
    });
}