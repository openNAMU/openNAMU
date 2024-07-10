"use strict";

function opennamu_list_old_page() {
    const url = window.location.pathname;
    const url_split = url.split('/');
    
    let set_type = '';
    let num = '';
    if(url_split.length === 4) {
        set_type = url_split[3];
        num = '1';
    } else {
        set_type = url_split[3];
        num = url_split[4];
    }

    fetch('/api/v2/list/document/' + set_type + '/' + num).then(function(res) {
        return res.json();
    }).then(function(data) {
        data = data["data"];

        let data_html = '';

        for(let for_a = 0; for_a < data.length; for_a++) {
            let doc_name = opennamu_do_url_encode(data[for_a][0]);

            let right = '<a href="/w/' + doc_name + '">' + opennamu_xss_filter(data[for_a][0]) + '</a> ';

            data_html += opennamu_make_list(right, data[for_a][1]);
        }

        data_html += opennamu_page_control('/list/document/' + set_type + '/{}', Number(num), data.length);

        document.getElementById('opennamu_list_old_page').innerHTML = data_html;
    });
}