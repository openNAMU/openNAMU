"use strict";

function do_insert_user_info() {
    if(document.getElementById('opennamu_get_user_info')) {
        let name = document.getElementById('opennamu_get_user_info').innerHTML;
        
        let xhr = new XMLHttpRequest();
        xhr.open("POST", "/api/user_info/" + opennamu_do_url_encode(name));
        xhr.send();

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                let get_data = JSON.parse(this.responseText);
                
                // 한글 지원 필요
                let get_data_auth = get_data[name]['auth'];
                if(get_data_auth === '0') {
                    get_data_auth = 'IP';
                } else if(get_data_auth === '1') {
                    get_data_auth = 'USER';
                } else {
                    get_data_auth = get_data[name]['auth'];
                }
                
                let get_data_ban = get_data[name]['ban'];
                if(get_data_ban === '0') {
                    get_data_ban = 'NORMAL';
                } else {
                    get_data_ban = 'BAN';
                    get_data_ban += '<br>';
                    
                    get_data_ban += 'TYPE : ';
                    if(get_data[name]['ban']['type'] === 'normal') {
                        get_data_ban += 'NORMAL'; 
                    } else {
                    	get_data_ban += 'REGEX';
                    }
                    get_data_ban += '<br>';
                    
                    get_data_ban += 'PERIOD : ';
                    if(get_data[name]['ban']['period'] === '0') {
                        get_data_ban += 'INF'; 
                    } else {
                    	get_data_ban += get_data[name]['ban']['period'];
                    }
                    get_data_ban += '<br>';
                    
                    get_data_ban += 'LOGIN ABLE : ';
                    if(get_data[name]['ban']['login_able'] === '1') {
                        get_data_ban += 'YES'; 
                    } else {
                    	get_data_ban += 'NO';
                    }
                    get_data_ban += '<br>';
                    
                    get_data_ban += 'REASON : ' + get_data[name]['ban']['reason'];
                }
                
                let data = '' +
                    '<table class="user_info_table">' +
                        '<tr>' +
                            '<td>NAME</td>' +
                            '<td>' + get_data[name]['render'] + '</td>' +
                        '</tr>' +
                        '<tr>' +
                            '<td>AUTH</td>' +
                            '<td>' + get_data_auth + '</td>' +
                        '</tr>' +
                        '<tr>' +
                            '<td>STATE</td>' +
                            '<td>' + get_data_ban + '</td>' +
                        '</tr>' +
                    '</table>' +
                '';
                
                document.getElementById('opennamu_get_user_info').innerHTML = data;
                
                opennamu_do_ip_parser();
            }
        }
    }
}

do_insert_user_info();