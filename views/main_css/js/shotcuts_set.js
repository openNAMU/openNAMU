all_list = [
    ['F', '/'],
    ['C', '/recent_changes'],
    ['D', '/recent_discuss'],
    ['A', '/random']
];

all_list.forEach(function(element) {
    shortcut.add(element[0], function() {
        window.location.href = element[1];
    }, {
        'disable_in_input' : true
    }); 
});