"use strict";

function opennamu_send_render(i = 0) {
    let get_class = document.getElementsByClassName('opennamu_send_content')[i];
    if(get_class) {
        opennamu_send_render(i + 1);
        
        let data = get_class.innerHTML;
        if(data === '&lt;br&gt;' || data === '') {
            document.getElementsByClassName('opennamu_send_content')[i].innerHTML = '<br>';
        } else {
            data = data.replace(/javascript:/i, '');
            data = data.replace(/&lt;a(?:(?:(?!&gt;).)*)&gt;((?:(?!&lt;\/a&gt;).)+)&lt;\/a&gt;/g, function(x, x_1) {
                let x_1_org = x_1.replace('&lt;', '<').replace('&gt;', '>');
                
                return '<a href="/w/' + encodeURIComponent(x_1_org) + '">' + x_1 + '</a>';
            });
            
            document.getElementsByClassName('opennamu_send_content')[i].innerHTML = data;
        }
    }
}

let opennamu_send_render_url = [
    '/alarm',
    '/recent_change',
    '/recent_changes'
];
if(opennamu_send_render_url.includes(window.location.pathname)) {
    opennamu_send_render();
}