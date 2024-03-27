"use strict";

function opennamu_bbs_main() {
    let lang_data = new FormData();
    // user_document -> 8
    lang_data.append('data', 'thread_base comment_base');

    fetch('/api/lang', {
        method : 'post',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];
    
        fetch('/api/bbs').then(function(res) {
            return res.json();
        }).then(function(bbs_list) {
            fetch('/api/bbs/main').then(function(res) {
                return res.json();
            }).then(function(data) {
                let data_html = '<ul class="opennamu_ul">';
                let bbs_id_to_name = {};

                for(let key in bbs_list) {
                    bbs_id_to_name[bbs_list[key][0]] = key;

                    data_html += '<li>';

                    data_html += '<a href="/bbs/w/' + bbs_list[key][0] + '">';
                    data_html += opennamu_xss_filter(key);
                    data_html += '</a>';

                    data_html += ' (';
                    if(bbs_list[key][1] === 'comment') {
                        data_html += lang[1];
                    } else {
                        data_html += lang[0];
                    }
                    data_html += ')';

                    data_html += ' (' + bbs_list[key][2] + ')';

                    data_html += '</li>';
                }
        
                data_html += '</ul>';
                data_html += '<hr class="main_hr">';

                for(let for_a = 0; for_a < data.length; for_a++) {
                    data_html += '<div class="opennamu_recent_change">';

                    data_html += '<a href="/bbs/w/' + data[for_a]['set_id'] + '/' + data[for_a]['set_code'] + '">' + opennamu_xss_filter(data[for_a]['title']) + '</a>';

                    data_html += '<div style="float: right;">';

                    data_html += '<a href="/bbs/w/' + data[for_a]['set_id'] + '">';
                    data_html += bbs_id_to_name[data[for_a]['set_id']]
                    data_html += '</a>';
                    data_html += ' | ';

                    data_html += data[for_a]['user_id_render'] + ' | ';
                    data_html += data[for_a]['date'];

                    data_html += '</div>'
                    data_html += '<div style="clear: both;"></div>';

                    data_html += '</div>';
                    data_html += '<hr class="main_hr">';
                }
                
                document.getElementById('opennamu_bbs_main').innerHTML = data_html;
            });
        });
    });
}