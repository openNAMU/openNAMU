"use strict";

// https://css-tricks.com/how-to-animate-the-details-element/
class Accordion {
    constructor(el) {
        this.el = el;
        this.summary = el.querySelector('summary');
        this.content = el.querySelector('.opennamu_folding');
    
        this.animation = null;
        this.isClosing = false;
        this.isExpanding = false;
        this.summary.addEventListener('click', (e) => this.onClick(e));
    }
  
    onClick(e) {
        e.preventDefault();
        this.el.style.overflow = 'hidden';
        if(this.isClosing || !this.el.open) {
            this.open();
        } else if(this.isExpanding || this.el.open) {
            this.shrink();
        }
    }
  
    shrink() {
        this.isClosing = true;
        
        const startHeight = `${this.el.offsetHeight}px`;
        const endHeight = `${this.summary.offsetHeight}px`;
        
        if(this.animation) {
            this.animation.cancel();
        }
        
        this.animation = this.el.animate({
            height: [startHeight, endHeight]
        }, {
            duration: 200,
            easing: 'ease-out'
        });
        
        this.animation.onfinish = () => this.onAnimationFinish(false);
        this.animation.oncancel = () => this.isClosing = false;
    }
  
    open() {
        this.el.style.height = `${this.el.offsetHeight}px`;
        this.el.open = true;
        window.requestAnimationFrame(() => this.expand());
    }
  
    expand() {
        this.isExpanding = true;
        const startHeight = `${this.el.offsetHeight}px`;
        const endHeight = `${this.summary.offsetHeight + this.content.offsetHeight}px`;
        
        if(this.animation) {
            this.animation.cancel();
        }
        
        this.animation = this.el.animate({
            height: [startHeight, endHeight]
        }, {
            duration: 200,
            easing: 'ease-out'
        });
        this.animation.onfinish = () => this.onAnimationFinish(true);
        this.animation.oncancel = () => this.isExpanding = false;
    }
  
    onAnimationFinish(open) {
        this.el.open = open;
        this.animation = null;
        this.isClosing = false;
        this.isExpanding = false;
        this.el.style.height = this.el.style.overflow = '';
    }
}

window.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('details').forEach((el) => {
        new Accordion(el);
    });
});

function opennamu_do_id_check(data) {
    if(data.match(/\.|\:/)) {
        return 0;
    } else {
        return 1;
    }
}

function opennamu_do_url_encode(data) {
    return encodeURIComponent(data);
}

function opennamu_cookie_split_regex(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function opennamu_get_main_skin_set(set_name) {
    return fetch("/api/setting/" + opennamu_do_url_encode(set_name)).then(function(res) {
        return res.json();
    }).then(function(text) {
        if(
            document.cookie.match(opennamu_cookie_split_regex(set_name)) &&
            document.cookie.match(opennamu_cookie_split_regex(set_name))[1] !== '' &&
            document.cookie.match(opennamu_cookie_split_regex(set_name))[1] !== 'default'
        ) {
            return document.cookie.match(opennamu_cookie_split_regex(set_name))[1];
        } else {
            if(text[set_name]) {
                return text[set_name][0][0];
            } else {
                return '';
            }
        }
    });
}

function opennamu_insert_v(name, data) {
    document.getElementById(name).value = data;
}