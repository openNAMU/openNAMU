function opennamu_user_rankup() {
    fetch('/api/v2/user/rankup').then(function(res) {
        return res.json();
    }).then(function(data) {
        console.log(data);
    });
}