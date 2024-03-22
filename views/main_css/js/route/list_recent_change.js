"use strict";

function opennamu_list_recent_change() {
    fetch('/api/lang/tool').then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];

        fetch('/api/recent_change/50').then(function(res) {
            return res.json();
        }).then(function(data) {
            /*
                data_list = append(data_list, []string{
                    id,
                    title,
                    date,
                    tool.IP_preprocess(db, db_set, ip, other_set["ip"])[0],
                    send,
                    leng,
                    hide,
                    tool.IP_parser(db, db_set, ip, other_set["ip"]),
                    type_data,
                })
            */

            let data_html = '<ul class="opennamu_ul">';
            for(let for_a = 0; for_a < data.length; for_a++) {
                let doc_name = opennamu_do_url_encode(data[for_a][1]);

                data_html += '<li>';
                data_html += '<a href="/w/' + doc_name + '">' + data[for_a][1] + '</a> ';

                let rev = Number(data[for_a][0]);
                if(rev <= 1) {
                    data_html += '<a href="/history/' + doc_name + '">(r' + data[for_a][0] + ')</a> ';
                } else {
                    data_html += '<a href="/diff/' + String(rev - 1) + '/' + data[for_a][0] + '/' + doc_name + '">(r' + data[for_a][0] + ')</a> ';
                }
                
                if(data[for_a][5] === '0') {
                    data_html += '<span style="color: gray;">(' + data[for_a][5] + ')</span> ';
                } else if(data[for_a][5].match(/\+/)) {
                    data_html += '<span style="color: green;">(' + data[for_a][5] + ')</span> ';
                } else {
                    data_html += '<span style="color: red;">(' + data[for_a][5] + ')</span> ';
                }
                
                data_html += '<a href="/history_tool/' + data[for_a][0] + '/' + doc_name + '">(' + lang[0] + ')</a> | ';
                data_html += data[for_a][7] + ' | ';
                data_html += data[for_a][2];
                data_html += '<br>'
                data_html += data[for_a][4];
                data_html += '</li>';
            }

            data_html += '</ul>';

            document.getElementById('opennamu_list_recent_change').innerHTML = data_html;
        });
    });
}