"use strict";

function opennamu_w_watch_list() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    let do_type = url_split[1];
    let page = url_split[2];
    let doc_name = url_split.slice(3, undefined).join('/');

    fetch('/api/v2/' + url).then(function(res) {
        return res.json();
    }).then(function(data) {
        let lang = data["language"];
        data = data["data"];

        let data_html = '<a href="/doc_watch_list/1/' + doc_name + '">(' + lang['watchlist'] + ')</a> <a href="/doc_star_doc/1/' + doc_name + '">(' + lang['star_doc'] + ')</a>';
        data_html += '<hr class="main_hr">'

        data_html += '<ul>';
        for(let for_a = 0; for_a < data.length; for_a++) {
            data_html += '<li>' + data[for_a][1] + '</li>';
        }

        data_html += '</ul>';
        data_html += '<hr class="main_hr">'
        
        data_html += opennamu_page_control('/' + do_type + '/{}/' + doc_name, Number(page), data.length);

        document.getElementById('opennamu_w_watch_list').innerHTML = data_html;
    });
}