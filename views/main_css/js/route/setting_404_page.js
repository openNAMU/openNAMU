"use strict";

function opennamu_setting_404_page_post() {
    let select = document.getElementById("opennamu_setting_404_page_select").value;

    let put_data_select = new FormData();
    put_data_select.append('data', select);
    
    fetch('/api/v2/setting/manage_404_page', {
        method : 'PUT',
        body : put_data_select,
    }).then(function(data) {
        let content = document.getElementById('opennamu_setting_404_page_textarea').value;

        let put_data_content = new FormData();
        put_data_content.append('data', content);
    
        fetch('/api/v2/setting/manage_404_page_content', {
            method : 'PUT',
            body : put_data_content,
        });

        history.go(0);
    });
}

function opennamu_setting_404_page_preview() {
    let content = document.getElementById('opennamu_setting_404_page_textarea').value;
    document.getElementById('opennamu_setting_404_page_preview').innerHTML = content;
}

function opennamu_setting_404_page() {
    let data = [];

    let lang_data = new FormData();
    lang_data.append('data', 'save 404_file 404_page preview');

    fetch('/api/lang', {
        method : 'POST',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(ajax_data) {
        data.push(ajax_data);
        return fetch('/api/v2/setting/manage_404_page');
    }).then(function(res) {
        return res.json();
    }).then(function(ajax_data) {
        data.push(ajax_data);
        return fetch('/api/v2/setting/manage_404_page_content');
    }).then(function(res) {
        return res.json();
    }).then(function(ajax_data) {
        data.push(ajax_data);

        let data_html = '';
        let select_list = [
            ['404_page', data[0]['data'][1]],
            ['404_file', data[0]['data'][2]],
        ];

        data_html += '<select id="opennamu_setting_404_page_select">';
        for(let for_a = 0; for_a < select_list.length; for_a++) {
            let selected = '';
            if(data[1]['data'][0] === select_list[for_a][0]) {
                selected = 'selected';
            }

            data_html += '<option value="' + select_list[for_a][0] + '" ' + selected + '>' + select_list[for_a][1] + '</option>';
        }
        data_html += '</select>';
        data_html += '<hr class="main_hr">';

        let set_data = '';
        if(data[2]['data'].length > 0) {
            set_data = data[2]['data'][0][0];
        }

        data_html += '<textarea class="opennamu_textarea_500" id="opennamu_setting_404_page_textarea">' + set_data + '</textarea>';
        data_html += '<hr class="main_hr">';

        data_html += '<button id="opennamu_save_button" onclick="opennamu_setting_404_page_post();">' + data[0]['data'][0] + '</button> ';
        data_html += '<button onclick="opennamu_setting_404_page_preview();">' + data[0]['data'][3] + '</button>';
        data_html += '<hr class="main_hr">';

        data_html += '<div id="opennamu_setting_404_page_preview"></div>';

        return data_html;
    }).then(function(end_data) {
        document.getElementById('opennamu_setting_404_page').innerHTML = end_data;
    });
}