"use strict";

// 폐지하고 다시 SSR 방식으로 전환 예정
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
            'ban'
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
                        
                        // 한글 지원 필요
                        let get_data_auth = get_data[name]['auth'];
                        if(get_data_auth === '0') {
                            get_data_auth = lang_data['ip'];
                        } else if(get_data_auth === '1') {
                            get_data_auth = lang_data['member'];
                        } else {
                            get_data_auth = get_data[name]['auth'];
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