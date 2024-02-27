"use strict";

function do_insert_user_info() {
    if(document.getElementById('opennamu_get_user_info')) {
        let name = document.getElementById('opennamu_get_user_info').innerHTML;

        fetch("/api/user_info/" + opennamu_do_url_encode(name)).then(function(res) {
            return res.json();
        }).then(function(data) {
            let lang_data = data["language"];

            let get_data_auth = data['data']['auth'];
            if(get_data_auth === '0') {
                get_data_auth = lang_data['ip'];
            } else if(get_data_auth === '1') {
                get_data_auth = lang_data['member'];
            } else {
                get_data_auth = data['data']['auth'];
            }

            let get_data_auth_date = data['data']['auth_date'];
            if(get_data_auth_date !== '0') {
                get_data_auth += ' (~' + get_data_auth_date + ')'
            }
            
            let get_data_ban = data['data']['ban'];
            if(get_data_ban === '0') {
                get_data_ban = lang_data['normal'];
            } else {
                get_data_ban = lang_data['ban'];
                get_data_ban += '<br>';
                
                get_data_ban += lang_data['type'] + ' : ';
                if(data['data']['ban']['type'] === 'normal') {
                    get_data_ban += '<a href="/block_log/user/' + opennamu_do_url_encode(name) + '">' + lang_data['normal'] + '</a>'; 
                } else {
                    get_data_ban += '<a href="/block_log/regex">' + lang_data['regex'] + '</a>';
                }
                get_data_ban += '<br>';
                
                get_data_ban += lang_data['period'] + ' : ';
                if(data['data']['ban']['period'] === '0') {
                    get_data_ban += '~ ' + lang_data['limitless']; 
                } else {
                    get_data_ban += '~ ' + data['data']['ban']['period'];
                }
                get_data_ban += '<br>';
                
                get_data_ban += lang_data['option'] + ' : ';
                if(data['data']['ban']['login_able'] === '2') {
                    get_data_ban += lang_data['edit_request_able']; 
                } else if(data['data']['ban']['login_able'] === '1') {
                    get_data_ban += lang_data['login_able']; 
                } else {
                    get_data_ban += 'X';
                }
                get_data_ban += '<br>';
                
                if(data['data']['ban']['reason'] === 'edit filter') {
                    get_data_ban += lang_data['why'] + ' : <a href="/edit_filter/' + opennamu_do_url_encode(name) + '">' + data['data']['ban']['reason'] + '</a>';
                } else {
                    get_data_ban += lang_data['why'] + ' : ' + data['data']['ban']['reason'];
                }
            }
            
            let end_data = '' +
                '<table class="user_info_table">' +
                    '<tr>' +
                        '<td>' + lang_data['user_name'] + '</td>' +
                        '<td>' + data['data']['render'] + '</td>' +
                    '</tr>' +
                    '<tr>' +
                        '<td>' + lang_data['authority'] + '</td>' +
                        '<td>' + get_data_auth + '</td>' +
                    '</tr>' +
                    '<tr>' +
                        '<td>' + lang_data['state'] + '</td>' +
                        '<td>' + get_data_ban + '</td>' +
                    '</tr>' +
                    '<tr>' +
                        '<td>' + lang_data['level'] + '</td>' +
                        '<td>' + data['data']['level'] + ' (' + data['data']['exp'] + ' / ' + data['data']['max_exp'] + ')</td>' +
                    '</tr>' +
                '</table>' +
            '';
            
            document.getElementById('opennamu_get_user_info').innerHTML = end_data;
        });
    }
}

do_insert_user_info();