"use strict";

function do_insert_user_info() {
    if(document.getElementById('opennamu_get_user_info')) {
        let name = document.getElementById('opennamu_get_user_info').innerHTML;
        let lang_data_list = [
            'user_name',
            'authority',
            'state',
            'member',
            'normal',
            'blocked',
            'type',
            'regex',
            'period',
            'limitless',
            'login_able',
            'why',
            'band_blocked',
            'ip',
            'ban',
            'level'
        ];

        let data_form = new FormData();
        data_form.append('title_list', JSON.stringify(lang_data_list));

        let xhr = new XMLHttpRequest();
        xhr.open("POST", "/api/lang/Test");
        xhr.send(data_form);

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                let lang_data = JSON.parse(this.responseText);
                
                let xhr_2 = new XMLHttpRequest();
                xhr_2.open("POST", "/api/user_info/" + opennamu_do_url_encode(name));
                xhr_2.send();

                xhr_2.onreadystatechange = function() {
                    if(this.readyState === 4 && this.status === 200) {
                        let get_data = JSON.parse(this.responseText);
                        
                        let get_data_auth = get_data[name]['auth'];
                        if(get_data_auth === '0') {
                            get_data_auth = lang_data['ip'];
                        } else if(get_data_auth === '1') {
                            get_data_auth = lang_data['member'];
                        } else {
                            get_data_auth = get_data[name]['auth'];
                        }

                        let get_data_auth_date = get_data[name]['auth_date'];
                        if(get_data_auth_date !== '0') {
                            get_data_auth += ' (~' + get_data_auth_date + ')'
                        }
                        
                        let get_data_ban = get_data[name]['ban'];
                        if(get_data_ban === '0') {
                            get_data_ban = lang_data['normal'];
                        } else {
                            get_data_ban = lang_data['ban'];
                            get_data_ban += '<br>';
                            
                            get_data_ban += lang_data['type'] + ' : ';
                            if(get_data[name]['ban']['type'] === 'normal') {
                                get_data_ban += lang_data['normal']; 
                            } else {
                                get_data_ban += lang_data['regex'];
                            }
                            get_data_ban += '<br>';
                            
                            get_data_ban += lang_data['period'] + ' : ';
                            if(get_data[name]['ban']['period'] === '0') {
                                get_data_ban += lang_data['limitless']; 
                            } else {
                                get_data_ban += get_data[name]['ban']['period'];
                            }
                            get_data_ban += '<br>';
                            
                            get_data_ban += lang_data['login_able'] + ' : ';
                            if(get_data[name]['ban']['login_able'] === '1') {
                                get_data_ban += 'O'; 
                            } else {
                                get_data_ban += 'X';
                            }
                            get_data_ban += '<br>';
                            
                            get_data_ban += lang_data['why'] + ' : ' + get_data[name]['ban']['reason'];
                        }

                        let level = '0';
                        let exp = '0';
                        let max_exp = '0';
                        if(get_data_auth !== lang_data['ip']) {
                            level = get_data[name]['level'];
                            exp = get_data[name]['exp'];
                            max_exp = String(500 + (Number(get_data[name]['level']) * 50));
                        }
                        
                        let data = '' +
                            '<table class="user_info_table">' +
                                '<tr>' +
                                    '<td>' + lang_data['user_name'] + '</td>' +
                                    '<td>' + get_data[name]['render'] + '</td>' +
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
                                    '<td>' + level + ' (' + exp + ' / ' + max_exp + ')</td>' +
                                '</tr>' +
                            '</table>' +
                        '';
                        
                        document.getElementById('opennamu_get_user_info').innerHTML = data;
                    }
                }
            }
        }
    }
}

do_insert_user_info();