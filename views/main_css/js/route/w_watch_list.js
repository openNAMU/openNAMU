"use strict";

function opennamu_w_watch_list(page = 1) {
    let lang_data = new FormData();
    lang_data.append('data', 'watchlist star_doc');

    fetch('/api/lang', {
        method : 'POST',
        body : lang_data
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];
        let url = window.location.pathname;

        let do_type = url.match('star_doc');
        if(do_type) {
            do_type = 'star_doc';
        } else {
            do_type = 'watch_list'
        }

        console.log();

        let split_url = url.split('/');
        let doc_name = split_url.slice(3, undefined);

        fetch('/api/' + url).then(function(res) {
            return res.json();
        }).then(function(data) {
            let data_html = '<a href="/doc_watch_list/1/' + doc_name + '">(' + lang[0] + ')</a> <a href="/doc_star_doc/1/' + doc_name + '">(' + lang[1] + ')</a>';
            data_html += '<hr class="main_hr">'

            data_html += '<ul class="opennamu_ul">';
            for(let for_a = 0; for_a < data.length; for_a++) {
                data_html += '<li><span class="opennamu_render_ip">' + data[for_a] + '</span></li>';
            }

            data_html += '</ul>';
            data_html += '<hr class="main_hr">'
            
            data_html += opennamu_page_control('/doc_' + do_type + '/{}/' + doc_name, page, data.length);

            document.getElementById('opennamu_w_watch_list').innerHTML = data_html;

            opennamu_do_ip_render();
        });
    });
}