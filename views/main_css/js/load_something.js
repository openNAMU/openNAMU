function do_twofa_check(init = 0) {
    let twofa_option = document.getElementById('twofa_check_input');
    let twofa_option_num = twofa_option.options.selectedIndex;
    let twofa_select_data = twofa_option.options[twofa_option_num].value;
    
    if(twofa_select_data === 'on') {
        document.getElementById('fa_plus_content').style.display = "block";
    } else {
        document.getElementById('fa_plus_content').style.display = "none";
    }
}

function send_render(i = 0) {
    let get_class = document.getElementsByClassName('send_content')[i];
    if(get_class) {
        send_render(i + 1);
        
        let data = get_class.innerHTML;
        if(data === '&lt;br&gt;') {
            document.getElementsByClassName('send_content')[i].innerHTML = '<br>';
        } else {
            data = data.replace(/javascript:/i, '');
            data = data.replace(/&lt;a(?:(?:(?!&gt;).)*)&gt;((?:(?!&lt;\/a&gt;).)+)&lt;\/a&gt;/g, function(x, x_1) {
                x_1_org = x_1.replace('&lt;', '<').replace('&gt;', '>');
                return '<a href="/w/' + encodeURIComponent(x_1_org) + '">' + x_1 + '</a>';
            });
            
            document.getElementsByClassName('send_content')[i].innerHTML = data;
        }
    }
}