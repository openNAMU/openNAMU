"use strict";

function opennamu_xss_filter(str) {
    return str.replace(/[&<>"']/g, function(match) {
        switch(match) {
            case '&':
                return '&amp;';
            case '<':
                return '&lt;';
            case '>':
                return '&gt;';
            case "'":
                return '&#x27;';
            case '"':
                return '&quot;';
        }
    });
}

function opennamu_list_hidden_remove() {
    const style = document.querySelector('#opennamu_list_hidden_style');
    if(style !== null) {
        if(style.innerHTML !== "") {
            style.innerHTML = '';
        } else {
            style.innerHTML = '.opennamu_list_hidden { display: none; }';
        }
    }
}

function opennamu_do_footnote_popover(set_name, load_name, sub_obj = undefined, do_type = 'open') {
    if(document.getElementById(set_name + '_load')) {
        document.getElementById(set_name + '_load').onclick = function (event) {
            event.stopPropagation();
        };
        if(do_type === 'open') {
            if(sub_obj !== undefined) {
                document.getElementById(set_name + '_load').innerHTML = document.getElementById(sub_obj).innerHTML;
            } else {
                document.getElementById(set_name).title = '';
                document.getElementById(set_name + '_load').innerHTML = '<a href="#' + load_name + '">(Go)</a> ' + document.getElementById(load_name + '_title').innerHTML;
            }
            document.getElementById(set_name + '_load').style.display = "inline-block";
            document.getElementById(set_name + '_load').count = 0;

            let width = document.getElementById(set_name + '_load').clientWidth;
            let screen_width = window.innerWidth;
            let left = document.getElementById(set_name).getBoundingClientRect().left;
            let top = window.pageYOffset + document.getElementById(set_name).getBoundingClientRect().top;

            document.getElementById(set_name + '_load').style.top = String(top) + "px";
            if(screen_width - (left + width) < 50) {
                if(left > 350) {
                    document.getElementById(set_name + '_load').style.left = String(left - 300) + "px";
                } else {
                    document.getElementById(set_name + '_load').style.left = "0px";
                }

                left = document.getElementById(set_name + '_load').getBoundingClientRect().left;
                width = document.getElementById(set_name + '_load').clientWidth;
                if(300 > width) {
                    document.getElementById(set_name + '_load').style.left = String(left + (300 - width)) + "px";
                } else {
                    document.getElementById(set_name + '_load').style.marginTop = "20px";
                }
            }
        } else {
            if(document.getElementById(set_name + '_load').count === 1) {
                document.getElementById(set_name + '_load').style.display = "none";
            } else {
                document.getElementById(set_name + '_load').count = 1;
            }
        }
    }
}

document.addEventListener("click", function () {
    document.querySelectorAll('span[id$="_over"] > .opennamu_popup_footnote').forEach(function (obj) {
        opennamu_do_footnote_popover(obj.id.slice(0, -5), '', undefined, 'close');
    });
});

