"use strict";

function opennamu_bbs_in() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    let bbs_num = url_split[3];
    let page;
    if(url_split.length > 4) {
        page = url_split[4];
    } else {
        page = '1';
    }

    fetch('/api/v2/bbs/in/' + bbs_num + '/' + page).then(function(res) {
        return res.json();
    }).then(function(data) {
        data = data["data"];

        let data_html = '';

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

        data_html += opennamu_page_control('/bbs/in/' + bbs_num + '/{}', Number(page), data.length);
        
        document.getElementById('opennamu_bbs_in').innerHTML = data_html;
    });
}