"use strict";

function opennamu_topic_list() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    let num;
    if(url_split[1] === 'topic') {
        num = '1';
    } else {
        num = url_split[2];
    }
    
    let doc_name;
    if(url_split[1] === 'topic') {
        doc_name = url_split.slice(2, undefined).join('/');
    } else {
        doc_name = url_split.slice(3, undefined).join('/');
    }

    fetch('/api/v2/topic/' + num + '/normal/' + doc_name).then(function(res) {
        return res.json();
    }).then(function(data) {
        let lang = data['language'];
        data = data['data'];

        let data_html = '';

        for(let for_a = 0; for_a < data.length; for_a++) {
            let left = '<a href="/thread/' + data[for_a][0] + '">' + opennamu_xss_filter(data[for_a][1]) + '</a>';
            
            let right = '';
            if(data[for_a][2] === 'O') {
                right += lang['closed'] + ' | ';
            } else if(data[for_a][2] === 'S') {
                right += lang['stop'] + ' | ';
            }

            if(data[for_a][3] !== '') {
                right += lang['agreed_discussion'] + ' | ';
            }

            right += '<a href="/thread/' + data[for_a][0] + '#' + data[for_a][7] + '">#' + data[for_a][7] + '</a> | ';
            right += data[for_a][5] + ' | ';
            right += data[for_a][6];

            data_html += opennamu_make_list(left, right);
        }

        if(data_html !== '') {
            data_html += '<hr class="main_hr">';
        }

        data_html += '<a href="/thread/0/' + doc_name + '">(' + lang['make_new_topic'] + ')</a>';

        data_html += opennamu_page_control('/topic_page/{}/normal/' + doc_name, Number(num), data.length);

        document.getElementById('opennamu_topic_list').innerHTML = data_html;
    });
}