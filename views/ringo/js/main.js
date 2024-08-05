"use strict";

let ringo_save_data = '';
let ringo_open = 0;
let ringo_menu_list = [
    'recent_cel',
    'other_cel',
    'user_cel',
    'add_cel'
];

function ringo_opening(data) {
    let element = [data];
    
    for(let for_a in ringo_menu_list) {
        if(ringo_menu_list[for_a] + '_in' !== data) {
            element.push(ringo_menu_list[for_a] + '_in');
        }
    }

    if((document.getElementById(element[0]).style.display === 'none' && ringo_open === 0) || ringo_save_data !== data) {
        document.getElementById(element[0]).style.display = 'block';

        for(let for_a in element) {
            if(for_a !== '0') { 
                if(document.getElementById(element[for_a]) !== null) {
                    document.getElementById(element[for_a]).style.display = 'none';
                }
            }
        }

        ringo_open = 1;
        ringo_save_data = data;

        setTimeout(function() { ringo_open = 2; }, 100);
    } else {
        document.getElementById(element[0]).style.display = 'none';

        ringo_open = 0
    }
}

document.addEventListener("click", function() {
    let cel_list = [];
    for(let for_a in ringo_menu_list) {
        cel_list.push(document.getElementById(ringo_menu_list[for_a]));
    }

    if(ringo_save_data !== '' && ringo_open === 2) {
        document.getElementById(ringo_save_data).style.display = 'none';

        setTimeout(function() { ringo_open = 0; }, 100);
    }
});