"use strict";

function opennamu_do_render_simple(name_ele) {
    let skin_set_data = document.getElementById(name_ele).innerHTML;
    
    // 목차 구현
    let toc_all_data = '<div id="toc"><span id="toc_title">TOC</span><br>';
    let split_toc;
    let toc_data;
    
    let i = 1;
    while(1) {
        toc_data = skin_set_data.match(/<h[1-6]>([^<>]+)<\/h[1-6]>/);
        if(toc_data) {
            split_toc = toc_data[1].match(/^([^ ]+)(.+)/);
            toc_all_data += '' + 
                '<br>' +
                '<span style="margin-left: ' + String(((toc_data[1].match(/\./g) || []).length - 1) * 10) + 'px;">' +
                    '<a href="#toc_' + String(i) + '">' + split_toc[1] + '</a>' + split_toc[2] +
                '</span>' +
            '';

            skin_set_data = skin_set_data.replace(
                /<(h[1-6])>([^<>]+)<\/h[1-6]>/, 
                '<$1 id="toc_' + String(i) + '"><a href="#toc">' + split_toc[1] + '</a>' + split_toc[2] + '</$1>'
            );
            
            i += 1;
        } else {
            break;
        }
    }
    
    skin_set_data = toc_all_data + '</div>' + skin_set_data;

    // 각주 구현
    var note_list = {};
    var plus_note;
    
    i = 1;
    while(1) {
        toc_data = skin_set_data.match(/<sup>([^<>]+)<\/sup>/);
        if(toc_data) {
            if(!note_list[toc_data[1]]) {
                note_list[toc_data[1]] = [String(i), 0];
            } else {
                note_list[toc_data[1]][1] += 1;
            }

            if(note_list[toc_data[1]][1] != 0) {
                plus_note = '_' + String(note_list[toc_data[1]][1]);
            } else {
                plus_note = '';
            }
            
            skin_set_data = skin_set_data.replace(
                /<sup>([^<>]+)<\/sup>/, 
                '<sup><a id="note_' + note_list[toc_data[1]][0] + plus_note + '" href="#note_' + note_list[toc_data[1]][0] + '_end">$1</a></sup>'
            );
            
            i += 1;
        } else {
            break;
        }
    }
    
    document.getElementById(name_ele).innerHTML = skin_set_data;
}

let opennamu_do_render_simple_url = [
    '/manager/1',
    '/manager',
    '/other',
    '/setting/phrase',
    '/setting/main',
    '/setting/external'
];
if(opennamu_do_render_simple_url.includes(window.location.pathname)) {
    opennamu_do_render_simple('opennamu_simple_render');
}