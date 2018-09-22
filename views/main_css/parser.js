hljs.initHighlightingOnLoad(); 

function folding(num) { 
    var fol = document.getElementById('folding_' + num); 
    if(fol.style.display == 'inline-block' || fol.style.display == 'block') { 
        fol.style.display = 'none';
    } else {
        if(num % 2 == 0) { 
            fol.style.display = 'block'; 
        } else { 
            fol.style.display = 'inline-block'; 
        } 
    } 
}