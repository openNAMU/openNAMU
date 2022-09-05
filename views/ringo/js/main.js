let ringo_save_data = '';
let ringo_open = 0;

function ringo_opening(data) {
    ringo_save_data = data;

    console.log(data);
    let element = [data];
    
    let menu_list = [
        'recent_cel_in',
        'other_cel_in',
        'user_cel_in'
    ];
    for(for_a in menu_list) {
        if(menu_list[for_a] !== data) {
            element.push(menu_list[for_a]);
        }
    }

    console.log(element);

    if(document.getElementById(element[0]).style.display == 'none') {
        document.getElementById(element[0]).style.display = 'block';

        for(for_a in element) {
            if(for_a !== '0') {
                console.log(for_a);
                document.getElementById(element[for_a]).style.display = 'none';
            }
        }
    } else {
        document.getElementById(element[0]).style.display = 'none';
    }

    ringo_open = 1;
    setTimeout(function() { ringo_open = 0; }, 100);
}

document.onclick = function(event) {
    if(ringo_save_data !== '' && ringo_open == 0) {
        document.getElementById(ringo_save_data).style.display = 'none';
    }
}