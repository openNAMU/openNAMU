"use strict";

function opennamu_give_auth_submit() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    if(url_split.length === 3) {
        if(url_split[2] === 'give_total') {
            let auth = document.getElementById('opennamu_give_auth_select_org').value;
            let change_auth = document.getElementById('opennamu_give_auth_select').value;

            let send_data = new FormData();

            send_data.append('auth', auth);
            send_data.append('change_auth', change_auth);

            fetch('/api/v2/auth/give', {
                method : 'PATCH',
                body : send_data,
            }).then(() => {
                history.go(0);
            }).catch(err => {
                console.error(err);
            });;
        } else {
            let change_auth = document.getElementById('opennamu_give_auth_select').value;
            let user_name_data = document.getElementById('opennamu_give_auth_user_name').value;
            user_name_data = user_name_data.replace('\r', '');

            let user_name_arr = user_name_data.split("\n");
            for(let for_a = 0; for_a < user_name_arr.length; for_a++) {
                let send_data = new FormData();

                send_data.append('user_name', user_name_arr[for_a]);
                send_data.append('change_auth', change_auth);

                fetch('/api/v2/auth/give', {
                    method : 'PATCH',
                    body : send_data,
                }).then(() => {
                    history.go(0);
                }).catch(err => {
                    console.error(err);
                });;
            }
        }
    } else {
        let user_name = url_split[3];
        let change_auth = document.getElementById('opennamu_give_auth_select').value;

        let send_data = new FormData();

        send_data.append('user_name', user_name);
        send_data.append('change_auth', change_auth);

        fetch('/api/v2/auth/give', {
            method : 'PATCH',
            body : send_data,
        }).then(() => {
            history.go(0);
        }).catch(err => {
            console.error(err);
        });
    }
}

function opennamu_give_auth() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    let mode = 0;
    if(url_split.length === 3) {
        if(url_split[2] === 'give_total') {
            mode = 1;
        }
    }

    let html_data = '';
    let user_name = '';

    fetch('/api/v2/list/auth').then(function(res) {
        return res.json();
    }).then(function(data) {
        let lang = data["language"];
        let data_list = data["data"];

        if(mode === 0) {
            if(url_split.length === 3) {
                html_data += '<textarea id="opennamu_give_auth_user_name" class="opennamu_textarea_100" placeholder="' + lang["many_delete_help"] + '"></textarea>';
                html_data += '<hr class="main_hr">';
            } else {
                user_name = url_split[3];
    
                html_data += '<div id="opennamu_get_user_info">' + opennamu_xss_filter(user_name) + '</div>';
                html_data += '<hr class="main_hr">';
            }
        } else {
            html_data += '<select id="opennamu_give_auth_select_org">';
            for(let for_a = 0; for_a < data_list.length; for_a++) {
                html_data += '<option value="' + data_list[for_a] + '">' + data_list[for_a] + '</option>';
            }
            html_data += '</select>';
            html_data += '<hr class="main_hr">';
        }

        html_data += '<select id="opennamu_give_auth_select">';
        for(let for_a = 0; for_a < data_list.length; for_a++) {
            html_data += '<option value="' + data_list[for_a] + '">' + data_list[for_a] + '</option>';
        }
        html_data += '</select>';
        html_data += '<hr class="main_hr">';

        html_data += '<button onclick="opennamu_give_auth_submit();">' + lang["send"] + '</button>';

        document.getElementById('opennamu_give_auth').innerHTML = html_data;

        if(user_name !== '') {
            fetch('/api/v2/auth/' + opennamu_do_url_encode(user_name)).then(function(res) {
                return res.json();
            }).then(function(data) {
                document.getElementById('opennamu_give_auth_select').value = data["name"];
            });
        }

        do_insert_user_info();
    });
}