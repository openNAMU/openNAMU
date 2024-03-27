"use strict";

function opennamu_bbs_main() {
    fetch('/api/bbs').then(function(res) {
        return res.json();
    }).then(function(bbs_list) {
        fetch('/api/bbs/main').then(function(res) {
            return res.json();
        }).then(function(data) {
            let data_html = '<ul class="opennamu_ul">';

            for(let key in bbs_list) {
                data_html += '<li><a href="/bbs/w/' + bbs_list[key] + '">' + opennamu_xss_filter(key) + '</a></li>';
            }
    
            data_html += '</ul>';
            data_html += '<hr class="main_hr">';

            for(let for_a = 0; for_a < data.length; for_a++) {
                data_html += '<div class="opennamu_recent_change">';

                data_html += '<a href="/bbs/w/' + data[for_a]['set_id'] + '/' + data[for_a]['set_code'] + '">' + opennamu_xss_filter(data[for_a]['title']) + '</a>';

                data_html += '<div style="float: right;">';

                data_html += data[for_a]['user_id_render'] + ' | ';
                data_html += data[for_a]['date'];

                data_html += '</div>'
                data_html += '<div style="clear: both;"></div>';

                data_html += '</div>';
                data_html += '<hr class="main_hr">';
            }
            
            document.getElementById('opennamu_bbs_main').innerHTML = data_html;
        });
    })
}